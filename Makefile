all: lint test

TOOLS_DIR := $(HOME)/.cache/deepseek-go/bin
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint

$(GOLANGCI_LINT):
	mkdir -p $(TOOLS_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLS_DIR) v1.64.8

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

test:
	go test -v ./...

test-short:
	go test -v -short ./...

test-race:
	go test -v -race ./...

test-integration:
	DEEPSEEK_LIVE_TESTS=1 go test -v -tags=integration ./...
