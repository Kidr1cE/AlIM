package session

import (
	"AlIM/pkg/room"
	"AlIM/pkg/tcp"
	"fmt"
)

type MessageHandler func(session *Session, message *tcp.Message)

const (
	GroupMessage = iota + 1
	PrivateMessage
)

type Session struct {
	ID        string
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
		s.TcpServer.Close()
	}()

	connectMessage, err := s.TcpServer.Receive()
	if err != nil {
		return
	}

	fmt.Println(connectMessage.String())

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
