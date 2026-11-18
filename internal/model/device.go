package model

import (
	"time"

	"gorm.io/gorm"
)

type Vendor string

const (
	VendorCisco  Vendor = "cisco"
	VendorH3C    Vendor = "h3c"
	VendorHuawei Vendor = "huawei"
	VendorRuijie Vendor = "ruijie"
)

func (v Vendor) IsValid() bool {
	switch v {
	case VendorCisco, VendorH3C, VendorHuawei, VendorRuijie:
		return true
	}
	return false
}

type DeviceModel string

const (
	DeviceModelIOS     DeviceModel = "ios"
	DeviceModelVRP     DeviceModel = "vrp"
	DeviceModelComware DeviceModel = "comware"
	DeviceModelRGOS    DeviceModel = "rg-os"
)

func (m DeviceModel) IsValid() bool {
	switch m {
	case DeviceModelIOS, DeviceModelVRP, DeviceModelComware, DeviceModelRGOS:
		return true
	}
	return false
}

type ConnProtocol string

const (
	ConnProtocolSSH    ConnProtocol = "ssh"
	ConnProtocolTelnet ConnProtocol = "telnet"
)

func (c ConnProtocol) IsValid() bool {
	switch c {
	case ConnProtocolSSH, ConnProtocolTelnet:
		return true
	}
	return false
}

// CollectionMode describes how the platform obtains a device's configuration.
type CollectionMode string

const (
	// CollectionModePull means the platform logs in and pulls the configuration.
	CollectionModePull CollectionMode = "pull"
	// CollectionModeSFTPPush means the device pushes its configuration to the
	// platform over SFTP and must never be pulled.
	CollectionModeSFTPPush CollectionMode = "sftp_push"
)

func (m CollectionMode) IsValid() bool {
	switch m {
	case CollectionModePull, CollectionModeSFTPPush:
		return true
	}
	return false
}

type Device struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	Name              string         `gorm:"size:128;not null" json:"name"`
	IP                string         `gorm:"size:64;not null;column:ip" json:"ip"`
	Vendor            Vendor         `gorm:"size:32;not null" json:"vendor"`
	Model             DeviceModel    `gorm:"size:128" json:"model"`
	Protocol          ConnProtocol   `gorm:"size:16;not null;default:ssh;column:protocol" json:"protocol"`
	Port              int            `gorm:"not null;default:22" json:"port"`
	Username          string         `gorm:"size:64;not null" json:"username"`
	Password          string         `gorm:"size:255;not null" json:"-"`
	SSHKey            string         `gorm:"size:4096" json:"-"`
	EnablePassword    string         `gorm:"size:255" json:"-"`
	Enabled           bool           `gorm:"default:true" json:"enabled"`
	CollectionMode    CollectionMode `gorm:"size:16;not null;default:pull" json:"collection_mode"`
	GroupID           *uint          `gorm:"index" json:"group_id"`
	Group             *Group         `gorm:"foreignKey:GroupID" json:"-"`
	GroupName         string         `gorm:"-" json:"group_name"`
	LastBackupStatus  string         `gorm:"size:32" json:"last_backup_status"`
	LastBackupAt      *time.Time     `json:"last_backup_at"`
	LastFailureReason string         `gorm:"size:512" json:"last_failure_reason"`
	RetryCount        int            `gorm:"default:0" json:"retry_count"`
	BackupInterval    string         `gorm:"size:16" json:"backup_interval"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Device) TableName() string {
	return "devices"
}

// IsPullable reports whether the platform may collect this device's
// configuration itself. SFTP push devices are excluded from every pull path.
func (d *Device) IsPullable() bool {
	return d.CollectionMode != CollectionModeSFTPPush
}

type Group struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:128;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"size:512" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Group) TableName() string {
	return "groups"
}
