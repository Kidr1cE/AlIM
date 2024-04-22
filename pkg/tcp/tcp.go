package tcp

import (
	"context"
	"fmt"
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
	content, err := message.Marshal()
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return err
	}
	_, err = ts.tcpConnection.Write(content)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}

	// TODO: handle error
	return nil
}

func (ts *TcpServer) Receive() (*Message, error) {
	var data = make([]byte, 1024)
	n, err := ts.tcpConnection.Read(data)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return nil, err
	}

	message := &Message{}
	err = message.Unmarshal(data[:n])
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil, err
	}

	fmt.Println("Received message:", message.String())
	return message, nil
}

func (ts *TcpServer) Close() {
	ts.tcpConnection.Close()
}
