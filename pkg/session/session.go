package session

import (
	"AlIM/pkg/mailbox"
	"AlIM/pkg/tcp"
	"fmt"
)

type Session struct {
	ID        string
	TcpServer *tcp.TcpServer
	mailBox   *mailbox.Mailbox
}

func NewSession(tcpServer *tcp.TcpServer) *Session {
	return &Session{
		TcpServer: tcpServer,
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
	mailBoxID := connectMessage.MailBoxID
	s.mailBox = mailbox.GetMailbox(mailBoxID)
	s.mailBox.AddClient(connectMessage.UserID, s.TcpServer)
	s.mailBox.BroadcastMessage(tcp.Message{
		MailBoxID: mailBoxID,
		Type:      1,
		UserID:    connectMessage.UserID,
		UserName:  connectMessage.UserName,
		Content:   []byte("User connected"),
	})

	for {
		message, err := s.TcpServer.Receive()
		if err != nil {
			fmt.Println("Error receiving message:", err)
			break
		}
		s.mailBox.BroadcastMessage(*message)
	}
}
