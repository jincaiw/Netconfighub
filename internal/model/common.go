package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 公共基础模型，包含 ID、CreatedAt、UpdatedAt、DeletedAt
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
