default: run

run:
	@go run playground/main.go
test:
	@go test ./core -cover
format:
	@gofmt -w .