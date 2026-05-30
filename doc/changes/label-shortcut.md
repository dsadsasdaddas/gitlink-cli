# Label shortcut

新增 `label` Shortcut 组，补齐 GitLink Issue 标签（项目标记 / `issue_tags`）OpenAPI 的常用操作封装：

- `label +list`
- `label +create`
- `label +update`
- `label +delete`

实现要点：

- 列表支持 `--keyword` 关键词过滤、`--only-name` 精简返回、`--sort-by` / `--sort-direction` 排序，映射到 API 的 `order_by` / `order_direction`。
- `+create` 的 `--color` 缺省为 `#1E90FF`；颜色统一做十六进制（`#RGB` / `#RRGGBB`）客户端校验，非法颜色在调用 API 前即报错。
- `+update` 先从列表接口取标签当前值并与传入字段合并，避免漏传字段被清空（更新接口要求 `name`/`description`/`color` 同时提交）；无任何变更字段时直接报错。
- 路径使用 `/api/v1/{owner}/{repo}/issue_tags`，与 webhook/milestone 等组保持一致的 `/v1/` 前缀约定。
- 补充单元测试覆盖各命令的 HTTP 方法、路径、查询参数、payload，以及颜色校验和 id 归一化逻辑。

背景：在此之前，Issue 标签只能通过 Raw API（`issue_tags`）手工管理；`gitlink-code-review`、`gitlink-insight` 等 Skill 在做 Issue 分拣 / 打标签时都需要拼接原始请求。`label` 组将其提升为一等命令，并配套 `skills/gitlink-label/` Skill 文档，方便人类与 AI Agent 直接复用。
