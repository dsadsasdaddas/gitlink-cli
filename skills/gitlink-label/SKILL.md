---
name: gitlink-label
version: 1.0.0
description: "Issue label management: list, create, update, and delete GitLink issue labels (项目标记). Triggered when a user needs to manage labels, set up a triage taxonomy, or tag issues."
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli label --help"
---

# gitlink-label

**CRITICAL**: Read [`../gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md) before starting. It covers authentication, permissions, global flags, and GitLink API behavior.
**CRITICAL**: Confirm user intent before running write or destructive operations such as `+create`, `+update`, or `+delete`.
**CRITICAL**: Use `gitlink-cli` for GitLink resources. Do not use GitHub-only tools such as `gh`.

## Shortcuts

| Shortcut | Description | Operation |
|----------|-------------|-----------|
| `label +list` | List issue labels | Read |
| `label +create` | Create an issue label | Write |
| `label +update` | Update a label, preserving unspecified fields | Write |
| `label +delete` | Delete an issue label | Destructive |

## Examples

```bash
# List all labels
gitlink-cli label +list --owner Gitlink --repo forgeplus

# Filter labels by keyword, sorted by issue count
gitlink-cli label +list --owner Gitlink --repo forgeplus -k bug --sort-by issues_count --sort-direction desc

# Create a label (color defaults to #1E90FF when omitted)
gitlink-cli label +create --owner Gitlink --repo forgeplus -n bug -d "Something is broken" -c "#FF0000"

# Update only the color; name and description are preserved
gitlink-cli label +update --owner Gitlink --repo forgeplus -i 42 -c "#00FF00"

# Delete a label
gitlink-cli label +delete --owner Gitlink --repo forgeplus -i 42
```

## Parameters

| Command | Key parameters |
|---------|----------------|
| `+list` | `--keyword`, `--only-name`, `--sort-by` (updated_on / created_on / issues_count), `--sort-direction` (asc / desc) |
| `+create` | `--name` (required), `--description`, `--color` (hex, default `#1E90FF`) |
| `+update` | `--id` (required) plus at least one of `--name`, `--description`, `--color` |
| `+delete` | `--id` (required) |

## API Notes

- Labels map to the GitLink "项目标记" / `issue_tags` API: `/api/v1/{owner}/{repo}/issue_tags`.
- `--color` must be a hex value (`#RGB` or `#RRGGBB`); it is validated client-side before the API call.
- `+update` first fetches the label's current values from the list endpoint and merges the requested changes, so fields you do not pass are preserved (the API requires `name`, `description`, and `color` together).
- To attach a label to an issue, pass its id via the issue update API field `issue_tag_ids` (see `gitlink-issue`); use `label +list --only-name true` to resolve label ids quickly.

## Typical workflow: bootstrap a triage taxonomy

```bash
# Create a consistent label set for issue triage
gitlink-cli label +create -n bug         -c "#D73A4A" -d "Confirmed defect"
gitlink-cli label +create -n enhancement -c "#A2EEEF" -d "Feature request"
gitlink-cli label +create -n question     -c "#D876E3" -d "Needs clarification"
gitlink-cli label +create -n security     -c "#B60205" -d "Security-sensitive"

# Verify the taxonomy
gitlink-cli label +list --only-name true --format json
```
