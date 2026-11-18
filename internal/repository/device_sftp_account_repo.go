package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type deviceSFTPAccountRepoImpl struct {
	db *gorm.DB
}

func NewDeviceSFTPAccountRepository(db *gorm.DB) DeviceSFTPAccountRepository {
	return &deviceSFTPAccountRepoImpl{db: db}
}

func (r *deviceSFTPAccountRepoImpl) Create(ctx context.Context, account *model.DeviceSFTPAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *deviceSFTPAccountRepoImpl) FindByUsername(ctx context.Context, username string) (*model.DeviceSFTPAccount, error) {
	var account model.DeviceSFTPAccount
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *deviceSFTPAccountRepoImpl) FindByDeviceID(ctx context.Context, deviceID uint) (*model.DeviceSFTPAccount, error) {
	var account model.DeviceSFTPAccount
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *deviceSFTPAccountRepoImpl) Update(ctx context.Context, account *model.DeviceSFTPAccount) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// Revoke soft-deletes the account and marks it disabled in one statement, so a
// revoked account can never authenticate even if a lookup bypasses the soft
// delete scope.
// Revoke removes the device's account. The row is hard-deleted rather than
// soft-deleted so the (device_id) and (username) unique indexes stay free for a
// replacement account; a revoked credential must never authenticate again, and
// the audit log keeps the history of who revoked it.
func (r *deviceSFTPAccountRepoImpl) Revoke(ctx context.Context, deviceID uint) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Where("device_id = ?", deviceID).
		Delete(&model.DeviceSFTPAccount{}).Error
}

func (r *deviceSFTPAccountRepoImpl) ListByDeviceID(ctx context.Context, deviceID uint) ([]model.DeviceSFTPAccount, error) {
	var accounts []model.DeviceSFTPAccount
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
