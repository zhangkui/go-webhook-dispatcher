# Bug Reproduction

- Case: BUG-004
- Task type: diagnosis
- Feature area: retry-scheduling
- Affected flow: 失败投递的首次及后续指数退避

## Observed Behavior

退避指数直接使用 attempt，首次失败 attempt=1 时计算 2^1 倍基础延迟，整个序列比约定多移一位；最大延迟封顶会掩盖较后次数的差异。

## Reproduction Steps

Run from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestRetryPolicyStartsAtBaseDelay$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The focused command is the reproduction. The full test command confirms repository impact, and the build command confirms compilation.