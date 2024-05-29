package message

import (
	"github.com/gofiber/websocket/v2"
	"sync"
	"time"
)

type AIResponse struct {
	Status     int    `json:"status"` // 1 for output, 0 for end, -1 for error
	StatusCode int    `json:"status_code,omitempty"`
	Output     string `json:"output,omitempty"`
	Stage      string `json:"stage,omitempty"`
}

type MossResponse struct {
	Status int    `json:"status"`
	Output string `json:"output"`
	Stage  string `json:"stage"`
}

type AddRecordsRequest struct {
	//ChatID   int
	UserID    int       `json:"user_id"`
	RoomID    string    `json:"room_id"`
	Type      string    `json:"type"`
	ToID      int       `json:"to_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type VideoMessage struct {
	C       *websocket.Conn
	Message []byte
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan VideoMessage)
	mu        sync.Mutex
)

//type VideoMessage struct {
//	Data string `json:"data"`
//}
