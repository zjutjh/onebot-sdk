# NapCat Go SDK 协作指南

## 项目现状

- 本仓库是面向 NapCat/OneBot 11 的纯 Go SDK，不是机器人框架。
- `api/` 中的 170 个 action 绑定由 NapCat 官方 OpenAPI 4.18.13 生成。
- 生成器使用 `libopenapi` 解析规范并生成模型；请求、响应和组件类型统一写入 `api/types_gen.go`。
- HTTP transport 使用 Resty；WebSocket transport 使用 coder/websocket；JSON 使用 Sonic 的 `ConfigStd`。
- 事件层当前细分私聊消息和群消息，其他事件通过 `UnknownEvent` 保留原始 JSON。
- 消息层提供常用 OneBot 消息段构造、文本提取和类型过滤。
- 根包统一暴露 HTTP、正向 WebSocket、反向 WebSocket 和原始 action 调用。

## 目标

- 保持 SDK 小、稳定、可直接嵌入 Go 项目。
- 强类型 API 与仓库内固定版本的官方 OpenAPI 保持一致。
- 对官方 schema 无法可靠表达的字段保留 `any` 或 map，不猜测协议类型。
- 未识别事件和新版 action 必须仍可通过 raw 入口使用。

## 约束

- 不手改 `api/*_gen.go`；修改 `internal/gen/` 后运行 `go generate ./...`。
- 当前规范唯一来源是 NapCatDocs 的 `src/api/4.18.13/openapi.json`，仓库副本位于 `internal/openapi/4.18.13/openapi.json`。升级时同时更新路径、生成标头和 README。
- 优先使用标准库和现有依赖；新增依赖必须解决已验证的问题。
- 保持 `transport.Caller` 为生成层与传输层之间的最小接口，不引入框架、插件或路由抽象。
- WebSocket 关闭必须唤醒 pending 调用并返回 `ErrClosed`，且只能由读循环关闭事件 channel。
- WebSocket 事件缓冲满时允许丢弃新事件，避免阻塞同一连接上的 API 响应；配置 logger 时必须记录该情况。
- 生成 API 以正确反映官方 schema 为先，不为旧生成结果保留兼容层。
- 文档与注释使用中文，Go 标识符和实现遵循标准 Go 风格。
- 修改后至少运行 `gofmt`、`go test ./...`、`go vet ./...`、`go generate ./...` 和 `git diff --check`。
