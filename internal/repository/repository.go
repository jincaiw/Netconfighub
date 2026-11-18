package repository

import (
	"context"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *model.Admin) error
	FindByID(ctx context.Context, id uint) (*model.Admin, error)
	FindByUsername(ctx context.Context, username string) (*model.Admin, error)
	List(ctx context.Context, offset, limit int) ([]model.Admin, int64, error)
	Update(ctx context.Context, admin *model.Admin) error
	Delete(ctx context.Context, id uint) error
}

type DeviceRepository interface {
	Create(ctx context.Context, device *model.Device) error
	FindByID(ctx context.Context, id uint) (*model.Device, error)
	FindByName(ctx context.Context, name string) (*model.Device, error)
	FindByIP(ctx context.Context, ip string) (*model.Device, error)
	List(ctx context.Context, groupID *uint, vendor *string, status *string, search *string, offset, limit int) ([]model.Device, int64, error)
	Update(ctx context.Context, device *model.Device) error
	Delete(ctx context.Context, id uint) error
	FindByGroupID(ctx context.Context, groupID uint) ([]model.Device, error)
	ListFailed(ctx context.Context, offset, limit int) ([]model.Device, int64, error)
	UpdateBackupState(ctx context.Context, id uint, state BackupStateUpdate) error
}

// BackupStateUpdate carries the four backup-related device columns that the
// scheduler updates after a collection attempt. It exists so those updates can
// be applied by ID without rewriting (and re-encrypting) the whole device row.
type BackupStateUpdate struct {
	Status        string
	At            *time.Time
	FailureReason string
	RetryCount    int
}

type GroupRepository interface {
	Create(ctx context.Context, group *model.Group) error
	FindByID(ctx context.Context, id uint) (*model.Group, error)
	List(ctx context.Context, offset, limit int) ([]model.Group, int64, error)
	ListAllNames(ctx context.Context, limit int) ([]model.Group, error)
	Update(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, id uint) error
}

type BackupTaskRepository interface {
	Create(ctx context.Context, task *model.BackupTask) error
	FindByID(ctx context.Context, id uint) (*model.BackupTask, error)
	List(ctx context.Context, deviceID *uint, status *model.BackupStatus, offset, limit int) ([]model.BackupTask, int64, error)
	Update(ctx context.Context, task *model.BackupTask) error
	UpdateStatus(ctx context.Context, id uint, status model.BackupStatus, failureReason string, errorType string) error
	FindByStatus(ctx context.Context, status model.BackupStatus) ([]model.BackupTask, error)
}

type BaselineRepository interface {
	Create(ctx context.Context, baseline *model.Baseline) error
	FindByID(ctx context.Context, id uint) (*model.Baseline, error)
	List(ctx context.Context, scope *model.BaselineScope, offset, limit int) ([]model.Baseline, int64, error)
	Update(ctx context.Context, baseline *model.Baseline) error
	Delete(ctx context.Context, id uint) error
	FindByDeviceID(ctx context.Context, deviceID uint) ([]model.Baseline, error)
	FindByGroupID(ctx context.Context, groupID uint) ([]model.Baseline, error)
}

type APITokenRepository interface {
	Create(ctx context.Context, token *model.APIToken) error
	FindByID(ctx context.Context, id uint) (*model.APIToken, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.APIToken, error)
	List(ctx context.Context, adminID uint) ([]model.APIToken, error)
	Delete(ctx context.Context, id uint) error
	UpdateLastUsed(ctx context.Context, id uint) error
}

type DeviationRepository interface {
	Create(ctx context.Context, deviation *model.Deviation) error
	FindByID(ctx context.Context, id uint) (*model.Deviation, error)
	List(ctx context.Context, baselineID *uint, offset, limit int) ([]model.Deviation, int64, error)
	FindByDeviceID(ctx context.Context, deviceID uint, offset, limit int) ([]model.Deviation, int64, error)
	FindByDeviceIDAndBaselineID(ctx context.Context, deviceID uint, baselineID uint, offset, limit int) ([]model.Deviation, int64, error)
	CountByBaselineID(ctx context.Context, baselineID uint) (int64, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	List(ctx context.Context, offset, limit int) ([]model.AuditLog, int64, error)
}

// DeviceSFTPAccountRepository persists the per-device credentials used by the
// inbound SFTP server. Only bcrypt hashes are ever stored.
type DeviceSFTPAccountRepository interface {
	Create(ctx context.Context, account *model.DeviceSFTPAccount) error
	FindByUsername(ctx context.Context, username string) (*model.DeviceSFTPAccount, error)
	FindByDeviceID(ctx context.Context, deviceID uint) (*model.DeviceSFTPAccount, error)
	Update(ctx context.Context, account *model.DeviceSFTPAccount) error
	Revoke(ctx context.Context, deviceID uint) error
	ListByDeviceID(ctx context.Context, deviceID uint) ([]model.DeviceSFTPAccount, error)
}

// SFTPUploadRepository persists inbound upload records and their two-stage
// reception/archiving state.
type SFTPUploadRepository interface {
	Create(ctx context.Context, upload *model.SFTPUpload) error
	FindByID(ctx context.Context, id uint) (*model.SFTPUpload, error)
	ListByDevice(ctx context.Context, deviceID uint, offset, limit int) ([]model.SFTPUpload, int64, error)
	ListByStage(ctx context.Context, stages []string, limit int) ([]model.SFTPUpload, error)
	// UpdateStage applies the given column updates only when the row is still in
	// one of fromStages, returning the number of affected rows so callers can
	// detect a concurrent archiver.
	UpdateStage(ctx context.Context, id uint, fromStages []string, updates map[string]interface{}) (int64, error)
	IncrementRetry(ctx context.Context, id uint) error
}
