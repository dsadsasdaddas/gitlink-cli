---
name: gitlink-code
version: 1.0.0
description: "Repository code browsing: files, entries, git trees, blobs, commits, blame, and tags. Use when users need to inspect GitLink repository source code metadata."
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli code --help"
---

# gitlink-code

Read `../gitlink-shared/SKILL.md` first for authentication, global flags, and safety rules.

Write operations must be confirmed by the user before execution. `code +delete-tag`
supports `--dry-run` and should be previewed before deleting a tag.

## Shortcuts

| Shortcut | Purpose | Auth |
|----------|---------|------|
| `code +files` | Search repository files | No for public repos |
| `code +entries` | List root code entries | No for public repos |
| `code +sub-entries` | Show a directory or file entry | No for public repos |
| `code +tree` | List a git tree by branch, tag, commit, or tree SHA | No for public repos |
| `code +blob` | Read a git blob by SHA | No for public repos |
| `code +commits` | List commits | No for public repos |
| `code +commit-files` | List files changed by one commit | No for public repos |
| `code +commit-diff` | Show one commit diff | No for public repos |
| `code +blame` | Show file blame data | No for public repos |
| `code +tags` | List repository tags | No for public repos |
| `code +tag` | Show tag details | No for public repos |
| `code +delete-tag` | Delete a tag | Yes |

## Examples

```bash
gitlink-cli code +files --owner Gitlink --repo forgeplus --search README --ref master
gitlink-cli code +entries --owner Gitlink --repo forgeplus --ref master
gitlink-cli code +sub-entries --owner Gitlink --repo forgeplus --path docs --ref master

gitlink-cli code +tree --owner Gitlink --repo forgeplus --sha master --recursive
gitlink-cli code +blob --owner Gitlink --repo forgeplus --sha <blob-sha>

gitlink-cli code +commits --owner Gitlink --repo forgeplus --sha master
gitlink-cli code +commit-files --owner Gitlink --repo forgeplus --sha <commit-sha>
gitlink-cli code +commit-diff --owner Gitlink --repo forgeplus --sha <commit-sha>
gitlink-cli code +blame --owner Gitlink --repo forgeplus --sha master --path README.md

gitlink-cli code +tags --owner Gitlink --repo forgeplus
gitlink-cli code +tag --owner Gitlink --repo forgeplus --name v1.0.0
gitlink-cli code +delete-tag --owner Gitlink --repo forgeplus --name v1.0.0 --dry-run
```

## API Mapping

| Shortcut | API |
|----------|-----|
| `code +files` | `GET /api/{owner}/{repo}/files.json` |
| `code +entries` | `GET /api/{owner}/{repo}/entries.json` |
| `code +sub-entries` | `GET /api/{owner}/{repo}/sub_entries.json` |
| `code +tree` | `GET /api/v1/{owner}/{repo}/git/trees/{sha}.json` |
| `code +blob` | `GET /api/v1/{owner}/{repo}/git/blobs/{sha}.json` |
| `code +commits` | `GET /api/v1/{owner}/{repo}/commits.json` |
| `code +commit-files` | `GET /api/v1/{owner}/{repo}/commits/{sha}/files.json` |
| `code +commit-diff` | `GET /api/v1/{owner}/{repo}/commits/{sha}/diff.json` |
| `code +blame` | `GET /api/v1/{owner}/{repo}/blame.json` |
| `code +tags` | `GET /api/v1/{owner}/{repo}/tags.json` |
| `code +tag` | `GET /api/v1/{owner}/{repo}/tags/{name}.json` |
| `code +delete-tag` | `DELETE /api/v1/{owner}/{repo}/tags/{tag}.json` |
