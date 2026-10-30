GO ?= go
.PHONY: run build test vet check fmt
run:
	$(GO) run ./cmd/outbox-consumer
build:
	$(GO) build -trimpath -o bin/outbox-consumer ./cmd/outbox-consumer
test:
	$(GO) test -race ./...
vet:
	$(GO) vet ./...
check: vet test build
fmt:
	gofmt -w cmd internal
