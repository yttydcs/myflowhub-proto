# 变更说明：Management 新增 node_info（节点自描述信息）

## 背景 / 目标
Win 的 `Devices` 页面需要在点击节点时展示“设备基本信息”（平台、版本号等）。为满足“数据来自节点本身”，需要在 management 子协议新增一个动作，由目标节点本地采集并回包。

## 具体变更内容
- 新增：
  - `protocol/management/types.go`
    - `ActionNodeInfo / ActionNodeInfoResp`
    - `NodeInfoReq / NodeInfoResp`（`items: map[string]string`，用于承载可扩展 KV）
- 更新：
  - `docs/protocol_map.md`：协议映射文档生成结果同步更新

## 对应 plan.md 任务映射
- `worktrees/node-info/MyFlowHub-Win/plan.md`
  - T2. Proto：新增 management `node_info`（并发布 `v0.1.1`）

## 关键设计决策与权衡
- 选择新增独立 action（`node_info`），而不是复用 `config_get`：
  - 优点：语义清晰、可扩展、不会污染配置域；更符合“节点自描述信息”的用途。
  - 代价：wire 新增 action（但不影响既有 action，属于向前兼容增强）。
- `items` 使用 KV（`map[string]string`）：
  - 优点：前端展示不依赖固定字段，未来扩展无需破坏旧 UI。
  - 代价：缺少强类型约束；但当前需求以“通用展示”优先。

## 测试与验证
- 已在本地执行：`go test ./...`

## 潜在影响与回滚方案
- 影响：仅新增 action 与类型，不改变既有 wire；旧节点不受影响。
- 回滚：
  - 若未发布 tag：可直接 revert 本提交。
  - 若已发布 tag：不建议删除远端 tag；改用补丁版本（如 `v0.1.2`）修复并同步下游依赖。

