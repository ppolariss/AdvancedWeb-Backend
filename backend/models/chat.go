package models

import (
	"gorm.io/gorm"
	"time"
)

type Chat struct {
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index:idx_chat_user_deleted,priority:2"`
	//Records   Records        `json:"records,omitempty"`
	//Count             int            `json:"count"` // Record 条数
	//MaxLengthExceeded bool           `json:"max_length_exceeded"`
}

type Chats []Chat

type Record struct {
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index:idx_record_chat_deleted,priority:2"`
	UserID    int            `json:"user_id" gorm:"index:idx_chat_user_deleted,priority:1"`
	Name      string         `json:"name"`
	ChatID    int            `json:"chat_id" gorm:"index:idx_record_chat_deleted,priority:1"`
	RoomID    string         `json:"room_id" gorm:"index:idx_room_id,priority:1"`
	Type      string         `json:"type"`
	ToID      int            `json:"to_id" gorm:"index:idx_to_id,priority:1"`
	Message   string
	//Duration           float64        `json:"duration"` // 处理时间，单位 s
	//Request            string         `json:"request"`
	//Response           string         `json:"response"`
	//Prefix             string         `json:"-"`
	//RawContent         string         `json:"raw_content"`
	//ExtraData          any            `json:"-" gorm:"serializer:json"` //`json:"extra_data" gorm:"serializer:json"`
	//ProcessedExtraData any            `json:"processed_extra_data" gorm:"serializer:json"`
	//LikeData           int            `json:"like_data"` // 1 like, -1 dislike
	//Feedback           string         `json:"feedback"`
	//RequestSensitive   bool           `json:"request_sensitive"`
	//ResponseSensitive  bool           `json:"response_sensitive"`
	//InnerThoughts      string         `json:"inner_thoughts"`
}

type Records []Record

func GetAllChatsByUserID(userID int) (chats []Chat, err error) {
	err = DB.Transaction(func(tx *gorm.DB) error {
		return tx.
			Model(&Chat{}).
			Preload("Records").
			Where("user_id=?", userID).
			Find(&chats).
			Error
	})
	return
}
