package session

import (
	"AlIM/pkg/room"
	"AlIM/pkg/tcp"
	"fmt"
)

type MessageHandler func(session *Session, message *tcp.Message)

type Session struct {
	ID        int
	Name      string
	TcpServer *tcp.TcpServer
	Room      *room.Room
	handlers  map[int]MessageHandler
}

func NewSession(tcpServer *tcp.TcpServer) *Session {
	return &Session{
		TcpServer: tcpServer,
		handlers:  make(map[int]MessageHandler),
	}
}

func (s *Session) Start() {
	defer func() {
		if s.Room != nil {
			s.Room.BroadcastMessage(tcp.Message{
				UserName: "AlIM Server",
				UserID:   0,
				Type:     0,
				Content:  []byte(fmt.Sprintf("%s#%d has left the room", s.Name, s.ID)),
			})
			s.TcpServer.Close()
		}
	}()

	for {
		message, err := s.TcpServer.Receive()
		if err != nil {
			fmt.Println("Error receiving message:", err)
			break
		}

		handler, ok := s.handlers[message.Type]
		if !ok {
			fmt.Println("Handler not found")
			continue
		}
		handler(s, message)
	}
}

func (s *Session) Handle(messageType int, handler MessageHandler) {
	s.handlers[messageType] = handler
}
