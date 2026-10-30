GO ?= go
.PHONY: run build test vet check fmt
run:
	$(GO) run ./cmd/outbox-relay
build:
	$(GO) build -trimpath -o bin/outbox-relay ./cmd/outbox-relay
test:
	$(GO) test -race ./...
vet:
	$(GO) vet ./...
check: vet test build
fmt:
	gofmt -w cmd internal
