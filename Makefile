.PHONY: build, test
build:
	go build -o ~/tools/bin/gnote

test:
	go test ./...
