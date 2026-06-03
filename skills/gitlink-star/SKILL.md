# GitLink Star Skill

Use this skill when the user wants to inspect or maintain starred/pinned projects for a GitLink user profile.

## Commands

### List starred projects

```bash
gitlink-cli star +list --login <user-login> --format json
```

OpenAPI mapping:

```text
GET /api/users/{owner}/is_pinned_projects.json
```

Notes:

- `--login` is the GitLink user login in the profile URL.
- The response contains `projects[].id` and `projects[].project_id`.
- `projects[].id` is the starred/pinned record ID used by `star +reorder`.
- `projects[].project_id` is the repository/project ID used by `star +set`.

### Set starred project IDs

```bash
gitlink-cli star +set --login <user-login> --project-ids 17,42 --dry-run --format json
gitlink-cli star +set --login <user-login> --project-ids 17,42 --format json
```

OpenAPI mapping:

```text
POST /api/users/{owner}/is_pinned_projects/pin.json
Body: { "is_pinned_project_ids": [17, 42] }
```

Safety rules:

1. For write operations, run `--dry-run` first.
2. Show the dry-run payload to the user.
3. Execute the non-dry-run command only after explicit confirmation.
4. Treat `--project-ids` as the desired saved starred project ID list from the OpenAPI field `is_pinned_project_ids`.

### Reorder a starred project

```bash
gitlink-cli star +reorder --login <user-login> --pinned-id 9 --position 10 --dry-run --format json
gitlink-cli star +reorder --login <user-login> --pinned-id 9 --position 10 --format json
```

OpenAPI mapping:

```text
PUT /api/users/{owner}/is_pinned_projects/{id}.json
Body: { "pinned_project": { "position": 10 } }
```

Notes:

- `--pinned-id` is the starred record ID from `star +list` (`projects[].id`), not `project_id`.
- Larger `--position` values rank higher according to the OpenAPI documentation.
- `--position` accepts non-negative integers.

## Typical workflow

```bash
# 1. Inspect current starred projects
gitlink-cli star +list --login wangyue111 --format json

# 2. Preview a new starred project list
gitlink-cli star +set --login wangyue111 --project-ids 17,42 --dry-run --format json

# 3. Apply after user confirmation
gitlink-cli star +set --login wangyue111 --project-ids 17,42 --format json

# 4. Preview reordering a starred record
gitlink-cli star +reorder --login wangyue111 --pinned-id 9 --position 10 --dry-run --format json
```

## Failure handling

- If `star +reorder` fails with 404, re-run `star +list` and confirm the `projects[].id` value.
- If `star +set` receives invalid IDs, verify `project_id` values from `star +list` or project detail output.
- Do not guess starred record IDs from project IDs.
