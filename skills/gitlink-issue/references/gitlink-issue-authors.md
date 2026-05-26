# Issue Authors

Use `issue +authors` to list users who have authored Issues in a repository.
This helps users find author IDs for Issue filtering and reporting workflows.

## Usage

```bash
gitlink-cli issue +authors --owner Gitlink --repo forgeplus
gitlink-cli issue +authors --owner Gitlink --repo forgeplus --keyword bob
```

## Options

| Option | Description |
|--------|-------------|
| `--owner` | Repository owner |
| `--repo` | Repository name |
| `--keyword`, `-k` | Optional search keyword |

## API

```http
GET /api/v1/{owner}/{repo}/issue_authors.json
```
