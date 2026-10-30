GO ?= go
GOLANGCI_LINT ?= golangci-lint
.PHONY: run build test vet check fmt lint db-up db-down migrate
run:
	$(GO) run ./cmd/api
build:
	$(GO) build -trimpath -o bin/api ./cmd/api
test:
	$(GO) test -race ./...
vet:
	$(GO) vet ./...
lint:
	$(GOLANGCI_LINT) run ./...
check: lint vet test build
fmt:
	gofmt -w cmd internal pkg
db-up:
	docker compose up -d --wait postgres
db-down:
	docker compose down
