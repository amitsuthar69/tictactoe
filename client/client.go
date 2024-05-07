package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var serverAddress = "ws://localhost:8080/ws"

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(serverAddress, nil)
	if err != nil {
		log.Fatal("Error connecting to game server:", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Tic Tac Toe!")
	fmt.Println("Enter 'create' to create a new room or '<roomid>' to join an existing room:")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		err := conn.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			log.Println("Error sending message to server:", err)
			continue
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading response from server:", err)
			continue
		}

		log.Println(":", string(msg))

		if strings.TrimSpace(string(msg)) == "joined" {
			fmt.Println("you joined the room")
		}
		if strings.TrimSpace(string(msg)) == "failed" {
			fmt.Println("failed to join")
		}
	}
}
