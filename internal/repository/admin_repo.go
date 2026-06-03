package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type adminRepoImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepoImpl{db: db}
}

func (r *adminRepoImpl) Create(ctx context.Context, admin *model.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *adminRepoImpl) FindByID(ctx context.Context, id uint) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.WithContext(ctx).First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepoImpl) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepoImpl) List(ctx context.Context, offset, limit int) ([]model.Admin, int64, error) {
	var admins []model.Admin
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Admin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&admins).Error; err != nil {
		return nil, 0, err
	}
	return admins, total, nil
}

func (r *adminRepoImpl) Update(ctx context.Context, admin *model.Admin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}

func (r *adminRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Admin{}, id).Error
}
