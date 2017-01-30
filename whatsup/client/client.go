package client

import (
	"../whatsup"
	"fmt"
	"bufio"
	"os"
	"strings"
)

func Start(user string, serverPort string, serverAddr string) {
	// Connect to chat server
	chatConn, err := whatsup.ServerConnect(user, serverAddr, serverPort)
	if err != nil {
		fmt.Printf("unable to connect to server: %s\n", err)
		return
	}
	// TODO: Receive input from the user and use the first return value of whatsup.ServerConnect
	// (currently ignored so the stencil will compile) to talk to the server.

	// This channel is used to stop the message `Enter <cmd>...` from printing until the server output is processed.
	// A slight modification is to set a timer in case server never sends a request or something goes wrong while
	// processing the response
	waitFoInput := make(chan bool)

	go func() {
		for {
			msg, err := whatsup.RecvMsg(chatConn)
			if err == nil {
				fmt.Printf("Got from server: %v", msg)
				switch msg.Action {
				case whatsup.LIST:
					fmt.Println(msg.Body)
				case whatsup.DISCONNECT:
					chatConn.Conn.Close()
					return
				case whatsup.MSG:
					fmt.Printf("Got %s from %s\n", msg.Body, msg.Username)
				default:
					fmt.Printf("%v\n", msg)
				}
				<-waitFoInput

			} else {
				fmt.Println(err)
				break
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter <cmd> [<username>] [<body>]: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		parsed := strings.Split(text, " ")
		var cmd, username, body string
		if len(parsed) > 0 {
			cmd = parsed[0]
			if len(parsed) > 1 {
				username = parsed[1]
			}
			if len(parsed) > 2 {
				body = parsed[2]
			}
			fmt.Println("Entered: ", cmd, username, body)
			var action whatsup.Purpose
			switch cmd {
			case "msg":
				action = whatsup.MSG
			case "list":
				action = whatsup.LIST
			case "disconnect":
				action = whatsup.DISCONNECT
			default:
				fmt.Printf("Invalid command \"%v\"\n", cmd)
				continue
			}
			msg := whatsup.WhatsUpMsg{Body: body, Username: username, Action: action}
			chatConn.Enc.Encode(&msg)
			waitFoInput <- true
		}
	}
}
