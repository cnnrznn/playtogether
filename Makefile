all: test
	go build -v

.PHONY: test
test:
	go test -v ./...
