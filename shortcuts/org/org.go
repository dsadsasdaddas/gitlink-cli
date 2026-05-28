package org

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gitlink-org/gitlink-cli/internal/i18n"
	"github.com/gitlink-org/gitlink-cli/shortcuts/common"
)

func Shortcuts(translators ...*i18n.Translator) []*common.Shortcut {
	tr := shortcutTranslator(translators...)
	return []*common.Shortcut{
		{
			Name:        "list",
			Description: tr.T("cmd.org.list.short"),
			Flags: []common.Flag{
				{Name: "page", Short: "p", Usage: tr.T("flag.page"), Default: "1"},
				{Name: "limit", Short: "l", Usage: tr.T("flag.limit"), Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				q := url.Values{}
				q.Set("page", ctx.Arg("page"))
				q.Set("limit", ctx.Arg("limit"))
				env, err := ctx.CallAPIWithQuery("GET", "/organizations", q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "info",
			Description: tr.T("cmd.org.info.short"),
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: tr.T("flag.org.id_or_login"), Required: true},
			},
			Run: func(ctx *common.RuntimeContext) error {
				id, _ := ctx.RequireArg("id")
				env, err := ctx.CallAPI("GET", fmt.Sprintf("/organizations/%s", id), nil)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "members",
			Description: tr.T("cmd.org.members.short"),
			Flags: []common.Flag{
				{Name: "id", Short: "i", Usage: tr.T("flag.org.id"), Required: true},
				{Name: "page", Short: "p", Usage: tr.T("flag.page"), Default: "1"},
				{Name: "limit", Short: "l", Usage: tr.T("flag.limit"), Default: "20"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				id, _ := ctx.RequireArg("id")
				q := url.Values{}
				q.Set("page", ctx.Arg("page"))
				q.Set("limit", ctx.Arg("limit"))
				env, err := ctx.CallAPIWithQuery("GET", fmt.Sprintf("/organizations/%s/organization_users", id), q)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "create",
			Description: tr.T("cmd.org.create.short"),
			Flags: []common.Flag{
				{Name: "name", Short: "n", Usage: tr.T("flag.org.name"), Required: true},
				{Name: "description", Short: "d", Usage: tr.T("flag.description")},
			},
			Run: func(ctx *common.RuntimeContext) error {
				name, _ := ctx.RequireArg("name")
				payload := map[string]interface{}{
					"name": name,
				}
				if d := ctx.Arg("description"); d != "" {
					payload["description"] = d
				}
				env, err := ctx.CallAPI("POST", "/organizations", payload)
				if err != nil {
					return err
				}
				return ctx.Output(env)
			},
		},
		{
			Name:        "team-projects-add-all",
			Description: "Add all organization projects to a team",
			Flags: []common.Flag{
				{Name: "organization", Short: "o", Usage: "Organization identifier", Required: true},
				{Name: "team-id", Short: "t", Usage: "Team ID", Required: true},
				{Name: "dry-run", Usage: "Preview the request without changing team projects", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				return runTeamProjectsAll(ctx, "POST", "add_all_team_projects")
			},
		},
		{
			Name:        "team-projects-remove-all",
			Description: "Remove all projects from an organization team",
			Flags: []common.Flag{
				{Name: "organization", Short: "o", Usage: "Organization identifier", Required: true},
				{Name: "team-id", Short: "t", Usage: "Team ID", Required: true},
				{Name: "dry-run", Usage: "Preview the request without changing team projects", Bool: true, Default: "false"},
			},
			Run: func(ctx *common.RuntimeContext) error {
				return runTeamProjectsAll(ctx, "DELETE", "remove_all_team_projects")
			},
		},
	}
}

func shortcutTranslator(translators ...*i18n.Translator) *i18n.Translator {
	if len(translators) > 0 && translators[0] != nil {
		return translators[0]
	}
	return i18n.Default()
}

func runTeamProjectsAll(ctx *common.RuntimeContext, method, action string) error {
	organization, err := ctx.RequireArg("organization")
	if err != nil {
		return err
	}
	teamID, err := ctx.RequireArg("team-id")
	if err != nil {
		return err
	}

	operation := "create_all"
	if method == "DELETE" {
		operation = "destroy_all"
	}
	path := fmt.Sprintf("/organizations/%s/teams/%s/team_projects/%s", organization, teamID, operation)
	if parseBool(ctx.Arg("dry-run")) {
		return ctx.OutputData(map[string]interface{}{
			"dry_run":      true,
			"action":       action,
			"method":       method,
			"path":         path,
			"organization": organization,
			"team_id":      teamID,
		})
	}
	env, err := ctx.CallAPI(method, path, nil)
	if err != nil {
		return err
	}
	return ctx.Output(env)
}

func parseBool(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}
