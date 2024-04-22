package server

import (
	"AlIM/pkg/room"
	"AlIM/pkg/session"
	"AlIM/pkg/tcp"
)

func GroupHandler(session *session.Session, message *tcp.Message) {
	if session.Room == nil {
		// group chat
		session.Room = room.GetMailbox(mailBoxID)
		session.Room.AddClient(connectMessage.UserID, s.TcpServer)
	}
	session.Room.BroadcastMessage(*message)
}
