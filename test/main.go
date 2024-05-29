package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrade   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mu sync.Mutex
)

// Message 定义消息结构
type Message struct {
	Data string `json:"data"`
}

//func main() {
//	http.HandleFunc("/video-chat.js", serveJs)
//	http.HandleFunc("/", serveHtml)
//	http.HandleFunc("/ws", handleConnections)
//
//	go handleMessages()
//
//	log.Println("Server is running on http://localhost:3000")
//	err := http.ListenAndServe(":3000", nil)
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "video-chat.js"))
}

func serveHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "index.html"))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	log.Println("connection")

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
		log.Println("message")
		log.Printf("Received message => %s", msg.Data)
		log.Printf("Clients count => %d", len(clients))
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				_ = client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}
