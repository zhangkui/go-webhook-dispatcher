# Bug Reproduction

- Case: BUG-005
- Task type: bugfix
- Feature area: dead-letter-replay
- Affected flow: POST /deliveries/{id}/replay 的审计记录保留

## Observed Behavior

重放直接把原死信 Delivery 改回 pending 并清零尝试次数，覆盖原始终态，历史 Attempts 与当前状态语义分裂，也无法区分重放任务。

## Reproduction Steps

Run from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestReplayCreatesNewDeliveryAndPreservesDeadLetter$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The focused command is the reproduction. The full test command confirms repository impact, and the build command confirms compilation.