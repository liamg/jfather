.PHONY: lint_install
lint_install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint
lint: lint_install
	golangci-lint run ./...

.PHONY: lint_fix
lint_fix: lint_install
	golangci-lint run ./... --fix

test:
	go test ./...