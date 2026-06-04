package todo

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

const (
	transferResource = "applied_transfer_projects"
	joinResource     = "applied_projects"
)

// Shortcuts returns user todo and request approval shortcuts.
func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "transfer-list",
			Description: "List pending project transfer requests for a user",
			Flags:       listFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				return listRequests(ctx, transferResource)
			},
		},
		{
			Name:        "transfer-accept",
			Description: "Accept a project transfer request",
			Flags:       actionFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				return runRequestAction(ctx, transferResource, "accept")
			},
		},
		{
			Name:        "transfer-refuse",
			Description: "Refuse a project transfer request",
			Flags:       actionFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				return runRequestAction(ctx, transferResource, "refuse")
			},
		},
		{
			Name:        "join-list",
			Description: "List pending project join requests for a user",
			Flags:       listFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				return listRequests(ctx, joinResource)
			},
		},
		{
			Name:        "join-accept",
			Description: "Accept a project join request",
			Flags:       actionFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				return runRequestAction(ctx, joinResource, "accept")
			},
		},
		{
			Name:        "join-refuse",
			Description: "Refuse a project join request",
			Flags:       actionFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				return runRequestAction(ctx, joinResource, "refuse")
			},
		},
	}
}

func listFlags() []common.Flag {
	return []common.Flag{
		{Name: "login", Short: "l", Usage: "User login that owns the todo queue", Required: true},
		{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
		{Name: "per-page", Usage: "Items per page", Default: "20"},
	}
}

func actionFlags() []common.Flag {
	return []common.Flag{
		{Name: "login", Short: "l", Usage: "User login that owns the todo queue", Required: true},
		{Name: "id", Short: "i", Usage: "Request ID", Required: true},
		{Name: "dry-run", Usage: "Preview the request without changing data", Bool: true, Default: "false"},
	}
}

func listRequests(ctx *common.RuntimeContext, resource string) error {
	login, err := requiredLogin(ctx)
	if err != nil {
		return err
	}
	q := url.Values{}
	if page := strings.TrimSpace(ctx.Arg("page")); page != "" {
		q.Set("page", page)
	}
	if perPage := strings.TrimSpace(ctx.Arg("per-page")); perPage != "" {
		q.Set("per_page", perPage)
	}
	env, err := ctx.CallAPIWithQuery("GET", todoCollectionPath(login, resource), q)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func runRequestAction(ctx *common.RuntimeContext, resource, action string) error {
	login, err := requiredLogin(ctx)
	if err != nil {
		return err
	}
	id, err := parseRequestID(ctx.Arg("id"))
	if err != nil {
		return err
	}
	path := todoActionPath(login, resource, id, action)
	if parseDryRun(ctx.Arg("dry-run")) {
		return ctx.OutputData(map[string]interface{}{
			"dry_run":  true,
			"method":   "POST",
			"path":     path,
			"login":    login,
			"id":       id,
			"resource": resource,
			"action":   action,
		})
	}
	env, err := ctx.CallAPI("POST", path, nil)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func requiredLogin(ctx *common.RuntimeContext) (string, error) {
	login, err := ctx.RequireArg("login")
	if err != nil {
		return "", err
	}
	login = strings.TrimSpace(login)
	if login == "" {
		return "", fmt.Errorf("required flag --login is missing")
	}
	return login, nil
}

func todoCollectionPath(login, resource string) string {
	return fmt.Sprintf("/users/%s/%s", url.PathEscape(login), resource)
}

func todoActionPath(login, resource, id, action string) string {
	return fmt.Sprintf("%s/%s/%s", todoCollectionPath(login, resource), id, action)
}

func parseRequestID(value string) (string, error) {
	value = strings.TrimSpace(value)
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return "", fmt.Errorf("invalid request ID %q", value)
	}
	return strconv.Itoa(id), nil
}

func parseDryRun(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}
