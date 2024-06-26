package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	ConnectMessage         = 3
	AddFriendMessage       = 4
	RoomChangeMessage      = 5
	SendMessage            = 6
	ListPublicRoomMessage  = 7
	RecommendFriendMessage = 8
)

// Message header: 32+100 132 bytes; body: length
type Message struct {
	UserID   int
	RoomID   int
	RoomType int
	Type     int
	Gap      [4]byte // 4bytes gap
	UserName string  // 25words max 100 bytes
	Content  []byte  // length
}

func (m *Message) String() string {
	return fmt.Sprintf("UserID: %d, RoomID: %d, RoomType: %d, Type: %d, UserName: %s, Content: %s", m.UserID, m.RoomID, m.RoomType, m.Type, m.UserName, m.Content)
}

func (m *Message) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write UserID
	if err := binary.Write(buf, binary.BigEndian, int32(m.UserID)); err != nil {
		return nil, fmt.Errorf("failed to write UserID: %v", err)
	}

	// Write RoomID
	if err := binary.Write(buf, binary.BigEndian, int32(m.RoomID)); err != nil {
		return nil, fmt.Errorf("failed to write RoomID: %v", err)
	}

	// Write RoomType
	if err := binary.Write(buf, binary.BigEndian, int32(m.RoomType)); err != nil {
		return nil, fmt.Errorf("failed to write RoomType: %v", err)
	}

	// Write Type
	if err := binary.Write(buf, binary.BigEndian, int32(m.Type)); err != nil {
		return nil, fmt.Errorf("failed to write Type: %v", err)
	}

	// Write Gap
	if err := binary.Write(buf, binary.BigEndian, m.Gap); err != nil {
		return nil, fmt.Errorf("failed to write Gap: %v", err)
	}

	// Write UserName
	userNameBytes := make([]byte, 25) // Create a byte array of length 25
	copy(userNameBytes, m.UserName)   // Copy the UserName into the byte array
	if err := binary.Write(buf, binary.BigEndian, userNameBytes); err != nil {
		return nil, fmt.Errorf("failed to write UserName: %v", err)
	}

	// Write Content
	if err := binary.Write(buf, binary.BigEndian, m.Content); err != nil {
		return nil, fmt.Errorf("failed to write Content: %v", err)
	}

	return buf.Bytes(), nil
}

func (m *Message) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	// Read UserID
	var userID int32
	if err := binary.Read(buf, binary.BigEndian, &userID); err != nil {
		return fmt.Errorf("failed to read UserID: %v", err)
	}
	m.UserID = int(userID)

	// Read RoomID
	var roomID int32
	if err := binary.Read(buf, binary.BigEndian, &roomID); err != nil {
		return fmt.Errorf("failed to read RoomID: %v", err)
	}
	m.RoomID = int(roomID)

	// Read RoomType
	var roomType int32
	if err := binary.Read(buf, binary.BigEndian, &roomType); err != nil {
		return fmt.Errorf("failed to read RoomType: %v", err)
	}
	m.RoomType = int(roomType)

	// Read Type
	var msgType int32
	if err := binary.Read(buf, binary.BigEndian, &msgType); err != nil {
		return fmt.Errorf("failed to read Type: %v", err)
	}
	m.Type = int(msgType)

	// Read Gap
	if _, err := buf.Read(m.Gap[:]); err != nil {
		return fmt.Errorf("failed to read Gap: %v", err)
	}

	// Read UserName
	userNameBytes := make([]byte, 25) // Assuming UserName is always 25 bytes
	if _, err := buf.Read(userNameBytes); err != nil {
		return fmt.Errorf("failed to read UserName: %v", err)
	}
	m.UserName = string(bytes.Trim(userNameBytes, "\x00")) // Remove trailing zero bytes

	// Read Content
	contentBytes := make([]byte, buf.Len())
	if _, err := buf.Read(contentBytes); err != nil {
		return fmt.Errorf("failed to read Content: %v", err)
	}
	m.Content = contentBytes

	return nil
}
