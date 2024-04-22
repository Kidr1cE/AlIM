package mailbox

import (
	"AlIM/pkg/tcp"
	"fmt"
	"sync"
)

var Mailboxes = make(map[int]*Mailbox)

type Mailbox struct {
	clients map[int]*tcp.TcpServer
	mu      sync.RWMutex
}

func NewMailbox() *Mailbox {
	return &Mailbox{
		clients: make(map[int]*tcp.TcpServer),
	}
}

func GetMailbox(mailboxID int) *Mailbox {
	// TODO get from cache
	if mailbox, ok := Mailboxes[mailboxID]; ok {
		fmt.Printf("Mailbox %d exists\n", mailboxID)
		return mailbox
	} else {
		fmt.Printf("Mailbox %d does not exist\n", mailboxID)
		mailbox := NewMailbox()
		Mailboxes[mailboxID] = mailbox
		return mailbox
	}
}

func (m *Mailbox) AddClient(id int, conn *tcp.TcpServer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[id] = conn
}

func (m *Mailbox) RemoveClient(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, id)
}

func (m *Mailbox) BroadcastMessage(message tcp.Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, conn := range m.clients {
		err := conn.Send(&message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
