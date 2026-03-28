# Change Archive

## Purpose

存放 `MyFlowHub-Proto` 已完成 workflow 的结果、验证方式和回滚说明。

## Rules

- 文件名使用 `YYYY-MM-DD_topic.md`
- 变更完成后补索引
- `protocol/*/types.go` 仍是长期协议真相；`docs/change/` 只记录结果

## Current Docs

- 最新归档：
  - [2026-03-28_stream-wire.md](2026-03-28_stream-wire.md)
  - `stream` 子协议 canonical wire：新增 `SubProto=8`、source/consumer/delivery 动作、内建 kind 与 DATA/ACK 小头。
  - [2026-03-26_auth-semi-central-authority.md](2026-03-26_auth-semi-central-authority.md)
  - [2026-03-25_proto-flow-delete-release.md](2026-03-25_proto-flow-delete-release.md)
  - [2026-03-22_management-node-display-name.md](2026-03-22_management-node-display-name.md)
- 历史归档保留在当前目录中，按文件名日期倒序浏览
