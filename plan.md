# Plan - Proto：新增 Flow 删除部署协议

## Workflow 信息
- Repo：`MyFlowHub-Proto`
- Branch：`feat/proto-deploy-delete`
- Worktree：`D:/project/MyFlowHub3/repo/MyFlowHub-Proto/repo/MyFlowHub-Proto/worktrees/feat-proto-deploy-delete`
- Base：`main`

## 项目目标与当前状态
- 目标：为 Flow 子协议新增“删除部署”能力，作为 Win“删除部署”与 SubProto 执行链路的协议基础。
- 当前状态：`protocol/flow/types.go` 仅有 `set/run/status/list/get`，无 delete 动作与请求/响应结构。

## 可执行任务清单（Checklist）
- [x] PROTO-DEL-1 新增 Flow delete action 与权限常量
- [x] PROTO-DEL-2 新增 DeleteReq / DeleteResp 协议结构
- [x] PROTO-DEL-3 同步协议映射文档（如适用）
- [x] PROTO-DEL-4 回归测试与自检

## 任务明细

### PROTO-DEL-1 新增 action/权限常量
- 目标：在 Flow 协议常量中新增 `ActionDelete`、`ActionDeleteResp`、`PermFlowDelete`。
- 涉及模块/文件：
  - `protocol/flow/types.go`
- 验收条件：
  - 常量名称与字符串值稳定可复用。
- 测试点：
  - 编译通过；下游可引用常量。
- 回滚点：
  - 回退 `types.go` 新增常量。

### PROTO-DEL-2 新增请求/响应结构
- 目标：新增 delete 请求/响应结构，字段风格与 set/get/run/status 一致。
- 涉及模块/文件：
  - `protocol/flow/types.go`
- 验收条件：
  - `DeleteReq` 包含 `req_id/origin_node/executor_node/flow_id`。
  - `DeleteResp` 包含 `req_id/code/msg/flow_id`。
- 测试点：
  - `go test ./... -count=1` 通过。
- 回滚点：
  - 回退新增结构体定义。

### PROTO-DEL-3 同步协议映射文档（如适用）
- 目标：确保协议映射文档可体现新 action/type。
- 涉及模块/文件：
  - `docs/protocol_map.md`（若由生成工具维护则按工具更新）
- 验收条件：
  - Flow actions/types 文档包含 delete/delete_resp 与新增 payload type。
- 测试点：
  - 文档生成命令可执行或文档静态检查无格式问题。
- 回滚点：
  - 回退文档变更。

### PROTO-DEL-4 回归测试与自检
- 目标：保证 Proto 仓无计划外变更。
- 涉及模块/文件：
  - 全仓
- 验收条件：
  - `go test ./... -count=1` 通过。
- 测试点：
  - 编译与单测。
- 回滚点：
  - 回退本 workflow 所有提交。

## 依赖关系
- 本 workflow 先于 SubProto/Win 实现。

## 风险与注意事项
- 风险：常量命名与下游实现不一致会导致 action 无法匹配。
- 注意：仅做协议字典变更，不引入运行时逻辑。

