package pipeline

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func Shortcuts() []*common.Shortcut {
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: "List platform pipelines",
			Flags: []common.Flag{
				{Name: "owner-id", Usage: "Owner user or organization ID"},
				{Name: "page", Short: "p", Usage: "Page number", Default: "1"},
				{Name: "limit", Short: "l", Usage: "Items per page", Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				q := pageLimitQuery(ctx)
				setQueryIfPresent(q, ctx, "owner-id", "owner_id")
				env, err := ctx.CallAPIWithQuery("GET", "/pm/pipelines", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "runs",
			Description: "List pipeline run records",
			Flags:       runFilterFlags(),
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				env, err := ctx.CallAPIWithQuery("GET", pipelineV1RepoPath(ctx)+"/actions/runs", runFilterQuery(ctx))
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "run",
			Description: "Run a pipeline workflow",
			Flags: append(runFilterFlags(),
				common.Flag{Name: "dry-run", Usage: "Preview the run request without starting a pipeline", Bool: true, Default: "false"},
			),
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				q := runFilterQuery(ctx)
				path := pipelineV1RepoPath(ctx) + "/actions/runs"
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(map[string]interface{}{
						"repository": fmt.Sprintf("%s/%s", ctx.Owner, ctx.Repo),
						"dry_run":    true,
						"action":     "run_pipeline",
						"method":     "POST",
						"path":       path,
						"query":      q,
					})
				}
				env, err := ctx.CallAPIWithQuery("POST", path, q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "view",
			Description: "Show pipeline details",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Pipeline ID", Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := requiredPositiveInt(ctx, "id")
				if err != nil {
					return err
				}
				env, err := ctx.CallAPI("GET", fmt.Sprintf("%s/pipelines/%d", pipelineV1RepoPath(ctx), id), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "delete",
			Description: "Delete a pipeline",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Pipeline ID", Required: true},
				{Name: "dry-run", Usage: "Preview the delete request without changing pipeline state", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := requiredPositiveInt(ctx, "id")
				if err != nil {
					return err
				}
				path := fmt.Sprintf("%s/pipelines/%d", pipelineV1RepoPath(ctx), id)
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(map[string]interface{}{
						"repository": fmt.Sprintf("%s/%s", ctx.Owner, ctx.Repo),
						"dry_run":    true,
						"action":     "delete_pipeline",
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
		{
			Name:        "save-yaml",
			Description: "Save a visual pipeline YAML graph",
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: "Pipeline ID", Required: true},
				{Name: "pipeline-json", Usage: "Pipeline graph JSON object or string", Required: true},
				{Name: "dry-run", Usage: "Preview the save request without changing pipeline state", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				id, err := requiredPositiveInt(ctx, "id")
				if err != nil {
					return err
				}
				pipelineJSON, err := ctx.RequireArg("pipeline-json")
				if err != nil {
					return err
				}
				payload := map[string]interface{}{
					"id":            id,
					"pipeline_json": parseJSONValue(pipelineJSON),
				}
				path := pipelineV1RepoPath(ctx) + "/pipelines/save_yaml"
				if ctx.Arg("dry-run") == "true" {
					return ctx.OutputData(map[string]interface{}{
						"repository": fmt.Sprintf("%s/%s", ctx.Owner, ctx.Repo),
						"dry_run":    true,
						"action":     "save_pipeline_yaml",
						"method":     "POST",
						"path":       path,
						"payload":    payload,
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
			Name:        "enable",
			Description: "Enable a pipeline workflow",
			Flags:       workflowStateFlags(true),
			Run: func(ctx *common.RuntimeContext) error {
				return runWorkflowState(ctx, "enable")
			},
		},
		{
			Name:        "disable",
			Description: "Disable a pipeline workflow",
			Flags:       workflowStateFlags(true),
			Run: func(ctx *common.RuntimeContext) error {
				return runWorkflowState(ctx, "disable")
			},
		},
		{
			Name:        "logs",
			Description: "Query pipeline run logs",
			Flags: []common.Flag{
				{Name: "run-id", Short: "r", Usage: "Pipeline run ID", Required: true},
				{Name: "id", Short: "i", Usage: "Pipeline ID", Required: true},
				{Name: "index", Usage: "Run index", Required: true},
				{Name: "job", Short: "j", Usage: "Job index", Default: "0"},
				{Name: "cursor", Usage: "Log cursor"},
				{Name: "step", Usage: "Log step", Default: "1"},
				{Name: "expanded", Usage: "Whether the log cursor is expanded", Bool: true, Default: "true"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				runID, err := ctx.RequireArg("run-id")
				if err != nil {
					return err
				}
				id, err := requiredPositiveInt(ctx, "id")
				if err != nil {
					return err
				}
				index, err := ctx.RequireArg("index")
				if err != nil {
					return err
				}
				job, err := requiredNonNegativeIntWithDefault(ctx, "job", 0)
				if err != nil {
					return err
				}
				step, err := requiredPositiveIntWithDefault(ctx, "step", 1)
				if err != nil {
					return err
				}
				payload := map[string]interface{}{
					"id":    id,
					"index": index,
					"job":   job,
					"owner": ctx.Owner,
					"repo":  ctx.Repo,
					"log_cursors": []map[string]interface{}{
						{
							"cursor":   ctx.Arg("cursor"),
							"expanded": ctx.Arg("expanded") == "true",
							"step":     step,
						},
					},
				}
				env, err := ctx.CallAPI("POST", fmt.Sprintf("%s/actions/runs/%s/jobs/0", pipelineV1RepoPath(ctx), url.PathEscape(runID)), payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "results",
			Description: "Show pipeline run report results",
			Flags: []common.Flag{
				{Name: "run-id", Short: "r", Usage: "Pipeline run ID"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				if err := ctx.ResolveOwnerRepo(); err != nil {
					return err
				}
				q := url.Values{}
				setQueryIfPresent(q, ctx, "run-id", "run_id")
				env, err := ctx.CallAPIWithQuery("GET", pipelineV1RepoPath(ctx)+"/pipelines/run_results", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
	}
}

func pipelineV1RepoPath(ctx *common.RuntimeContext) string {
	return "/v1" + ctx.RepoPath()
}

func runFilterFlags() []common.Flag {
	return []common.Flag{
		{Name: "ref", Short: "r", Usage: "Branch, tag, or commit SHA"},
		{Name: "workflow", Short: "w", Usage: "Workflow file name"},
	}
}

func runFilterQuery(ctx *common.RuntimeContext) url.Values {
	q := url.Values{}
	setQueryIfPresent(q, ctx, "ref", "ref")
	setQueryIfPresent(q, ctx, "workflow", "workflow")
	return q
}

func workflowStateFlags(includeDryRun bool) []common.Flag {
	flags := []common.Flag{
		{Name: "id", Short: "i", Usage: "Pipeline ID", Required: true},
		{Name: "workflow", Short: "w", Usage: "Workflow file name", Required: true},
	}
	if includeDryRun {
		flags = append(flags, common.Flag{Name: "dry-run", Usage: "Preview the request without changing pipeline state", Bool: true, Default: "false"})
	}
	return flags
}

func runWorkflowState(ctx *common.RuntimeContext, action string) error {
	if err := ctx.ResolveOwnerRepo(); err != nil {
		return err
	}
	id, err := requiredPositiveInt(ctx, "id")
	if err != nil {
		return err
	}
	workflow, err := ctx.RequireArg("workflow")
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"id":       id,
		"workflow": workflow,
	}
	path := fmt.Sprintf("%s/actions/%s", pipelineV1RepoPath(ctx), action)
	if ctx.Arg("dry-run") == "true" {
		return ctx.OutputData(map[string]interface{}{
			"repository": fmt.Sprintf("%s/%s", ctx.Owner, ctx.Repo),
			"dry_run":    true,
			"action":     action + "_pipeline",
			"method":     "POST",
			"path":       path,
			"payload":    payload,
		})
	}
	env, err := ctx.CallAPI("POST", path, payload)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func pageLimitQuery(ctx *common.RuntimeContext) url.Values {
	q := url.Values{}
	setQueryIfPresent(q, ctx, "page", "page")
	setQueryIfPresent(q, ctx, "limit", "limit")
	return q
}

func setQueryIfPresent(q url.Values, ctx *common.RuntimeContext, flagName, queryName string) {
	if value := ctx.Arg(flagName); value != "" {
		q.Set(queryName, value)
	}
}

func parseJSONValue(value string) interface{} {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return value
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
		return parsed
	}
	return value
}

func requiredPositiveInt(ctx *common.RuntimeContext, flagName string) (int, error) {
	value, err := ctx.RequireArg(flagName)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("--%s must be a positive integer", flagName)
	}
	return id, nil
}

func requiredPositiveIntWithDefault(ctx *common.RuntimeContext, flagName string, defaultValue int) (int, error) {
	value := ctx.Arg(flagName)
	if value == "" {
		return defaultValue, nil
	}
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("--%s must be a positive integer", flagName)
	}
	return id, nil
}

func requiredNonNegativeIntWithDefault(ctx *common.RuntimeContext, flagName string, defaultValue int) (int, error) {
	value := ctx.Arg(flagName)
	if value == "" {
		return defaultValue, nil
	}
	id, err := strconv.Atoi(value)
	if err != nil || id < 0 {
		return 0, fmt.Errorf("--%s must be a non-negative integer", flagName)
	}
	return id, nil
}
