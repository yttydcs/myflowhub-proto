# flow-graph-contract-canonical-sync

## Summary

- 当 Flow graph/node contract 同时在 Proto、SubProto、Win、Server 维护多份 truth 时，advanced node、binding source、branch op 和 legacy compatibility 很容易各自漂移。
- 正确的收束方式是：Proto 先拥有 canonical typed contract 和 generated artifact，SubProto 消费 Proto wire type，Win 消费 generated artifact，Server 只镜像 spec 和行为级验证。

## Lookup Hints

- `flow_contract.ts`
- `Node.Kind`
- `BindingSourceKind`
- `BranchMatchOp`
- `call.args`
- `args_template`
- `local`
- `exec`
- `edge.case`
- `flow graph/node contract`

## Symptoms

- Win 侧手写的 `FlowNodeKind` / `FlowBindingSourceKind` / `FlowBranchMatchOp` 与 Proto 不一致。
- SubProto 仍保留私有重复 wire contract struct，或者仍兜底接受 `call.args` / `local` / `exec`。
- Server spec 仍描述 legacy compatibility，但 runtime 或 canonical doc 已经不是这样。
- `flow.get` 无法稳定 round-trip advanced graph/node payload，特别是 `branch / foreach / subflow / edge.case`。

## Impact

- 同一份 Flow payload 在 authoring、runtime、spec、integration test 之间会出现分叉。
- 新节点或新字段上线前需要反复人工对齐，多仓变更风险高。
- legacy compatibility 会持续拖住 runtime 和 spec 的收束。

## Trigger Conditions

- 新增 node/source/op 时，先改 Win 或 SubProto，再回头补 Proto。
- 修改 Proto contract 后，没有同步 generated TS artifact 或下游 spec。
- 在项目尚未上线阶段仍保留历史兼容字段。

## Root Cause

- 没有 single source of truth 来定义 Flow graph/node wire contract。
- Proto、SubProto、Win、Server 之间混淆了“协议真相”“runtime 语义”“UI draft 投影”和“稳定文档镜像”的边界。

## Investigation Trail

- 先检查 `MyFlowHub-Proto/protocol/flow/types.go` 是否已经覆盖当前支持的 node/source/op contract。
- 对比 SubProto 是否仍有 repo-local duplicate wire structs。
- 对比 Win 是否仍在 store / inspector / dialog 中手写完整 literal unions 或 option arrays。
- 对比 Server stable spec 是否仍写 legacy compatibility。
- 用 behavior-level round-trip test 验证 `flow.set/get` 是否真的保留 advanced graph/node payload。

## Resolution

- 在 Proto 增加 typed Flow graph/node contract、canonical doc、TS generated artifact。
- 让 SubProto 直接消费 Proto types，并显式拒绝 legacy `call.args` / `local` / `exec`。
- 让 Win 改为消费 generated artifact，只保留本地 draft / filter / UX 适配。
- 让 Server spec 指回 Proto canonical source，并用 integration test 验证 advanced graph/node round-trip。

## Prevention / Guardrails

- 以后新增 Flow node/source/op，第一步必须是 Proto。
- 修改 Proto contract 后，必须同步 canonical doc 和 TS artifact，再做下游仓适配。
- Win 不得重新手写完整协议 truth，只能在 canonical truth 之上做 UI draft 映射。
- Server spec 不是 source of truth，只能镜像 Proto contract 和已验证行为。
- 若项目尚未上线，不要为了“也许将来用到”的旧 payload 继续保留 legacy compatibility。

## Related Requirements / Specs / Changes

- Requirement: `D:/project/MyFlowHub3/repo/MyFlowHub-Server/docs/requirements/flow_data_dag.md`
- Specs:
  - `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/docs/flow_contract.md`
  - `D:/project/MyFlowHub3/worktrees/server-graph-contract-closure/docs/specs/flow.md`
  - `D:/project/MyFlowHub3/worktrees/win-graph-contract-closure/docs/specs/flow-editor-visual-form.md`
- Change:
  - `D:/project/MyFlowHub3/worktrees/proto-graph-contract-closure/docs/change/2026-04-05_proto-flow-graph-contract-closure.md`
