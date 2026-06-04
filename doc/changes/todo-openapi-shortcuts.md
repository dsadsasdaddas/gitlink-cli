# Todo OpenAPI Shortcuts

This change adds `gitlink-cli todo` shortcuts for GitLink user todo queues.

## New Commands

```bash
gitlink-cli todo +transfer-list --login wangyue111
gitlink-cli todo +transfer-accept --login wangyue111 --id 7 --dry-run
gitlink-cli todo +transfer-refuse --login wangyue111 --id 7 --dry-run
gitlink-cli todo +join-list --login wangyue111 --page 1 --per-page 20
gitlink-cli todo +join-accept --login wangyue111 --id 11 --dry-run
gitlink-cli todo +join-refuse --login wangyue111 --id 11 --dry-run
```

## Rationale

Project transfer and project join requests are real maintainer workflows, but previously they required raw API calls. The new shortcuts make request review easier for humans and safer for Agents.

## Safety

Write operations support `--dry-run`, returning the method, path, login, request ID, resource, and action without changing data.

## Tests

- Unit tests cover list paths/query parameters, accept/refuse paths, dry-run behavior, and request ID validation.
- Verified with `go test ./...` and `go vet ./...`.
