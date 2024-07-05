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

const (
	defaultPort = "8989"
	maxClients  = 10
)

type User struct {
	username   string
	ipaddress  string
	joinedAt   time.Time
	connection net.Conn
}

type Message struct {
	content       string
	timeStamp     time.Time
	clientName    string
	systemMessage bool
}

var (
	messages        []Message
	users           = make(map[net.Conn]User)
	usersMutex      sync.Mutex
	messagesMutex   sync.Mutex
	activeClients   int
	activeClientsMu sync.Mutex
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
		// conn, err := listener.Accept()
		// if err != nil {
		// 	log.Printf("Error accepting connection from client: %s", err)
		// } else {
		// 	fmt.Println("Connection accepted!")
		// 	go processClient(conn)
		// }
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
			continue
		}

		activeClientsMu.Lock()
		if activeClients >= maxClients {
			activeClientsMu.Unlock()
			conn.Write([]byte("Server is full. Please try again later.\n"))
			conn.Close()
			continue
		}

		activeClients++
		activeClientsMu.Unlock()

		go ProcessClient(conn)
	}
}

func NotifyAll(message string, sender User, isSystemMessage bool) {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	for conn, user := range users {
		if sender.username != "" && sender.username == user.username {
			continue
		}
		prompt := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), user.username)

		conn.Write([]byte("\n"))
		if isSystemMessage {
			// conn.Write([]byte(message + "\n"))
			conn.Write([]byte(message))
		} else {
			formattedMessage := fmt.Sprintf("[%s][%s]:%s", time.Now().Format("2006-01-02 15:04:05"), sender.username, message)
			conn.Write([]byte(formattedMessage + "\n" + prompt))
		}
	}
}

func BroadcastMessage(user User, content string, isSystemMessage bool) {
	msg := Message{
		content:       content,
		timeStamp:     time.Now(),
		clientName:    user.username,
		systemMessage: isSystemMessage,
	}

	messagesMutex.Lock()
	defer messagesMutex.Unlock()
	messages = append(messages, msg)

	if isSystemMessage {
		NotifyAll(content, user, true)
	} else {
		NotifyAll(content, user, false)
	}
}

func SendPreviousMessages(conn net.Conn) {
	messagesMutex.Lock()
	defer messagesMutex.Unlock()

	for _, msg := range messages {
		if msg.systemMessage {
			conn.Write([]byte(msg.content + "\n"))
		} else {
			fmt.Fprintf(conn, "[%s][%s]:%s\n", msg.timeStamp.Format("2006-01-02 15:04:05"), msg.clientName, msg.content)
		}
	}
}

func ProcessClient(conn net.Conn) {
	fmt.Println("Processing client connection...")

	defer func() {
		conn.Close()
		usersMutex.Lock()
		delete(users, conn)
		usersMutex.Unlock()

		activeClientsMu.Lock()
		activeClients--
		activeClientsMu.Unlock()
	}()

	logo, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Printf("Error reading logo file: %v", err)
		conn.Write([]byte("Welcome to TCP-Chat!\n[ENTER YOUR NAME]: "))
	} else {
		conn.Write([]byte("Welcome to TCP-Chat!\n"))
		conn.Write(logo)
		conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	}

	scanner := bufio.NewScanner(conn)
	var username string
	if scanner.Scan() {
		username = scanner.Text()
	}
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
	BroadcastMessage(user, joinMessage, true)
	SendPreviousMessages(conn)

	displayPrompt := func() {
		prompt := fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), username)
		conn.Write([]byte(prompt))
	}

	displayPrompt()

	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "exit" {
			log.Printf("%s has requested to close the connection.", username)
			exitMessage := fmt.Sprintf("%s has left our chat...", username)
			BroadcastMessage(user, exitMessage, true)
			fmt.Fprintln(conn, "Goodbye!")
			conn.Close()
			break
		}
		if msg == "" {
			fmt.Fprintln(conn, "Empty messages are not allowed!")
			displayPrompt()
			continue
		}

		BroadcastMessage(user, msg, false)
		displayPrompt()
	}

	usersMutex.Lock()
	delete(users, conn)
	usersMutex.Unlock()
}
