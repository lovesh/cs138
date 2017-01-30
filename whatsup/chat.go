package main

import (
	"./client"
	"./server"

	"flag"
	"fmt"
)

func main() {
	isServerPtr := flag.Bool("server", false, "run as a server?")
	isClientPtr := flag.Bool("client", false, "run as a client?")
	userPtr := flag.String("user", "", "username used for this chat client")
	serverPortPtr := flag.String("port", "", "chat server port to connect to")
	serverAddrPtr := flag.String("addr", "", "address to chat server")
	flag.Parse()

	if *isClientPtr && *isServerPtr {
		fmt.Println("Please rerun with only one: -server or -client")
		return
	}

	if *isServerPtr {
		fmt.Println("Starting server")
		server.Start()
	} else if *isClientPtr {
		fmt.Println("Starting client")
		client.Start(*userPtr, *serverPortPtr, *serverAddrPtr)
	} else {
		fmt.Println("Please rerun with either -server or -client")
	}
}
