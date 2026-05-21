# PR Draft: Add workflow agent commands for issue triage and repository health analysis

## Summary

This PR adds `workflow +triage` and `workflow +health` with three execution modes:

- local flags
- local JSON input
- read-only GitLink fetch mode

It also adds stable `json`, `table`, and `markdown` rendering for Agent consumption.

## Motivation

- Help maintainers triage issues faster
- Provide a structured repository health overview
- Give AI Agents stable machine-readable output
- Keep the workflow local-first and safe by default
- Avoid any dependency on external LLM APIs

## Changes

- New `shortcuts/workflow` rule engine and DTOs
- Local command layer for `workflow +triage` and `workflow +health`
- Workflow-local renderer for `json`, `table`, and `markdown`
- Read-only GitLink fetch and normalization helpers
- Unit tests for rules, fetch normalization, rendering, and command wiring
- Competition and test documentation updates

## Safety

- Remote mode is read-only
- No comment, label, close, merge, or release write actions
- Health scoring tolerates unknown or unavailable metrics
- Test fixtures do not contain secrets or tokens

## Tests

- `gofmt -w shortcuts/workflow/*.go shortcuts/register.go`
- `go test ./shortcuts/workflow`
- `go test ./...`
- `httptest` coverage for API normalization and fetch tolerance
- Manual command examples in local and remote read-only modes

## Known Limitations

- Real API response shapes may still require minor normalization tweaks
- Write operations are intentionally deferred to a later PR
- `pr-summary`, `release-notes`, and `stale` are planned next

## Screenshots or Examples

```bash
gitlink-cli workflow +triage --title "Token leaked in logs" --body "The access token appears in command output" --format json
gitlink-cli workflow +health --repository Gitlink/gitlink-cli --open-issues 3 --open-prs 1 --has-readme --has-license --has-contributing --agent-readiness-known --agent-readiness-score 9 --format table
gitlink-cli workflow +triage --owner Gitlink --repo gitlink-cli --state open --limit 5 --format table
```
