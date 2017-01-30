package whatsup

import (
	"encoding/gob"
	"fmt"
	"net"
)

type ChatConn struct {
	Enc  *gob.Encoder
	Dec  *gob.Decoder
	Conn net.Conn
}

// Connects a chat client to a chat server
func ServerConnect(username string, serverAddr string, serverPort string) (ChatConn, error) {
	chatConn := ChatConn{}
	fmt.Printf("Connecting to %s:%s\n", serverAddr, serverPort)
	conn, err := net.Dial("tcp", serverAddr+":"+serverPort)
	if err != nil {
		return chatConn, err
	}
	chatConn.Conn = conn
	chatConn.Enc = gob.NewEncoder(conn)
	chatConn.Dec = gob.NewDecoder(conn)

	msg := WhatsUpMsg{Username: username, Action: CONNECT}
	chatConn.Enc.Encode(&msg)

	return chatConn, nil
}

func SendMsg(chatConn ChatConn, msg WhatsUpMsg) {
	chatConn.Enc.Encode(&msg)
}

// Receive next WhatsUpMsg from a ChatConn (blocks)
func RecvMsg(chatConn ChatConn) (WhatsUpMsg, error) {
	var chatMsg WhatsUpMsg
	err := chatConn.Dec.Decode(&chatMsg)
	return chatMsg, err
}

func (msg WhatsUpMsg) String() string {
	return fmt.Sprintf("{Username: \"%v\", Body: \"%v\", Action: %v}\n", msg.Username, msg.Body, msg.Action)
}

func (purpose Purpose) String() string {
	switch purpose {
	case CONNECT:
		return "CONNECT"
	case MSG:
		return "MSG"
	case LIST:
		return "LIST"
	case ERROR:
		return "ERROR"
	case DISCONNECT:
		return "DISCONNECT"
	default:
		return "Unknown Purpose!"
	}
}
