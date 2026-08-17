# Bug Reproduction

- Case: BUG-003
- Task type: bugfix
- Feature area: signature-key-rotation
- Affected flow: 按投递记录密钥版本生成 HMAC

## Observed Behavior

函数虽然验证请求版本存在，却继续扫描密钥并用更高版本覆盖 secret；历史任务保留旧 KeyVersion，但签名实际使用新密钥，接收方按记录版本验签失败。

## Reproduction Steps

Run from /app:

`ash
go test -buildvcs=false -count=1 -run "^TestSignPayloadUsesRequestedKeyVersion$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

The focused command is the reproduction. The full test command confirms repository impact, and the build command confirms compilation.