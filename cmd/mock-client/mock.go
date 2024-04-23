package main

import (
	"AlIM/pkg/tcp"
	"fmt"
	"net"
)

var userName string
var messageType int
var userID, roomID, roomType int

func main() {
	// client config init
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	tcpServer := tcp.NewTcpServer(conn)
	defer tcpServer.Close()

	// User init
	userInit(tcpServer)

	// Start listener
	go listen(tcpServer)

	// Start sender
	sender(tcpServer)
}

// Connect
func userInit(tcpServer *tcp.TcpServer) {
	fmt.Println("Init your user, set Username")
	_, _ = fmt.Scan(&userName)

	err := Send(tcpServer, &tcp.Message{
		UserName: userName,
		RoomID:   roomID,
		UserID:   userID,
		Type:     tcp.ConnectMessage,
	})
	if err != nil {
		fmt.Println("Error sending message", err)
		return
	}

	initResp, err := Receive(tcpServer)
	if err != nil {
		fmt.Println("Error receiving message", err)
		return
	}
	userID = initResp.UserID

	fmt.Println("Connected to server, type /menu to see commands")

	messageType = tcp.SendMessage
}

func sender(tcpServer *tcp.TcpServer) {
	for {
		var input, content string
		_, _ = fmt.Scanln(&input)
		switch input {
		case "/change": // Change room
			fmt.Print("Set your RoomID, RoomType \nRoomType:\t\n1 - Group,\t\n2 - Private\n")
			_, _ = fmt.Scan(&roomID, &roomType)
			fmt.Println("RoomID:", roomID, "RoomType:", roomType)
			messageType = tcp.RoomChangeMessage
		case "/add": // Add friend
			fmt.Print("Set your Friend UserID\n")
			_, _ = fmt.Scan(&roomID)
			messageType = tcp.AddFriendMessage
		case "/list":
			messageType = tcp.ListPublicRoomMessage
		case "/recommend":
			messageType = tcp.RecommendFriendMessage
		case "/menu":
			fmt.Println("Commands:\n" +
				"\t/change - Change room\n" +
				"\t/add - Add friend\n" +
				"\t/list - List public rooms\n" +
				"\t/recommend - Recommend friend\n" +
				"\t/menu - Show commands")
			continue
		default:
			content = input
		}

		err := Send(tcpServer, &tcp.Message{
			UserName: userName,
			RoomID:   roomID,
			RoomType: roomType,
			UserID:   userID,
			Type:     messageType,
			Content:  []byte(content),
		})
		if err != nil {
			fmt.Println("Error sending message", err)
			return
		}

		messageType = tcp.SendMessage
	}
}

func listen(tcpServer *tcp.TcpServer) {
	for {
		message, err := Receive(tcpServer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}
		_ = message
		fmt.Printf("%s#%d : %s\n", message.UserName, message.UserID, message.Content)
	}
}

func Receive(tcpServer *tcp.TcpServer) (*tcp.Message, error) {
	message, err := tcpServer.Receive()
	if err != nil {
		fmt.Println("Error receiving message", err)
		return nil, err
	}
	return message, nil
}

func Send(tcpServer *tcp.TcpServer, message *tcp.Message) error {
	err := tcpServer.Send(message)
	if err != nil {
		fmt.Println("Error sending message", err)
		return err
	}
	return nil
}
