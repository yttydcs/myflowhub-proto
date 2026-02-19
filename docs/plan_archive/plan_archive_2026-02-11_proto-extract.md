# Plan - 新协议仓库 MyFlowHub-Proto（PR1）

## Workflow 信息
- Repo：`MyFlowHub-Proto`
- 分支：`refactor/proto-extract`
- Worktree：`d:\project\MyFlowHub3\worktrees\proto-extract\MyFlowHub-Proto`
- 目标 PR：PR1（跨多个 repo 同步提交/合并）

## 项目目标（PR1）
1) 将协议定义（types/常量/子协议号）从 Server 仓库抽离到独立仓库 `MyFlowHub-Proto`。
2) wire 保持不变：action 名称、消息 JSON 结构、SubProto 值均不变。
3) Proto 作为“纯协议字典”：只依赖标准库，不依赖 Core/Server/Win。

## 范围
### 必须（PR1）
- 初始化 `go.mod`（go/toolchain 与现有仓库对齐）
- 提供 `protocol/auth|varstore|topicbus|file|flow|management|exec` 包：
  - const（SubProto、Action、Visibility 等）
  - data struct（请求/响应结构）
  - Message 类型（Action + Data）
- `go test ./...` 通过（至少保证可编译）

### 不做（本 PR）
- 协议 wire 版本化/协商
- assist/up/notify 的 wire 重命名或合并（本轮策略 A）
- 额外 SDK（如 client helper、统一 Encode/Decode 工具）

## 已确认信息
- Go module path：`github.com/yttydcs/myflowhub-proto`
- GitHub remote：`https://github.com/yttydcs/myflowhub-proto.git`

## 问题清单（阻塞：否）
- 无

## 任务清单（Checklist）

### P1 - 初始化 go.mod
- 目标：建立可被 Server/Win replace 引用的 Go module。
- 涉及模块/文件：
  - `go.mod`
- 验收条件：
  - `go test ./...` 通过。
- 回滚点：
  - revert `go.mod`。

### P2 - 迁移现有协议包（从 Server 复制并保持 wire 不变）
- 目标：迁移 `auth/varstore/topicbus/file/flow/management`。
- 涉及模块/文件（预期）：
  - `protocol/*/types.go`
- 验收条件：
  - 对应包可被 `myflowhub-server` 兼容壳引用。
- 测试点：
  - `go test ./... -count=1`
- 回滚点：
  - 逐包提交，便于单独 revert。

### P3 - 新增 exec 协议包
- 目标：将 Server 内 `exec` 子协议的 types/常量抽到 Proto（用户已确认）。
- 涉及模块/文件：
  - `protocol/exec/types.go`
- 验收条件：
  - Server 可改为依赖 `protocol/exec`（不再让 flow import exec handler 包）。

### P4 - 全量回归
- 目标：确保 Proto 仓库本身可编译、可被上层依赖。
- 验收条件：
  - `go test ./...` 通过。

## 依赖关系
- Proto 完成后，Server/Win 才能切换 import/依赖。

## 风险与注意事项
- module path 一旦落地，会影响所有依赖方 import；必须先确认再写代码。
- 严禁在 Proto 中引入 Server/Core 的实现依赖（避免循环依赖）。
