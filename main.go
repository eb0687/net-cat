package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type ipAddr string
type user string

type User struct {
	username   string
	ipaddress  string
	joinedAt   time.Time
	connection net.Conn
}

type Message struct {
	messageBody string
	sentAt      time.Time
	sentBy      User
}

var Messages = []Message{}
var Users = []User{}

func main() {
	addr := fmt.Sprintf("%s:%v", "10.1.204.172", "3333")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("port is already in use")
		fmt.Printf("err: %v\n", err)
		return
	}

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

func processClient(conn net.Conn) {
	fmt.Println("Processing client connection...")
	msg := "Hello from server!" + "\n" + conn.RemoteAddr().String() + "\n"
	conn.Write([]byte(msg))
	fmt.Printf("msg: %v\n", msg)

	s := bufio.NewScanner(conn)
	conn.Write([]byte("Please enter a username: "))

	isfirstmsg := true

	username := ""

	for s.Scan() {

		msg := s.Text()
		if isfirstmsg {
			username = msg
			Users = append(Users, User{
				username:   username,
				ipaddress:  conn.RemoteAddr().String(),
				joinedAt:   time.Now(),
				connection: conn,
			})
			isfirstmsg = false
		}
		if msg == "exit" {
			fmt.Println("Client requested to close the connection.")
			fmt.Fprintln(conn, "Goodbye!")
			break
		}
		if msg == "" {
			fmt.Fprintln(conn, "Empty messages are not allowed!")
			continue
		}

		for _, c := range Users {
			fmt.Printf("c: %+v\n", c)
			if username == c.username {
				continue
			}
			fmt.Fprintf(c.connection, "[%v][%v]:%v\n[%v][%v]:",
				time.Now().Format("2020-01-20 15:48:41"),
				username,
				msg,
				time.Now().Format("2020-01-20 15:48:41"),
				c.username)
		}

	}

	fmt.Println("finished processing client connection.")
	conn.Close()
}
