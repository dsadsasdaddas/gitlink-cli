---
name: gitlink-pipeline
version: 1.0.0
description: "Pipeline workflow operations: list pipelines, run workflows, inspect runs, fetch logs, save visual YAML, enable, disable, and delete pipelines."
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli pipeline --help"
---

# gitlink-pipeline

Read [`../gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md) first for authentication, global flags, and API behavior.

**CRITICAL**: confirm user intent before running write or delete commands. Prefer `--dry-run` before `pipeline +run`, `pipeline +save-yaml`, `pipeline +enable`, `pipeline +disable`, or `pipeline +delete`.

## Shortcuts

| Shortcut | Description |
|----------|-------------|
| `pipeline +list` | List platform pipelines |
| `pipeline +runs` | List repository pipeline run records |
| `pipeline +run` | Start a pipeline workflow |
| `pipeline +view` | Show pipeline details |
| `pipeline +delete` | Delete a pipeline |
| `pipeline +save-yaml` | Save a visual pipeline YAML graph |
| `pipeline +enable` | Enable a pipeline workflow |
| `pipeline +disable` | Disable a pipeline workflow |
| `pipeline +logs` | Query pipeline run logs |
| `pipeline +results` | Show pipeline run report results |

## Examples

```bash
# List platform pipelines
gitlink-cli pipeline +list --owner-id 123 --page 1 --limit 20

# List and run workflow records for a repository
gitlink-cli pipeline +runs --owner Gitlink --repo forgeplus --ref master --workflow build.yml
gitlink-cli pipeline +run --owner Gitlink --repo forgeplus --ref master --workflow build.yml --dry-run

# Inspect a pipeline and its logs
gitlink-cli pipeline +view --owner Gitlink --repo forgeplus --id 7
gitlink-cli pipeline +logs --owner Gitlink --repo forgeplus --run-id 99 --id 7 --index 43
gitlink-cli pipeline +results --owner Gitlink --repo forgeplus --run-id 99

# Save a visual pipeline graph
gitlink-cli pipeline +save-yaml --owner Gitlink --repo forgeplus \
  --id 7 --pipeline-json '{"nodes":[]}' --dry-run

# Toggle or delete pipeline workflows
gitlink-cli pipeline +enable --owner Gitlink --repo forgeplus --id 7 --workflow build.yml --dry-run
gitlink-cli pipeline +disable --owner Gitlink --repo forgeplus --id 7 --workflow build.yml --dry-run
gitlink-cli pipeline +delete --owner Gitlink --repo forgeplus --id 7 --dry-run
```

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
