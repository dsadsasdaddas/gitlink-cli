package mirror

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

// Shortcuts returns repository mirror management shortcuts.
func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "create",
			Description: "Create a GitLink mirror repository from a remote clone URL",
			Flags: []common.Flag{
				{Name: "user-id", Short: "u", Usage: "Target user or organization ID", Required: true},
				{Name: "name", Short: "n", Usage: "Project display name", Required: true},
				{Name: "repository-name", Short: "r", Usage: "Repository identifier", Required: true},
				{Name: "clone-addr", Short: "c", Usage: "Remote clone URL", Required: true},
				{Name: "description", Short: "d", Usage: "Project description"},
				{Name: "category-id", Usage: "Project category ID"},
				{Name: "language-id", Usage: "Project language ID"},
				{Name: "private", Usage: "Create as a private project", Bool: true, Default: "false"},
				{Name: "auth-username", Usage: "Username for private mirror source"},
				{Name: "auth-password", Usage: "Password or token for private mirror source"},
				{Name: "dry-run", Usage: "Preview the request without creating the mirror", Bool: true, Default: "false"},
			},
			Run: runCreate,
		},
		{
			Name:        "sync",
			Description: "Trigger manual synchronization for a mirror repository",
			Flags: []common.Flag{
				{Name: "repo-id", Short: "i", Usage: "GitLink repository ID", Required: true},
				{Name: "dry-run", Usage: "Preview the request without synchronizing the mirror", Bool: true, Default: "false"},
			},
			Run: runSync,
		},
	}
}

func runCreate(ctx *common.RuntimeContext) error {
	payload, err := mirrorPayload(ctx)
	if err != nil {
		return err
	}
	if parseDryRun(ctx.Arg("dry-run")) {
		return ctx.OutputData(map[string]interface{}{
			"dry_run": true,
			"method":  "POST",
			"path":    mirrorCreatePath(),
			"payload": redactMirrorPayload(payload),
		})
	}
	env, err := ctx.CallAPI("POST", mirrorCreatePath(), payload)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func runSync(ctx *common.RuntimeContext) error {
	repoID, err := parsePositiveInt("repo-id", ctx.Arg("repo-id"))
	if err != nil {
		return err
	}
	path := mirrorSyncPath(repoID)
	if parseDryRun(ctx.Arg("dry-run")) {
		return ctx.OutputData(map[string]interface{}{
			"dry_run": true,
			"method":  "POST",
			"path":    path,
			"repo_id": repoID,
		})
	}
	env, err := ctx.CallAPI("POST", path, nil)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func mirrorPayload(ctx *common.RuntimeContext) (map[string]interface{}, error) {
	userID, err := parsePositiveInt("user-id", ctx.Arg("user-id"))
	if err != nil {
		return nil, err
	}
	name, err := requireTrimmed(ctx, "name")
	if err != nil {
		return nil, err
	}
	repositoryName, err := requireTrimmed(ctx, "repository-name")
	if err != nil {
		return nil, err
	}
	cloneAddr, err := requireTrimmed(ctx, "clone-addr")
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"user_id":         userID,
		"name":            name,
		"repository_name": repositoryName,
		"clone_addr":      cloneAddr,
		"is_mirror":       true,
		"private":         parseDryRun(ctx.Arg("private")),
	}
	if description := strings.TrimSpace(ctx.Arg("description")); description != "" {
		payload["description"] = description
	}
	if categoryID := strings.TrimSpace(ctx.Arg("category-id")); categoryID != "" {
		id, err := parsePositiveInt("category-id", categoryID)
		if err != nil {
			return nil, err
		}
		payload["project_category_id"] = id
	}
	if languageID := strings.TrimSpace(ctx.Arg("language-id")); languageID != "" {
		id, err := parsePositiveInt("language-id", languageID)
		if err != nil {
			return nil, err
		}
		payload["project_language_id"] = id
	}
	authUsername := strings.TrimSpace(ctx.Arg("auth-username"))
	authPassword := strings.TrimSpace(ctx.Arg("auth-password"))
	if authPassword != "" && authUsername == "" {
		return nil, fmt.Errorf("--auth-username is required when --auth-password is provided")
	}
	if authUsername != "" {
		payload["auth_username"] = authUsername
	}
	if authPassword != "" {
		payload["auth_password"] = authPassword
	}
	return payload, nil
}

func requireTrimmed(ctx *common.RuntimeContext, name string) (string, error) {
	value, err := ctx.RequireArg(name)
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("required flag --%s is missing", name)
	}
	return value, nil
}

func mirrorCreatePath() string {
	return "/projects/migrate"
}

func mirrorSyncPath(repoID int) string {
	return fmt.Sprintf("/repositories/%s/sync_mirror", url.PathEscape(strconv.Itoa(repoID)))
}

func parsePositiveInt(name, value string) (int, error) {
	value = strings.TrimSpace(value)
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid --%s value %q", name, value)
	}
	return id, nil
}

func parseDryRun(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}

func redactMirrorPayload(payload map[string]interface{}) map[string]interface{} {
	redacted := make(map[string]interface{}, len(payload))
	for k, v := range payload {
		if k == "auth_password" {
			redacted[k] = "***REDACTED***"
			continue
		}
		redacted[k] = v
	}
	return redacted
}
