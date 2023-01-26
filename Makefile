CONFIG_PATH=./examples/config.yaml

fmt:
	go fmt ./...
	gofmt -s -w .


vet:
	go vet ./...

golangci-lint:
	golangci-lint run

linter: fmt vet golangci-lint ## Run all linters once

build:
	go build

run: build
	BULKRD_CONFIGPATH=./examples/config.yaml ./bulkrd