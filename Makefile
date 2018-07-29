
BUILD_DIR=./build

.PHONY: test
test:
	@go test

.PHONY: build
build:
	@go build -o $(BUILD_DIR)/markdown-table-formatter

.PHONY: run
run:
	@go run cmd/main.go
