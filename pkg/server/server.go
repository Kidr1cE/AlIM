package server

import (
	"AlIM/pkg/session"
	"AlIM/pkg/tcp"
	"context"
	"fmt"
	"net"
)

const (
	GroupMessage = iota + 1
	PrivateMessage
	ConnectMessage
	AddFriendMessage
	RoomChangeMessage
	SendMessage
)

var connectNum int

type HandlerFunc func(ctx context.Context, conn net.Conn)

type MailServer struct {
	Address string
	Handler map[int]HandlerFunc
}

func NewMailServer(address string) *MailServer {
	return &MailServer{
		Address: address,
	}
}

func InitSession(tcpServer *tcp.TcpServer) *session.Session {
	newSession := session.NewSession(tcpServer)
	newSession.Handle(ConnectMessage, ConnectHandler)
	newSession.Handle(AddFriendMessage, AddFriendHandler)
	newSession.Handle(RoomChangeMessage, RoomChangeHandler)
	newSession.Handle(SendMessage, SendMessageHandler)

	return newSession
}

func (ms *MailServer) Start() {
	conn, err := net.Listen("tcp", ms.Address)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// TODO: limit the number of connections
		connectNum++

		tcpServer := tcp.NewTcpServer(conn)
		newSession := InitSession(tcpServer)

		go func() {
			newSession.Start()
		}()
	}
}
