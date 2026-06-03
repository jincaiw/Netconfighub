package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type backupTaskRepoImpl struct {
	db *gorm.DB
}

func NewBackupTaskRepository(db *gorm.DB) BackupTaskRepository {
	return &backupTaskRepoImpl{db: db}
}

func (r *backupTaskRepoImpl) Create(ctx context.Context, task *model.BackupTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *backupTaskRepoImpl) FindByID(ctx context.Context, id uint) (*model.BackupTask, error) {
	var task model.BackupTask
	if err := r.db.WithContext(ctx).Preload("Device").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *backupTaskRepoImpl) List(ctx context.Context, deviceID *uint, status *model.BackupStatus, offset, limit int) ([]model.BackupTask, int64, error) {
	var tasks []model.BackupTask
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.BackupTask{})
	if deviceID != nil {
		countQuery = countQuery.Where("device_id = ?", *deviceID)
	}
	if status != nil {
		countQuery = countQuery.Where("status = ?", *status)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.BackupTask{})
	if deviceID != nil {
		findQuery = findQuery.Where("device_id = ?", *deviceID)
	}
	if status != nil {
		findQuery = findQuery.Where("status = ?", *status)
	}
	if err := findQuery.Preload("Device").Offset(offset).Limit(limit).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *backupTaskRepoImpl) Update(ctx context.Context, task *model.BackupTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *backupTaskRepoImpl) UpdateStatus(ctx context.Context, id uint, status model.BackupStatus, failureReason string, errorType string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}
	if errorType != "" {
		updates["error_type"] = errorType
	}
	return r.db.WithContext(ctx).Model(&model.BackupTask{}).Where("id = ?", id).Updates(updates).Error
}

func (r *backupTaskRepoImpl) FindByStatus(ctx context.Context, status model.BackupStatus) ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
