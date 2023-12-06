lint:
	golangci-lint run
	
fmt:
	go fmt ecommerce/...

run:
	go run cmd/main.go


all: fmt run