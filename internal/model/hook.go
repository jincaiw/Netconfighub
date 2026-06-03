package model

import "time"

type Hook struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Type      string    `gorm:"size:32;not null" json:"type"`
	Config    string    `gorm:"size:4096;not null" json:"config"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	Events    string    `gorm:"size:256" json:"events"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Hook) TableName() string {
	return "hooks"
}
