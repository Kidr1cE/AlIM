package room

import (
	"AlIM/pkg/tcp"
	"fmt"
	"sync"
)

const (
	PublicRoom = iota + 1
	PrivateRoom
)

// PublicRooms ID is the group ID
var PublicRooms = make(map[int]*Room)

// PrivateRooms ID is the user ID
var PrivateRooms = make(map[int]*Room)

type Room struct {
	ID      int
	UserNum int
	clients map[int]*tcp.TcpServer
	mu      sync.RWMutex
}

func NewRoom(roomID int) *Room {
	return &Room{
		ID:      roomID,
		clients: make(map[int]*tcp.TcpServer),
	}
}

func GetRoom(roomID int) *Room {
	// TODO get from cache
	if room, ok := PublicRooms[roomID]; ok {
		fmt.Printf("Room %d exists\n", roomID)
		return room
	} else {
		fmt.Printf("Room %d does not exist\n", roomID)
		room := NewRoom(roomID)
		PublicRooms[roomID] = room
		return room
	}
}

func GetPrivateRoom(roomID int) *Room {
	// TODO get from cache
	if room, ok := PrivateRooms[roomID]; ok {
		fmt.Printf("Room %d exists\n", roomID)
		return room
	} else {
		fmt.Printf("Room %d does not exist\n", roomID)
		room := NewRoom(roomID)
		PrivateRooms[roomID] = room
		return room
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
