# Plan Archive - 2026-03-29 Proto stream 主线合并

## Workflow 信息
- Repo：`MyFlowHub-Proto`
- Branch：`fix/proto-stream-mainline-merge`
- Base：`main @ 28e827d`
- Worktree：`D:\project\MyFlowHub3\worktrees\fix-proto-stream-mainline-merge`

## 目标
- 将 `feat/proto-stream-subproto` 中已经完成的 `stream` canonical wire 合并回 `repo/MyFlowHub-Proto/main`。
- 让根 workspace 默认命中的 Proto 主线路径直接提供 `protocol/stream`。

## 已确认约束
- `MyFlowHub-Proto` 仍保持纯协议字典边界，不引入运行时逻辑。
- `docs/protocol_map.md` 为受生成器门禁保护的 canonical 文档。
- 来源 worktree 的 `.gitignore` 未跟踪改动不能带入主线。
- 本轮不承担远端 semver tag 发布。

## 任务映射
- `PROTO-STRM-M1`
  - 合并 `feat/proto-stream-subproto`，解决 `docs/change/README.md` 冲突。
- `PROTO-STRM-M2`
  - 执行 `GOWORK=off go test ./... -count=1 -p 1`。
- `PROTO-STRM-M3`
  - 归档 `docs/change/2026-03-29_proto-stream-mainline-merge.md`。
- `PROTO-STRM-M4`
  - 将验证通过的分支合回 `repo/MyFlowHub-Proto/main`。

## 验收标准
- 主线包含 `protocol/stream/types.go`。
- `docs/protocol_map.md` 包含 `SubProto=8 / protocol/stream`。
- Proto 仓库自身 `GOWORK=off go test ./... -count=1 -p 1` 通过。
