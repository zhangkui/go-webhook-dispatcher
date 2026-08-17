# go-webhook-dispatcher

## 项目说明

go-webhook-dispatcher 是一个使用 Go 标准库实现的内存 Webhook 投递调度服务。它支持订阅端点注册与停用、事件过滤、HMAC 签名、投递尝试记录、指数退避、死信查询和安全重放。网络投递由可替换的 `Sender` 接口完成，默认示例使用内存模拟器。

## 标准命令

```bash
go build ./...
go test ./...
go vet ./...
go run ./cmd/webhookd
```

## 使用方式

服务默认监听 `:8080`。主要接口：

- `POST /subscriptions` 注册订阅端点。
- `POST /subscriptions/{id}/disable` 停用订阅。
- `POST /events` 发布领域事件。
- `GET /deliveries` 按状态、端点和事件筛选投递。
- `POST /deliveries/{id}/replay` 重放死信投递。

请求和响应均使用 JSON。可通过 `WEBHOOK_ADDR` 环境变量修改监听地址。
