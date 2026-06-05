# GitLink Access Skill

Use this skill when an Agent needs to help a user apply to join a GitLink project or quit a project.

## Commands

```bash
# Preview a project join application
gitlink-cli access +join --code <invite_or_application_code> --role developer --dry-run

# Submit a project join application
gitlink-cli access +join --code <invite_or_application_code> --role reporter

# Preview leaving a project
gitlink-cli access +quit --owner <owner> --repo <repo> --dry-run

# Leave a project after explicit confirmation
gitlink-cli access +quit --owner <owner> --repo <repo>
```

## Safety Contract

- Always run `access +join` with `--dry-run` first and confirm the code and requested role.
- Always run `access +quit` with `--dry-run` first and confirm the owner/repo before leaving a project.
- Valid roles are `manager`, `developer`, and `reporter`; default is `developer`.
- Do not submit or quit without explicit user confirmation.

## API Mapping

| Command | Method | API path |
|---|---|---|
| `access +join` | POST | `/api/applied_projects.json` |
| `access +quit` | POST | `/api/{owner}/{repo}/quit.json` |
