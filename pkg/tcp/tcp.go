package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type HandlerFunc func(ctx context.Context, conn net.Conn)

type TcpServer struct {
	MessageID     uint32
	Address       string
	tcpConnection net.Conn
	handlerFunc   map[uint32]HandlerFunc
}

func NewTcpServer(conn net.Conn) *TcpServer {
	return &TcpServer{
		tcpConnection: conn,
	}
}

func (ts *TcpServer) Send(message *Message) error {
	// Marshal the message
	content, err := message.Marshal()
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return err
	}

	contentLength := uint32(len(content))
	buf := new(bytes.Buffer)

	// Write the length to the buffer
	if err := binary.Write(buf, binary.BigEndian, contentLength); err != nil {
		fmt.Println("Error writing length to buffer:", err)
		return err
	}

	// Write the content to the buffer
	if err := binary.Write(buf, binary.BigEndian, content); err != nil {
		fmt.Println("Error writing content to buffer:", err)
		return err
	}

	// Send the buffer
	_, err = ts.tcpConnection.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}

	return nil
}

func (ts *TcpServer) Receive() (*Message, error) {
	// Create a buffer to hold the length
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(ts.tcpConnection, lenBuf)
	if err != nil {
		fmt.Println("Error reading length from server:", err)
		return nil, err
	}

	// Convert the length to an int
	length := binary.BigEndian.Uint32(lenBuf)

	// Create a buffer to hold the data
	data := make([]byte, length)
	_, err = io.ReadFull(ts.tcpConnection, data)
	if err != nil {
		fmt.Println("Error reading data from server:", err)
		return nil, err
	}

	// Unmarshal the data
	message := &Message{}
	err = message.Unmarshal(data)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil, err
	}

	fmt.Println("Received message:", message.String())
	return message, nil
}

//func (ts *TcpServer) Receive() (*Message, error) {
//	var data = make([]byte, 1024)
//	n, err := ts.tcpConnection.Read(data)
//	if err != nil {
//		fmt.Println("Error reading from server:", err)
//		return nil, err
//	}
//
//	message := &Message{}
//	err = message.Unmarshal(data[:n])
//	if err != nil {
//		fmt.Println("Error unmarshalling message:", err)
//		return nil, err
//	}
//
//	fmt.Println("Received message:", message.String())
//	return message, nil
//}

func (ts *TcpServer) Close() {
	ts.tcpConnection.Close()
}
