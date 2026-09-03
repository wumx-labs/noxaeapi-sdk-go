package main

import (
	"fmt"
	"log"

	"github.com/NoxlyDev/noxaeapi-sdk-go"
	"github.com/NoxlyDev/noxaeapi-sdk-go/ws"
)

func main() {
	client, err := noxaeapi.NewClient(noxaeapi.ClientOptions{
		BaseURL: "http://localhost:8080",
		APIKey:  "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	socket, err := ws.Connect(ws.Options{
		BaseURL: client.BaseURL(),
		APIKey:  client.APIKey(),
		Route:   "console",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	for event := range socket.Events {
		switch event.Type {
		case ws.EventOpen:
			fmt.Println("connected")
		case ws.EventConsole:
			fmt.Println("[console]", string(event.Payload))
		case ws.EventError:
			fmt.Println("error:", event.Err)
		case ws.EventClose:
			fmt.Println("disconnected")
		}
	}
}
