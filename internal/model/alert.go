package model

import (
	"time"
)

type Alert struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Type       string    `gorm:"size:32;not null;index" json:"type"`
	DeviceID   *uint     `gorm:"index" json:"device_id"`
	DeviceName string    `gorm:"size:128" json:"device_name"`
	Title      string    `gorm:"size:256;not null" json:"title"`
	Message    string    `gorm:"size:1024" json:"message"`
	Severity   string    `gorm:"size:16;not null;default:warning" json:"severity"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Alert) TableName() string {
	return "alerts"
}
