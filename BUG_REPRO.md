# Bug Reproduction

- Case: BUG-002
- Task type: diagnosis
- Feature area: event-idempotency
- Affected flow: POST /events 重复事件 ID 冲突处理

## Observed Behavior

发布流程先枚举订阅并写入投递任务，最后才调用 AddEvent 检查事件 ID；重复 ID 返回冲突时，副作用已经发生，导致重复投递任务残留。

## Reproduction Steps

Run from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestDuplicateEventDoesNotCreateAdditionalDeliveries$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The focused command is the reproduction. The full test command confirms repository impact, and the build command confirms compilation.