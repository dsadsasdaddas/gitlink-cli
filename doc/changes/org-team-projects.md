# Organization Team Project Binding Shortcuts

## Summary

This change adds high-level shortcuts for two organization team project binding endpoints:

- `org +team-projects-add-all` → `POST /api/organizations/{organization}/teams/{id}/team_projects/create_all.json`
- `org +team-projects-remove-all` → `DELETE /api/organizations/{organization}/teams/{id}/team_projects/destroy_all.json`

## User Value

Organization maintainers can quickly bind all organization projects to a team, or clear all project bindings from a team, without manually using Raw API.

## Safety

Both write/destructive operations support `--dry-run`, returning the method, path, organization, and team ID before any remote mutation.

## Validation

Unit tests cover POST/DELETE request paths and dry-run no-network behavior. README and `gitlink-org` Skill docs include command examples and safety workflow guidance.
