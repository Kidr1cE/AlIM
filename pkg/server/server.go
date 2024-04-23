package server

import (
	"AlIM/pkg/session"
	"AlIM/pkg/tcp"
	"context"
	"fmt"
	"net"
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

func (ms *MailServer) Start() {
	conn, err := net.Listen("tcp", ms.Address)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		conn, err := conn.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// TODO: limit the number of connections
		connectNum++

		tcpServer := tcp.NewTcpServer(conn)
		newSession := session.NewSession(tcpServer)

		go func() {
			newSession.Start(ctx)
		}()
	}
}
