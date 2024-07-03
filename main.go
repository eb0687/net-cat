package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "8989"

type ipAddr string
type user string

type User struct {
	username   string
	ipaddress  string
	joinedAt   time.Time
	connection net.Conn
}

type Message struct {
	content   string
	timeStamp time.Time
	clienName User
}

var Messages = []Message{}
var Users = []User{}

func main() {

	port := defaultPort
	arg := os.Args[1:]

	if len(arg) == 1 {
		port = arg[0]
	} else if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $host")
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

func getPort() string {
	arg := os.Args[1:]
	if len(arg) == 1 {
		return os.Args[1]
	} else if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
	return ""
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
