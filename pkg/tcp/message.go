package tcp

import (
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	MailBoxID int
	Type      int
	UserID    int
	UserName  string
	Content   []byte
}

func (m *Message) String() string {
	return fmt.Sprintf("ID: %d, Type: %d, UserID: %d, UserName: %s, Content: %s", m.MailBoxID, m.Type, m.UserID, m.UserName, m.Content)
}

func (m *Message) Marshal() ([]byte, error) {
	sep := "\\sep"
	messageBytes := []byte(fmt.Sprintf("%d%s%d%s%d%s%s%s%s", m.MailBoxID, sep, m.Type, sep, m.UserID, sep, m.UserName, sep, m.Content))
	return messageBytes, nil
}

func (m *Message) Unmarshal(content []byte) error {
	parts := strings.Split(string(content), "\\sep")
	if len(parts) != 5 {
		return fmt.Errorf("error unmarshalling message")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	m.MailBoxID = id

	messageType, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	m.Type = messageType

	userID, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	m.UserID = userID

	m.UserName = parts[3]
	m.Content = []byte(parts[4])

	return nil
}
