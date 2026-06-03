package model

import (
	"time"
)

type AuditLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Username   string    `gorm:"size:64" json:"username"`
	Action     string    `gorm:"size:64;not null" json:"action"`
	TargetType string    `gorm:"size:64;not null" json:"target_type"`
	TargetID   uint      `gorm:"index" json:"target_id"`
	Detail     string    `gorm:"size:1024" json:"detail"`
	ClientIP   string    `gorm:"size:64" json:"client_ip"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
