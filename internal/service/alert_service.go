package service

import (
	"context"
	"fmt"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type AlertService interface {
	Create(ctx context.Context, alert *model.Alert) error
	GetByID(ctx context.Context, id uint) (*model.Alert, error)
	List(ctx context.Context, alertType *string, severity *string, deviceID *uint, offset, limit int) ([]model.Alert, int64, error)
	CountUnread(ctx context.Context) (int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context) error
	Delete(ctx context.Context, id uint) error
}

type alertServiceImpl struct {
	alertRepo repository.AlertRepository
}

func NewAlertService(alertRepo repository.AlertRepository) AlertService {
	return &alertServiceImpl{alertRepo: alertRepo}
}

func (s *alertServiceImpl) Create(ctx context.Context, alert *model.Alert) error {
	return s.alertRepo.Create(ctx, alert)
}

func (s *alertServiceImpl) GetByID(ctx context.Context, id uint) (*model.Alert, error) {
	alert, err := s.alertRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("告警不存在")
	}
	return alert, nil
}

func (s *alertServiceImpl) List(ctx context.Context, alertType *string, severity *string, deviceID *uint, offset, limit int) ([]model.Alert, int64, error) {
	return s.alertRepo.List(ctx, alertType, severity, deviceID, offset, limit)
}

func (s *alertServiceImpl) CountUnread(ctx context.Context) (int64, error) {
	return s.alertRepo.CountUnread(ctx)
}

func (s *alertServiceImpl) MarkAsRead(ctx context.Context, id uint) error {
	return s.alertRepo.MarkAsRead(ctx, id)
}

func (s *alertServiceImpl) MarkAllAsRead(ctx context.Context) error {
	return s.alertRepo.MarkAllAsRead(ctx)
}

func (s *alertServiceImpl) Delete(ctx context.Context, id uint) error {
	if _, err := s.alertRepo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("告警不存在")
	}
	return s.alertRepo.Delete(ctx, id)
}
