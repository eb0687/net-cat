- [x] TCP connection between server and multiple clients (relation of 1 to many).
- [x] A name requirement to the client.
- [ ] Control connections quantity.
- [x] Clients must be able to send messages to the chat.
- [x] Do not broadcast EMPTY messages from a client.
- [x] Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : [2020-01-20 15:48:41][client.name]:[client.message]
- [x] If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
- [x] If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
- [x] If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
- [x] All Clients must receive the messages sent by other Clients.
- [x] If a Client leaves the chat, the rest of the Clients must not disconnect.
- [x] If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: [USAGE]: ./TCPChat $port0

## Instructions

- [x] Start TCP server, listen and accept connections
- [x] Your project must have Go-routines
- [x] Your project must have channels or Mutexes
- [x] Maximum 10 connections
- [ ] It is recommended to have test files for unit testing both the server connection and the client.

## Bonus

- [ ] Terminal UI (you are allowed to use only this package : https://github.com/jroimartin/gocui).
- [x] Find a way to save all the logs into a file.
- [ ] Creating more than 1 group chat.
