package dataset

import (
	"fmt"
	"net/url"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "view",
			Description: "Show repository dataset details and attachments",
			Flags: []common.Flag{
				{Name: "page", Short: "p", Usage: "Attachment page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Attachments per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				env, err := ctx.CallAPIWithQuery("GET", datasetRepoPath(ctx), datasetPageQuery(ctx))
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "list",
			Description: "List project datasets",
			Flags: []common.Flag{
				{Name: "ids", Usage: "Comma-separated dataset IDs to query"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				q := url.Values{}
				setDatasetQueryIfPresent(q, ctx, "ids", "ids")
				env, err := ctx.CallAPIWithQuery("GET", "/v1/project_datasets", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "create",
			Description: "Create a repository dataset",
			Flags:       datasetWriteFlags(true),
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				payload, err := datasetPayload(ctx)
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("POST", datasetRepoPath(ctx), payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "update",
			Description: "Update a repository dataset",
			Flags:       datasetWriteFlags(true),
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				payload, err := datasetPayload(ctx)
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("PUT", datasetRepoPath(ctx), payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "delete-attachment",
			Description: "Delete a dataset attachment by UUID",
			Flags: []common.Flag{
				{Name: "uuid", Short: "u", Usage: "Attachment UUID", Required: true},
				{Name: "dry-run", Usage: "Preview the delete request without changing dataset state", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				uuid, err := ctx.RequireArg("uuid")
				if err != nil {
					return err
				}
				path := fmt.Sprintf("/attachments/%s", url.PathEscape(uuid))
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(map[string]interface{}{
						"dry_run": true,
						"action":  "delete_dataset_attachment",
						"uuid":    uuid,
						"method":  "DELETE",
						"path":    path,
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

func datasetRepoPath(ctx *common.RuntimeContext) string {
	return "/v1" + ctx.RepoPath() + "/dataset"
}

func datasetWriteFlags(required bool) []common.Flag {
	return []common.Flag{
		{Name: "title", Short: "t", Usage: "Dataset title", Required: required},
		{Name: "description", Short: "d", Usage: "Dataset description", Required: required},
		{Name: "license-id", Usage: "License ID"},
		{Name: "paper-content", Usage: "Research paper content"},
	}
}

func datasetPayload(ctx *common.RuntimeContext) (map[string]interface{}, error) {
	title, err := ctx.RequireArg("title")
	if err != nil {
		return nil, err
	}
	description, err := ctx.RequireArg("description")
	if err != nil {
		return nil, err
	}
	payload := map[string]interface{}{
		"title":       title,
		"description": description,
	}
	setDatasetPayloadIfPresent(payload, ctx, "license-id", "license_id")
	setDatasetPayloadIfPresent(payload, ctx, "paper-content", "paper_content")
	return payload, nil
}

func datasetPageQuery(ctx *common.RuntimeContext) url.Values {
	q := url.Values{}
	setDatasetQueryIfPresent(q, ctx, "page", "page")
	setDatasetQueryIfPresent(q, ctx, "limit", "limit")
	return q
}

func setDatasetQueryIfPresent(q url.Values, ctx *common.RuntimeContext, flagName, queryName string) {
	if value := ctx.Arg(flagName); value != "" {
		q.Set(queryName, value)
	}
}

func setDatasetPayloadIfPresent(payload map[string]interface{}, ctx *common.RuntimeContext, flagName, payloadName string) {
	if value := ctx.Arg(flagName); value != "" {
		payload[payloadName] = value
	}
}
