package webhook

import (
	"fmt"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

var allowedWebhookTypes = map[string]bool{
	"gitea": true, "slack": true, "discord": true, "dingtalk": true, "telegram": true,
	"msteams": true, "feishu": true, "matrix": true, "jianmu": true, "softbot": true,
}

var allowedWebhookContentTypes = map[string]bool{"json": true, "form": true}
var allowedWebhookMethods = map[string]bool{"GET": true, "POST": true}

var allowedWebhookEvents = map[string]bool{
	"push": true, "create": true, "delete": true,
	"issues_only": true, "issue_assign": true, "issue_label": true, "issue_comment": true,
	"pull_request_only": true, "pull_request_assign": true, "pull_request_comment": true,
}

// Shortcuts returns webhook management shortcuts.
func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: "List repository webhooks",
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", webhookPath(ctx), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "view",
			Description: "View webhook details",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Webhook ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", webhookItemPath(ctx, id), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "create",
			Description: "Create a repository webhook",
			Flags: []common.Flag{
				{Name: "url", Short: "u", Usage: "Webhook target URL", Required: true},
				{Name: "events", Short: "e", Usage: "Comma-separated events, for example: push,issues_only", Required: true},
				{Name: "type", Short: "t", Usage: "Webhook type: gitea/slack/discord/dingtalk/telegram/msteams/feishu/matrix/jianmu/softbot", Default: "gitea"},
				{Name: "content-type", Usage: "Payload content type: json or form", Default: "json"},
				{Name: "method", Short: "m", Usage: "HTTP method: POST or GET", Default: "POST"},
				{Name: "secret", Short: "s", Usage: "Webhook secret"},
				{Name: "branch-filter", Usage: "Branch glob filter for push/create/delete events", Default: "*"},
				{Name: "active", Usage: "Whether the webhook is active: true or false", Default: "true"},
			},
			Run: runCreate,
		},
		{
			Name:        "update",
			Description: "Update a repository webhook while preserving unspecified fields when available",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Webhook ID", Required: true},
				{Name: "url", Short: "u", Usage: "Webhook target URL"},
				{Name: "events", Short: "e", Usage: "Comma-separated events, for example: push,issues_only"},
				{Name: "content-type", Usage: "Payload content type: json or form"},
				{Name: "method", Short: "m", Usage: "HTTP method: POST or GET"},
				{Name: "secret", Short: "s", Usage: "Webhook secret. Pass it again if the server does not return existing secrets."},
				{Name: "branch-filter", Usage: "Branch glob filter for push/create/delete events"},
				{Name: "active", Usage: "Whether the webhook is active: true or false"},
			},
			Run: runUpdate,
		},
		{
			Name:        "delete",
			Description: "Delete a repository webhook",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Webhook ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("DELETE", webhookItemPath(ctx, id), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "test",
			Description: "Trigger a test delivery for a webhook",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Webhook ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("POST", fmt.Sprintf("%s/tests", webhookItemPath(ctx, id)), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "tasks",
			Description: "List webhook delivery tasks",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Webhook ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", fmt.Sprintf("%s/hooktasks", webhookItemPath(ctx, id)), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
	}
}

func runCreate(ctx *common.RuntimeContext) error {
	if err := ctx.ResolveOwnerRepo(); err != nil {
		return err
	}
	payload, err := webhookPayloadFromArgs(ctx, nil)
	if err != nil {
		return err
	}
	env, err := ctx.CallAPI("POST", webhookPath(ctx), payload)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func runUpdate(ctx *common.RuntimeContext) error {
	if err := ctx.ResolveOwnerRepo(); err != nil {
		return err
	}
	id, err := ctx.RequireArg("id")
	if err != nil {
		return err
	}

	current, err := fetchWebhook(ctx, id)
	if err != nil {
		return fmt.Errorf("fetch webhook: %w", err)
	}
	payload, err := webhookPayloadFromArgs(ctx, current)
	if err != nil {
		return err
	}
	env, err := ctx.CallAPI("PUT", webhookItemPath(ctx, id), payload)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func webhookPath(ctx *common.RuntimeContext) string {
	return fmt.Sprintf("/v1/%s/%s/webhooks", ctx.Owner, ctx.Repo)
}

func webhookItemPath(ctx *common.RuntimeContext, id string) string {
	return fmt.Sprintf("%s/%s", webhookPath(ctx), id)
}

func fetchWebhook(ctx *common.RuntimeContext, id string) (map[string]interface{}, error) {
	env, err := ctx.CallAPI("GET", webhookItemPath(ctx, id), nil)
	if err != nil {
		return nil, err
	}
	data, ok := env.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse webhook data")
	}
	return data, nil
}

func webhookPayloadFromArgs(ctx *common.RuntimeContext, current map[string]interface{}) (map[string]interface{}, error) {
	url := firstNonEmpty(ctx.Arg("url"), stringFromMap(current, "url"))
	if url == "" {
		return nil, fmt.Errorf("required flag --url is missing")
	}

	eventValue := ctx.Arg("events")
	var events []string
	var err error
	if eventValue != "" {
		events, err = parseWebhookEvents(eventValue)
		if err != nil {
			return nil, err
		}
	} else {
		events, err = eventsFromMap(current)
		if err != nil {
			return nil, err
		}
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("required flag --events is missing")
	}

	webhookType := strings.ToLower(firstNonEmpty(ctx.Arg("type"), stringFromMap(current, "type"), "gitea"))
	if err := validateOneOf("type", webhookType, allowedWebhookTypes); err != nil {
		return nil, err
	}
	contentType := strings.ToLower(firstNonEmpty(ctx.Arg("content-type"), stringFromMap(current, "content_type"), "json"))
	if err := validateOneOf("content-type", contentType, allowedWebhookContentTypes); err != nil {
		return nil, err
	}
	method := strings.ToUpper(firstNonEmpty(ctx.Arg("method"), stringFromMap(current, "http_method"), "POST"))
	if err := validateOneOf("method", method, allowedWebhookMethods); err != nil {
		return nil, err
	}
	branchFilter := firstNonEmpty(ctx.Arg("branch-filter"), stringFromMap(current, "branch_filter"), "*")
	active, err := activeFromArgs(ctx.Arg("active"), current)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"type":          webhookType,
		"active":        active,
		"content_type":  contentType,
		"http_method":   method,
		"url":           url,
		"branch_filter": branchFilter,
		"events":        events,
	}
	if secret := firstNonEmpty(ctx.Arg("secret"), stringFromMap(current, "secret")); secret != "" {
		payload["secret"] = secret
	}
	return payload, nil
}

func parseWebhookEvents(value string) ([]string, error) {
	parts := strings.Split(value, ",")
	events := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, part := range parts {
		event := strings.TrimSpace(part)
		if event == "" {
			continue
		}
		if !allowedWebhookEvents[event] {
			return nil, fmt.Errorf("invalid --events value %q", event)
		}
		if seen[event] {
			continue
		}
		seen[event] = true
		events = append(events, event)
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("required flag --events is missing")
	}
	return events, nil
}

func eventsFromMap(values map[string]interface{}) ([]string, error) {
	if values == nil {
		return nil, nil
	}
	raw, ok := values["events"]
	if !ok || raw == nil {
		return nil, nil
	}
	switch events := raw.(type) {
	case []interface{}:
		result := make([]string, 0, len(events))
		for _, event := range events {
			name, ok := event.(string)
			if !ok {
				return nil, fmt.Errorf("failed to parse webhook events")
			}
			result = append(result, name)
		}
		return result, nil
	case []string:
		return events, nil
	default:
		return nil, fmt.Errorf("failed to parse webhook events")
	}
}

func activeFromArgs(value string, current map[string]interface{}) (bool, error) {
	if value != "" {
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, fmt.Errorf("invalid --active value %q: use true or false", value)
		}
	}
	if current != nil {
		if active, ok := current["active"].(bool); ok {
			return active, nil
		}
		if active, ok := current["is_active"].(bool); ok {
			return active, nil
		}
	}
	return true, nil
}

func stringFromMap(values map[string]interface{}, key string) string {
	if values == nil {
		return ""
	}
	value, _ := values[key].(string)
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func validateOneOf(name, value string, allowed map[string]bool) error {
	if allowed[value] {
		return nil
	}
	return fmt.Errorf("invalid --%s value %q", name, value)
}
