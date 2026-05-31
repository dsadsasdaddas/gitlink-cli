# Pipeline OpenAPI Shortcuts

Submitter: Wang Yue

This change adds a dedicated `pipeline` shortcut group for GitLink Pipeline OpenAPI coverage.

## Commands

- `pipeline +list`
- `pipeline +runs`
- `pipeline +run`
- `pipeline +view`
- `pipeline +delete`
- `pipeline +save-yaml`
- `pipeline +enable`
- `pipeline +disable`
- `pipeline +logs`
- `pipeline +results`

## API Mapping

| Shortcut | Method | API path |
|----------|--------|----------|
| `pipeline +list` | GET | `/api/pm/pipelines.json` |
| `pipeline +runs` | GET | `/api/v1/{owner}/{repo}/actions/runs.json` |
| `pipeline +run` | POST | `/api/v1/{owner}/{repo}/actions/runs.json` |
| `pipeline +view` | GET | `/api/v1/{owner}/{repo}/pipelines/{id}.json` |
| `pipeline +delete` | DELETE | `/api/v1/{owner}/{repo}/pipelines/{id}.json` |
| `pipeline +save-yaml` | POST | `/api/v1/{owner}/{repo}/pipelines/save_yaml` |
| `pipeline +enable` | POST | `/api/v1/{owner}/{repo}/actions/enable.json` |
| `pipeline +disable` | POST | `/api/v1/{owner}/{repo}/actions/disable.json` |
| `pipeline +logs` | POST | `/api/v1/{owner}/{repo}/actions/runs/{run_id}/jobs/0` |
| `pipeline +results` | GET | `/api/v1/{owner}/{repo}/pipelines/run_results.json` |

## Verification

- Unit tests cover request methods, paths, query parameters, request bodies, dry-run behavior, and invalid ID validation.
- Help documentation is available through `gitlink-cli pipeline --help` and command-specific help.
- Write and delete commands support `--dry-run` to preview requests before changing pipeline state.
