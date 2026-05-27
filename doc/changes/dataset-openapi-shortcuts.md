# Dataset OpenAPI Shortcuts

## Summary

This change adds a dedicated `dataset` shortcut group for GitLink dataset APIs.
It covers repository dataset details, global dataset lookup, create/update, and
dataset attachment deletion.

## Added shortcuts

| Shortcut | API |
|----------|-----|
| `dataset +view` | `GET /api/v1/{owner}/{repo}/dataset` |
| `dataset +list` | `GET /api/v1/project_datasets.json` |
| `dataset +create` | `POST /api/v1/{owner}/{repo}/dataset` |
| `dataset +update` | `PUT /api/v1/{owner}/{repo}/dataset` |
| `dataset +delete-attachment` | `DELETE /api/attachments/{uuid}.json` |

`dataset +delete-attachment` supports `--dry-run` so users can preview the exact
delete request before changing dataset state.

## Submitter

Wang Yue
