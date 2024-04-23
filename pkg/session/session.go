package session

import (
	"AlIM/pkg/room"
	"AlIM/pkg/tcp"
	"context"
	"fmt"
	"io"
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
	session := &Session{
		TcpServer: tcpServer,
		handlers:  make(map[int]MessageHandler),
	}

	session.Handle(tcp.ConnectMessage, ConnectHandler)
	session.Handle(tcp.AddFriendMessage, AddFriendHandler)
	session.Handle(tcp.RoomChangeMessage, RoomChangeHandler)
	session.Handle(tcp.SendMessage, SendMessageHandler)
	session.Handle(tcp.ListPublicRoomMessage, ListPublicRoomHandler)
	session.Handle(tcp.RecommendFriendMessage, RecommendFriendHandler)

	return session
}

func (s *Session) Start(ctx context.Context) {
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
		select {
		case <-ctx.Done():
			return
		default:
			message, err := s.TcpServer.Receive()
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Println("Error receiving message:", err)
				continue
			}

			handler, ok := s.handlers[message.Type]
			if !ok {
				fmt.Println("Handler not found")
				continue
			}
			handler(s, message)
		}
	}
}

func (s *Session) Handle(messageType int, handler MessageHandler) {
	s.handlers[messageType] = handler
}
