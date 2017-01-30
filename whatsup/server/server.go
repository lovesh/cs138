package server

import (
	"../whatsup"

	"fmt"
	"net"
	"encoding/gob"
	"strings"
)


type Client struct {
	conn net.Conn
	ch   chan whatsup.WhatsUpMsg
}

var clientConnections map[string]Client = make(map[string]Client)

func handleConnection(conn net.Conn) {
	fmt.Println("Connection from ", conn.RemoteAddr())
	chatConn := whatsup.ChatConn{}
	chatConn.Conn = conn
	chatConn.Enc = gob.NewEncoder(conn)
	chatConn.Dec = gob.NewDecoder(conn)
	defer conn.Close()
	var name string
	for {
		msg, err := whatsup.RecvMsg(chatConn)
		if err == nil {
			fmt.Printf("Got from client %v", msg)
			switch msg.Action {
			case whatsup.CONNECT:
				ch := make(chan whatsup.WhatsUpMsg)
				clientConnections[msg.Username] = Client{conn, ch}
				//clientNames = append(clientNames, msg.Username)
				name = msg.Username
				go func(msgs chan whatsup.WhatsUpMsg) {
					for {
						msg, ok := <- msgs
						if !ok {
							break
						}
						fmt.Printf("Sending to %s ", name)
						fmt.Println("msg: ", msg)
						chatConn.Enc.Encode(&msg)
					}
					fmt.Println("Exiting go func")
				}(ch)
			case whatsup.LIST:
				var clientNames []string
				for k := range clientConnections {
					clientNames = append(clientNames, k)
				}
				reply := whatsup.WhatsUpMsg{Body: "Connected users: [" + strings.Join(clientNames, ", ") + "]",
					Username: "", Action: whatsup.LIST}
				clientConnections[name].ch <- reply
			case whatsup.MSG:
				reply := whatsup.WhatsUpMsg{Body: msg.Body, Username: name, Action: whatsup.MSG}
				clientConnections[msg.Username].ch <- reply
			case whatsup.DISCONNECT:
				delete(clientConnections, name)
				close(clientConnections[name].ch)
				break
			default:
				fmt.Printf("%v", msg)
			}

		} else {
			fmt.Println(err)
			break
		}
	}
}

func Start() {
	listen, port, err := whatsup.OpenListener()
	fmt.Printf("Listening on port %d\n", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listen.Accept() // this blocks until connection or error
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
