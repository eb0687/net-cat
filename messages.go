package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Message struct {
	content       string
	timeStamp     time.Time
	clientName    string
	systemMessage bool
}

var (
	messages      []Message
	messagesMutex sync.Mutex
)

func isAllowedChar(c rune) bool {
	allowedSpecialChars := "!@#$%^&*()_+-={}[]|\\:;\"'<>,.?/ "
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
		return true
	}
	if strings.ContainsRune(allowedSpecialChars, c) {
		return true
	}
	return false
}

func isASCII(s string) bool {
	for _, v := range s {
		if !isAllowedChar(v) {
			return false
		}
	}
	return true
}

func NotifyAll(message string, sender User, isSystemMessage bool) {
	if !isASCII(message) {
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	for conn, user := range users {
		if sender.username != "" && sender.username == user.username {
			continue
		}
		prompt := fmt.Sprintf("[%s][%s]:",
			time.Now().Format("2006-01-02 15:04:05"),
			user.username)

		formattedMessage := fmt.Sprintf(
			"[%s][%s]:%s",
			time.Now().Format("2006-01-02 15:04:05"),
			sender.username,
			message,
		)

		conn.Write([]byte("\n"))
		if isSystemMessage {
			conn.Write([]byte(message))
			conn.Write([]byte("\n" + prompt))
		} else {
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
			fmt.Fprintf(conn, "[%s][%s]:%s\n",
				msg.timeStamp.Format("2006-01-02 15:04:05"),
				msg.clientName,
				msg.content)
		}
	}
}
