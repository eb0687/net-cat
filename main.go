package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const defaultPort = "8989"

type User struct {
	username   string
	ipaddress  string
	joinedAt   time.Time
	connection net.Conn
}

type Message struct {
	content    string
	timeStamp  time.Time
	clientName string
}

var (
	messages      []Message
	users         = make(map[net.Conn]User)
	usersMutex    sync.Mutex
	messagesMutex sync.Mutex
)

func main() {
	port := defaultPort
	arg := os.Args[1:]

	if len(arg) == 1 {
		port = arg[0]
	} else if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	addr := fmt.Sprintf("%s:%v", "localhost", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("port is already in use")
		fmt.Printf("err: %v\n", err)
		return
	}
	defer listener.Close()

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			fmt.Println("Connection accepted!")
			go processClient(conn)
		}
	}
}

// func promptUsername(conn net.Conn) string {
// 	conn.Write([]byte("Please enter a username: "))
// 	scanner := bufio.NewScanner(conn)
// 	if scanner.Scan() {
// 		return scanner.Text()
// 	}
// 	return ""
// }

func notifyAll(message string, sender User) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	for conn, user := range users {
		if sender.username != "" && sender.username == user.username {
			continue
		}
		conn.Write([]byte(message + "\n"))
	}
}

func broadcastMessage(user User, content string) {
	msg := Message{
		content:    content,
		timeStamp:  time.Now(),
		clientName: user.username,
	}

	messagesMutex.Lock()
	defer messagesMutex.Unlock()
	messages = append(messages, msg)

	formattedMessage := fmt.Sprintf("[%s][%s]:%s", msg.timeStamp.Format("2006-01-02 15:04:05"), msg.clientName, msg.content)
	notifyAll(formattedMessage, user)

}

func sendPreviousMessages(conn net.Conn) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()

	for _, msg := range messages {
		conn.Write([]byte(fmt.Sprintf("[%s][%s]:%s\n", msg.timeStamp.Format("2006-01-02 15:04:05"), msg.clientName, msg.content)))
	}
}

func processClient(conn net.Conn) {
	fmt.Println("Processing client connection...")
	// msg := "Hello & welcome to net-cat server!" + "\n" + conn.RemoteAddr().String() + "\n"
	// conn.Write([]byte(msg))
	// fmt.Printf("SERVER: %v\n", msg)

	logo, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Printf("Error reading logo file: %v", err)
		conn.Write([]byte("Welcome to TCP-Chat!\n[ENTER YOUR NAME]: "))
	} else {
		conn.Write([]byte("Welcome to TCP-Chat!\n"))
		conn.Write(logo)
		conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	}

	// defer conn.Close()
	scanner := bufio.NewScanner(conn)
	var username string
	if scanner.Scan() {
		username = scanner.Text()
	}
	// username := promptUsername(conn)
	if username == "" {
		fmt.Fprintln(conn, "Empty username is not allowed!")
		fmt.Fprintln(conn, "Reconnect to the server and try again.")
		conn.Close()
		return
	}

	user := User{
		username:   username,
		ipaddress:  conn.RemoteAddr().String(),
		joinedAt:   time.Now(),
		connection: conn,
	}

	usersMutex.Lock()
	users[conn] = user
	usersMutex.Unlock()

	joinMessage := fmt.Sprintf("%s has joined our chat...", username)
	notifyAll(joinMessage, user)
	sendPreviousMessages(conn)

	// scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "exit" {
			fmt.Println("Client requested to close the connection.")
			exitMessage := fmt.Sprintf("%s has left our chat...", username)
			notifyAll(exitMessage, user)
			fmt.Fprintln(conn, "Goodbye!")
			conn.Close()
			break
		}
		if msg == "" {
			fmt.Fprintln(conn, "Empty messages are not allowed!")
			continue
		}
		broadcastMessage(user, msg)
	}

	usersMutex.Lock()
	delete(users, conn)
	usersMutex.Unlock()
}
