package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type auditLogRepoImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepoImpl{db: db}
}

func (r *auditLogRepoImpl) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepoImpl) List(ctx context.Context, offset, limit int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
