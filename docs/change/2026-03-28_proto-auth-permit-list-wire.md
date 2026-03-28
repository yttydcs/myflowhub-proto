# 2026-03-28 Proto：auth register permit list wire contract

## 变更背景 / 目标

- permit issuance UI 需要从 auth runtime 读取当前活动 permit 列表，因此 `MyFlowHub-Proto` 需要先补齐 canonical wire 定义和 protocol map。
- 目标是保持 Proto 仍然只承担“协议字典与映射”的职责，不引入运行时语义。

## 具体变更内容

- 修改 `protocol/auth/types.go`
  - 新增 `ActionListRegisterPermits`
  - 新增 `ActionListRegisterPermitsResp`
  - 新增 `RegisterPermitInfo`
  - 新增 `ListRegisterPermitsReq`
  - 新增 `ListRegisterPermitsResp`
- 修改 `docs/protocol_map.md`
  - 使用 `protocolmapgen` 重新生成 auth action / payload 映射
  - 纳入 permit list 的 request / response / item 结构

## Impact

- Requirements impact: `none`
- Specs impact: `updated`
- Lessons impact: `none`
- Related requirements: `none`
- Related specs:
  - `docs/protocol_map.md`
- Related lessons: `none`

## 对应 plan.md 任务映射

- `PROTO-PERMIT-1`
- `PROTO-PERMIT-2`
- `REVIEW-PROTO-PERMIT-1`
- `ARCHIVE-PROTO-PERMIT-1`

## 经验 / 教训摘要

- 只要新增 action 或 payload，Proto 仓就必须立刻重跑 `protocolmapgen`；否则 canonical wire 和映射文档会立刻失真。
- permit list 这类读接口的稳定性同样依赖 canonical struct 命名，不应让下游各自定义临时 JSON 结构。

## 可复用排查线索

- 症状
  - 下游报 `undefined: protocol.ActionListRegisterPermits`
  - `docs/protocol_map.md` 校验失败
- 触发条件
  - 新增 permit list 契约后没有同步 Proto
- 关键词
  - `list_register_permits`
  - `RegisterPermitInfo`
  - `protocolmapgen`
- 快速检查
  - 查看 `protocol/auth/types.go` 是否包含 permit list action 和结构
  - 运行 `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md`

## 关键设计决策与权衡

- 仅新增 wire action 与 payload，不在 Proto 中表达“活动 permit”“权限判断”等运行时语义
  - 优点：职责清晰，继续由 SubProto/Server 解释业务行为
  - 代价：阅读者需要到下游 spec/runtime 看更完整语义

## 测试与验证方式 / 结果

- `GOWORK=off go test ./... -count=1 -p 1`
  - 结果：通过
- `GOWORK=off go run ./cmd/protocolmapgen -check -out docs/protocol_map.md`
  - 结果：通过

## 潜在影响与回滚方案

- 潜在影响
  - 下游若仍停留在旧版 Proto，将无法编译或识别 permit list action
- 回滚方案
  - 回退 `protocol/auth/types.go`
  - 回退 `docs/protocol_map.md`

## 子Agent执行轨迹

- 本轮未使用子Agent
