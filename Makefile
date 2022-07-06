deps:
	go mod tidy
	go mod vendor

test: deps
	go test ./...
