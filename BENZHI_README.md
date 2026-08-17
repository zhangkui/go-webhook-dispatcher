# BENZHI Evaluation Guide

## Project

- Repository: $repo
- Purpose: in-memory webhook filtering, signing, delivery retry, and dead-letter replay service.
- Go toolchain: golang:1.22.
- Frontend toolchain: none.

## Standard Build, Run, and Test Commands

Run inside the container:

`ash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd/webhookd
cd '/app' && GOTOOLCHAIN=local go test -buildvcs=false -count=1 ./...
`

## Docker Build and Shell

`ash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-webhook-dispatcher-bug3-candidate linux/amd64
./build_benzhi_docker.sh go-webhook-dispatcher-bug3-candidate-arm64 linux/arm64
docker run --rm -it --platform linux/amd64 go-webhook-dispatcher-bug3-candidate bash
docker run --rm -it --platform linux/arm64 go-webhook-dispatcher-bug3-candidate-arm64 bash
`

## Task Verification Commands

`ash
go test -buildvcs=false -count=1 -run "^TestSignPayloadUsesRequestedKeyVersion$" ./internal/webhook
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
`

Task type: $taskType. For diagnosis tasks, the target test remains red by design and serves as the reproduction evidence; the build command remains green.

## Bug Reproduction

See BUG_REPRO.md for the observed behavior and expected evidence.