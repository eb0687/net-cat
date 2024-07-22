package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	isSystemMessage = true
)

type User struct {
	username   string
	ipAddress  string
	joinedAt   time.Time
	connection net.Conn
}

var (
	users      = make(map[net.Conn]User)
	usersMutex sync.Mutex
)

func ProcessClient(conn net.Conn) {
	fmt.Print("new client has joined: ")

	defer func() {
		log.Printf("%s has requested to close the connection.", users[conn].username)
		exitMessage := fmt.Sprintf("%s has left our chat...", users[conn].username)
		BroadcastMessage(users[conn], exitMessage, true)
		fmt.Fprintln(conn, "Goodbye!")

		conn.Close()
		usersMutex.Lock()
		delete(users, conn)
		usersMutex.Unlock()

		activeClientsMu.Lock()
		activeClients--
		activeClientsMu.Unlock()
	}()

	DisplayLogo(conn)

	scanner := bufio.NewScanner(conn)
	var username string
	for scanner.Scan() {
		username = scanner.Text()
		if username == "" {
			conn.Write([]byte("Empty username is not allowed! Please enter a valid username: "))
			continue
		}

		usersMutex.Lock()
		duplicate := false
		for _, user := range users {
			if user.username == username {
				duplicate = true
				break
			}
		}
		usersMutex.Unlock()

		if duplicate {
			conn.Write([]byte("Username is already taken. Please enter a different username: "))
			log.Printf("Username already taken: %s", username)
		} else {
			break
		}
	}

	user := User{
		username:   username,
		ipAddress:  conn.RemoteAddr().String(),
		joinedAt:   time.Now(),
		connection: conn,
	}

	usersMutex.Lock()
	users[conn] = user
	fmt.Println(username)
	usersMutex.Unlock()

	joinMessage := fmt.Sprintf("%s has joined our chat...", username)
	log.Printf("%s [%s] has joined the server", username, user.ipAddress)

	BroadcastMessage(user, joinMessage, isSystemMessage)
	SendPreviousMessages(conn)
	DisplayPrompt(conn, username)

	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "exit" {
			break
		}
		if msg == "" {
			fmt.Fprintln(conn, "Empty messages are not allowed!")
			DisplayPrompt(conn, username)
			continue
		}

		BroadcastMessage(user, msg, false)
		DisplayPrompt(conn, username)
	}
}
