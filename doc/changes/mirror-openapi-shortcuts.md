# Mirror OpenAPI Shortcuts

This change adds `gitlink-cli mirror` shortcuts for mirror repository creation and manual synchronization.

## New Commands

```bash
gitlink-cli mirror +create --user-id 42 --name demo-mirror \
  --repository-name demo-mirror --clone-addr https://example.com/demo.git --dry-run

gitlink-cli mirror +sync --repo-id 99 --dry-run
```

## Rationale

GitLink OpenAPI supports mirror migration and manual mirror synchronization, but users previously had to call raw API endpoints. These shortcuts make repository migration and mirroring easier to use from both humans and Agents.

## Safety

- `mirror +create` and `mirror +sync` support `--dry-run`.
- `auth_password` is redacted in dry-run output.
- IDs and optional category/language fields are validated before making API calls.

## Tests

- Unit tests cover create payloads, dry-run behavior, credential validation, sync path generation, and ID validation.
- Verified with `go test ./...` and `go vet ./...`.
