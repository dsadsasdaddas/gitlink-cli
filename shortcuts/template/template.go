package template

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: "List repository project templates",
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", templatePath(ctx), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "view",
			Description: "Show project template details",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Project template ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", templateItemPath(ctx, id), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "create",
			Description: "Create a project template",
			Flags: []common.Flag{
				{Name: "type", Short: "t", Usage: "Template type, for example: ProjectTemplates::Issue", Default: "ProjectTemplates::Issue"},
				{Name: "name", Short: "n", Usage: "Template name", Required: true},
				{Name: "content", Short: "c", Usage: "Template content"},
				{Name: "content-file", Usage: "Path to a template content file"},
				{Name: "dry-run", Usage: "Preview the create request without changing templates", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				name, err := ctx.RequireArg("name")
				if err != nil {
					return err
				}
				content, err := templateContentFromArgs(ctx)
				if err != nil {
					return err
				}
				payload := map[string]interface{}{
					"type":    firstTemplateValue(ctx.Arg("type"), "ProjectTemplates::Issue"),
					"name":    name,
					"content": content,
				}
				path := templatePath(ctx)
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(templateDryRun(ctx, "create_project_template", "POST", path, payload))
				}
				env, err := ctx.CallAPI("POST", path, payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "update",
			Description: "Update a project template while preserving unspecified fields",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Project template ID", Required: true},
				{Name: "type", Short: "t", Usage: "Template type"},
				{Name: "name", Short: "n", Usage: "Template name"},
				{Name: "content", Short: "c", Usage: "Template content"},
				{Name: "content-file", Usage: "Path to a template content file"},
				{Name: "dry-run", Usage: "Preview the update request without changing templates", Bool: true, Default: "false"},
			},
			Run: runUpdate,
		},
		{
			Name:        "delete",
			Description: "Delete a project template",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Project template ID", Required: true},
				{Name: "dry-run", Usage: "Preview the delete request without changing templates", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				path := templateItemPath(ctx, id)
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(templateDryRun(ctx, "delete_project_template", "DELETE", path, nil))
				}
				env, err := ctx.CallAPI("DELETE", path, nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
	}
}

func runUpdate(ctx *common.RuntimeContext) error {
	if err := ctx.ResolveOwnerRepo(); err != nil {
		return err
	}
	id, err := ctx.RequireArg("id")
	if err != nil {
		return err
	}
	if !hasTemplateUpdateArgs(ctx) {
		return fmt.Errorf("at least one of --type, --name, --content, or --content-file is required")
	}
	if ctx.Arg("content") != "" && ctx.Arg("content-file") != "" {
		return fmt.Errorf("use only one of --content or --content-file")
	}
	current, err := fetchTemplate(ctx, id)
	if err != nil {
		return fmt.Errorf("fetch project template: %w", err)
	}
	payload, err := templateUpdatePayload(ctx, current)
	if err != nil {
		return err
	}
	path := templateItemPath(ctx, id)
	if ctx.Arg("dry-run") == "true" {
		return ctx.OutputData(templateDryRun(ctx, "update_project_template", "PUT", path, payload))
	}
	env, err := ctx.CallAPI("PUT", path, payload)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func templatePath(ctx *common.RuntimeContext) string {
	return "/v1" + ctx.RepoPath() + "/project_templates"
}

func templateItemPath(ctx *common.RuntimeContext, id string) string {
	return fmt.Sprintf("%s/%s", templatePath(ctx), url.PathEscape(id))
}

func fetchTemplate(ctx *common.RuntimeContext, id string) (map[string]interface{}, error) {
	env, err := ctx.CallAPI("GET", templateItemPath(ctx, id), nil)
	if err != nil {
		return nil, err
	}
	data, ok := env.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse project template response")
	}
	if template, ok := data["project_template"].(map[string]interface{}); ok {
		return template, nil
	}
	return data, nil
}

func templateUpdatePayload(ctx *common.RuntimeContext, current map[string]interface{}) (map[string]interface{}, error) {
	content := templateString(current, "content")
	if ctx.Arg("content") != "" || ctx.Arg("content-file") != "" {
		value, err := templateContentFromArgs(ctx)
		if err != nil {
			return nil, err
		}
		content = value
	}
	payload := map[string]interface{}{
		"type":    firstTemplateValue(ctx.Arg("type"), templateString(current, "type")),
		"name":    firstTemplateValue(ctx.Arg("name"), templateString(current, "name")),
		"content": content,
	}
	if payload["type"] == "" {
		return nil, fmt.Errorf("required template type is missing; pass --type")
	}
	if payload["name"] == "" {
		return nil, fmt.Errorf("required template name is missing; pass --name")
	}
	return payload, nil
}

func templateContentFromArgs(ctx *common.RuntimeContext) (string, error) {
	content := ctx.Arg("content")
	contentFile := strings.TrimSpace(ctx.Arg("content-file"))
	if content != "" && contentFile != "" {
		return "", fmt.Errorf("use only one of --content or --content-file")
	}
	if contentFile != "" {
		data, err := os.ReadFile(contentFile)
		if err != nil {
			return "", fmt.Errorf("read --content-file: %w", err)
		}
		content = string(data)
	}
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("required flag --content or --content-file is missing")
	}
	return content, nil
}

func templateDryRun(ctx *common.RuntimeContext, action, method, path string, payload interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"repository": fmt.Sprintf("%s/%s", ctx.Owner, ctx.Repo),
		"dry_run":    true,
		"action":     action,
		"method":     method,
		"path":       path,
	}
	if payload != nil {
		result["payload"] = payload
	}
	return result
}

func hasTemplateUpdateArgs(ctx *common.RuntimeContext) bool {
	for _, name := range []string{"type", "name", "content", "content-file"} {
		if ctx.Arg(name) != "" {
			return true
		}
	}
	return false
}

func templateString(values map[string]interface{}, key string) string {
	if values == nil {
		return ""
	}
	value, _ := values[key].(string)
	return value
}

func firstTemplateValue(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
