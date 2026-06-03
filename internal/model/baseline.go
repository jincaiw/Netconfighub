package model

import (
	"time"
)

// BaselineScope 基线作用域枚举
type BaselineScope string

const (
	BaselineScopeDevice BaselineScope = "device"
	BaselineScopeGroup  BaselineScope = "group"
)

// Baseline 基线模型
type Baseline struct {
	ID        uint          `gorm:"primarykey" json:"id"`
	DeviceID  *uint         `gorm:"index" json:"device_id"`
	Device    *Device       `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	GroupID   *uint         `gorm:"index" json:"group_id"`
	Group     *Group        `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Scope     BaselineScope `gorm:"size:16;not null" json:"scope"`
	Content   string        `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// TableName 指定表名
func (Baseline) TableName() string {
	return "baselines"
}

// Deviation 偏差模型
type Deviation struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	BackupTaskID uint      `gorm:"index;not null" json:"backup_task_id"`
	BackupTask   BackupTask `gorm:"foreignKey:BackupTaskID" json:"backup_task,omitempty"`
	BaselineID   uint      `gorm:"index;not null" json:"baseline_id"`
	Baseline     Baseline  `gorm:"foreignKey:BaselineID" json:"baseline,omitempty"`
	DiffContent  string    `gorm:"type:text;not null" json:"diff_content"`
	DetectedAt   time.Time `json:"detected_at"`
}

// TableName 指定表名
func (Deviation) TableName() string {
	return "deviations"
}
