package room

import (
	"AlIM/pkg/tcp"
	"fmt"
	"sync"
)

// Mailboxes ID is the group ID
var Rooms = make(map[int]*Room)

// PrivateMailboxes ID is the user ID
var PrivateMailboxes = make(map[int]*Room)

type Room struct {
	UserNum int
	clients map[int]*tcp.TcpServer
	mu      sync.RWMutex
}

func NewMailbox() *Room {
	return &Room{
		clients: make(map[int]*tcp.TcpServer),
	}
}

func GetMailbox(roomID int) *Room {
	// TODO get from cache
	if room, ok := Rooms[roomID]; ok {
		fmt.Printf("Room %d exists\n", roomID)
		return room
	} else {
		fmt.Printf("Room %d does not exist\n", roomID)
		room := NewMailbox()
		Rooms[roomID] = room
		return room
	}
}

func SetPrivateMailbox(userID int, conn *tcp.TcpServer) *Room {
	// TODO get from cache
	if room, ok := PrivateMailboxes[userID]; ok {
		fmt.Printf("Private room %d exists\n", userID)
		room.AddClient(userID, conn)
	} else {
		fmt.Printf("Private room %d does not exist\n", userID)
		room := NewMailbox()
		room.AddClient(userID, conn)
		PrivateMailboxes[userID] = room
	}
	return PrivateMailboxes[userID]
}

func GetPrivateMailbox(userID int) *Room {
	// TODO get from cache
	if room, ok := PrivateMailboxes[userID]; ok {
		fmt.Printf("Private room %d exists\n", userID)
		return room
	} else {
		fmt.Printf("Private room %d does not exist\n", userID)
		return nil
	}
}

func (m *Room) AddClient(id int, conn *tcp.TcpServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.UserNum++

	m.clients[id] = conn
}

func (m *Room) RemoveClient(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.UserNum--

	delete(m.clients, id)
}

func (m *Room) BroadcastMessage(message tcp.Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, conn := range m.clients {
		err := conn.Send(&message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			// TODO remove client
			delete(m.clients, message.UserID)
		}
	}
}
