package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type baselineRepoImpl struct {
	db *gorm.DB
}

func NewBaselineRepository(db *gorm.DB) BaselineRepository {
	return &baselineRepoImpl{db: db}
}

func (r *baselineRepoImpl) Create(ctx context.Context, baseline *model.Baseline) error {
	return r.db.WithContext(ctx).Create(baseline).Error
}

func (r *baselineRepoImpl) FindByID(ctx context.Context, id uint) (*model.Baseline, error) {
	var baseline model.Baseline
	if err := r.db.WithContext(ctx).Preload("Device").Preload("Group").First(&baseline, id).Error; err != nil {
		return nil, err
	}
	return &baseline, nil
}

func (r *baselineRepoImpl) List(ctx context.Context, scope *model.BaselineScope, offset, limit int) ([]model.Baseline, int64, error) {
	var baselines []model.Baseline
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.Baseline{})
	if scope != nil {
		countQuery = countQuery.Where("scope = ?", *scope)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.Baseline{})
	if scope != nil {
		findQuery = findQuery.Where("scope = ?", *scope)
	}
	if err := findQuery.Preload("Device").Preload("Group").Offset(offset).Limit(limit).Order("created_at DESC").Find(&baselines).Error; err != nil {
		return nil, 0, err
	}

	return baselines, total, nil
}

func (r *baselineRepoImpl) Update(ctx context.Context, baseline *model.Baseline) error {
	return r.db.WithContext(ctx).Save(baseline).Error
}

func (r *baselineRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Baseline{}, id).Error
}

func (r *baselineRepoImpl) FindByDeviceID(ctx context.Context, deviceID uint) ([]model.Baseline, error) {
	var baselines []model.Baseline
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Find(&baselines).Error; err != nil {
		return nil, err
	}
	return baselines, nil
}

func (r *baselineRepoImpl) FindByGroupID(ctx context.Context, groupID uint) ([]model.Baseline, error) {
	var baselines []model.Baseline
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&baselines).Error; err != nil {
		return nil, err
	}
	return baselines, nil
}
