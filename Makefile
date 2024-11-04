.PHONY: test

run:
	go mod tidy && go run cmd/main.go

test:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

clean:
	rm -rf coverage.*