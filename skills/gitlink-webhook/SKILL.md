---
name: gitlink-webhook
version: 1.0.0
description: "Webhook 配置：列出、查看、创建、更新、删除、测试仓库 Webhook，并查看推送历史。当用户需要配置 GitLink Webhook 或外部系统集成时触发。"
metadata:
  requires:
    bins: ["gitlink-cli"]
  cliHelp: "gitlink-cli webhook --help"
---

# gitlink-webhook（Webhook 配置）

**CRITICAL — 开始前必须先阅读 [`../gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md)，其中包含认证、权限处理和 API 注意事项。**
**CRITICAL — 创建、更新、删除、测试 Webhook 都可能触发外部系统或修改仓库配置，执行前必须确认用户意图。**
**CRITICAL — GitLink 操作只能用 `gitlink-cli`。禁止用 `gh`（GitHub CLI）操作 GitLink 资源。`gh` 仅适用于 GitHub 平台。**

## Shortcuts

| Shortcut | 说明 | 需要认证 |
|----------|------|----------|
| `webhook +list` | 列出仓库 Webhook | 是 |
| `webhook +view` | 查看 Webhook 详情 | 是 |
| `webhook +create` | 创建 Webhook | 是 |
| `webhook +update` | 更新 Webhook | 是 |
| `webhook +delete` | 删除 Webhook | 是 |
| `webhook +test` | 触发测试推送 | 是 |
| `webhook +tasks` | 查看推送历史 | 是 |

## `--events` 支持的事件

`--events` 使用逗号分隔，例如 `push,issues_only`。当前 CLI 会校验以下事件名：

| 事件 | 说明 |
|------|------|
| `push` | 代码推送 |
| `create` | 创建分支或标签 |
| `delete` | 删除分支或标签 |
| `issues_only` | Issue 创建/变更主事件 |
| `issue_assign` | Issue 指派变更 |
| `issue_label` | Issue 标签变更 |
| `issue_comment` | Issue 评论 |
| `pull_request_only` | Pull Request 创建/变更主事件 |
| `pull_request_assign` | Pull Request 指派变更 |
| `pull_request_comment` | Pull Request 评论 |

示例：

```text
push,create,delete,issues_only,issue_assign,issue_label,issue_comment,pull_request_only,pull_request_assign,pull_request_comment
```

## 使用示例

```bash
# 查看 Webhook 列表
gitlink-cli webhook +list --owner myuser --repo myrepo

# 创建 Webhook
gitlink-cli webhook +create --owner myuser --repo myrepo \
  --url https://example.com/gitlink-hook \
  --events push,issues_only

# 更新 Webhook
gitlink-cli webhook +update --owner myuser --repo myrepo -i 123 \
  --url https://example.com/new-hook \
  --events push,issue_comment

# 测试 Webhook 并查看推送历史
gitlink-cli webhook +test --owner myuser --repo myrepo -i 123
gitlink-cli webhook +tasks --owner myuser --repo myrepo -i 123

# 删除 Webhook
gitlink-cli webhook +delete --owner myuser --repo myrepo -i 123
```

## 安全流程

1. 确认目标仓库 `owner/repo`。
2. 创建或更新前确认 URL、事件列表、secret 和分支过滤规则。
3. 删除前先执行 `webhook +view`，确认 ID 对应的 Webhook。
4. `webhook +test` 可能向外部系统发送请求，执行前先确认用户意图。
