default: run

run:
	@go run playground/main.go
test:
	@go test ./core ./prompts -cover
snap:
	@UPDATE_SNAPSHOTS=true go test ./prompts
format:
	@gofmt -w .