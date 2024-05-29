default: run

APP_NAME=go-clack

run:
	@go run playground/main.go
build:
	@go build -a -o $(APP_NAME) cmd/api/main.go
test:
	@go test ./ ...
clean:
	@rm -f $(APP_NAME)
format:
	@gofmt -w .