package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type deviationRepoImpl struct {
	db *gorm.DB
}

func NewDeviationRepository(db *gorm.DB) DeviationRepository {
	return &deviationRepoImpl{db: db}
}

func (r *deviationRepoImpl) Create(ctx context.Context, deviation *model.Deviation) error {
	return r.db.WithContext(ctx).Create(deviation).Error
}

func (r *deviationRepoImpl) FindByID(ctx context.Context, id uint) (*model.Deviation, error) {
	var deviation model.Deviation
	if err := r.db.WithContext(ctx).Preload("BackupTask.Device").Preload("Baseline").First(&deviation, id).Error; err != nil {
		return nil, err
	}
	return &deviation, nil
}

func (r *deviationRepoImpl) List(ctx context.Context, baselineID *uint, offset, limit int) ([]model.Deviation, int64, error) {
	var deviations []model.Deviation
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.Deviation{})
	if baselineID != nil {
		countQuery = countQuery.Where("baseline_id = ?", *baselineID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.Deviation{})
	if baselineID != nil {
		findQuery = findQuery.Where("baseline_id = ?", *baselineID)
	}
	if err := findQuery.Preload("BackupTask.Device").Preload("Baseline").Offset(offset).Limit(limit).Order("detected_at DESC").Find(&deviations).Error; err != nil {
		return nil, 0, err
	}

	return deviations, total, nil
}

func (r *deviationRepoImpl) FindByDeviceID(ctx context.Context, deviceID uint, offset, limit int) ([]model.Deviation, int64, error) {
	var deviations []model.Deviation
	var total int64

	subQuery := r.db.WithContext(ctx).Model(&model.BackupTask{}).Select("id").Where("device_id = ?", deviceID)

	countQuery := r.db.WithContext(ctx).Model(&model.Deviation{}).Where("backup_task_id IN (?)", subQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.Deviation{}).Where("backup_task_id IN (?)", subQuery)
	if err := findQuery.Preload("BackupTask.Device").Preload("Baseline").Offset(offset).Limit(limit).Order("detected_at DESC").Find(&deviations).Error; err != nil {
		return nil, 0, err
	}

	return deviations, total, nil
}

func (r *deviationRepoImpl) FindByDeviceIDAndBaselineID(ctx context.Context, deviceID uint, baselineID uint, offset, limit int) ([]model.Deviation, int64, error) {
	var deviations []model.Deviation
	var total int64

	subQuery := r.db.WithContext(ctx).Model(&model.BackupTask{}).Select("id").Where("device_id = ?", deviceID)

	countQuery := r.db.WithContext(ctx).Model(&model.Deviation{}).Where("backup_task_id IN (?)", subQuery).Where("baseline_id = ?", baselineID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.Deviation{}).Where("backup_task_id IN (?)", subQuery).Where("baseline_id = ?", baselineID)
	if err := findQuery.Preload("BackupTask.Device").Preload("Baseline").Offset(offset).Limit(limit).Order("detected_at DESC").Find(&deviations).Error; err != nil {
		return nil, 0, err
	}

	return deviations, total, nil
}

func (r *deviationRepoImpl) CountByBaselineID(ctx context.Context, baselineID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Deviation{}).Where("baseline_id = ?", baselineID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
