package access

import (
	"fmt"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

var allowedRoles = map[string]string{
	"manager":   "manager",
	"developer": "developer",
	"reporter":  "reporter",
}

// Shortcuts returns self-service project access shortcuts.
func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "join",
			Description: "Apply to join a project with an invitation code",
			Flags: []common.Flag{
				{Name: "code", Short: "c", Usage: "Project invitation or application code", Required: true},
				{Name: "role", Short: "r", Usage: "Requested role: manager, developer, or reporter", Default: "developer"},
				{Name: "dry-run", Usage: "Preview the request without submitting the join application", Bool: true, Default: "false"},
			},
			Run: runJoin,
		},
		{
			Name:        "quit",
			Description: "Quit the current repository project",
			Flags: []common.Flag{
				{Name: "dry-run", Usage: "Preview the request without leaving the project", Bool: true, Default: "false"},
			},
			Run: runQuit,
		},
	}
}

func runJoin(ctx *common.RuntimeContext) error {
	code, err := requireTrimmed(ctx, "code")
	if err != nil {
		return err
	}
	role, err := normalizeRole(ctx.Arg("role"))
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"applied_project": map[string]interface{}{
			"code": code,
			"role": role,
		},
	}
	path := joinPath()
	if parseDryRun(ctx.Arg("dry-run")) {
		return ctx.OutputData(map[string]interface{}{
			"dry_run": true,
			"method":  "POST",
			"path":    path,
			"payload": payload,
		})
	}
	env, err := ctx.CallAPI("POST", path, payload)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func runQuit(ctx *common.RuntimeContext) error {
	if err := ctx.ResolveOwnerRepo(); err != nil {
		return err
	}
	path := quitPath(ctx)
	if parseDryRun(ctx.Arg("dry-run")) {
		return ctx.OutputData(map[string]interface{}{
			"dry_run": true,
			"method":  "POST",
			"path":    path,
			"owner":   ctx.Owner,
			"repo":    ctx.Repo,
		})
	}
	env, err := ctx.CallAPI("POST", path, nil)
	if err != nil {
		return err
	}
	return ctx.Output(env)
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

func normalizeRole(value string) (string, error) {
	role := strings.ToLower(strings.TrimSpace(value))
	if role == "" {
		role = "developer"
	}
	if normalized, ok := allowedRoles[role]; ok {
		return normalized, nil
	}
	return "", fmt.Errorf("invalid --role value %q: use manager, developer, or reporter", value)
}

func joinPath() string {
	return "/applied_projects"
}

func quitPath(ctx *common.RuntimeContext) string {
	return ctx.RepoPath() + "/quit"
}

func parseDryRun(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}
