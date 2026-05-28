# org team project binding shortcuts

> **前置条件：** 先阅读 [`../gitlink-shared/SKILL.md`](../../gitlink-shared/SKILL.md) 了解认证、全局参数和安全规则。

用于批量维护组织团队与项目的绑定关系。

## 命令

```bash
# 预览：将组织下全部项目加入指定团队
gitlink-cli org +team-projects-add-all --organization Gitlink --team-id 7 --dry-run --format json

# 执行：将组织下全部项目加入指定团队
gitlink-cli org +team-projects-add-all --organization Gitlink --team-id 7

# 预览：从指定团队移除全部项目
gitlink-cli org +team-projects-remove-all --organization Gitlink --team-id 7 --dry-run --format json

# 执行：从指定团队移除全部项目
gitlink-cli org +team-projects-remove-all --organization Gitlink --team-id 7
```

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `--organization` / `-o` | 是 | 组织标识，例如 `Gitlink` |
| `--team-id` / `-t` | 是 | 团队 ID |
| `--dry-run` | 否 | 只预览 method/path/organization/team_id，不修改线上数据 |
| `--format` | 否 | 输出格式：json / table / yaml |

## 安全流程

1. 先确认组织标识和团队 ID。
2. 执行 `--dry-run --format json` 展示将要调用的接口。
3. 用户确认后再执行真实命令。
4. 如果是删除/移除操作，必须再次提醒该操作会影响团队下全部项目绑定。

## OpenAPI 对齐

- `POST /api/organizations/{organization}/teams/{id}/team_projects/create_all.json`
- `DELETE /api/organizations/{organization}/teams/{id}/team_projects/destroy_all.json`

## References

- [gitlink-org](../SKILL.md)
- [gitlink-shared](../../gitlink-shared/SKILL.md)
