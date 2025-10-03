package model

import "time"

type LoggingAPI struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URI       string    `json:"uri"`
	Method    string    `json:"method"`
	Request   string    `json:"request"`
	Response  string    `json:"response"`
	IP        string    `json:"ip"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (LoggingAPI) TableName() string {
	return "logging_api"
}