package service

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type ConfigService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	List(ctx context.Context) ([]model.SystemConfig, error)
}

type configServiceImpl struct {
	configRepo repository.ConfigRepository
}

func NewConfigService(configRepo repository.ConfigRepository) ConfigService {
	return &configServiceImpl{configRepo: configRepo}
}

func (s *configServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return s.configRepo.Get(ctx, key)
}

func (s *configServiceImpl) Set(ctx context.Context, key, value string) error {
	return s.configRepo.Set(ctx, key, value)
}

func (s *configServiceImpl) List(ctx context.Context) ([]model.SystemConfig, error) {
	return s.configRepo.List(ctx)
}
