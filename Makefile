build:
	go build -o bin/TCPChat
run: build
	go run main.go
test:
	go test -v ./...
