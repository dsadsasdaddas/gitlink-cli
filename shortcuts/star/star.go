package star

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

// Shortcuts returns user starred/pinned project shortcut commands.
func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: "List a user's starred projects",
			Flags: []common.Flag{
				{Name: "login", Short: "l", Usage: "User login name", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				login, err := ctx.RequireArg("login")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", starListPath(login), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "set",
			Description: "Set starred project IDs for a user",
			Flags: []common.Flag{
				{Name: "login", Short: "l", Usage: "User login name", Required: true},
				{Name: "project-ids", Short: "p", Usage: "Comma-separated project IDs to keep starred", Required: true},
				{Name: "dry-run", Usage: "Preview the request without changing starred projects", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				login, err := ctx.RequireArg("login")
				if err != nil {
					return err
				}
				idsArg, err := ctx.RequireArg("project-ids")
				if err != nil {
					return err
				}
				ids, err := parsePositiveIntCSV(idsArg, "project-ids")
				if err != nil {
					return err
				}
				path := starSetPath(login)
				payload := map[string]interface{}{
					"is_pinned_project_ids": ids,
				}
				if parseBool(ctx.Arg("dry-run")) {
					return ctx.OutputData(map[string]interface{}{
						"dry_run": true,
						"action":  "set_starred_projects",
						"method":  "POST",
						"path":    path,
						"login":   login,
						"payload": payload,
					})
				}
				env, err := ctx.CallAPI("POST", path, payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "reorder",
			Description: "Update a starred project's display position",
			Flags: []common.Flag{
				{Name: "login", Short: "l", Usage: "User login name", Required: true},
				{Name: "pinned-id", Short: "i", Usage: "Starred project record ID from star +list", Required: true},
				{Name: "position", Short: "p", Usage: "Display position; larger numbers rank higher", Required: true},
				{Name: "dry-run", Usage: "Preview the request without changing order", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				login, err := ctx.RequireArg("login")
				if err != nil {
					return err
				}
				pinnedID, err := ctx.RequireArg("pinned-id")
				if err != nil {
					return err
				}
				positionArg, err := ctx.RequireArg("position")
				if err != nil {
					return err
				}
				position, err := parseNonNegativeInt(positionArg, "position")
				if err != nil {
					return err
				}
				path := starReorderPath(login, pinnedID)
				payload := map[string]interface{}{
					"pinned_project": map[string]interface{}{
						"position": position,
					},
				}
				if parseBool(ctx.Arg("dry-run")) {
					return ctx.OutputData(map[string]interface{}{
						"dry_run":   true,
						"action":    "reorder_starred_project",
						"method":    "PUT",
						"path":      path,
						"login":     login,
						"pinned_id": pinnedID,
						"payload":   payload,
					})
				}
				env, err := ctx.CallAPI("PUT", path, payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
	}
}

func starListPath(login string) string {
	return fmt.Sprintf("/users/%s/is_pinned_projects", url.PathEscape(login))
}

func starSetPath(login string) string {
	return fmt.Sprintf("/users/%s/is_pinned_projects/pin", url.PathEscape(login))
}

func starReorderPath(login, pinnedID string) string {
	return fmt.Sprintf("/users/%s/is_pinned_projects/%s", url.PathEscape(login), url.PathEscape(pinnedID))
}

func parsePositiveIntCSV(value, flagName string) ([]int, error) {
	parts := strings.Split(value, ",")
	ids := make([]int, 0, len(parts))
	seen := map[int]bool{}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.Atoi(part)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("invalid --%s value %q: use positive integer IDs", flagName, part)
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("--%s must include at least one project ID", flagName)
	}
	return ids, nil
}

func parseNonNegativeInt(value, flagName string) (int, error) {
	value = strings.TrimSpace(value)
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("invalid --%s value %q: use a non-negative integer", flagName, value)
	}
	return parsed, nil
}

func parseBool(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}
