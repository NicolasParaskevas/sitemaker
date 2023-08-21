build:
	@go build -o bin/sitemaker cmd/main.go

run: build
	@./bin/sitemaker

test:
	@go test -v ./...