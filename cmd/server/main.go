package main

import (
	"AlIM/pkg/server"
)

func main() {
	mailServer := server.NewMailServer(":8080")
	mailServer.Start()
}
