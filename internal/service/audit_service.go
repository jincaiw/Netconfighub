package service

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type auditServiceImpl struct {
	auditRepo repository.AuditLogRepository
}

func NewAuditService(auditRepo repository.AuditLogRepository) AuditService {
	return &auditServiceImpl{auditRepo: auditRepo}
}

func (s *auditServiceImpl) Log(ctx context.Context, userID uint, username, action, targetType string, targetID uint, detail, clientIP string) error {
	log := &model.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
		ClientIP:   clientIP,
	}
	return s.auditRepo.Create(ctx, log)
}

func (s *auditServiceImpl) List(ctx context.Context, offset, limit int) ([]model.AuditLog, int64, error) {
	return s.auditRepo.List(ctx, offset, limit)
}
