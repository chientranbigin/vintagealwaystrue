package models

import "time"

type CopilotConversation struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	PSID         string     `gorm:"uniqueIndex;size:128" json:"psid"`
	CustomerName string     `gorm:"size:255" json:"customer_name"`
	AvatarURL    string     `gorm:"type:text" json:"avatar_url"`
	LastOrderAt  *time.Time `json:"last_order_at"`
	LastOrderID  *uint      `json:"last_order_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CopilotMessageRecord struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID uint      `gorm:"index" json:"conversation_id"`
	FBMessageID    string    `gorm:"uniqueIndex;size:191" json:"fb_message_id"`
	Sender         string    `gorm:"size:32" json:"sender"` // "customer" or "shop"
	Type           string    `gorm:"size:32" json:"type"`   // "text" or "image"
	Text           string    `gorm:"type:longtext" json:"text"`
	ImageURL       string    `gorm:"type:text" json:"image_url"`
	OCRText        string    `gorm:"type:longtext" json:"ocr_text"`
	MsgTimestamp   time.Time `gorm:"index" json:"msg_timestamp"`
	CreatedAt      time.Time `json:"created_at"`
}
