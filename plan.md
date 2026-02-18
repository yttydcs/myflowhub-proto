# Plan - Proto：发布 v0.1.0（semver 依赖基线）（PR19-PROTO-SemVer）

## Workflow 信息
- Repo：`MyFlowHub-Proto`
- 分支：`chore/proto-semver`
- Worktree：`d:\project\MyFlowHub3\worktrees\pr19-semver-deps\MyFlowHub-Proto`
- Base：`main`
- 参考：
  - `d:\project\MyFlowHub3\target.md`
  - `d:\project\MyFlowHub3\repos.md`
  - `d:\project\MyFlowHub3\guide.md`（commit 信息中文）
- 目标：发布并推送 tag `v0.1.0`

## 约束（边界）
- Proto 必须保持“纯协议字典”：
  - 只依赖标准库；
  - 不引入运行时逻辑/网络逻辑/handler 逻辑；
  - wire 不变（SubProto 值、Action 名称、JSON struct 字段不变）。
- 验收必须使用 `GOWORK=off`。

## 当前状态（事实，可审计）
- Proto 当前尚未发布任何 tag；上游（Server/SDK/Win）通过 `replace ../MyFlowHub-Proto` 联调。
- 本 workflow 要发布 `v0.1.0`，作为上游去 `replace` 的可拉取依赖版本。

---

## 1) 需求分析

### 目标
1) 发布 `github.com/yttydcs/myflowhub-proto@v0.1.0`，供上游仓库以 semver 方式引用。
2) 归档发布说明，包含验证方式与回滚策略。

### 范围（必须 / 不做）
#### 必须
- 新增归档文档：`docs/change/2026-02-18_proto-v0.1.0.md`
- 回归验证：`GOWORK=off go test ./... -count=1 -p 1`
- 结束 workflow 且你确认后：
  - 在 `repo/` 合并分支到 `main` 并 push
  - 创建 annotated tag `v0.1.0` 并 push tag

#### 不做
- 不做 wire 重命名、assist/up/notify 的 wire 收敛（策略 A：wire 不改）
- 不引入 Core/SDK 依赖

### 验收标准
- `GOWORK=off go test ./... -count=1 -p 1` 通过
- 远端存在并可拉取：`github.com/yttydcs/myflowhub-proto@v0.1.0`

### 风险
- tag 一旦 push 不应删除/改写；如需修复，追加发布 `v0.1.1`。

---

## 2) 架构设计（分析）

### 总体方案（采用）
- 使用 semver tag 发布 Proto：
  - 首次发布：`v0.1.0`
- 上游通过 `require github.com/yttydcs/myflowhub-proto@v0.1.0` 引用，移除 `replace`。

### 测试策略
- `go test ./...`（主要验证可编译）

---

## 3.1) 计划拆分（Checklist）

## 问题清单（阻塞：否）
- 已确认版本：Proto 发布 `v0.1.0`；验收使用 `GOWORK=off`。

### PROTOSEM1 - 归档发布文档
- 目标：写清楚 `v0.1.0` 发布范围、验证命令、风险与回滚方案。
- 涉及文件：
  - `docs/change/2026-02-18_proto-v0.1.0.md`
- 验收条件：
  - 文档包含 tag、发布范围、`GOWORK=off` 验收命令、回滚策略。
- 回滚点：
  - revert 文档提交。

### PROTOSEM2 - 回归测试（GOWORK=off）
- 目标：确保 Proto 发布点可通过回归（不依赖 go.work）。
- 命令：
  - `$env:GOTMPDIR='d:\\project\\MyFlowHub3\\.tmp\\gotmp'`
  - `New-Item -ItemType Directory -Force -Path $env:GOTMPDIR | Out-Null`
  - `$env:GOWORK='off'`
  - `go test ./... -count=1 -p 1`
- 验收条件：通过。
- 回滚点：revert 本分支改动（若有）。

### PROTOSEM3 - Code Review（阶段 3.3）+ 归档（阶段 4）
- 目标：Review 覆盖需求/风险/测试，并在 docs/change 归档。
- 验收条件：Review 结论为“通过”。

### PROTOSEM4 - 合并与打 tag（你确认结束 workflow 后执行）
- 目标：将本分支合并到 `main`，并发布 `v0.1.0`。
- 步骤（在 `repo/MyFlowHub-Proto` 执行）：
  1) `git merge --ff-only origin/chore/proto-semver`（或等价方式）
  2) `git push origin main`
  3) `git tag -a v0.1.0 -m \"chore: 发布 v0.1.0\"`
  4) `git push origin v0.1.0`
- 验收条件：
  - tag 可在远端查看并被 `go` 拉取。
- 回滚方案：
  - 不删除 tag；如需修复，追加发布 `v0.1.1`。
