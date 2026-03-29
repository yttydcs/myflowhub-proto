# 2026-03-29 Proto：stream 主线合并

## 变更背景 / 目标

- 根 workspace 的 `go.work` 固定使用 `repo/MyFlowHub-Proto` 主线，而不是 `worktrees/proto-stream-subproto`。
- 只保留 feature worktree 会导致 `repo/MyFlowHub-Server` 主线路径在 workspace mode 下仍拿不到 `protocol/stream`。
- 本次目标是把 `feat/proto-stream-subproto` 已完成的 `stream` canonical wire 正式合回 `repo/MyFlowHub-Proto/main`。

## 具体变更内容

- 合并 `feat/proto-stream-subproto` 到当前 Proto 主线语境。
- 新增 [`protocol/stream/types.go`](D:/project/MyFlowHub3/worktrees/fix-proto-stream-mainline-merge/protocol/stream/types.go)
  - 提供 `SubProtoStream = 8`
  - 提供 `stream` 的公开 action / permission / descriptor / DATA/ACK 头定义
- 更新 [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/fix-proto-stream-mainline-merge/docs/protocol_map.md)
  - 生成结果纳入 `SubProto=8 / protocol/stream`
- 更新 [`docs/change/README.md`](D:/project/MyFlowHub3/worktrees/fix-proto-stream-mainline-merge/docs/change/README.md)
  - 补入 `stream` 相关归档入口
- 删除来源分支带入的根 [`plan.md`](D:/project/MyFlowHub3/worktrees/fix-proto-stream-mainline-merge/plan.md)
  - 改为保留本次 workflow 的 `todo.md` 作为 worktree 控制文档，并把可保留信息归档到 `docs/plan_archive/`

## Impact

- Requirements impact: `none`
- Specs impact: `updated`
- Lessons impact: `none`
- Related requirements:
  - none
- Related specs:
  - [`docs/protocol_map.md`](D:/project/MyFlowHub3/worktrees/fix-proto-stream-mainline-merge/docs/protocol_map.md)
- Related lessons:
  - none

## 对应 plan.md / todo.md 任务映射

- `PROTO-STRM-M1`
- `PROTO-STRM-M2`
- `PROTO-STRM-M3`
- `PROTO-STRM-M4`

## 经验 / 教训摘要

- 根 workspace 若已把上游协议仓库固定进 `go.work`，仅合业务仓主线并不能恢复真正的主线路径。
- `Proto` 这类 canonical 仓库在主线合并时，要优先确认 `docs/protocol_map.md` 这类受生成门禁保护的文件已经正确收敛。

## 可复用排查线索

- 症状
  - 主线 `Server` 已接入 `stream`，但 workspace mode 下仍提示缺少 `protocol/stream`
- 触发条件
  - 根 `go.work` 指向 `repo/MyFlowHub-Proto/main`
  - `stream` 只存在于 Proto feature worktree，尚未回到主线
- 关键词
  - `protocol/stream`
  - `go.work`
  - `SubProto=8`
  - `docs/protocol_map.md`
- 快速检查
  - 查看 [`protocol/stream/types.go`](D:/project/MyFlowHub3/worktrees/fix-proto-stream-mainline-merge/protocol/stream/types.go) 是否已在主线路径存在
  - 执行 `GOWORK=off go test ./... -count=1 -p 1`

## 关键设计决策与权衡

- 采用“先合回本地主线，再由下游 Server 主线路径消费”的收口方式，而不是继续让 root 脚本依赖 Proto worktree。
- 本轮不承担新的 semver tag 发布，只解决 workspace 主线路径缺口；发布链另行处理。

## 测试与验证方式 / 结果

- `Proto merge worktree`
  - 执行：`$env:GOWORK='off'; go test ./... -count=1 -p 1`
  - 结果：通过

## 潜在影响与回滚方案

- 潜在影响
  - 主线 `docs/protocol_map.md` 增加 `stream` 协议映射后，下游会开始默认看到 `SubProto=8`
- 回滚方案
  - 在主线对本次 merge commit 执行 `git revert`

## 子Agent执行轨迹

- 本轮未使用子Agent
