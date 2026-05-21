# Workflow Agent Test Report

## Scope

This phase covers:

- Issue triage rules
- health scoring rules
- PR summary rules
- local command execution
- API fetch boundary tests
- remote read-only manual verification
- `json` / `table` / `markdown` rendering
- language handling
- mock tests do not depend on the real remote API

## Environment

- OS: Windows
- Go version: `go1.26.1 windows/amd64`
- Go path: `E:\GitLinkCLI-Competition\tools\go1.26.1\go\bin\go.exe`
- gofmt path: `E:\GitLinkCLI-Competition\tools\go1.26.1\go\bin\gofmt.exe`

## Test Commands

Executed:

```bash
gofmt -w shortcuts/workflow/*.go shortcuts/register.go
go test ./shortcuts/workflow
go test ./...
```

Results:

- `go test ./shortcuts/workflow` passed.
- `go test ./...` passed.

## Unit Tests

- triage rules tests
- health score tests
- messages tests
- render tests
- command tests
- fetch boundary tests
- PR summary rules and fetch tests

## API Fetch Boundary Tests

- empty issue responses return a clear error instead of panicking
- missing issue titles still allow body-only issues to be normalized
- label normalization supports string arrays, object arrays, and title/name variants
- author normalization supports string, `user`, and `creator` shapes
- GitLink error-in-body responses return readable errors
- health activity timestamps accept `updated_at`, `updatedAt`, `last_activity_at`, `merged_at`, and `closed_at`
- release responses accept `releases`, `data`, and direct array shapes
- CI unavailability is recorded as `unknown` without failing the whole health run
- stale-days values `0` and negative values fall back to the default `30`
- PR summary fetch normalizes PR metadata, changed files, commits, authors, branches, and list limits
- PR summary tolerates partial files or commits fetch failures while keeping base PR metadata
- PR summary base PR error-in-body responses return readable errors

## Manual Command Examples

```bash
gitlink-cli workflow +triage --title "Install failed on Windows" --body "go install failed with error" --format table
gitlink-cli workflow +triage --title "Token leaked in logs" --body "The access token appears in command output" --format json
gitlink-cli workflow +triage --title "安装失败，无法登录" --body "运行命令时报错" --lang zh-CN --format markdown
gitlink-cli workflow +triage --from shortcuts/workflow/testdata/issue_bug.json --format json
gitlink-cli workflow +health --repository Gitlink/gitlink-cli --open-issues 3 --open-prs 1 --has-readme --has-license --has-contributing --agent-readiness-known --agent-readiness-score 9 --format table
gitlink-cli workflow +health --repository demo/repo --open-issues 60 --stale-issues 25 --open-prs 12 --stale-prs 6 --recent-activity-known --recent-activity-days 120 --release-known=false --format json
gitlink-cli workflow +health --repository Gitlink/gitlink-cli --open-issues 3 --open-prs 1 --has-readme --has-license --has-contributing --lang zh-CN --format markdown
gitlink-cli workflow +pr-summary --owner Gitlink --repo gitlink-cli --number 1 --format markdown
gitlink-cli workflow +pr-summary --from shortcuts/workflow/testdata/pr_summary.json --format json
```

## Remote Manual Verification

- Command: `gitlink-cli workflow +triage --owner Gitlink --repo gitlink-cli --state open --limit 5 --format table`
- Result: succeeded, returned five issues in table form.
- Command: `gitlink-cli workflow +health --owner Gitlink --repo gitlink-cli --stale-days 30 --lang zh-CN --format markdown`
- Result: succeeded, returned a markdown health report with score `58` and risk level `high`.
- Remote writes: `No`

## Known Limitations

- Current workflow commands support local analysis and read-only GitLink fetch mode.
- `workflow +triage` still supports local parameters or a local JSON file via `--from`.
- `workflow +health` still supports local parameters or a local JSON file via `--from`.
- `workflow +pr-summary` supports local JSON input and read-only GitLink fetch mode.
- `json/table/markdown` are rendered inside the workflow package, not by the global formatter.
- Fetch-layer tests use `httptest` and do not depend on the real remote API.

## Conclusion

The rule-based Agent Workflow prototype, including the read-only fetch layer, is implemented, tested, and locally runnable.
