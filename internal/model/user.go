package model

import (
	"time"
)

// Admin 管理员模型
type Admin struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

// APIToken API 令牌模型
type APIToken struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	Name       string     `gorm:"size:128;not null" json:"name"`
	TokenHash  string     `gorm:"size:255;uniqueIndex;not null" json:"-"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	AdminID    uint       `gorm:"index;not null" json:"admin_id"`
	Admin      Admin      `gorm:"foreignKey:AdminID" json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TableName 指定表名
func (APIToken) TableName() string {
	return "api_tokens"
}
