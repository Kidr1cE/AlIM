package server

import (
	"AlIM/pkg/room"
	"AlIM/pkg/session"
	"AlIM/pkg/store"
	"AlIM/pkg/tcp"
	"fmt"
)

func GroupHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("GroupHandler: ", message.String())
	if session.Room == nil {
		// group chat
		session.Room = room.GetRoom(message.RoomID)
		session.Room.AddClient(message.UserID, session.TcpServer)
	}
	session.Room.BroadcastMessage(*message)
}

// PrivateHandler :if the message type is private, message.RoomID is the receiver's UserID
func PrivateHandler(session *session.Session, message *tcp.Message) {
	// Check if the sender and receiver are friends
	if !store.IsFriend(message.UserID, message.RoomID) {
		fmt.Printf("User %d is not friends with user %d\n", message.UserID, message.RoomID)
		// If they are not friends, send an error message to the sender
		session.TcpServer.Send(&tcp.Message{
			UserName: "AlIM Server",
			RoomID:   message.RoomID,
			UserID:   message.UserID,
			Content:  []byte("You are not friends with this user."),
		})
		return
	}
	roomID := store.GenerateRoomID(message.UserID, message.RoomID)
	fmt.Printf("Private RoomID: %d\n", roomID)

	// Check if the room exists
	session.Room = room.GetRoom(roomID)
	session.Room.AddClient(message.UserID, session.TcpServer)

	session.Room.BroadcastMessage(*message)

	return
}

func ConnectHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("ConnectHandler", message.String())
	if session.Room == nil {
		session.Room = room.GetRoom(message.RoomID)
		session.Room.AddClient(message.UserID, session.TcpServer)
	}
	session.TcpServer.Send(&tcp.Message{
		UserName: "AlIM Server",
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Content:  []byte("Connect success!"),
	})
}

func AddFriendHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("AddFriendHandler", message.String())
	store.AddFriend(message.UserID, message.RoomID)
	session.TcpServer.Send(&tcp.Message{
		UserName: "AlIM Server",
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Content:  []byte("Friend added!"),
	})
}

func RoomChangeHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("AddFriendHandler", message.String())
	store.AddFriend(message.UserID, message.RoomID)
	session.TcpServer.Send(&tcp.Message{
		UserName: "AlIM Server",
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Content:  []byte("Friend added!"),
	})
}
