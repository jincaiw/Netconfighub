package model

import (
	"time"
)

type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"
	BackupStatusRunning   BackupStatus = "running"
	BackupStatusSuccess   BackupStatus = "success"
	BackupStatusFailed    BackupStatus = "failed"
	BackupStatusRetrying  BackupStatus = "retrying"
	BackupStatusTimeout   BackupStatus = "timeout"
	BackupStatusCancelled BackupStatus = "cancelled"
	BackupStatusSkipped   BackupStatus = "skipped"
)

type BackupTask struct {
	ID            uint         `gorm:"primarykey" json:"id"`
	DeviceID      uint         `gorm:"index;not null" json:"device_id"`
	Device        Device       `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	Status        BackupStatus `gorm:"size:32;not null;default:pending" json:"status"`
	StartedAt     *time.Time   `json:"started_at"`
	FinishedAt    *time.Time   `json:"finished_at"`
	FailureReason string       `gorm:"size:512" json:"failure_reason"`
	ErrorType     string       `gorm:"size:64" json:"error_type"`
	RetryCount    int          `gorm:"default:0" json:"retry_count"`
	MaxRetries    int          `gorm:"default:3" json:"max_retries"`
	ConfigHash    string       `gorm:"size:64" json:"config_hash"`
	CreatedAt     time.Time    `json:"created_at"`
}

func (BackupTask) TableName() string {
	return "backup_tasks"
}
