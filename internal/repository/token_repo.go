package repository

import (
	"context"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type apiTokenRepoImpl struct {
	db *gorm.DB
}

func NewAPITokenRepository(db *gorm.DB) APITokenRepository {
	return &apiTokenRepoImpl{db: db}
}

func (r *apiTokenRepoImpl) Create(ctx context.Context, token *model.APIToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *apiTokenRepoImpl) FindByID(ctx context.Context, id uint) (*model.APIToken, error) {
	var token model.APIToken
	if err := r.db.WithContext(ctx).First(&token, id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *apiTokenRepoImpl) FindByTokenHash(ctx context.Context, tokenHash string) (*model.APIToken, error) {
	var token model.APIToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *apiTokenRepoImpl) List(ctx context.Context, adminID uint) ([]model.APIToken, error) {
	var tokens []model.APIToken
	if err := r.db.WithContext(ctx).Where("admin_id = ?", adminID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *apiTokenRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.APIToken{}, id).Error
}

func (r *apiTokenRepoImpl) UpdateLastUsed(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.APIToken{}).Where("id = ?", id).Update("last_used_at", now).Error
}
