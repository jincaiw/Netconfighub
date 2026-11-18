package model

import (
	"time"
)

// SFTP upload stages. Reception and archiving are tracked separately so that
// "transfer succeeded" (stage=received) never implies "backup succeeded"
// (stage=archived, plus the related backup task status).
const (
	SFTPUploadStageReceived  = "received"
	SFTPUploadStageArchiving = "archiving"
	SFTPUploadStageArchived  = "archived"
	SFTPUploadStageFailed    = "failed"
	SFTPUploadStageRejected  = "rejected"
)

// SFTP upload error types persisted in sftp_uploads.error_type.
const (
	SFTPUploadErrorGitCommit = "git_commit"
	SFTPUploadErrorDB        = "db"
	SFTPUploadErrorContent   = "content"
)

// SFTPUpload records one configuration file pushed by a device over SFTP.
type SFTPUpload struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	DeviceID       uint       `gorm:"index;not null;index:idx_sftp_uploads_device_created,priority:1" json:"device_id"`
	AccountID      uint       `gorm:"index" json:"account_id"`
	BackupTaskID   *uint      `gorm:"index" json:"backup_task_id"`
	Stage          string     `gorm:"size:16;not null;index:idx_sftp_uploads_stage" json:"stage"`
	ReceivedAt     time.Time  `json:"received_at"`
	ReceivedBytes  int64      `json:"received_bytes"`
	ReceivedHash   string     `gorm:"size:64" json:"received_hash"`
	ArchivedAt     *time.Time `json:"archived_at"`
	ConfigHash     string     `gorm:"size:64" json:"config_hash"`
	StoredPath     string     `gorm:"size:512" json:"stored_path"`
	ClientFilename string     `gorm:"size:255" json:"client_filename"`
	RemoteAddr     string     `gorm:"size:64" json:"remote_addr"`
	ErrorType      string     `gorm:"size:64" json:"error_type"`
	FailureReason  string     `gorm:"size:512" json:"failure_reason"`
	RetryCount     int        `gorm:"default:0" json:"retry_count"`
	CreatedAt      time.Time  `gorm:"index:idx_sftp_uploads_device_created,priority:2" json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (SFTPUpload) TableName() string {
	return "sftp_uploads"
}

// IsArchived reports whether the upload's content reached the Git store.
func (u *SFTPUpload) IsArchived() bool {
	return u.Stage == SFTPUploadStageArchived
}

// CanRetry reports whether the upload may be re-queued for archiving.
func (u *SFTPUpload) CanRetry() bool {
	return u.Stage == SFTPUploadStageFailed || u.Stage == SFTPUploadStageReceived
}
