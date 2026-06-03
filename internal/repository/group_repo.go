package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type groupRepoImpl struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepoImpl{db: db}
}

func (r *groupRepoImpl) Create(ctx context.Context, group *model.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *groupRepoImpl) FindByID(ctx context.Context, id uint) (*model.Group, error) {
	var group model.Group
	if err := r.db.WithContext(ctx).First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepoImpl) List(ctx context.Context, offset, limit int) ([]model.Group, int64, error) {
	var groups []model.Group
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Group{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (r *groupRepoImpl) Update(ctx context.Context, group *model.Group) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *groupRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Group{}, id).Error
}
