# gitlink-gatekeeper — 故障排查（TROUBLESHOOTING）

**CRITICAL — 开始前请先阅读 [`gitlink-shared/SKILL.md`](../gitlink-shared/SKILL.md)（认证、全局参数、真实 API 坑）与 [`./SKILL.md`](./SKILL.md)（工作流）、[`./REFERENCE.md`](./REFERENCE.md)（策略字段与评分算法）。**
**CRITICAL — gatekeeper 默认 dry-run，绝不自动合并。本文档中任何「修复」都不应放宽该安全默认，除非用户明确要求。**
**CRITICAL — GitLink 资源只能用 `gitlink-cli` 操作，禁止用 `gh`/`glab`。**

本文档列出运行 gatekeeper 时的常见问题，按「症状 / 原因 / 解决」三段式给出。先看下表速查，再到对应小节读细节。

## 速查表

| # | 症状 | 根因 | 一句话解决 |
|---|------|------|-----------|
| 1 | 提示「未找到策略文件」 | 仓库根无 `gatekeeper.yaml` 且未传 `--policy` | 回退内置默认策略，或 `--policy <path>` 指定 |
| 2 | 报「weights 之和 ≠ 100」 | 五维权重总和不是 100 | 调整权重让 `review+test+hygiene+commit+ci = 100` |
| 3 | `401`（请登录）/ `403`（无权限） | Token 过期 / owner-repo 错或无权限 | `gitlink-cli auth login` / 核对 owner、repo、权限 |
| 4 | `pr +diff` 输出巨大、超出处理窗口 | 大型 PR 的 diff 一次性返回太大 | 改用 `pr +files` 选关键文件 + `version-diff -f` 按文件分段 |
| 5 | 评分卡里 CI 显示 `unknown` | `ci +builds` 取不到与该 PR 对应的构建 | 按 unknown 计 0.5×权重，不要当失败 |
| 6 | 如何表达 REQUEST_CHANGES 裁决 | 自动门禁刻意不替人按 `approved`/`rejected` | 用建议性 `common` 评论 + 标题标注裁决 + 打 `gatekeeper:needs-changes` 标签表达 |
| 7 | 打标签报「标签不存在」 | 该 label 尚未在仓库创建 | 先 `label +list` 查，缺则 `label +create` 再挂到 issue |
| 8 | `auto_merge` 开了却没合并 | 三条件未同时满足 | 需 `verdict==PASS` + 命令带 `--apply` + 策略 `auto_merge: true` |
| 9 | 评分卡对 draft PR 给出裁决 | 草稿 PR 不应被门禁裁决 | 检测 draft 标志，提示先转 Ready for Review |
| 10 | 大仓库只看到前 20 条 PR/构建/标签 | 列表接口默认分页 limit=20 | 用 `--page` / `--limit` 翻页，按 `meta.total_count` 判完整 |
| 11 | `--state open` 仍返回已合并/关闭的 PR | GitLink `--state` 仅影响计数 | 用返回里的 `pull_request_status` 客户端过滤 |
| 12 | `pr +merge` 报 `--do` 不是已知 flag | flag 名记错（`do` 是底层 API 字段） | CLI flag 是 `--method`/`-m`，不是 `--do` |

---

## 1. 找不到 `gatekeeper.yaml`（回退默认）

**症状**：运行时提示「未在仓库根目录找到 `gatekeeper.yaml`」，或评分卡页脚 `policy:` 显示为内置默认而非你期望的文件。

**原因**：策略文件默认从**仓库根目录**的 `gatekeeper.yaml` 读取。当根目录没有该文件、且命令未带 `--policy <path>` 时，gatekeeper 不会报错中止，而是按设计**回退到内置默认策略**（见 `REFERENCE.md` 的默认值：`pass=85 / request_changes=60`，权重 `40/20/15/15/10`）。

**解决**：
- 若你确实想用默认策略 —— 这是正常行为，无需处理；评分卡页脚会标 `policy: built-in@v1`。
- 若你有自定义策略 —— 确认文件确实在仓库根，或显式指定路径：

```bash
# 显式指定策略文件（路径相对当前工作目录）
gitlink-gatekeeper --policy ./policies/gatekeeper.strict.yaml --pr 42

# 把示例策略复制到仓库根，让默认查找命中
cp skills/gitlink-gatekeeper/examples/gatekeeper.yaml ./gatekeeper.yaml
```

> 注意：策略文件路径区分大小写；`Gatekeeper.yaml`、`gatekeeper.yml`（`.yml` 后缀）都不会被默认查找命中。

---

## 2. `weights` 之和不等于 100

**症状**：加载策略时报错「weights 之和必须为 100，当前为 `<n>`」，gatekeeper 拒绝按该策略评分。

**原因**：评分算法要求五个维度权重之和**严格等于 100**，这样 `total ∈ 0..100` 才有可比性与可复现性。常见错误是只改了一两个维度、忘了让其余维度补平。

**解决**：调整 `weights`，使五项相加正好 100。

```yaml
weights:
  review_findings: 40
  test_coverage:   20
  pr_hygiene:      15
  commit_quality:  15
  ci_status:       10   # 40+20+15+15+10 = 100 ✅
```

| 错误示例 | 和 | 问题 |
|----------|:--:|------|
| `40/20/15/15/5` | 95 | 少 5，需补到某一维度 |
| `50/20/15/15/10` | 110 | 多 10，把 review 降回 40 |

> 注意：维度名必须是这五个固定键（`review_findings`/`test_coverage`/`pr_hygiene`/`commit_quality`/`ci_status`），多写或漏写键也会校验失败。某维度想「不计分」应把权重设 0 并把差额加到别处，而不是删除该键。

---

## 3. `401` / `403` 认证与权限错误

**症状**：任意数据采集命令返回
```json
{"ok": false, "error": {"code": 401, "message": "请登录后再操作", "suggestion": "请先运行 gitlink-cli auth login 登录"}}
```
或 `403`（拒绝访问）。

**原因**（详见 `gitlink-shared/SKILL.md`「认证错误处理」）：
- `401`：未登录或 Token 已过期。GitLink Token 有效期 **7 天**，过期需重新登录。
- `403`：已登录但对该 `owner/repo` 无权限，或 owner/repo 解析错误。

**解决**：

```bash
# 401：重新登录
gitlink-cli auth login
gitlink-cli auth status        # 确认登录态

# 403：核对上下文（在仓库目录下可从 git remote 自动解析）
gitlink-cli repo +info --owner <owner> --repo <repo> --format json
```

> 注意：gatekeeper 全程**不回显 Token**。若 `auth status` 正常但仍 `403`，多半是 owner/repo 写错或你只有读权限——读权限足够采集与生成评分卡（dry-run），但写评论、打标签、合并需要写权限。

---

## 4. PR diff 过大需分段

**症状**：`gitlink-cli pr +diff -i <id> --format json` 返回内容极大，超出单次处理窗口，或采集很慢/报频率限制。

**原因**：大型 PR 的全量 diff 一次返回会非常大（`gitlink-shared` 与 `gitlink-code-review` 均提示 `pr +diff` 输出可能很大，需分段处理；files/diff 接口有频率限制）。

**解决**：先用 `pr +files` 拿到文件清单，只对**与评分相关**的源码/测试文件按文件取 diff：

```bash
# 1. 先取文件清单（轻量），用于 test_coverage 维度与体量判断
gitlink-cli pr +files -i <id> --format json

# 2. 对单个文件取差异（避免一次拉全量）
gitlink-cli pr +version-diff -i <id> -v <version-id> -f path/to/file.go --format json
#   version-id 从 pr +versions -i <id> 获取（patchset 版本）
```

> 注意：`test_coverage`、`pr_hygiene` 的体量判定只依赖**文件清单**（`changed_src`/`changed_tests`/`changed_files`），不需要全量 diff；只有 `review_findings`（AI 审查）才需要看 diff 内容。所以分段时优先保证文件清单完整，diff 可按文件懒加载。避免短时间内重复请求 files/diff 接口。

---

## 5. CI 状态取不到（按 `unknown` 处理）

**症状**：评分卡 `CI status` 行显示 `unknown`，得分为权重的一半（默认 `10 → 5`）。

**原因**：`gitlink-cli ci +builds` 返回的是仓库的构建列表，并不保证能定位到**正好对应当前 PR 头部 commit** 的那次构建——仓库可能没配 CI、构建尚未触发、或无法把构建与该 PR 的 commit 关联。此时 CI 既非「通过」也非「失败」，而是**未知**。

**原因细节与评分映射**（见 `REFERENCE.md` 3.5）：

| CI 情况 | ci_status 得分（权重 10 时） |
|---------|:---:|
| 找到对应构建且通过 | 10 |
| 找到对应构建且失败 | 0 |
| 取不到 / 无 CI / 无法关联 | 5（`round(10 * 0.5)`） |

**解决**：
- 这是**预期的降级行为**，不是 bug：取不到就当 `unknown`，给一半分，**不要**当成失败而误触发 `require_ci_pass` 硬门禁。
- 若希望 `unknown` 也拦截，可在策略里收紧（但要清楚这会拦下没配 CI 的仓库）。默认 `require_ci_pass: true` 的语义是「**CI 明确失败**才拦」，`unknown` 不触发硬门禁。

```bash
# 排查：先看仓库到底有没有构建记录
gitlink-cli ci +builds --limit 20 --format json
```

> 注意：硬门禁 `require_ci_pass` 只在 CI **明确失败**时命中；`unknown` 不算失败、不触发硬门禁，仅按 0.5×权重计分。

---

## 6. 如何表达 REQUEST_CHANGES 裁决（为何用评论而非 rejected）

**症状**：纠结要不要用 GitLink review 的 `rejected` 状态把 gatekeeper 的 `REQUEST_CHANGES` 裁决"硬"标到 PR 上。

**原因 / 设计选择**：GitLink 原生 PR review 的 `status` 为 `common`/`approved`/`rejected`（`rejected` 即"请求修改"的等价）。但 `approved`/`rejected` 是**强语义的人工授意动作**。gatekeeper 是自动门禁，**刻意不替人按下 approved/rejected** —— 所有自动裁决（含 REQUEST_CHANGES）一律以**建议性的 `common` 评论**回写，把强语义留给维护者。

**解决**：用「评论 + 标签」组合表达裁决（这正是本作品复用 `label` 命令的原因）：
1. 评分卡以 `pr +comment`（或 `pr +review --status common`）回写，**标题里明确标注裁决**（`Verdict: ❌ REQUEST_CHANGES`），一眼可读；
2. 同时打上状态标签 `gatekeeper:needs-changes`，让状态可被列表/过滤识别。

```bash
# 评分卡回写（标题已含裁决，正文是评分卡）—— 与 workflow 脚本一致，走 pr +comment
gitlink-cli pr +comment -i <id> -b "$(cat scorecard.md)"
#   也可作为评审记录：gitlink-cli pr +review -i <id> --status common --content "$(cat scorecard.md)"
```

> 注意：gatekeeper **绝不**自动发 `approved`/`rejected` —— 那是维护者的权限。裁决的「拦截」语义靠标题文字 + `gatekeeper:needs-changes` 标签承载。

---

## 7. `label` 不存在，需先 `create`

**症状**：打标签时提示标签不存在，或 `issue_tag_ids` 里传了一个查不到的 ID。

**原因**：`gatekeeper.yaml` 里配置的 `labels`（`gatekeeper:pass` / `gatekeeper:needs-changes` / `gatekeeper:review`）只是名字，仓库里**未必已经创建**对应的 label 实体。给 PR/Issue 挂标签时用的是 label 的**数字 ID**，名字对不上数据库里就没有，自然挂不上。

**解决**：先查后建——先用 `label +list` 看标签是否存在并拿到 ID，缺的用 `label +create` 创建（注意 `label` 组**只有** `+list/+create/+update/+delete`，**没有** `+view`）：

```bash
# 1. 查标签是否已存在，拿 id（--only-name true 只返回 id+name，便于解析）
gitlink-cli label +list --keyword gatekeeper --only-name true --format json

# 2. 缺失则创建（color 为十六进制，缺省有内置默认色）
gitlink-cli label +create -n "gatekeeper:needs-changes" -d "Gatekeeper requested changes" -c "#D73A4A" --format json
gitlink-cli label +create -n "gatekeeper:pass"          -d "Gatekeeper passed"            -c "#0E8A16" --format json
gitlink-cli label +create -n "gatekeeper:review"        -d "Gatekeeper left comments"     -c "#FBCA04" --format json

# 3. 把标签 id 挂到 PR 背后的 issue（用 PR 关联的 issue.id，非 PR 号；先取：）
ISSUE_ID=$(gitlink-cli pr +view -i <pr> --format json | jq -r '.data.issue.id')
#    更新时需带 done_ratio/subject/description，见 gitlink-shared
gitlink-cli api POST /:owner/:repo/issues/$ISSUE_ID --body '{
  "issue_tag_ids": [<tag_id>],
  "done_ratio": 0,
  "subject": "<原始标题>",
  "description": "<原始描述>"
}'
#    也可用 gitlink-cli issue +update，它会自动保留 subject/description，免手动回传。
```

> 注意：标签是**幂等创建**——重复 `+create` 同名标签前应先 `+list` 检查，避免产生重名标签。更新 issue 挂标签时务必带上当前 `subject`/`description`，否则可能清空描述（`gitlink-shared` 已警示）。

---

## 8. `auto_merge` 想生效但没合并（需三者同时满足）

**症状**：策略里写了 `auto_merge: true`，跑完却没有合并 PR。

**原因**：这是**故意的安全设计**，不是 bug。gatekeeper 绝不轻易合并，合并必须**三个条件同时成立**：

| 条件 | 来源 | 缺了会怎样 |
|------|------|-----------|
| `verdict == PASS` | 评分结果 | 非 PASS（COMMENT/REQUEST_CHANGES）一律不合并 |
| 命令带 `--apply` | 运行参数 | 不带则全程 dry-run，只打印不写不合并 |
| 策略 `auto_merge: true` | `gatekeeper.yaml` | 默认 `false`，不会合并 |

三者缺一不可。最常见的是**忘了 `--apply`**（默认 dry-run），或裁决其实不是 PASS。

**解决**：

```bash
# 先 dry-run 看裁决是不是 PASS（默认就是 dry-run，不写任何东西）
gitlink-gatekeeper --pr <id>

# 确认 PASS 且策略 auto_merge: true 后，显式带 --apply 才会合并
gitlink-gatekeeper --pr <id> --apply

# 底层合并命令（CLI flag 是 --method，不是 --do，详见 #12）
gitlink-cli pr +merge -i <id> --method squash --format json
```

> 注意：`merge_method`（`merge`/`rebase`/`squash`）由策略 `behavior.merge_method` 决定，映射到 `pr +merge --method`。即便三条件都满足，合并前也应向用户复述「将以 `<method>` 合并 PR #<id>」。

---

## 9. 草稿（draft）PR

**症状**：对一个还是草稿状态的 PR 跑出了正式裁决；作者反馈「我还没写完」。

**原因**：草稿 PR 通常尚未完成（描述未补、测试未加、commit 未整理），此时打门禁裁决意义不大，且容易误伤——`gitlink-code-review` 也建议对 draft PR 先提示作者转 Ready for Review。

**解决**：采集 PR 元信息后检查草稿标志（`pr +view` 返回里的 draft / WIP 标记，或标题含 `WIP`/`[Draft]`），命中则**不出裁决**，只提示：

```bash
gitlink-cli pr +view -i <id> --format json   # 检查是否 draft / 标题含 WIP
```

> 注意：可在 dry-run 下对 draft PR 生成「预览评分卡」帮作者自查，但**不回写评论、不打标签、不合并**，直到 PR 转为 Ready for Review。

---

## 10. 分页 / 大仓库（列表只看到一页）

**症状**：大仓库里 `pr +list` / `ci +builds` / `label +list` 只返回 20 条，漏掉了你要找的 PR、构建或标签。

**原因**：列表类接口默认分页，`limit` 默认 20、`page` 默认 1。返回 Envelope 的 `meta` 里有 `page`/`limit`/`total_count`，但**不会自动翻页**。

**解决**：按 `meta.total_count` 判断是否还有下一页，循环翻页直到取全：

```bash
# 第一页，先看 meta.total_count
gitlink-cli pr +list --state open --page 1 --limit 50 --format json
# 还有更多则继续翻页
gitlink-cli pr +list --state open --page 2 --limit 50 --format json

# 构建、标签同理
gitlink-cli ci +builds --page 1 --limit 50 --format json
gitlink-cli label +list --page 1 --limit 50 --format json
```

> 注意：把 `limit` 适当调大（如 50）可减少请求轮次，但 files/diff 接口有频率限制，翻页时不要过于密集。评分只针对**单个 PR**，分页主要用于「先在列表里定位到目标 PR 号」这一步。

---

## 11. 用 `--state` 过滤 PR 不准

**症状**：`gitlink-cli pr +list --state open` 返回的列表里混进了已合并 / 已关闭的 PR。

**原因**：GitLink 的真实行为（`gitlink-shared` 已记录）——`--state` 参数**仅影响统计计数**，返回的列表可能包含所有状态。不能只信 `--state`。

**解决**：在客户端用每条 PR 的 `pull_request_status` 字段二次过滤：

| `pull_request_status` | 含义 |
|:---:|------|
| `0` | open |
| `1` | merged |
| `2` | closed |

gatekeeper 只对 `pull_request_status == 0`（open）的 PR 做裁决；对已 merged/closed 的应跳过并提示。

> 注意：这是平台行为，不是 CLI bug。任何依赖「PR 是否仍 open」的逻辑（如批量门禁扫描）都必须以 `pull_request_status` 为准，而非 `--state`。

---

## 12. `pr +merge` 的 flag 是 `--method`，不是 `--do`

**症状**：执行 `gitlink-cli pr +merge -i <id> --do squash` 报错「未知 flag `--do`」。

**原因**：`do` 是 GitLink 合并 API（`POST .../pulls/:id/pr_merge`）的**底层请求字段名**；而 `gitlink-cli` 暴露给用户的 **flag 名是 `--method`（短选项 `-m`）**，默认值 `merge`。两者不要混淆——文档/笔记里若看到 `--do` 是记错了。

**解决**：

```bash
# 正确：用 --method / -m
gitlink-cli pr +merge -i <id> --method squash --format json
gitlink-cli pr +merge -i <id> -m rebase       --format json

# 合并方式取值：merge | rebase | squash（缺省 merge）
```

> 注意：gatekeeper 的 `behavior.merge_method` 直接对应 `--method`。合并仍受 #8 的三条件约束（PASS + `--apply` + `auto_merge: true`）。

---

## 还没解决？

1. 加 `--debug` 重跑，看原始请求/响应（`gitlink-cli ... --debug`）。
2. 用 `--format json` 拿结构化输出，核对 `error.code` / `error.suggestion` 与本文对照。
3. 回到 [`SKILL.md`](./SKILL.md) 重走工作流，确认每步命令与参数；字段/算法/阈值一律以 [`REFERENCE.md`](./REFERENCE.md)（SSOT）与 [`REFERENCE.md`](./REFERENCE.md) 为准。
4. 始终遵守安全默认：默认 dry-run、绝不自动合并、写操作前复述意图、不回显 Token。
