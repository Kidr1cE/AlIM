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

	// Check if the message is received from a private room
	if message.RoomType == room.PrivateRoom {
		privateRoomID := store.GenerateRoomID(session.ID, message.RoomID)
		session.Room = room.GetPrivateRoom(privateRoomID)
	}
}

// ListPublicRoomHandler
func ListPublicRoomHandler(session *session.Session, message *tcp.Message) {
	fmt.Println("ListPublicRoomHandler", message.String())

	// Send the list of public rooms
	for roomID := range room.PublicRooms {
		err := session.TcpServer.Send(&tcp.Message{
			UserName: "AlIM Server",
			RoomID:   roomID,
			UserID:   0,
			Content:  []byte(fmt.Sprintf("Room %d", roomID)),
		})
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}

// RecommendFriendHandler :Send friend recommendations
func RecommendFriendHandler(session *session.Session, message *tcp.Message) {
	ownFriends := store.GetFriends(message.UserID)

	unaddedFriends := make(map[int]struct{})
	for _, friendID := range ownFriends {
		for _, friendID2 := range store.GetFriends(friendID) {
			if _, ok := unaddedFriends[friendID2]; !ok && friendID2 != message.UserID {
				unaddedFriends[friendID2] = struct{}{}
			}
		}
	}

	friendList := make([]int, 0, len(unaddedFriends))
	for friendID := range unaddedFriends {
		friendList = append(friendList, friendID)
	}

	recommandContent := fmt.Sprintf("Recommend friends: %v", friendList)
	err := session.TcpServer.Send(&tcp.Message{
		UserName: "AlIM Server",
		RoomID:   message.RoomID,
		UserID:   message.UserID,
		Content:  []byte(recommandContent),
	})
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}
