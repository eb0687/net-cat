build:
    go build -o bin/
run: build
    go run main.go
test:
    go test -v ./...
