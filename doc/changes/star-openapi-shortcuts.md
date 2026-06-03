# Starred Project Shortcuts

Submitter: Wang Yue

This change adds a medium-size OpenAPI shortcut group for user starred/pinned projects.

## Commands

- `star +list`
- `star +set`
- `star +reorder`

## API Mapping

| Shortcut | Method | API path |
|----------|--------|----------|
| `star +list` | GET | `/api/users/{owner}/is_pinned_projects.json` |
| `star +set` | POST | `/api/users/{owner}/is_pinned_projects/pin.json` |
| `star +reorder` | PUT | `/api/users/{owner}/is_pinned_projects/{id}.json` |

## Design Notes

- The feature is implemented as a standalone `star` command group instead of extending `user`, reducing conflicts with other user-account PRs.
- `star +set` accepts `--project-ids` and sends the OpenAPI field `is_pinned_project_ids`.
- `star +reorder` uses `--pinned-id`, which is the starred record ID returned by `star +list` as `projects[].id`; it is intentionally not named `--project-id`.
- Write operations support `--dry-run` to preview request method, path, and payload before changing profile data.

## Verification

```bash
git diff --check
GOPROXY=https://goproxy.cn,direct go test ./...
go vet ./...
go run . star --help
go run . star +list --help
go run . star +set --help
go run . star +reorder --help
go run . star +set --login wangyue111 --project-ids 17,42 --dry-run --format json
go run . star +reorder --login wangyue111 --pinned-id 9 --position 10 --dry-run --format json
```

Unit tests cover:

- API method/path mapping.
- `star +set` JSON payload and duplicate ID normalization.
- `star +set --dry-run` no-network behavior.
- `star +reorder` JSON payload and dry-run no-network behavior.
- Required flag validation and invalid numeric input rejection.
