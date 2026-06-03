package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *model.Alert) error
	FindByID(ctx context.Context, id uint) (*model.Alert, error)
	List(ctx context.Context, alertType *string, severity *string, deviceID *uint, offset, limit int) ([]model.Alert, int64, error)
	CountUnread(ctx context.Context) (int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context) error
	Delete(ctx context.Context, id uint) error
}

type alertRepoImpl struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepoImpl{db: db}
}

func (r *alertRepoImpl) Create(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *alertRepoImpl) FindByID(ctx context.Context, id uint) (*model.Alert, error) {
	var alert model.Alert
	if err := r.db.WithContext(ctx).First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *alertRepoImpl) List(ctx context.Context, alertType *string, severity *string, deviceID *uint, offset, limit int) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.Alert{})
	if alertType != nil && *alertType != "" {
		countQuery = countQuery.Where("type = ?", *alertType)
	}
	if severity != nil && *severity != "" {
		countQuery = countQuery.Where("severity = ?", *severity)
	}
	if deviceID != nil {
		countQuery = countQuery.Where("device_id = ?", *deviceID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.Alert{})
	if alertType != nil && *alertType != "" {
		findQuery = findQuery.Where("type = ?", *alertType)
	}
	if severity != nil && *severity != "" {
		findQuery = findQuery.Where("severity = ?", *severity)
	}
	if deviceID != nil {
		findQuery = findQuery.Where("device_id = ?", *deviceID)
	}
	if err := findQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&alerts).Error; err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

func (r *alertRepoImpl) CountUnread(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Alert{}).Where("is_read = ?", false).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *alertRepoImpl) MarkAsRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Alert{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *alertRepoImpl) MarkAllAsRead(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&model.Alert{}).Where("is_read = ?", false).Update("is_read", true).Error
}

func (r *alertRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Alert{}, id).Error
}
