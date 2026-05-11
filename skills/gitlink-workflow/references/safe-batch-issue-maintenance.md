# Safe Batch Issue Maintenance Workflow（安全批量 Issue 维护工作流）

> **前置条件：** 先阅读 [`../../gitlink-shared/SKILL.md`](../../gitlink-shared/SKILL.md) 了解认证、全局参数、安全规则和写操作确认要求。

本工作流用于指导 AI Agent 安全地批量关闭 stale、duplicate 或 resolved 状态的 GitLink Issues，底层命令为 `gitlink-cli issue +batch-close`。

> 依赖说明：本工作流依赖 PR [#12](https://gitlink.org.cn/Gitlink/gitlink-cli/pulls/12) 中新增的 `gitlink-cli issue +batch-close` 命令，或后续已经包含该命令的版本。

## 适用场景

当用户提出以下需求时，Agent 可以使用本工作流：

- 批量关闭 stale / inactive Issues
- 批量关闭 duplicate Issues
- 批量关闭已经修复或已完成的 Issues
- 根据 CSV、表格、脚本输出批量处理 Issues
- 在 Release / Sprint 收尾时批量关闭 Issues

## 安全合约

批量 Issue 维护属于写操作。Agent 必须遵守以下安全合约：

1. 明确目标仓库，并向用户说明 `--owner` 和 `--repo`。
2. 明确 Issue 编号来源：`--numbers` 或 `--from <csv>`。
3. 必须先执行 `gitlink-cli issue +batch-close ... --dry-run --format json`。
4. 必须向用户展示 dry-run 汇总，包括 `total`、`succeeded`、`failed` 和计划处理的 Issue 编号。
5. 必须等待用户明确确认后，才能去掉 `--dry-run`。
6. 只有在用户确认后，才能执行真实批量关闭命令。
7. 执行后必须报告最终汇总和失败项。

> [!CAUTION]
> Agent 不得在未执行 dry-run、未获得用户明确确认的情况下直接执行真实批量关闭。

## 工作流步骤

### Step 1：确认目标仓库

如果当前目录是 GitLink 仓库，owner/repo 可以自动解析。否则应显式指定：

```bash
gitlink-cli repo +info --owner <owner> --repo <repo> --format json
```

Agent 需要向用户说明即将操作的仓库，例如：

```text
目标仓库：<owner>/<repo>
```

### Step 2：准备 Issue 输入

可以直接使用 Issue 编号列表：

```bash
ISSUE_NUMBERS="101,102,103"
```

也可以准备 CSV 文件：

```csv
number,reason
101,stale
102,duplicate
103,resolved
```

CSV 支持 `number`、`issue_number`、`project_issues_index` 作为 Issue 编号列名。如果没有表头，则默认读取第一列。

### Step 3：先执行 dry-run

Issue 编号列表方式：

```bash
gitlink-cli issue +batch-close \
  --owner <owner> \
  --repo <repo> \
  --numbers "$ISSUE_NUMBERS" \
  --dry-run \
  --format json
```

CSV 方式：

```bash
gitlink-cli issue +batch-close \
  --owner <owner> \
  --repo <repo> \
  --from issues.csv \
  --dry-run \
  --format json
```

### Step 4：展示 dry-run 结果

Agent 需要在真实执行前向用户展示计划，例如：

```text
仓库：<owner>/<repo>
Dry-run：true
计划处理数量：3
计划关闭的 Issue：101, 102, 103
预检失败数量：0

请确认是否继续执行真实批量关闭。
```

### Step 5：用户确认后执行

只有在用户明确确认后，才能执行真实命令。

Issue 编号列表方式：

```bash
gitlink-cli issue +batch-close \
  --owner <owner> \
  --repo <repo> \
  --numbers "$ISSUE_NUMBERS" \
  --format json
```

CSV 方式：

```bash
gitlink-cli issue +batch-close \
  --owner <owner> \
  --repo <repo> \
  --from issues.csv \
  --format json
```

### Step 6：输出最终报告

最终报告应包含：

- 目标仓库
- Issue 总数
- 成功数量
- 失败数量
- 每个失败 Issue 编号及错误原因
- 是否需要后续处理

## 输出结构示例

Dry-run 输出：

```json
{
  "ok": true,
  "data": {
    "repository": "owner/repo",
    "dry_run": true,
    "total": 3,
    "succeeded": 3,
    "failed": 0,
    "results": [
      {"id": "101", "action": "close", "status": "planned"},
      {"id": "102", "action": "close", "status": "planned"},
      {"id": "103", "action": "close", "status": "planned"}
    ]
  }
}
```

真实执行输出：

```json
{
  "ok": true,
  "data": {
    "repository": "owner/repo",
    "dry_run": false,
    "total": 3,
    "succeeded": 3,
    "failed": 0,
    "results": [
      {"id": "101", "action": "close", "status": "closed"},
      {"id": "102", "action": "close", "status": "closed"},
      {"id": "103", "action": "close", "status": "closed"}
    ]
  }
}
```

## Agent 应答模板

当用户说“帮我关闭这些 Issue”或“批量关闭 stale Issues”时，Agent 应先说明计划：

```text
我会先对仓库 <owner>/<repo> 的 Issue 编号 <numbers> 执行 dry-run 预览，不会修改任何数据。你确认 dry-run 结果后，我再执行真实批量关闭。
```

dry-run 完成后，Agent 应询问用户：

```text
Dry-run 已完成。计划关闭 3 个 Issue：101、102、103。未发现失败项。请确认是否现在执行真实关闭。
```

## 失败处理

如果 `failed > 0`：

1. 不要声称全部成功。
2. 列出失败的 Issue 编号 和错误信息。
3. 建议用户检查权限、Issue 是否存在、owner/repo 是否正确。
4. 如果只有部分失败，需要分别报告成功项和失败项。

## 参考

- [`gitlink-workflow`](../SKILL.md)
- [`gitlink-issue`](../../gitlink-issue/SKILL.md)
- [`issue +batch-close` PR #12](https://gitlink.org.cn/Gitlink/gitlink-cli/pulls/12)
- [`gitlink-shared`](../../gitlink-shared/SKILL.md)
