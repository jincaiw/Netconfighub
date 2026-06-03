package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type HookRepository interface {
	Create(ctx context.Context, hook *model.Hook) error
	FindByID(ctx context.Context, id uint) (*model.Hook, error)
	List(ctx context.Context, offset, limit int) ([]model.Hook, int64, error)
	Update(ctx context.Context, hook *model.Hook) error
	Delete(ctx context.Context, id uint) error
	FindEnabledByEvent(ctx context.Context, event string) ([]model.Hook, error)
}

type hookRepoImpl struct {
	db *gorm.DB
}

func NewHookRepository(db *gorm.DB) HookRepository {
	return &hookRepoImpl{db: db}
}

func (r *hookRepoImpl) Create(ctx context.Context, hook *model.Hook) error {
	return r.db.WithContext(ctx).Create(hook).Error
}

func (r *hookRepoImpl) FindByID(ctx context.Context, id uint) (*model.Hook, error) {
	var hook model.Hook
	if err := r.db.WithContext(ctx).First(&hook, id).Error; err != nil {
		return nil, err
	}
	return &hook, nil
}

func (r *hookRepoImpl) List(ctx context.Context, offset, limit int) ([]model.Hook, int64, error) {
	var hooks []model.Hook
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Hook{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&hooks).Error; err != nil {
		return nil, 0, err
	}
	return hooks, total, nil
}

func (r *hookRepoImpl) Update(ctx context.Context, hook *model.Hook) error {
	return r.db.WithContext(ctx).Save(hook).Error
}

func (r *hookRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Hook{}, id).Error
}

func (r *hookRepoImpl) FindEnabledByEvent(ctx context.Context, event string) ([]model.Hook, error) {
	var hooks []model.Hook
	if err := r.db.WithContext(ctx).Where("enabled = ? AND events LIKE ?", true, "%"+event+"%").Find(&hooks).Error; err != nil {
		return nil, err
	}
	return hooks, nil
}
