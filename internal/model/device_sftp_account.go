package model

import (
	"time"

	"gorm.io/gorm"
)

// Device SFTP account status values. Revoking an account soft-deletes the row
// and additionally marks it disabled.
const (
	SFTPAccountStatusActive   = "active"
	SFTPAccountStatusDisabled = "disabled"
)

// DeviceSFTPAccount is the credential a device uses to log in to the inbound
// SFTP server. Only the bcrypt hash is stored: the plaintext password is
// returned exactly once by the service layer and is never persisted or logged.
type DeviceSFTPAccount struct {
	ID uint `gorm:"primarykey" json:"id"`
	// The unique indexes are scoped by deleted_at so that revoking an account
	// (soft delete) frees the device and username for a new account.
	DeviceID     uint           `gorm:"not null;uniqueIndex:idx_sftp_account_device_active,where:deleted_at IS NULL" json:"device_id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex:idx_sftp_account_username_active,where:deleted_at IS NULL" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Status       string         `gorm:"size:16;not null;default:active" json:"status"`
	LastUsedAt   *time.Time     `json:"last_used_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DeviceSFTPAccount) TableName() string {
	return "device_sftp_accounts"
}

// IsActive reports whether the account may currently authenticate.
func (a *DeviceSFTPAccount) IsActive() bool {
	return a.Status == SFTPAccountStatusActive
}
