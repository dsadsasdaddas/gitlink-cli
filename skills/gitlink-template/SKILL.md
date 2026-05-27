---
name: gitlink-template
version: 1.0.0
description: "Repository project template operations: list, view, create, update, and delete GitLink project templates."
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli template --help"
---

# gitlink-template

Read [`../gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md) first for authentication, global flags, and API behavior.

**CRITICAL**: confirm user intent before creating, updating, or deleting templates. Prefer `--dry-run` before write or delete operations.

## Shortcuts

| Shortcut | Description |
|----------|-------------|
| `template +list` | List repository project templates |
| `template +view` | Show project template details |
| `template +create` | Create a project template |
| `template +update` | Update a project template while preserving unspecified fields |
| `template +delete` | Delete a project template |

## Examples

```bash
# List and inspect templates
gitlink-cli template +list --owner Gitlink --repo forgeplus
gitlink-cli template +view --owner Gitlink --repo forgeplus --id 7

# Create an issue template from inline content
gitlink-cli template +create --owner Gitlink --repo forgeplus \
  --type ProjectTemplates::Issue \
  --name "Bug report" \
  --content "## Problem\n\n## Steps\n\n## Expected" \
  --dry-run

# Create from a file
gitlink-cli template +create --owner Gitlink --repo forgeplus \
  --name "Feature request" --content-file .gitlink/FEATURE_TEMPLATE.md --dry-run

# Update only the fields that changed
gitlink-cli template +update --owner Gitlink --repo forgeplus \
  --id 7 --content-file .gitlink/ISSUE_TEMPLATE.md --dry-run

# Delete a template
gitlink-cli template +delete --owner Gitlink --repo forgeplus --id 7 --dry-run
```

## API Mapping

| Shortcut | Method | API path |
|----------|--------|----------|
| `template +list` | GET | `/api/v1/{owner}/{repo}/project_templates.json` |
| `template +view` | GET | `/api/v1/{owner}/{repo}/project_templates/{id}.json` |
| `template +create` | POST | `/api/v1/{owner}/{repo}/project_templates.json` |
| `template +update` | PUT | `/api/v1/{owner}/{repo}/project_templates/{id}.json` |
| `template +delete` | DELETE | `/api/v1/{owner}/{repo}/project_templates/{id}.json` |

## Notes

- `template +create` supports `--content` and `--content-file`; use only one at a time.
- `template +update` reads the current template first and preserves unspecified `type`, `name`, or `content` fields.
- `template +create`, `template +update`, and `template +delete` support `--dry-run`.
