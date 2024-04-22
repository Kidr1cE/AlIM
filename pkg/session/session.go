package session

import (
	"AlIM/pkg/mailbox"
	"AlIM/pkg/tcp"
	"fmt"
)

type Session struct {
	ID             string
	TcpServer      *tcp.TcpServer
	mailBox        *mailbox.Mailbox
	privateMailBox *mailbox.Mailbox
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
	// group chat
	s.mailBox = mailbox.GetMailbox(mailBoxID)
	s.mailBox.AddClient(connectMessage.UserID, s.TcpServer)
	s.mailBox.BroadcastMessage(tcp.Message{
		MailBoxID: mailBoxID,
		Type:      1,
		UserID:    connectMessage.UserID,
		UserName:  connectMessage.UserName,
		Content:   []byte("User connected"),
	})
	// private chat
	s.privateMailBox = mailbox.SetPrivateMailbox(connectMessage.UserID, s.TcpServer)

	for {
		message, err := s.TcpServer.Receive()
		if err != nil {
			fmt.Println("Error receiving message:", err)
			break
		}
		switch message.Type {
		case tcp.GroupMessage:
			s.mailBox.BroadcastMessage(*message)
		case tcp.PrivateMessage:
			friendMailBoxID := message.MailBoxID
			friendMailBox := mailbox.GetPrivateMailbox(friendMailBoxID)
			if friendMailBox == nil {
				fmt.Println("Friend mailbox not found")
				break
			}
			friendMailBox.BroadcastMessage(*message)
		}
	}
}
