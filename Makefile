.PHONY: test run

test:
	go test -v ./...

run:
	go run .
	go fmt ./..
	goimports -w .
