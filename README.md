# Net-Cat

## Description

This project is a recreation of the NetCat (nc) command-line utility in a Server-Client architecture.

NetCat (nc) is a versatile networking tool that reads and writes data across network connections using TCP or UDP. This implementation focuses on TCP connections, enabling features like opening TCP connections, listening on arbitrary TCP ports, and facilitating a group chat.

## Features

- TCP connection between server and multiple clients (1 to many).
- Maximum number of simultaneous connections to 10.
- Message Broadcast: Clients can send messages to the chat.
- Message Filtering: Do not broadcast empty messages.
- Message Format: Messages are timestamped and include the sender's name.  
  `[YYYY-MM-DD HH:MM:SS][client.name]: [client.message]`
- Message History: New clients receive all previous messages upon joining.
- Join/Leave Notifications: Clients are informed when others join or leave.
- Continuous Operation: Clients leaving do not disrupt the connection for others.
- This project utilizes go-routines and mutexes to handle concurrent client connections and ensure thread-safe access to shared resources.

## Build

In a terminal navigate to the root of the project directory & run the below command to build the binary.

```
$ make build
```

## Usage

```
$ cd bin/
$ ./TCPChat <port>
```

If no port is specified, the server will be running on the default port of **8989**.

To test the connection, launch another terminal on the host machine or another device and run the below command:

```
nc <ipaddress> <port>
```

## Authors

- [nhussain](https://learn.reboot01.com/git/nhussain)
- [yabuzuha](https://learn.reboot01.com/git/yabuzuha)
- [etarada](https://learn.reboot01.com/git/etarada)
