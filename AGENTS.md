# onebot-sdk 协作指南

## 项目现状

- 本仓库是面向 OneBot 11 的纯 Go 统一接入层 SDK，不是机器人框架；当前支持 NapCat 与 SnowLuma 两种后端，未知后端降级为通用 OneBot 11 模式。
- `api/` 中的 178 个 action 绑定由 NapCat 官方 OpenAPI 4.18.33 生成。
- 生成器使用 `libopenapi` 解析规范并生成模型；请求、响应和组件类型统一写入 `api/types_gen.go`。
- HTTP transport 使用标准库 `net/http`；WebSocket transport 使用 coder/websocket；JSON 使用 Sonic 的 `ConfigStd`。
- `dialect/` 是方言行为表（NapCat/SnowLuma/Generic），`Generic` 全 false，代表纯 OB11 核心能力、零补拉零扩展。
- 事件层覆盖标准 OB11 核心集（消息/通知/请求/元事件），ID 兼容 JSON number/string；未识别事件通过 `UnknownEvent` 保留原始 JSON，强类型解析失败通过 `RawEvent` 交付。
- 消息层是强类型段（`SegmentData` 接口 + 类型化实现，全部值类型）；`ForwardData` 内含 `ID/Content/FetchErr`。
- 归一化管线（根包 `pipeline.go`）与手写方言门面（`facade.go`）解决两端协议差异。
- 根包统一暴露 HTTP、正向 WebSocket、反向 WebSocket、方言检测和原始 action 调用。

## 目标

- 保持 SDK 小、稳定、可直接嵌入 Go 项目。
- 强类型 API 与仓库内固定版本的官方 OpenAPI 保持一致。
- 对官方 schema 无法可靠表达的字段保留 `any` 或 map，不猜测协议类型。
- 未识别事件和新版 action 必须仍可通过 raw 入口使用。
- 方言差异在 SDK 内归一，上游不感知后端身份。

## 约束

- 不手改 `api/*_gen.go`；修改 `internal/gen/` 后运行 `go generate ./...`。
- 当前规范唯一来源是 NapCatDocs 的 `src/api/4.18.33/openapi.json`，仓库副本位于 `internal/openapi/4.18.33/openapi.json`。升级时同时更新路径、生成标头和 README。
- 方言表扩展规则：新增后端差异时先在 `dialect/` 加行为位，再在管线或门面统一处理；行为位必须取两端并集语义，缺席为零值。新增后端（如 Lagrange）只加方言表行，不改生成物。
- 方言差异做双路径归一时，宁缺毋滥：补拉失败保留原始 id 并暴露错误（`FetchErr`），不伪造内容、不吞事件。
- 优先使用标准库和现有依赖；新增依赖必须解决已验证的问题。
- 保持 `transport.Caller` 为 api/方言层与传输层之间的最小调用接口，不引入框架、插件或路由抽象。
- WebSocket 关闭必须唤醒 pending 调用并返回首次终止原因；读循环是传输原始队列（`TakeRawEvents`）唯一的关闭者，用户事件通道由归一化管线唯一关闭，两条通道各自只有一个写入者和一个关闭者。管线在交付阻塞期间通过 `Done()` 观察连接终止，连接一停就立即退出，不等满交付超时。
- 管线保序：事件归一化由单 goroutine 串行处理（解析 → 归一 → 投递），投递顺序与到达顺序一致；普通事件零 API 调用快速路径，不得引入并发投递。
- 读循环永不阻塞：它是 action(echo) 响应唯一的路由方，阻塞会让在途调用（含管线自身的补拉 API 往返）永远等不到响应。因此传输层到管线这一跳固定为非阻塞投递，队列满一律丢弃并记日志；`TakeRawEvents()` 返回的主队列由管线唯一消费。
- 背压只在管线到用户这一跳（`deliverEvent`）：事件队列必须有界；未配置 `WithEventDeliveryTimeout` 时队列满丢弃新事件并记日志，配置正数超时后用户持续不消费以 `ErrEventBackpressure` 终止连接。方言补拉的 API 往返不计入交付超时。
- `RawEvents()` 是旁路（transport 的 `TeeRawEvents`）：与 `Events()` 并存，每帧复制一份；首次调用才启用，须在事件到达前调用，之后不回放；旁路队列满丢弃并记日志，不影响主事件流。禁止把管线消费者切换到旁路通道。
- 生成 API 以正确反映官方 schema 为先，不为旧生成结果保留兼容层。数值 ID（user_id/group_id/message_id/self_id/operator_id/sender_id/target_id）在生成后处理中统一改写为 `message.ID` 值类型（解码兼容字符串与数字，编码为数字），改写规则维护在 `internal/gen` 的 `idFieldJSONNames`。
- message_id 只做等值比较，禁止排序、范围比较或当水位用（两端数值语义不可比：NapCat 31 位正数 vs SnowLuma 有符号 int32）。
- 禁止照抄 SnowLuma 源码中的 schema 描述文案（源码可见非商业许可），文档自写。
- 文档与注释使用中文，Go 标识符和实现遵循标准 Go 风格。
- 双后端事件样本位于 `testdata/napcat/` 与 `testdata/snowluma/`，字段级 golden 回归必须全绿；修改事件解析时同步补充两后端样本。
- 修改后至少运行 `gofmt`、`go test -race ./...`、`go vet ./...`、`go generate ./...` 和 `git diff --check`（CI 会执行同等检查，生成物零漂移以 `go generate` 后 `git status --porcelain` 为空为准）。
