# Bug Reproduction

- Case: $bug
- Task type: $taskType
- Feature area: $feature
- Affected flow: $flow

## Observed Behavior

只使用 secrets[version] 对应的密钥，不被其他版本覆盖。

## Reproduction Steps

Run the commands below from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestSignPayloadUsesRequestedKeyVersion$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The first command is the focused reproduction. The full test command confirms whether the same behavior affects the repository regression suite. The build command confirms the project still compiles.

## Recorded Baseline Error

See the preserved pre_fix.jsonl trajectory and cases/BUG-003/data/branch-base-red.txt outside the repository for the complete command output.