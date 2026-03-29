# 2026-03-28 Proto：stream canonical wire

## 变更背景 / 目标

- 为新的 `stream` 子协议补齐 canonical wire，供 `SubProto` 与 `Server` 共享。
- 保持 `MyFlowHub-Proto` 继续只承载协议字典，不引入运行时协调逻辑。

## 具体变更内容

- 新增 [`protocol/stream/types.go`](D:/project/MyFlowHub3/worktrees/proto-stream-subproto/protocol/stream/types.go)
  - 新增 `SubProtoStream = 8`
  - 新增公开 CTRL action：
    - `announce / withdraw / list_sources / get_source`
    - `announce_consumer / withdraw_consumer / list_consumers / get_consumer`
    - `subscribe / unsubscribe / connect / disconnect / signal`
  - 新增权限：
    - `stream.publish`
    - `stream.consume`
    - `stream.subscribe`
    - `stream.connect`
  - 锁定内建 `kind`：
    - `music`
    - `video`
    - `text`
    - `custom`
  - 新增 `SourceDescriptor` / `ConsumerDescriptor`
  - 新增 `StreamDataHeaderV1` / `StreamAckHeaderV1`
  - DATA flag 使用显式常量值 `1/2/4/8`，避免 `protocolmapgen` 解析位移表达式失败
- 修改 [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/proto-stream-subproto/docs/protocol_map.md)
  - 重新生成协议映射，纳入 `SubProto=8 / protocol/stream`

## Impact

- Requirements impact: `none`
- Specs impact: `updated`
- Lessons impact: `none`
- Related requirements:
  - [`docs/requirements/stream.md`](D:/project/MyFlowHub3/worktrees/server-stream-subproto-design/docs/requirements/stream.md)
- Related specs:
  - [`docs/specs/stream.md`](D:/project/MyFlowHub3/worktrees/server-stream-subproto-design/docs/specs/stream.md)
  - [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/proto-stream-subproto/docs/protocol_map.md)
- Related lessons: `none`

## 对应 plan.md 任务映射

- `STRM-PROTO-1`

## 经验 / 教训摘要

- `stream` 的公开 wire 只描述稳定对外合同；`delivery_prepare/activate/abort/close` 这类实现内协调动作不进入 Proto。
- `protocolmapgen` 对常量表达式有约束，位标志常量用显式字面量最稳妥。

## 可复用排查线索

- 症状
  - `undefined: protocol.SubProtoStream`
  - `undefined: protocol.ActionAnnounceConsumer`
  - `docs/protocol_map.md is out of date`
- 触发条件
  - 下游开始引入 `stream` 子协议，但 Proto 尚未补齐 canonical wire
- 关键词
  - `protocol/stream`
  - `SubProtoStream`
  - `music video text custom`
  - `protocolmapgen`
- 快速检查
  - 运行 `GOWORK=off go test ./... -count=1 -p 1`
  - 检查 [`protocol/stream/types.go`](D:/project/MyFlowHub3/worktrees/proto-stream-subproto/protocol/stream/types.go) 是否已存在

## 关键设计决策与权衡

- 只把稳定公开动作写入 Proto，不把 handler 私有协调动作公开化。
- 把 `text` 作为一级内建 `kind`，而不是让文本场景全部退回 `custom`。
- DATA / ACK 采用固定二进制小头，避免热路径依赖 JSON。

## 测试与验证方式 / 结果

- 执行：`GOWORK=off go test ./... -count=1 -p 1`
- 结果：通过

## 潜在影响与回滚方案

- 潜在影响
  - 下游若仍锁定旧版 Proto，将无法编译 `protocol/stream`
- 回滚方案
  - 回退 [`protocol/stream/types.go`](D:/project/MyFlowHub3/worktrees/proto-stream-subproto/protocol/stream/types.go)
  - 回退 [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/proto-stream-subproto/docs/protocol_map.md)

## 子Agent执行轨迹

- 本轮未使用子Agent
