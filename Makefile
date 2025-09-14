LOCAL_BIN := $(CURDIR)/bin

.PHONY: bin-deps
bin-deps:
	$(info Installing binary dependencies...)

	mkdir -p $(LOCAL_BIN)

	ls $(LOCAL_BIN)/golangci-lint &> /dev/null || GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

.PHONY: lint
lint:
	$(info Running lint...)
	$(LOCAL_BIN)/golangci-lint run ./...

.PHONY: test
test:
	go test -count=1 -race ./...

.PHONY: run-sever
run-sever:
	go run ./cmd/server/main.go -verbose -address 127.0.0.1:8881

.PHONY: run-client
run-client:
	go run ./cmd/client/main.go -address 127.0.0.1:8881