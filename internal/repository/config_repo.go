package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConfigRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	List(ctx context.Context) ([]model.SystemConfig, error)
}

type configRepoImpl struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepoImpl{db: db}
}

func (r *configRepoImpl) Get(ctx context.Context, key string) (string, error) {
	var cfg model.SystemConfig
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&cfg).Error; err != nil {
		return "", err
	}
	return cfg.Value, nil
}

func (r *configRepoImpl) Set(ctx context.Context, key, value string) error {
	cfg := model.SystemConfig{
		Key:   key,
		Value: value,
	}
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&cfg)
	return result.Error
}

func (r *configRepoImpl) List(ctx context.Context) ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	if err := r.db.WithContext(ctx).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}
