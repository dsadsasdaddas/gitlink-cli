---
name: gitlink-dataset
version: 1.0.0
description: "Dataset management: create, update, view repository datasets, list project datasets, and delete dataset attachments."
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli dataset --help"
---

# gitlink-dataset

Read `../gitlink-shared/SKILL.md` first for authentication, global flags, and safety rules.

Write operations must be confirmed by the user before execution. Use
`dataset +delete-attachment --dry-run` before deleting an attachment.

## Shortcuts

| Shortcut | Purpose | Auth |
|----------|---------|------|
| `dataset +view` | Show repository dataset details and attachments | No for public repos |
| `dataset +list` | List project datasets, optionally by IDs | No for public repos |
| `dataset +create` | Create a repository dataset | Yes |
| `dataset +update` | Update a repository dataset | Yes |
| `dataset +delete-attachment` | Delete a dataset attachment by UUID | Yes |

## Examples

```bash
gitlink-cli dataset +view --owner Gitlink --repo forgeplus
gitlink-cli dataset +list --ids 1,2

gitlink-cli dataset +create --owner Gitlink --repo forgeplus \
  --title "Research dataset" \
  --description "Dataset description" \
  --license-id 3 \
  --paper-content "Paper content"

gitlink-cli dataset +update --owner Gitlink --repo forgeplus \
  --title "Research dataset" \
  --description "Updated description"

gitlink-cli dataset +delete-attachment --uuid <attachment-uuid> --dry-run
```

## API Mapping

| Shortcut | API |
|----------|-----|
| `dataset +view` | `GET /api/v1/{owner}/{repo}/dataset` |
| `dataset +list` | `GET /api/v1/project_datasets.json` |
| `dataset +create` | `POST /api/v1/{owner}/{repo}/dataset` |
| `dataset +update` | `PUT /api/v1/{owner}/{repo}/dataset` |
| `dataset +delete-attachment` | `DELETE /api/attachments/{uuid}.json` |
