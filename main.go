package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	defaultPort = "8989"
	maxClients  = 10
)

var (
	activeClients   int
	activeClientsMu sync.Mutex
)

func main() {
	err := SetupLogging("chat.log")
	if err != nil {
		fmt.Printf("Error setting up logging: %v\n", err)
		return
	}

	port := defaultPort
	arg := os.Args[1:]

	if len(arg) == 1 {
		port = arg[0]
	} else if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	addr := fmt.Sprintf(":%v", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		log.Printf("err: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for connections on %s\n", listener.Addr().String())
	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
			continue
		}

		activeClientsMu.Lock()
		if activeClients >= maxClients {
			activeClientsMu.Unlock()
			conn.Write([]byte("Server is full. Please try again later.\n"))
			log.Printf("Maximum number of clients reached.")
			conn.Close()
			continue
		}

		activeClients++
		activeClientsMu.Unlock()

		go ProcessClient(conn)
	}
}
