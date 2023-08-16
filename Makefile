test:
	go test ./...

lint:
	golangci-lint run

run:
	go run cmd/cli/main.go -f contrib/pli07.txt