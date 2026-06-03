package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type deviceRepoImpl struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepoImpl{db: db}
}

func (r *deviceRepoImpl) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *deviceRepoImpl) FindByID(ctx context.Context, id uint) (*model.Device, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).Preload("Group").First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepoImpl) FindByName(ctx context.Context, name string) (*model.Device, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepoImpl) FindByIP(ctx context.Context, ip string) (*model.Device, error) {
	var device model.Device
	if err := r.db.WithContext(ctx).Where("ip = ?", ip).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepoImpl) List(ctx context.Context, groupID *uint, vendor *string, status *string, search *string, offset, limit int) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.Device{})
	if groupID != nil {
		countQuery = countQuery.Where("group_id = ?", *groupID)
	}
	if vendor != nil && *vendor != "" {
		countQuery = countQuery.Where("vendor = ?", *vendor)
	}
	if status != nil && *status != "" {
		countQuery = countQuery.Where("last_backup_status = ?", *status)
	}
	if search != nil && *search != "" {
		countQuery = countQuery.Where("name LIKE ? OR ip LIKE ?", "%"+*search+"%", "%"+*search+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.Device{})
	if groupID != nil {
		findQuery = findQuery.Where("group_id = ?", *groupID)
	}
	if vendor != nil && *vendor != "" {
		findQuery = findQuery.Where("vendor = ?", *vendor)
	}
	if status != nil && *status != "" {
		findQuery = findQuery.Where("last_backup_status = ?", *status)
	}
	if search != nil && *search != "" {
		findQuery = findQuery.Where("name LIKE ? OR ip LIKE ?", "%"+*search+"%", "%"+*search+"%")
	}
	if err := findQuery.Preload("Group").Offset(offset).Limit(limit).Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}

func (r *deviceRepoImpl) Update(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *deviceRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Device{}, id).Error
}

func (r *deviceRepoImpl) FindByGroupID(ctx context.Context, groupID uint) ([]model.Device, error) {
	var devices []model.Device
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepoImpl) ListFailed(ctx context.Context, offset, limit int) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Device{}).Where("last_backup_status IN ?", []string{"failed", "timeout", "retrying"}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Group").Where("last_backup_status IN ?", []string{"failed", "timeout", "retrying"}).Offset(offset).Limit(limit).Find(&devices).Error; err != nil {
		return nil, 0, err
	}

	return devices, total, nil
}
