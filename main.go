package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Room struct {
	ID          string
	Player1Conn *websocket.Conn
	Player2Conn *websocket.Conn
}

var rooms map[string]*Room = make(map[string]*Room)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	if string(msg) == "create" {
		HandleCreate(conn, string(msg))
	} else {
		HandleJoin(conn, string(msg))
	}
}

func HandleCreate(conn *websocket.Conn, msg string) {
	roomID := generateRoomID()
	rooms[roomID] = &Room{ID: roomID, Player1Conn: conn}
	err := conn.WriteMessage(websocket.TextMessage, []byte(roomID))
	if err != nil {
		log.Println(err)
	}
}

func HandleJoin(conn *websocket.Conn, msg string) {
	room, ok := rooms[msg]
	if !ok || room.Player2Conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, []byte("failed"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	room.Player2Conn = conn
	err := conn.WriteMessage(websocket.TextMessage, []byte("joined"))
	if err != nil {
		log.Println(err)
	}
	if room.Player1Conn != nil && room.Player2Conn != nil {
		startGame(room)
	}
}

func generateRoomID() string {
	return randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func startGame(room *Room) {
	fmt.Println("game started!", room)
	fmt.Println(room.Player1Conn.LocalAddr().String(), "joined")
	fmt.Println(room.Player2Conn.LocalAddr().String(), "joined")
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
