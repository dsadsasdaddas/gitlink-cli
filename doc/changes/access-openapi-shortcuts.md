# Access OpenAPI Shortcuts

This change adds `gitlink-cli access` shortcuts for project access self-service workflows.

## New Commands

```bash
gitlink-cli access +join --code <invite-code> --role developer --dry-run
gitlink-cli access +quit --owner Gitlink --repo forgeplus --dry-run
```

## Rationale

GitLink OpenAPI already supports project join applications and project quit actions, but users previously needed raw API calls. These shortcuts make the applicant side of the access workflow easier to automate and pair well with `todo` approval shortcuts.

## Safety

- `access +join` supports `--dry-run` and validates requested role.
- `access +quit` supports `--dry-run` and shows the resolved owner/repo before any write.
- Write commands should be executed only after user confirmation.

## Tests

- Unit tests cover join payload generation, default role, role validation, dry-run behavior, and quit path generation.
- Verified with `go test ./...` and `go vet ./...`.
