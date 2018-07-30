
BUILD_DIR=./build

.PHONY: dep
dep:
	dep ensure -v

.PHONY: dep-update
dep-update:
	dep ensure -v -update ./...

.PHONY: fmt
fmt:
	@gofmt ./...

.PHONY: test
test:
	@go test ./...

.PHONY: cover
cover:
	@GOPWT_OFF=1 go test -race -coverpkg=./... -coverprofile=coverage.txt ./...

.PHONY: test-run
test-run:
	@cat example.md | make run | diff -u example.md -; true

.PHONY: build
build:
	@go build -o $(BUILD_DIR)/markdown-table-formatter

.PHONY: run
run:
	@go run main.go
