# Bug Reproduction

- Case: $bug
- Task type: $taskType
- Feature area: $feature
- Affected flow: $flow

## Observed Behavior

按索引遍历并直接更新切片元素的 Active 字段。

## Reproduction Steps

Run the commands below from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestDisabledSubscriptionStopsFutureDeliveries$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The first command is the focused reproduction. The full test command confirms whether the same behavior affects the repository regression suite. The build command confirms the project still compiles.

## Recorded Baseline Error

See the preserved pre_fix.jsonl trajectory and cases/BUG-001/data/branch-base-red.txt outside the repository for the complete command output.