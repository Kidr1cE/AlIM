package main

import (
	"AlIM/pkg/server"
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
	fmt.Println("Init your user, set Username, UserID")
	_, _ = fmt.Scan(&userName, &userID)
	fmt.Println("Set user name:", userName, "UserID:", userID)

	err := Send(tcpServer, &tcp.Message{
		UserName: userName,
		RoomID:   roomID,
		UserID:   userID,
		Type:     server.ConnectMessage,
	})
	if err != nil {
		fmt.Println("Error sending message", err)
		return
	}

	fmt.Println("Set your RoomID, MessageType")
	_, _ = fmt.Scan(&roomID, &messageType)

	messageType = server.SendMessage
}

func sender(tcpServer *tcp.TcpServer) {
	for {
		var input, content string
		_, _ = fmt.Scanln(&input)
		switch input {
		case "/change": // Change room
			fmt.Print("Set your RoomID, RoomType\nRoomType: 1 - Group, 2 - Private\n")
			_, _ = fmt.Scan(&roomID, &roomType)
			messageType = server.RoomChangeMessage
		case "/add": // Add friend
			fmt.Print("Set your RoomID\n")
			_, _ = fmt.Scan(&roomID)
			messageType = server.AddFriendMessage
		case "/list":
			messageType = server.ListPublicRoomMessage
		case "/recommend":
			messageType = server.RecommendFriendMessage
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

		messageType = server.SendMessage
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
