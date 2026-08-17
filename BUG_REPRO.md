# Bug Reproduction

- Case: BUG-001
- Task type: bugfix
- Feature area: subscription-lifecycle
- Affected flow: POST /subscriptions/{id}/disable 后的事件匹配

## Observed Behavior

方法对 subscriptions 切片使用值拷贝遍历，只修改循环变量，返回成功但存储中的 Active 仍为 true，后续发布继续匹配该端点。

## Reproduction Steps

Run from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestDisabledSubscriptionStopsFutureDeliveries$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The focused command is the reproduction. The full test command confirms repository impact, and the build command confirms compilation.