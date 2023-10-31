.PHONY: sync
sync:
	go work sync

.PHONY: tidy.client
tidy.client:
	pushd ./client && go mod tidy && popd

.PHONY: test.client
test.client:
	go test ./client -cover

.PHONY: deps
deps:
	@command -v mockery >/dev/null 2>&1 || go install github.com/vektra/mockery/v2@latest