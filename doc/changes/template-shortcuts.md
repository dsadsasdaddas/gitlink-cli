# Project Template Shortcuts

Submitter: Wang Yue

This change adds a dedicated `template` shortcut group for GitLink repository project template OpenAPI coverage.

## Commands

- `template +list`
- `template +view`
- `template +create`
- `template +update`
- `template +delete`

## API Mapping

| Shortcut | Method | API path |
|----------|--------|----------|
| `template +list` | GET | `/api/v1/{owner}/{repo}/project_templates.json` |
| `template +view` | GET | `/api/v1/{owner}/{repo}/project_templates/{id}.json` |
| `template +create` | POST | `/api/v1/{owner}/{repo}/project_templates.json` |
| `template +update` | PUT | `/api/v1/{owner}/{repo}/project_templates/{id}.json` |
| `template +delete` | DELETE | `/api/v1/{owner}/{repo}/project_templates/{id}.json` |

## Behavior

- `template +create` supports inline `--content` and file-based `--content-file`.
- `template +update` fetches current template data first, then preserves unspecified fields.
- Write/delete operations support `--dry-run`.

## Verification

- Unit tests cover request methods, paths, JSON payloads, field preservation, dry-run behavior, and missing-argument validation.
