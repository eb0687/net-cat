package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func SetupLogging(filename string) error {
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)
	return nil
}

func DisplayPrompt(conn net.Conn, username string) {
	prompt := fmt.Sprintf("[%s][%s]:",
		time.Now().Format("2006-01-02 15:04:05"),
		username)
	conn.Write([]byte(prompt))
}

func DisplayLogo(conn net.Conn) {
	logo, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Printf("Error reading logo file: %v", err)
		conn.Write([]byte("Welcome to TCP-Chat!\n[ENTER YOUR NAME]: "))
	} else {
		conn.Write([]byte("Welcome to TCP-Chat!\n"))
		conn.Write(logo)
		conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	}
}
