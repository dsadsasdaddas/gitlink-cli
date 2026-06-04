# GitLink Mirror Skill

Use this skill when an Agent needs to create GitLink mirror repositories or trigger manual mirror synchronization.

## Commands

```bash
# Preview creating a mirror repository
gitlink-cli mirror +create --user-id <user_or_org_id> \
  --name <display_name> \
  --repository-name <identifier> \
  --clone-addr <remote_git_url> \
  --dry-run

# Create a mirror for a private source
gitlink-cli mirror +create --user-id <user_or_org_id> \
  --name <display_name> \
  --repository-name <identifier> \
  --clone-addr <remote_git_url> \
  --auth-username <username> \
  --auth-password <token_or_password>

# Preview or run a manual mirror sync
gitlink-cli mirror +sync --repo-id <repository_id> --dry-run
gitlink-cli mirror +sync --repo-id <repository_id>
```

## Safety Contract

- Always run `mirror +create` with `--dry-run` first and confirm `user_id`, `repository_name`, and `clone_addr`.
- Dry-run output redacts `auth_password` as `***REDACTED***`.
- Do not echo credentials in chat transcripts or PR descriptions.
- Always run `mirror +sync --dry-run` before triggering a real sync.

## API Mapping

| Command | Method | API path |
|---|---|---|
| `mirror +create` | POST | `/api/projects/migrate.json` |
| `mirror +sync` | POST | `/api/repositories/{id}/sync_mirror.json` |
