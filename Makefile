.PHONY: run mod-vendor linter gosec test validate

run:
	go run main.go

mod-vendor:
	go mod vendor

linter:
	@golangci-lint run

gosec:
	@gosec -quiet ./...

test:
	@go test ./test/ -p 32

validate: linter gosec test