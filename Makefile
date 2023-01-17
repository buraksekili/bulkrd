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
	./bulkrd