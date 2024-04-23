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
	History map[int][]tcp.Message
	clients map[int]*tcp.TcpServer
	mu      sync.RWMutex
}

func NewRoom(roomID int) *Room {
	return &Room{
		ID:      roomID,
		History: make(map[int][]tcp.Message),
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

func (r *Room) AddClient(id int, conn *tcp.TcpServer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.UserNum++

	r.clients[id] = conn
}

func (r *Room) RemoveClient(id int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.UserNum--

	delete(r.clients, id)
}

func (r *Room) BroadcastMessage(message tcp.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, conn := range r.clients {
		err := conn.Send(&message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			// TODO remove client
			delete(r.clients, message.UserID)
		}
	}
}

func (r *Room) SendToUser(userID int, message tcp.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if conn, ok := r.clients[userID]; ok {
		err := conn.Send(&message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			// TODO remove client
			delete(r.clients, message.UserID)
		}
	}
}
