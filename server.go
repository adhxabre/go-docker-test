package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	gubrak "github.com/novalagung/gubrak/v2"
)

type M map[string]interface{}

const (
	MESSAGE_NEW_USER = "New User"
	MESSAGE_CHAT     = "Chat"
	MESSAGE_LEAVE    = "Leave"
)

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections = make([]*WebSocketConnection, 0)
)

func main() {
	http.HandleFunc("/ws", chatWebSocket)

	fmt.Println("Server starting at :8080")
	http.ListenAndServe(":8080", nil)
}

func chatWebSocket(w http.ResponseWriter, r *http.Request) {
	currentGorillaConn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
	connections = append(connections, &currentConn)

	go handleIO(&currentConn, connections)
}

func handleIO(currentConn *WebSocketConnection, _ []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[ERR] Error: %v\n", r)
			return
		}
	}()

	broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				ejectConnection(currentConn)
				return
			}

			fmt.Printf("[ERR] Error: %v", err.Error())
			continue
		}

		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}
}

func ejectConnection(currentConn *WebSocketConnection) {
	filtered := gubrak.From(connections).Reject(func(each *WebSocketConnection) bool {
		return each == currentConn
	}).Result()
	connections = filtered.([]*WebSocketConnection)
}

func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}
