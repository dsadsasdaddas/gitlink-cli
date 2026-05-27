# Code OpenAPI Shortcuts

## Summary

This change adds a dedicated `code` shortcut group for repository source-code
inspection APIs. The group keeps code browsing separate from repository
administration commands, which reduces coupling with existing `repo` shortcuts.

## Added shortcuts

| Shortcut | API |
|----------|-----|
| `code +files` | `GET /api/{owner}/{repo}/files.json` |
| `code +entries` | `GET /api/{owner}/{repo}/entries.json` |
| `code +sub-entries` | `GET /api/{owner}/{repo}/sub_entries.json` |
| `code +tree` | `GET /api/v1/{owner}/{repo}/git/trees/{sha}.json` |
| `code +blob` | `GET /api/v1/{owner}/{repo}/git/blobs/{sha}.json` |
| `code +commits` | `GET /api/v1/{owner}/{repo}/commits.json` |
| `code +commit-files` | `GET /api/v1/{owner}/{repo}/commits/{sha}/files.json` |
| `code +commit-diff` | `GET /api/v1/{owner}/{repo}/commits/{sha}/diff.json` |
| `code +blame` | `GET /api/v1/{owner}/{repo}/blame.json` |
| `code +tags` | `GET /api/v1/{owner}/{repo}/tags.json` |
| `code +tag` | `GET /api/v1/{owner}/{repo}/tags/{name}.json` |
| `code +delete-tag` | `DELETE /api/v1/{owner}/{repo}/tags/{tag}.json` |

`code +delete-tag` supports `--dry-run` so users can preview the exact delete
request before changing repository state.

## Submitter

Wang Yue
