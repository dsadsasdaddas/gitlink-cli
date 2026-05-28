---
name: gitlink-org
version: 1.0.0
description: "组织管理：查看组织列表、详情、成员，创建组织，批量维护团队项目绑定。当用户需要操作 GitLink 组织时触发。"
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli org --help"
---

# gitlink-org（组织操作）

**CRITICAL — 开始前必须先阅读 [`../gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md)，其中包含认证、权限处理和 API 注意事项。**
**CRITICAL — 所有 Shortcuts 在执行写入/删除操作前，务必先确认用户意图。**
**CRITICAL — GitLink 操作只能用 `gitlink-cli`。禁止用 `gh`（GitHub CLI）操作 GitLink 资源。`gh` 仅适用于 GitHub 平台。**

> **前置条件：** 先阅读 [`../gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md)

## Shortcuts

| Shortcut | 说明 |
|----------|------|
| `org +list` | 组织列表 |
| `org +info` | 组织详情 |
| `org +members` | 成员列表 |
| `org +create` | 创建组织 |
| `org +team-projects-add-all` | 将组织下全部项目加入指定团队 |
| `org +team-projects-remove-all` | 从指定团队移除全部项目 |

## 使用示例

```bash
gitlink-cli org +list
gitlink-cli org +info --id Gitlink
gitlink-cli org +members --id Gitlink
gitlink-cli org +create --name my-org --description "我的组织"

# 写操作建议先 dry-run
gitlink-cli org +team-projects-add-all --organization Gitlink --team-id 7 --dry-run
gitlink-cli org +team-projects-remove-all --organization Gitlink --team-id 7 --dry-run
```

## 团队项目绑定工作流

> [!CAUTION]
> `team-projects-add-all` 和 `team-projects-remove-all` 都会批量修改团队项目绑定，真实执行前必须先使用 `--dry-run` 并确认组织标识与团队 ID。

```bash
# 预览：将组织所有项目加入团队
gitlink-cli org +team-projects-add-all --organization Gitlink --team-id 7 --dry-run --format json

# 执行
gitlink-cli org +team-projects-add-all --organization Gitlink --team-id 7

# 预览：从团队移除所有项目
gitlink-cli org +team-projects-remove-all --organization Gitlink --team-id 7 --dry-run --format json
```

## Raw API 补充

```bash
# 其他组织团队管理接口仍可通过 Raw API 访问
gitlink-cli api GET /organizations/:id/teams
gitlink-cli api POST /organizations/:id/teams --body '{"name":"dev-team"}'

# 移除成员
gitlink-cli api DELETE /organizations/:id/organization_users/:uid
```

## References

- [team project binding](references/gitlink-org-team-projects.md)
