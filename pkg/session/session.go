package session

import (
	"AlIM/pkg/room"
	"AlIM/pkg/tcp"
	"context"
	"fmt"
	"io"
)

type MessageHandler func(session *Session, message *tcp.Message)

// global map of message types to handlers
var tcpHandler map[int]MessageHandler

type Session struct {
	ID        int
	Name      string
	TcpServer *tcp.TcpServer
	Room      *room.Room
	handlers  map[int]MessageHandler
}

func init() {
	tcpHandler = make(map[int]MessageHandler)

	Handle(tcp.ConnectMessage, ConnectHandler)
	Handle(tcp.AddFriendMessage, AddFriendHandler)
	Handle(tcp.RoomChangeMessage, RoomChangeHandler)
	Handle(tcp.SendMessage, SendMessageHandler)
	Handle(tcp.ListPublicRoomMessage, ListPublicRoomHandler)
	Handle(tcp.RecommendFriendMessage, RecommendFriendHandler)
}

func NewSession(tcpServer *tcp.TcpServer) *Session {
	session := &Session{
		TcpServer: tcpServer,
		handlers:  tcpHandler,
	}

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

func Handle(messageType int, handler MessageHandler) {
	tcpHandler[messageType] = handler
}
