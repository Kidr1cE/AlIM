package main

import (
	"AlIM/pkg/tcp"
	"fmt"
	"net"
	"sync"
)

func main() {
	var userName string
	var id, userID int
	fmt.Scan(&userName, &id, &userID)
	fmt.Println("Set user name:", userName, "ID:", id, "UserID:", userID)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			message, err := Receive(conn)
			if err != nil {
				fmt.Println("Error receiving message:", err)
				return
			}
			fmt.Println("Received message:", message.String())
		}
	}()

	for {
		var content string
		fmt.Scanln(&content)
		message := &tcp.Message{
			UserName:  userName,
			MailBoxID: id,
			UserID:    userID,
			Type:      1,
			Content:   []byte(content),
		}

		err = Send(conn, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}

	wg.Wait()
}

func Send(conn net.Conn, message *tcp.Message) error {
	content, err := message.Marshal()
	fmt.Println("Content:", message.String())
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return err
	}
	_, err = conn.Write(content)
	if err != nil {
		fmt.Println("Error sending message", err)
		return err
	}
	return nil
}

func Receive(conn net.Conn) (*tcp.Message, error) {
	res := make([]byte, 1024)
	n, err := conn.Read(res)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return nil, err
	}

	resMessage := &tcp.Message{}
	err = resMessage.Unmarshal(res[:n])
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil, err
	}
	return resMessage, nil
}
