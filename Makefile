all: test
	go build -v

.PHONY: test
test:
	go clean -testcache
	go test -v ./...
