package schemas

import "time"

type AddRecordRequest struct {
	//ChatID   int
	UserID    int       `json:"user_id"`
	RoomID    string    `json:"room_id"`
	Type      string    `json:"type"`
	ToID      int       `json:"to_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type AddChatRequest struct {
	UserID   int
	RoomID   string
	CreateAt time.Time
}

type AddChatResponse struct {
	ChatID int
}

type AddRecordResponse struct {
}

//type AddChatAndRecord struct {
//}
//
//type AddChatAndRecord struct {
//}
