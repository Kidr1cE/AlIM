package server

import (
	"AlIM/pkg/room"
	"AlIM/pkg/session"
	"AlIM/pkg/store"
	"AlIM/pkg/tcp"
	"fmt"
)

// ConnectHandler :Set session ID and Name
func ConnectHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("ConnectHandler", message.String())

	session.ID = message.UserID
	session.Name = message.UserName

	if room.GetRoom(message.RoomID) != nil {
		session.Room = room.GetRoom(message.RoomID)
		session.Room.AddClient(message.UserID, session.TcpServer)
	}

	err := session.TcpServer.Send(&tcp.Message{
		UserName: "AlIM Server",
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Content:  []byte("Connected!"),
	})
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func AddFriendHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("AddFriendHandler", message.String())

	if !store.IsFriend(message.UserID, message.RoomID) { // If they are not friends, add them as friends
		store.AddFriend(message.UserID, message.RoomID)
	}

	privateRoomID := store.GenerateRoomID(message.UserID, message.RoomID)
	privateRoom := room.GetPrivateRoom(privateRoomID)
	session.Room = privateRoom

	err := session.TcpServer.Send(&tcp.Message{
		UserName: "AlIM Server",
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Content:  []byte("Friend added!"),
	})
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

// RoomChangeHandler :Change the room of the session. Required Message.RoomID Message.RoomType
func RoomChangeHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("RoomChangeHandler", message.String())

	oldRoom := session.Room
	if oldRoom != nil {
		oldRoom.RemoveClient(session.ID)
	}

	switch message.RoomType {
	case room.PublicRoom:
		session.Room = room.GetRoom(message.RoomID)
	case room.PrivateRoom:
		// Check if they are friends
		if !store.IsFriend(session.ID, message.RoomID) {
			fmt.Println("Not friends")
			return
		}
		privateRoomID := store.GenerateRoomID(session.ID, message.RoomID)
		session.Room = room.GetPrivateRoom(privateRoomID)
	}
	session.Room.AddClient(session.ID, session.TcpServer)
}

// SendMessageHandler :Change the room of the session. Required Content
func SendMessageHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("SendMessageHandler", message.String())

	if session.Room == nil {
		fmt.Println("No room found for user")
		return
	}
	session.Room.BroadcastMessage(*message)
}
