# 2026-03-26 Proto：auth semi-central authority wire contract

## 变更背景 / 目标
- 为 auth 半中心 authority 方案补齐 canonical wire 定义，供 `SubProto` 和 `Server docs` 复用。
- 保持 `MyFlowHub-Proto` 仍然只承担“协议字典”职责，不引入运行时逻辑。

## 具体变更内容
- 修改 [`protocol/auth/types.go`](D:/project/MyFlowHub3/worktrees/proto-auth-semi-central-authority/protocol/auth/types.go)
  - 新增 `ActionAuthorityPolicySync = "authority_policy_sync"`
  - 新增 `AuthorityPolicySyncData`
    - `mode`
    - `effective_authority_id`
    - `epoch`
    - `ttl_sec`
- 修改 [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/proto-auth-semi-central-authority/docs/protocol_map.md)
  - 重新生成 Auth action / payload 映射，纳入 `authority_policy_sync`

## Impact
- Requirements impact: `none`
- Specs impact: `updated`
- Lessons impact: `none`
- Related requirements: `none`
- Related specs:
  - [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/proto-auth-semi-central-authority/docs/protocol_map.md)
  - [`auth.md`](D:/project/MyFlowHub3/worktrees/server-auth-semi-central-authority/docs/specs/auth.md)
- Related lessons: `none`

## 对应 plan.md 任务映射
- `AUTHPOL-PROTO-1`

## 经验 / 教训摘要
- `protocol/*/types.go` 仍是 single source-of-truth；协议扩展后必须立刻重跑 `protocolmapgen`，否则 Proto 仓的门禁测试会失败。

## 可复用排查线索
- 症状
  - `undefined: protocol.ActionAuthorityPolicySync`
  - `undefined: protocol.AuthorityPolicySyncData`
  - `docs/protocol_map.md is out of date`
- 触发条件
  - 下游开始使用半中心 authority 同步动作，但 Proto 未同步新 wire
- 关键词
  - `authority_policy_sync`
  - `effective_authority_id`
  - `protocolmapgen`
- 快速检查
  - 运行 `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md`
  - 检查 [`protocol/auth/types.go`](D:/project/MyFlowHub3/worktrees/proto-auth-semi-central-authority/protocol/auth/types.go) 是否包含新 action / struct

## 关键设计决策与权衡
- 仅增加 action 常量与 payload struct，不把“租约过期 / 转发 / 降级”之类运行时语义放进 Proto。
- `effective_authority_id` 作为运行时下发字段存在于 wire 中，但不意味着要持久化到各节点配置。

## 测试与验证方式 / 结果
- 执行：`GOWORK=off go test ./... -count=1 -p 1`
- 结果：通过

## 潜在影响与回滚方案
- 潜在影响
  - 下游若仍引用旧版 Proto，将无法编译或识别 `authority_policy_sync`
- 回滚方案
  - 回退 [`protocol/auth/types.go`](D:/project/MyFlowHub3/worktrees/proto-auth-semi-central-authority/protocol/auth/types.go)
  - 回退 [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/proto-auth-semi-central-authority/docs/protocol_map.md)

## 子Agent执行轨迹
- 本轮未使用子Agent
