default: run

run:
	@go run playground/main.go
test:
	@go test ./core -cover; go test -p 1 ./prompts -cover
snap:
	@UPDATE_SNAPSHOTS=true go test -p 1 ./prompts
format:
	@gofmt -w .