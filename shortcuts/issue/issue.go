package issue

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

type existingIssue struct {
	Subject     string
	Description string
}

func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: "List issues",
			Flags: []common.Flag{
				{Name: "state", Short: "s", Usage: "Filter by state: open, closed, all", Default: "open"},
				{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Items per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				q := url.Values{}
				q.Set("page", ctx.Arg("page"))
				q.Set("limit", ctx.Arg("limit"))
				if s := ctx.Arg("state"); s != "" {
					q.Set("state", s)
				}
				env, err := ctx.CallAPIWithQuery("GET", ctx.RepoPath()+"/issues", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "create",
			Description: "Create a new issue",
			Flags: []common.Flag{
				{Name: "title", Short: "t", Usage: "Issue title", Required: true},
				{Name: "body", Short: "b", Usage: "Issue description"},
				{Name: "assignee", Short: "a", Usage: "Assignee login"},
				{Name: "milestone", Short: "m", Usage: "Milestone ID"},
				{Name: "label", Usage: "Label ID"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				title, err := ctx.RequireArg("title")
				if err != nil {
					return err
				}
				body := map[string]interface{}{
					"subject":     title,
					"done_ratio":  0,
					"priority_id": 2,
				}
				if desc := ctx.Arg("body"); desc != "" {
					body["description"] = desc
				}
				if a := ctx.Arg("assignee"); a != "" {
					body["assigned_to_id"] = a
				}
				if m := ctx.Arg("milestone"); m != "" {
					body["fixed_version_id"] = m
				}
				env, err := ctx.CallAPI("POST", ctx.RepoPath()+"/issues", body)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "view",
			Description: "View issue details",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Issue ID or number", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", fmt.Sprintf("%s/issues/%s", ctx.RepoPath(), id), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "close",
			Description: "Close an issue",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Issue ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				current, err := fetchExistingIssue(ctx, id)
				if err != nil {
					return err
				}

				body := map[string]interface{}{
					"subject":     current.Subject,
					"description": current.Description,
					"status_id":   5, // 5 = closed
				}
				env, err := ctx.CallAPI("PUT", fmt.Sprintf("%s/issues/%s", ctx.RepoPath(), id), body)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "update",
			Description: "Update an issue",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Issue ID", Required: true},
				{Name: "title", Short: "t", Usage: "New title"},
				{Name: "body", Short: "b", Usage: "New description"},
				{Name: "state", Short: "s", Usage: "New state: open, closed, or numeric status_id"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				title := ctx.Arg("title")
				description := ctx.Arg("body")
				state := ctx.Arg("state")
				if title == "" && description == "" && state == "" {
					return fmt.Errorf("at least one of --title, --body, or --state is required")
				}

				current, err := fetchExistingIssue(ctx, id)
				if err != nil {
					return err
				}

				body := map[string]interface{}{
					"subject":     current.Subject,
					"description": current.Description,
				}
				if t := ctx.Arg("title"); t != "" {
					body["subject"] = t
				}
				if b := ctx.Arg("body"); b != "" {
					body["description"] = b
				}
				if s := ctx.Arg("state"); s != "" {
					statusID, err := normalizeIssueStatus(s)
					if err != nil {
						return err
					}
					body["status_id"] = statusID
				}
				env, err := ctx.CallAPI("PUT", fmt.Sprintf("%s/issues/%s", ctx.RepoPath(), id), body)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "comment",
			Description: "Add a comment to an issue",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Issue ID", Required: true},
				{Name: "body", Short: "b", Usage: "Comment body", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := ctx.RequireArg("id")
				if err != nil {
					return err
				}
				body, err := ctx.RequireArg("body")
				if err != nil {
					return err
				}
				payload := map[string]interface{}{
					"content": body,
				}
				env, err := ctx.CallAPI("POST", fmt.Sprintf("/issues/%s/journals", id), payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
	}
}

func fetchExistingIssue(ctx *common.RuntimeContext, id string) (*existingIssue, error) {
	getEnv, err := ctx.CallAPI("GET", fmt.Sprintf("%s/issues/%s", ctx.RepoPath(), id), nil)
	if err != nil {
		return nil, err
	}
	issueData, ok := getEnv.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse issue data")
	}
	subject, _ := issueData["subject"].(string)
	if subject == "" {
		return nil, fmt.Errorf("failed to parse issue subject")
	}
	description, _ := issueData["description"].(string)
	return &existingIssue{
		Subject:     subject,
		Description: description,
	}, nil
}

func normalizeIssueStatus(state string) (interface{}, error) {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "open":
		return 1, nil
	case "closed":
		return 5, nil
	default:
		if id, err := strconv.Atoi(state); err == nil {
			return id, nil
		}
		return nil, fmt.Errorf("invalid --state %q: use open, closed, or a numeric status_id", state)
	}
}
