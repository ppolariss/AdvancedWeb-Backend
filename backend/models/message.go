package models

import "gorm.io/gorm"

type Message struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	User            *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID          int    `json:"user_id" gorm:"not null;index"`
	CommodityItemID int    `json:"commodity_item_id"` //如果不constraint join时要注意
	CreateAt        MyTime `json:"create_at" gorm:"autoCreateTime"`
	//PriceLimit   float64
}

type ByCreatedAt []Message

func (a ByCreatedAt) Len() int           { return len(a) }
func (a ByCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedAt) Less(i, j int) bool { return a[i].CreateAt.Time.Before(a[j].CreateAt.Time) }

func GetMessagesByUserID(userID int) (messages []Message, err error) {
	err = DB.Transaction(func(tx *gorm.DB) error {
		return tx.
			Preload("CommodityItem").
			Preload("CommodityItem.Commodity").
			Preload("CommodityItem.Platform").
			Preload("CommodityItem.Seller").
			Where("user_id=?", userID).
			Find(&messages).
			Error
	})
	return
}
