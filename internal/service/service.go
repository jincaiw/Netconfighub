package service

import (
	"context"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
)

type AdminService interface {
	Create(ctx context.Context, username, password string) (*model.Admin, error)
	Authenticate(ctx context.Context, username, password string) (*model.Admin, error)
	GetByID(ctx context.Context, id uint) (*model.Admin, error)
	List(ctx context.Context, offset, limit int) ([]model.Admin, int64, error)
	Update(ctx context.Context, id uint, username string, password *string) (*model.Admin, error)
	Delete(ctx context.Context, id uint) error
}

type DeviceService interface {
	Create(ctx context.Context, device *model.Device) (*model.Device, error)
	GetByID(ctx context.Context, id uint) (*model.Device, error)
	List(ctx context.Context, groupID *uint, vendor *string, status *string, search *string, offset, limit int) ([]model.Device, int64, error)
	Update(ctx context.Context, id uint, device *model.Device) (*model.Device, error)
	Delete(ctx context.Context, id uint) error
	TriggerBackup(ctx context.Context, deviceID uint) (*model.BackupTask, error)
	GetDeviceConfig(ctx context.Context, deviceID uint) (string, error)
	GetDeviceVersions(ctx context.Context, deviceID uint, limit int) ([]model.ConfigVersion, error)
	GetDeviceDeviations(ctx context.Context, deviceID uint, offset, limit int) ([]model.Deviation, int64, error)
	ListFailedDevices(ctx context.Context, offset, limit int) ([]model.Device, int64, error)
	GetConfigDiff(ctx context.Context, deviceID, fromHash, toHash string) (*model.DiffResult, error)
	GetLatestDiff(ctx context.Context, deviceID string) (*model.DiffResult, error)
	ImportDevices(ctx context.Context, csvData []byte) (*model.ImportResult, error)
	ExportDevices(ctx context.Context) ([]byte, error)
	EnableDevice(ctx context.Context, id uint) error
	DisableDevice(ctx context.Context, id uint) error
}

type GroupService interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	GetByID(ctx context.Context, id uint) (*model.Group, error)
	List(ctx context.Context, offset, limit int) ([]model.Group, int64, error)
	Update(ctx context.Context, id uint, group *model.Group) (*model.Group, error)
	Delete(ctx context.Context, id uint) error
	TriggerGroupBackup(ctx context.Context, groupID uint) ([]model.BackupTask, error)
}

type BackupService interface {
	GetByID(ctx context.Context, id uint) (*model.BackupTask, error)
	List(ctx context.Context, deviceID *uint, status *model.BackupStatus, offset, limit int) ([]model.BackupTask, int64, error)
	GetConfig(ctx context.Context, taskID uint) (string, error)
}

type BaselineService interface {
	Create(ctx context.Context, baseline *model.Baseline) (*model.Baseline, error)
	GetByID(ctx context.Context, id uint) (*model.Baseline, error)
	List(ctx context.Context, scope *model.BaselineScope, offset, limit int) ([]model.Baseline, int64, error)
	Update(ctx context.Context, id uint, baseline *model.Baseline) (*model.Baseline, error)
	Delete(ctx context.Context, id uint) error
	FindBaselineForDevice(ctx context.Context, deviceID uint) (*model.Baseline, error)
}

type DeviationService interface {
	Create(ctx context.Context, deviation *model.Deviation) error
	GetByID(ctx context.Context, id uint) (*model.Deviation, error)
	List(ctx context.Context, baselineID *uint, deviceID *uint, offset, limit int) ([]model.Deviation, int64, error)
	GetByDeviceID(ctx context.Context, deviceID uint) ([]model.Deviation, error)
	ExportDeviations(ctx context.Context, deviceID, baselineID *uint) ([]byte, error)
}

type APITokenService interface {
	Create(ctx context.Context, adminID uint, name string, expiresAt interface{}) (*model.APIToken, string, error)
	GetByToken(ctx context.Context, token string) (*model.APIToken, error)
	List(ctx context.Context, adminID uint) ([]model.APIToken, error)
	Delete(ctx context.Context, id uint, adminID uint) error
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, time.Time, error)
	Logout(ctx context.Context) error
	RefreshToken(ctx context.Context, tokenString string) (string, time.Time, error)
	InitAdmin(ctx context.Context, username, password string) (*model.Admin, error)
	ChangePassword(ctx context.Context, adminID uint, oldPassword, newPassword string) error
}

type AuditService interface {
	Log(ctx context.Context, userID uint, username, action, targetType string, targetID uint, detail, clientIP string) error
	List(ctx context.Context, offset, limit int) ([]model.AuditLog, int64, error)
}
