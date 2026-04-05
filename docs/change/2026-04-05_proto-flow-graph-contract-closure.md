# 2026-04-05 Proto Flow Graph Contract Closure

## 变更背景 / 目标

- 上一轮完整性审查后，剩余唯一结构性缺口是 `flow graph/node contract` 仍在 Proto、SubProto、Win、Server 四仓各自维护局部真相。
- 本轮目标是把 Flow graph/node wire contract 收束为 Proto canonical truth，并让 SubProto、Win、Server 变成消费方或镜像方，同时明确不保留旧数据兼容。

## 具体变更内容

### Proto canonical contract

- `protocol/flow/types.go`
  - 新增 `NodeKind`、`BindingSourceKind`、`BranchMatchOp`
  - 为 `call / compose / transform / set_var / branch / foreach / subflow` 建立 typed spec struct
  - `Trigger` 增加 `cron`
  - `Edge` 增加 `case`
  - `Node.Kind` 改为 typed enum
- `cmd/flowcontractgen/**`
  - 新增 canonical contract generator
- `internal/flowcontract/**`
  - 新增生成逻辑
- `docs/flow_contract.md`
  - 新增 human-readable canonical contract 文档
- `generated/flow_contract.ts`
  - 新增供下游同步的 TS canonical artifact
- `README.md`
  - 补 canonical contract generator 和用途说明
- `docs/protocol_map.md`
  - 同步 canonical contract 相关 flow protocol 映射

### SubProto runtime deduplication

- `flow/runtime_bindings.go`
  - 删除私有重复 wire contract struct，改为直接消费 Proto canonical types
- `flow/handler.go`
  - decode / validate / materialize 路径切到 typed normalize
  - legacy `call.args` 直接报错，要求改用 `args_template`
- `flow/types.go`
  - 暴露 Proto 类型别名和常量供本仓使用
- `flow/go.mod`
  - 增加本地 Proto `replace`，冻结到本 workflow 的 GC1 接口
- `flow/*_test.go`
  - 覆盖 strict decode 和 legacy payload rejection

### Win contract consumption

- `frontend/src/generated/flow_contract.ts`
  - 同步 Proto canonical TS artifact
- `frontend/src/stores/flow.ts`
  - 删除本地完整 `FlowNodeKind / FlowBindingSourceKind / FlowBranchMatchOp` truth
  - 改为从 generated artifact 派生 canonical option lists
  - 保留 Win 本地 draft 投影、root/body source kind 过滤、strict parse/export 和标签映射
- `frontend/src/components/flow/editor/FlowNodeInspector.vue`
- `frontend/src/components/flow/editor/FlowBodyNodeInspector.vue`
- `frontend/src/components/flow/editor/FlowFieldBindingDialog.vue`
  - inspector / dialog option lists 改为消费 canonical option exports
- `frontend/src/stores/flow.test.ts`
  - 增加 canonical option import regression test
- `docs/specs/flow-editor-visual-form.md`
  - 明确 Win 消费 canonical contract artifact 的边界

### Server spec and behavior closure

- `docs/specs/flow.md`
  - 增加 canonical contract source 指向 Proto
  - 去除 legacy `local` / `exec` 和 `call.args` 兼容措辞
- `tests/integration_flow_round_trip_test.go`
  - 扩展为 advanced graph/node round-trip 集成测试
  - 覆盖 `set_var / transform / branch / foreach / subflow / edge.case / detail`

## Requirements impact

- `none`

## Specs impact

- `updated`

## Lessons impact

- `updated`

## Related requirements

- `D:/project/MyFlowHub3/repo/MyFlowHub-Server/docs/requirements/flow_data_dag.md`

## Related specs

- `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/docs/flow_contract.md`
- `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/docs/protocol_map.md`
- `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure/docs/specs/flow.md`
- `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure/docs/specs/flow-editor-visual-form.md`

## Related lessons

- `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/docs/lessons/flow-graph-contract-canonical-sync.md`
- `D:/project/MyFlowHub3/docs/lessons/wails-binding-proto-drift.md`
- `D:/project/MyFlowHub3/docs/lessons/frontend-worktree-wailsjs-missing.md`
- `D:/project/MyFlowHub3/docs/lessons/frontend-build-empty-node-modules.md`

## 对应 plan.md 任务映射

- `GC1`
  - `protocol/flow/types.go`
  - `cmd/flowcontractgen/**`
  - `internal/flowcontract/**`
  - `docs/flow_contract.md`
  - `generated/flow_contract.ts`
  - `README.md`
  - `docs/protocol_map.md`
- `GC2`
  - `flow/runtime_bindings.go`
  - `flow/handler.go`
  - `flow/types.go`
  - `flow/go.mod`
  - `flow/graph_test.go`
  - `flow/local_capability_test.go`
  - `flow/capability_provider_test.go`
  - `flow/trigger_test.go`
- `GC3`
  - `frontend/src/generated/flow_contract.ts`
  - `frontend/src/stores/flow.ts`
  - `frontend/src/components/flow/editor/FlowNodeInspector.vue`
  - `frontend/src/components/flow/editor/FlowBodyNodeInspector.vue`
  - `frontend/src/components/flow/editor/FlowFieldBindingDialog.vue`
  - `frontend/src/stores/flow.test.ts`
  - `docs/specs/flow-editor-visual-form.md`
- `GC4`
  - `docs/specs/flow.md`
  - `tests/integration_flow_round_trip_test.go`
- `GC5`
  - `plan.md`
  - `docs/change/2026-04-05_proto-flow-graph-contract-closure.md`
  - `docs/lessons/flow-graph-contract-canonical-sync.md`
  - `docs/lessons/README.md`

## 经验 / 教训摘要

- Flow graph/node wire contract 必须先在 Proto 收束，再让 SubProto、Win、Server 去消费或镜像，不能继续让四仓各自维护局部 truth。
- Win 侧 canonical artifact 只应提供协议真相，不应吞掉 draft model、ordinary-mode 边界和本地 UX 适配职责。
- 在当前 Windows 环境里，Vitest 默认 fork worker 可能误报 OOM；Flow editor 相关验证需要显式单 worker 运行。

## 可复用排查线索

- 症状：
  - Win 侧 `FlowNodeKind` / `FlowBindingSourceKind` / `FlowBranchMatchOp` 与 Proto 不一致
  - SubProto 仍接受 `call.args` 或 `local` / `exec`
  - Server spec 仍写 legacy 兼容措辞
  - `flow.get` round-trip 对 advanced graph/node payload 回不来
- 触发条件：
  - 新增 node/source/op 时直接先改 Win 或 SubProto
  - 修改 Proto contract 后未同步 TS artifact 或稳定 spec
  - 继续保留历史兼容字段，导致 runtime 和 spec 双重分叉
- 关键词 / 错误文本：
  - `flow_contract.ts`
  - `Node.Kind string`
  - `BranchMatchOp`
  - `BindingSourceKind`
  - `args_template`
  - `call.args`
  - `local`
  - `exec`
  - `edge.case`
- 快速检查：
  1. 先看 `MyFlowHub-Proto/protocol/flow/types.go` 是否是新增或变更的唯一 canonical source
  2. 运行 Proto generator，确认 `docs/flow_contract.md` 与 `generated/flow_contract.ts` 已更新
  3. 对比 Win `frontend/src/generated/flow_contract.ts` 与 Proto 生成物是否一致
  4. 检查 SubProto 是否仍存在私有重复 wire contract struct
  5. 跑 Server `TestIntegrationFlowRoundTrip` 看 advanced graph/node payload 是否完整 round-trip

## 关键设计决策与权衡

- 决策：不把 runtime validation 移进 Proto
  - 好处：不违反 `repos.md` 和 Proto 仓边界
  - 代价：SubProto 仍需保留 runtime 语义层，但必须消费 Proto wire contract 类型
- 决策：完全移除 legacy `local` / `exec` 与 `call.args` 兼容
  - 好处：趁项目尚未上线，直接消除双轨语义和未来迁移成本
  - 代价：任何旧 payload 现在都会显式失败
- 决策：Win 用 generated artifact 消费 canonical enums，但保留本地 UI draft 适配
  - 好处：协议真相与 UI 适配边界清晰
  - 代价：未来 Proto 新增节点后，Win 仍需显式补 ordinary-mode UI 支持

## 测试与验证方式 / 结果

- `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure`
  - `GOWORK=off go test ./...`
  - 结果：通过
- `D:/project/MyFlowHub3/worktrees/subproto-graph-contract-closure/flow`
  - `GOWORK=off go test ./... -count=1 -p 1`
  - 结果：通过
- `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure`
  - `GOWORK=off go test ./tests -run TestIntegrationFlowRoundTrip -count=1 -p 1`
  - 结果：通过
- `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure`
  - `GOWORK=off wails generate module`
  - 结果：通过
  - `GOWORK=off go test ./internal/services/flow ./internal/mcp ./internal/mcpapp -count=1 -p 1`
  - 结果：通过
- `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure/frontend`
  - `npx vitest run src/stores/flow.test.ts --maxWorkers=1 --no-file-parallelism`
  - 结果：通过
  - `npx vitest run src/components/flow/editor/FlowNodeInspector.test.ts --maxWorkers=1 --no-file-parallelism`
  - 结果：通过
  - `npx vitest run src/components/flow/editor/FlowBodyNodeInspector.test.ts --maxWorkers=1 --no-file-parallelism`
  - 结果：通过
  - `npx vitest run src/components/flow/editor/FlowFieldBindingDialog.test.ts --pool threads --maxWorkers=1 --no-file-parallelism`
  - 结果：通过
  - `npx vitest run src/windows/FlowEditorWindow.test.ts --pool threads --maxWorkers=1 --no-file-parallelism`
  - 结果：通过
- Contract sync
  - 比较 Proto `generated/flow_contract.ts` 与 Win `frontend/src/generated/flow_contract.ts`
  - 结果：一致

## 潜在影响

- 后续新增 Flow node/source/op 时，Proto 会成为唯一 first-change 入口，下游需要通过 generator/sync 和显式适配完成闭环。
- 旧 payload 中若仍出现 `call.args`、`local`、`exec`，现在会明确失败。
- Win worktree 的前端测试在当前机器上需要避免默认 fork worker 策略，否则可能出现环境性 OOM 假失败。

## 回滚方案

1. Proto 回退 `protocol/flow/types.go`、`cmd/flowcontractgen/**`、`internal/flowcontract/**`、`docs/flow_contract.md`、`generated/flow_contract.ts`、`README.md`、`docs/protocol_map.md`
2. SubProto 回退 `flow/runtime_bindings.go`、`flow/handler.go`、`flow/types.go`、`flow/go.mod` 以及相关测试
3. Win 回退 `frontend/src/generated/flow_contract.ts` 以及 6 个相关 UI/store/spec 文件
4. Server 回退 `docs/specs/flow.md` 与 `tests/integration_flow_round_trip_test.go`

## 子Agent执行轨迹

- `GC2`
  - 子Agent：`019d5bb2-9629-7c00-9c5d-9aa780875417`
  - 负责仓库：`MyFlowHub-SubProto`
  - 结果：完成；主代理复核后再次运行 `GOWORK=off go test ./... -count=1 -p 1` 通过
- `GC3`
  - 子Agent：`019d5bb2-ace7-74f0-a435-245c868c2b31`
  - 负责仓库：`MyFlowHub-Win`
  - 结果：完成；主代理补齐 `npm install`、`GOWORK=off wails generate module` 和 Win targeted tests 后复核通过
- `GC4`
  - 子Agent：`019d5bb2-c3c3-7e40-abd2-cee93c010d34`
  - 负责仓库：`MyFlowHub-Server`
  - 结果：完成；主代理复核后再次运行 `TestIntegrationFlowRoundTrip` 通过
- `GC1` / `GC5`
  - 由主代理完成
