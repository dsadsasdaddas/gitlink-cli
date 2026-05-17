# 分支管理完整工作流

本示例演示一个完整的分支管理流程：查看分支 → 创建分支 → 开发推送 → 保护分支 → 清理旧分支。

## 场景描述

你正在维护一个 GitLink 仓库，需要开发一个新功能。开发完成后需要保护主分支，并清理已完成的功能分支。

## 步骤

### 1. 查看当前所有分支

```bash
gitlink-cli branch +list
```

输出示例：
```
Branch            Protected
master            Yes
develop           No
feature/old-login No
```

### 2. 创建新功能分支

```bash
gitlink-cli branch +create --name feature/new-auth --from develop
```

输出确认分支创建成功。

### 3. 本地开发并推送

```bash
git checkout -b feature/new-auth
# ... 编写代码 ...
git add .
git commit -m "feat: add new auth module"
git push gitlink feature/new-auth
```

### 4. 创建 PR 并合并

通过 `gitlink-cli pr +create` 创建 PR，审查后合并。合并后，功能分支 `feature/new-auth` 将自动删除（取决于仓库设置）。

### 5. 保护主分支

确保 `master` 分支有保护规则，防止误删或直接推送：

```bash
gitlink-cli branch +protect --name master
```

### 6. 清理旧分支（可选）

开发完成后，删除不再需要的旧分支：

```bash
# 确认旧分支已合并或无保留价值后
gitlink-cli branch +delete --name feature/old-login
```

## 注意事项

- `branch +delete` 是 **Destructive Operation**，操作前请确认分支内容已合并或无保留价值
- `branch +protect` 会对分支施加保护规则，影响后续的推送和合并流程
- `branch +unprotect` 仅支持简单分支名（如 `main`），含 `/` 的路径需通过 Web 页面操作
- 建议使用 `branch +list` 先查看分支状态，再执行写入/删除操作
- 建议使用 `branch +list` 先查看分支状态，再执行写入/删除操作
