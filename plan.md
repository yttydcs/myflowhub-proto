# Plan - Proto Stream Release Chain

## Workflow 信息
- Repo：`MyFlowHub-Proto`
- Branch：`feat/proto-stream-subproto`
- Base：`main`
- Worktree：`D:\project\MyFlowHub3\worktrees\proto-stream-subproto`
- 当前 Stage：`3.1`

## 当前状态
- `stream` canonical wire 已在本 worktree 落地，并已形成结果归档：
  - `D:\project\MyFlowHub3\worktrees\proto-stream-subproto\docs\change\2026-03-28_stream-wire.md`
- 已确认远端最新 Proto tag 为 `v0.1.3`；下游 `SubProto/stream` 与 `Server` 若切到 `GOWORK=off`，仍拿不到 `protocol/stream`
- 本轮不再扩展 wire，目标是把现有 `stream` 协议定义发布为可被下游真实消费的 patch 版本

## Stage 1 - 需求分析

### 目标
- 发布 `github.com/yttydcs/myflowhub-proto@v0.1.4`，为 `stream` 下游消费提供可拉取版本锚点。
- 保证发布基线仍是纯协议字典，不夹带运行时语义或计划外 wire 调整。

### 范围
- 必须：
  - 以当前 `protocol/stream` 和 `docs/protocol_map.md` 为发布基线
  - 在 `GOWORK=off` 下完成 Proto 自身验证
  - 对远端创建并校验 `v0.1.4` tag
  - 为本轮发布补充可审计的 change 归档
- 可选：
  - 无
- 不做：
  - 不新增或重命名任何 `stream` action / payload / 常量
  - 不引入任何运行时逻辑
  - 不改写已存在 tag

### 使用场景
- `MyFlowHub-SubProto/stream` 在 `GOWORK=off` 下依赖 `github.com/yttydcs/myflowhub-proto v0.1.4`
- `MyFlowHub-Server` 在 `GOWORK=off` 下通过 compat wrapper 引入 `protocol/stream`

### 功能需求
- `protocol/stream/types.go` 必须保持当前公开合同不变
- `docs/protocol_map.md` 必须继续和 `protocol/*/types.go` 保持一致
- 发布后远端必须存在 `v0.1.4`，供下游 `go list -m` / `go test` 拉取

### 非功能需求
- Proto 仍保持“纯协议字典”边界，只依赖标准库
- 发布验证必须使用 `GOWORK=off`
- 已 push tag 不允许重写；如发现问题，只能追加更高 patch

### 输入输出
- 输入：
  - `D:\project\MyFlowHub3\worktrees\proto-stream-subproto\protocol\stream\types.go`
  - `D:\project\MyFlowHub3\worktrees\proto-stream-subproto\docs\protocol_map.md`
  - `D:\project\MyFlowHub3\worktrees\server-stream-subproto-design\docs\requirements\stream.md`
  - `D:\project\MyFlowHub3\worktrees\server-stream-subproto-design\docs\specs\stream.md`
- 输出：
  - 远端 tag：`v0.1.4`
  - 本轮发布归档文档

### 边界异常
- 若远端已存在 `v0.1.4`，不得覆盖或重打
- 若 `docs/protocol_map.md` 与真源不一致，必须先修正后再发布
- 若下游在远端 tag 可见后仍无法拉取，需要继续核查 module cache / version chain，而不是回退为 `go.work`

### 验收标准
- `GOWORK=off go test ./... -count=1 -p 1` 通过
- `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md` 通过
- `git ls-remote --tags origin refs/tags/v0.1.4` 能看到远端 tag

### 风险
- tag 一旦 push 不能回收，只能发更高 patch 修复
- 若下游过早切换版本而远端 tag 尚未传播，`GOWORK=off` 可能短暂失败

## Stage 2 - 架构设计

### 总体方案
- 方案 A：继续只用本地 `go.work` 联调，不发布 Proto tag
  - 不选：无法满足下游 `GOWORK=off` 真实消费
- 方案 B：以当前已落地的 `stream` wire 为基线，发布 `v0.1.4`
  - 采用：改动面最小，且能直接解除下游 `unknown package protocol/stream` 问题
- 方案 C：在 Proto 再叠加一轮 wire 调整后统一发布
  - 不选：会扩大范围，破坏“先固定版本锚点，再让下游升级”的节奏

### 模块职责
- `MyFlowHub-Proto`
  - 提供唯一公开协议真源
  - 对外发布 `v0.1.4`
- `MyFlowHub-SubProto`
  - 依赖 `v0.1.4` 并发布 `stream/v0.1.0`
- `MyFlowHub-Server`
  - 依赖 `v0.1.4 + stream/v0.1.0`，完成 `GOWORK=off` 验证与集成测试

### 数据 / 调用流
1. Proto 在本地执行 `GOWORK=off` 验证
2. Proto 创建并推送 `v0.1.4`
3. SubProto `stream` module 升级依赖到 `v0.1.4`
4. Server 升级依赖到 `v0.1.4`

### 接口草案
- 发布版本：`v0.1.4`
- 验证命令：
  - `GOWORK=off go test ./... -count=1 -p 1`
  - `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md`

### 错误与安全
- 若远端 tag 已存在，则停止发布动作并改发更高 patch
- 不允许把 `docs/change` 当长期协议真相替代 `protocol/*/types.go`

### 性能与测试策略
- 以仓库级回归为主，不新增运行时测试
- 发布后由下游 `SubProto/Server` 再执行真实消费验证

### 可扩展性设计点
- 本轮只固定 `stream` v1 的版本锚点，不在 Proto 侧新增第二轮 `stream` 合同

## Stage 3.1 - 计划
- Requirements impact：`none`
- Specs impact：`none`
- Related requirements：
  - `D:\project\MyFlowHub3\worktrees\server-stream-subproto-design\docs\requirements\stream.md`
- Related specs：
  - `D:\project\MyFlowHub3\worktrees\server-stream-subproto-design\docs\specs\stream.md`
  - `D:\project\MyFlowHub3\worktrees\proto-stream-subproto\docs\protocol_map.md`
- Related lessons：
  - `D:\project\MyFlowHub3\docs\lessons\cross-repo-semver-release.md`

### 执行清单
- [ ] `PROTO-REL-1` 审计 Proto 发布基线并完成本仓 `GOWORK=off` 验证
- [ ] `PROTO-REL-2` 创建并推送 `v0.1.4`
- [ ] `PROTO-DOC-1` 归档本轮 Proto 发布结果

### 任务明细

#### PROTO-REL-1
- Owner：主 Agent
- Worktree：`D:\project\MyFlowHub3\worktrees\proto-stream-subproto`
- Files：
  - `protocol/stream/types.go`
  - `docs/protocol_map.md`
- Goal：
  - 确认当前 `stream` wire 和 canonical 文档就是 `v0.1.4` 的发布基线
- Acceptance：
  - `GOWORK=off go test ./... -count=1 -p 1` 通过
  - `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md` 通过
- Tests：
  - `GOWORK=off go test ./... -count=1 -p 1`
  - `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md`
- Rollback：
  - 若发现基线不一致，先修正文档或代码，再继续发布

#### PROTO-REL-2
- Owner：主 Agent
- Worktree：`D:\project\MyFlowHub3\worktrees\proto-stream-subproto`
- Files：
  - `none (git commit / tag / push)`
- Goal：
  - 创建并推送 `v0.1.4`
- Acceptance：
  - `git ls-remote --tags origin refs/tags/v0.1.4` 有结果
- Tests：
  - `git ls-remote --tags origin refs/tags/v0.1.4`
- Rollback：
  - 不删除 tag；若发布点有问题，追加更高 patch

#### PROTO-DOC-1
- Owner：主 Agent
- Worktree：`D:\project\MyFlowHub3\worktrees\proto-stream-subproto`
- Files：
  - `docs/change/2026-03-28_stream-wire-release.md`
  - `docs/change/README.md`
- Goal：
  - 归档本轮 `v0.1.4` 发布与验证证据
- Acceptance：
  - change 文档记录版本号、验证命令、风险和回滚策略
- Tests：
  - 人工核对文档与发布结果一致
- Rollback：
  - 回退本轮文档改动

### 依赖 / 风险 / 备注
- 依赖顺序：`Proto v0.1.4 -> SubProto stream/v0.1.0 -> Server v0.0.12`
- 远端 tag 状态已核对；当前 `v0.1.4` 尚不存在
- Proto 发布完成前，禁止在下游提交 `GOWORK=off` 依赖对齐结论

阻塞：否
进入 3.2
