# Lessons

存放 `MyFlowHub-Proto` workflow 中具备复用价值的协议、生成链路与跨仓同步经验。

## How To Use

- 先查这里，再去翻 `docs/change/`。
- 只保留可重复命中的结构性问题或排查规则，不把它当成逐次变更日志。

## Current Status

- [flow-graph-contract-canonical-sync.md](flow-graph-contract-canonical-sync.md)
  - 症状：Flow graph/node contract 在 Proto、SubProto、Win、Server 之间出现多份 truth，或 Win/Server/SubProto 先于 Proto 各自扩展字段。
  - 关键词：`flow_contract.ts`、`Node.Kind`、`BindingSourceKind`、`BranchMatchOp`、`call.args`、`edge.case`。

## Rules

- 文件名使用稳定主题名，不用日期前缀。
- 每条 lesson 都要回链到对应 `change` / `spec`。
