package code

import (
	"fmt"
	"net/url"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "files",
			Description: "Search repository files",
			Flags: []common.Flag{
				{Name: "search", Short: "s", Usage: "Search keyword"},
				{Name: "ref", Short: "r", Usage: "Branch, tag, or commit SHA"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				q := url.Values{}
				setQueryIfPresent(q, ctx, "search", "search")
				setQueryIfPresent(q, ctx, "ref", "ref")
				env, err := ctx.CallAPIWithQuery("GET", ctx.RepoPath()+"/files", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "entries",
			Description: "List repository root entries",
			Flags: []common.Flag{
				{Name: "ref", Short: "r", Usage: "Branch, tag, or commit SHA"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				q := url.Values{}
				setQueryIfPresent(q, ctx, "ref", "ref")
				env, err := ctx.CallAPIWithQuery("GET", ctx.RepoPath()+"/entries", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "sub-entries",
			Description: "Show a repository directory or file entry",
			Flags: []common.Flag{
				{Name: "path", Short: "p", Usage: "Repository-relative file or directory path", Required: true},
				{Name: "ref", Short: "r", Usage: "Branch, tag, or commit SHA"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				filepath, err := ctx.RequireArg("path")
				if err != nil {
					return err
				}
				q := url.Values{}
				q.Set("filepath", filepath)
				setQueryIfPresent(q, ctx, "ref", "ref")
				env, err := ctx.CallAPIWithQuery("GET", ctx.RepoPath()+"/sub_entries", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "tree",
			Description: "List a git tree",
			Flags: []common.Flag{
				{Name: "sha", Short: "s", Usage: "Branch, tag, commit SHA, or tree SHA", Required: true},
				{Name: "recursive", Short: "R", Usage: "Show recursive entries", Bool: true, Default: "false"},
				{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Items per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				sha, err := ctx.RequireArg("sha")
				if err != nil {
					return err
				}
				q := pageLimitQuery(ctx)
				if ctx.Arg("recursive") == "true" {
					q.Set("recursive", "true")
				}
				env, err := ctx.CallAPIWithQuery("GET", fmt.Sprintf("%s/git/trees/%s", codeV1RepoPath(ctx), url.PathEscape(sha)), q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "blob",
			Description: "Show a git blob",
			Flags: []common.Flag{
				{Name: "sha", Short: "s", Usage: "Blob SHA", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				sha, err := ctx.RequireArg("sha")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", fmt.Sprintf("%s/git/blobs/%s", codeV1RepoPath(ctx), url.PathEscape(sha)), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "commits",
			Description: "List repository commits",
			Flags: []common.Flag{
				{Name: "sha", Short: "s", Usage: "Branch, tag, or commit SHA"},
				{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Items per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				q := pageLimitQuery(ctx)
				setQueryIfPresent(q, ctx, "sha", "sha")
				env, err := ctx.CallAPIWithQuery("GET", codeV1RepoPath(ctx)+"/commits", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "commit-files",
			Description: "List files changed by a commit",
			Flags: []common.Flag{
				{Name: "sha", Short: "s", Usage: "Commit SHA", Required: true},
				{Name: "file", Short: "f", Usage: "Filter by file path"},
				{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Items per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				sha, err := ctx.RequireArg("sha")
				if err != nil {
					return err
				}
				q := pageLimitQuery(ctx)
				setQueryIfPresent(q, ctx, "file", "filepath")
				env, err := ctx.CallAPIWithQuery("GET", fmt.Sprintf("%s/commits/%s/files", codeV1RepoPath(ctx), url.PathEscape(sha)), q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "commit-diff",
			Description: "Show a commit diff",
			Flags: []common.Flag{
				{Name: "sha", Short: "s", Usage: "Commit SHA", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				sha, err := ctx.RequireArg("sha")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", fmt.Sprintf("%s/commits/%s/diff", codeV1RepoPath(ctx), url.PathEscape(sha)), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "blame",
			Description: "Show file blame information",
			Flags: []common.Flag{
				{Name: "sha", Short: "s", Usage: "Branch, tag, or commit SHA", Required: true},
				{Name: "path", Short: "p", Usage: "Repository-relative file path", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				sha, err := ctx.RequireArg("sha")
				if err != nil {
					return err
				}
				filepath, err := ctx.RequireArg("path")
				if err != nil {
					return err
				}
				q := url.Values{}
				q.Set("sha", sha)
				q.Set("filepath", filepath)
				env, err := ctx.CallAPIWithQuery("GET", codeV1RepoPath(ctx)+"/blame", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "tags",
			Description: "List repository tags",
			Flags: []common.Flag{
				{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Items per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				env, err := ctx.CallAPIWithQuery("GET", codeV1RepoPath(ctx)+"/tags", pageLimitQuery(ctx))
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "tag",
			Description: "Show repository tag details",
			Flags: []common.Flag{
				{Name: "name", Short: "n", Usage: "Tag name", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				name, err := ctx.RequireArg("name")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", fmt.Sprintf("%s/tags/%s", codeV1RepoPath(ctx), url.PathEscape(name)), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "delete-tag",
			Description: "Delete a repository tag",
			Flags: []common.Flag{
				{Name: "name", Short: "n", Usage: "Tag name", Required: true},
				{Name: "dry-run", Usage: "Preview the delete request without changing repository state", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				name, err := ctx.RequireArg("name")
				if err != nil {
					return err
				}
				path := fmt.Sprintf("%s/tags/%s", codeV1RepoPath(ctx), url.PathEscape(name))
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(map[string]interface{}{
						"repository": fmt.Sprintf("%s/%s", ctx.Owner, ctx.Repo),
						"dry_run":    true,
						"action":     "delete_tag",
						"tag":        name,
						"method":     "DELETE",
						"path":       path,
					})
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

func codeV1RepoPath(ctx *common.RuntimeContext) string {
	return "/v1" + ctx.RepoPath()
}

func setQueryIfPresent(q url.Values, ctx *common.RuntimeContext, flagName, queryName string) {
	if value := ctx.Arg(flagName); value != "" {
		q.Set(queryName, value)
	}
}

func pageLimitQuery(ctx *common.RuntimeContext) url.Values {
	q := url.Values{}
	setQueryIfPresent(q, ctx, "page", "page")
	setQueryIfPresent(q, ctx, "limit", "limit")
	return q
}
