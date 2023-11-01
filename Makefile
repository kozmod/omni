.PHONY: sync
sync:
	go work sync

.PHONY: tidy.client
tidy.client:
	pushd ./client && go mod tidy && popd

.PHONY: test.client
test.client:
	go test ./client -cover

.PHONY: lint
lint:
	@go version
	@golangci-lint --version
	pushd ./client && GOWORK=off golangci-lint run ./...

.PHONY: deps
deps:
	@command -v mockery >/dev/null 2>&1 || go install github.com/vektra/mockery/v2@latest
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest