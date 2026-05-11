# Safe Batch Issue Maintenance Example（安全批量 Issue 维护示例）

本示例展示人类用户或 AI Agent 如何通过 dry-run 确认流程，安全地批量关闭 GitLink Issues。

> 依赖说明：本示例依赖 PR [#12](https://gitlink.org.cn/Gitlink/gitlink-cli/pulls/12) 中新增的 `gitlink-cli issue +batch-close` 命令，或后续已经包含该命令的版本。

## 场景

维护者希望在确认批量计划后，关闭 stale、duplicate 或 already resolved 的 Issues。

## 文件

- `issues.csv`：示例 Issue 编号输入文件。

## Step 1：先使用 dry-run 预览

```bash
gitlink-cli issue +batch-close \
  --owner Gitlink \
  --repo forgeplus \
  --from issues.csv \
  --dry-run \
  --format json
```

预期行为：

- 不修改任何 Issue 状态。
- 输出结构化汇总结果。
- Agent 必须向用户展示结果并请求确认。

## Step 2：用户确认后再执行

只有在用户明确确认 dry-run 结果后，才能执行真实关闭：

```bash
gitlink-cli issue +batch-close \
  --owner Gitlink \
  --repo forgeplus \
  --from issues.csv \
  --format json
```

## Agent 检查清单

- [ ] 确认 `owner/repo`。
- [ ] 确认 Issue 编号或 CSV 来源。
- [ ] 先执行 `--dry-run`。
- [ ] 展示计划关闭的 Issue 编号和汇总。
- [ ] 等待用户明确确认。
- [ ] 用户确认后，才执行不带 `--dry-run` 的真实命令。
- [ ] 报告最终成功/失败数量。

## 注意事项

`issues.csv` 中的编号是占位示例。对真实仓库执行前，请替换为真实 Issue 编号。
