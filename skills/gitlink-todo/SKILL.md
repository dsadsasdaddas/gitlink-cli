# GitLink Todo Skill

Use this skill when an Agent needs to inspect or approve user todo queues on GitLink, especially project transfer requests and project join requests.

## Commands

```bash
# List project transfer requests for a user
gitlink-cli todo +transfer-list --login <user> --page 1 --per-page 20

# Accept or refuse a transfer request; preview writes first
gitlink-cli todo +transfer-accept --login <user> --id <request_id> --dry-run
gitlink-cli todo +transfer-refuse --login <user> --id <request_id> --dry-run

# List project join requests for a user
gitlink-cli todo +join-list --login <user> --page 1 --per-page 20

# Accept or refuse a join request; preview writes first
gitlink-cli todo +join-accept --login <user> --id <request_id> --dry-run
gitlink-cli todo +join-refuse --login <user> --id <request_id> --dry-run
```

## Safety Contract

- Always list requests before accepting or refusing them.
- For `transfer-accept`, `transfer-refuse`, `join-accept`, and `join-refuse`, run with `--dry-run` first and show the target `login`, `id`, `resource`, and `action` to the user.
- Do not approve a request if the project, requester, or target owner cannot be identified from the list output.

## API Mapping

| Command | Method | API path |
|---|---|---|
| `todo +transfer-list` | GET | `/api/users/{owner}/applied_transfer_projects.json` |
| `todo +transfer-accept` | POST | `/api/users/{owner}/applied_transfer_projects/{id}/accept.json` |
| `todo +transfer-refuse` | POST | `/api/users/{owner}/applied_transfer_projects/{id}/refuse.json` |
| `todo +join-list` | GET | `/api/users/{owner}/applied_projects.json` |
| `todo +join-accept` | POST | `/api/users/{owner}/applied_projects/{id}/accept.json` |
| `todo +join-refuse` | POST | `/api/users/{owner}/applied_projects/{id}/refuse.json` |
