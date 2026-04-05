# MyFlowHub-Proto

MyFlowHub 的协议定义仓库（types/常量/协议号等）。

- 目标：供 MyFlowHub-Server / MyFlowHub-Win 复用，减少重复与耦合。
- 说明：本仓库仅放协议定义，不放 Server 业务实现。

---

## Single source-of-truth（重要）

本仓库的 **single source-of-truth** 为：`protocol/*/types.go`。

- `SubProto*`：子协议编号
- `Action*`：wire action 字符串常量
- 请求/响应 payload struct（JSON tag / omitempty 等）
- 其它协议相关常量（例如 `Kind*`/`Op*`/`Perm*` 等）

> 约束：不得在本仓库引入运行时逻辑/网络逻辑/handler 逻辑；wire 不改（除非另起 workflow 明确破坏性变更）。

---

## 协议文档与生成产物

本仓库维护两类 canonical 衍生物：

- `docs/protocol_map.md`
  - 子协议、action、type、常量总览
- `docs/flow_contract.md`
  - Flow graph/node contract 的 typed 说明
- `generated/flow_contract.ts`
  - 供下游同步消费的 Flow contract TypeScript 类型产物

### Flow contract 生成 / 校验命令

- 更新 Flow contract：
  - `go run ./cmd/flowcontractgen -write -md-out docs/flow_contract.md -ts-out generated/flow_contract.ts`
- 校验 Flow contract 是否最新：
  - `go run ./cmd/flowcontractgen -check -md-out docs/flow_contract.md -ts-out generated/flow_contract.ts`

### 协议映射文档（半自动生成）

本仓库维护 canonical 文档：`docs/protocol_map.md`。

- 该文件为“半自动文档”：生成器只会覆盖 `<!-- BEGIN GENERATED -->` 与 `<!-- END GENERATED -->` 中间区块。
- 请不要手工修改生成区块内容；手写说明请写在区块外（例如 `Notes（Manual）`）。

### 生成 / 校验命令

- 更新文档（写入）：
  - `go run ./cmd/protocolmapgen -write -out docs/protocol_map.md`
- 校验文档是否最新（不写入，非 0 表示需要更新）：
  - `go run ./cmd/protocolmapgen -check -out docs/protocol_map.md`

### 强制校验范围（重要）

仅在 `MyFlowHub-Proto` 仓库内强制：`go test ./...` 会校验 `docs/protocol_map.md` 是否与 `protocol/*/types.go` 一致。

---

## 同步到 MyFlowHub3 meta workspace（副本）

`d:\project\MyFlowHub3` 会保留一份同步副本：`d:\project\MyFlowHub3\docs\protocol_map.md`。

- canonical 以本仓库的 `docs/protocol_map.md` 为准；meta workspace 副本不做门禁。
- 同步方式（整文件复制）：
  - `Copy-Item d:\\project\\MyFlowHub3\\repo\\MyFlowHub-Proto\\docs\\protocol_map.md d:\\project\\MyFlowHub3\\docs\\protocol_map.md -Force`
