package service

import (
	"context"
	"fmt"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type baselineServiceImpl struct {
	baselineRepo  repository.BaselineRepository
	deviceRepo    repository.DeviceRepository
	deviationRepo repository.DeviationRepository
}

func NewBaselineService(baselineRepo repository.BaselineRepository, deviceRepo repository.DeviceRepository, deviationRepo repository.DeviationRepository) BaselineService {
	return &baselineServiceImpl{
		baselineRepo:  baselineRepo,
		deviceRepo:    deviceRepo,
		deviationRepo: deviationRepo,
	}
}

func (s *baselineServiceImpl) Create(ctx context.Context, baseline *model.Baseline) (*model.Baseline, error) {
	if baseline.Scope != model.BaselineScopeDevice && baseline.Scope != model.BaselineScopeGroup {
		return nil, fmt.Errorf("无效的基线作用域")
	}
	if baseline.Content == "" {
		return nil, fmt.Errorf("基线内容不能为空")
	}
	if baseline.Scope == model.BaselineScopeDevice && baseline.DeviceID == nil {
		return nil, fmt.Errorf("设备级基线必须指定 device_id")
	}
	if baseline.Scope == model.BaselineScopeGroup && baseline.GroupID == nil {
		return nil, fmt.Errorf("分组级基线必须指定 group_id")
	}
	if baseline.Scope == model.BaselineScopeDevice && baseline.DeviceID != nil {
		existing, err := s.baselineRepo.FindByDeviceID(ctx, *baseline.DeviceID)
		if err == nil && len(existing) > 0 {
			return nil, fmt.Errorf("该设备已存在基线，请先删除旧基线或更新")
		}
	}
	if baseline.Scope == model.BaselineScopeGroup && baseline.GroupID != nil {
		existing, err := s.baselineRepo.FindByGroupID(ctx, *baseline.GroupID)
		if err == nil && len(existing) > 0 {
			return nil, fmt.Errorf("该分组已存在基线，请先删除旧基线或更新")
		}
	}
	if err := s.baselineRepo.Create(ctx, baseline); err != nil {
		return nil, err
	}
	return s.baselineRepo.FindByID(ctx, baseline.ID)
}

func (s *baselineServiceImpl) GetByID(ctx context.Context, id uint) (*model.Baseline, error) {
	baseline, err := s.baselineRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("基线不存在")
	}
	return baseline, nil
}

func (s *baselineServiceImpl) List(ctx context.Context, scope *model.BaselineScope, offset, limit int) ([]model.Baseline, int64, error) {
	return s.baselineRepo.List(ctx, scope, offset, limit)
}

func (s *baselineServiceImpl) Update(ctx context.Context, id uint, baseline *model.Baseline) (*model.Baseline, error) {
	existing, err := s.baselineRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("基线不存在")
	}
	if baseline.Content != "" {
		existing.Content = baseline.Content
	}
	if err := s.baselineRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.baselineRepo.FindByID(ctx, id)
}

func (s *baselineServiceImpl) Delete(ctx context.Context, id uint) error {
	if _, err := s.baselineRepo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("基线不存在")
	}
	count, err := s.deviationRepo.CountByBaselineID(ctx, id)
	if err != nil {
		return fmt.Errorf("检查关联偏差记录失败")
	}
	if count > 0 {
		return fmt.Errorf("该基线存在 %d 条关联偏差记录，无法删除", count)
	}
	return s.baselineRepo.Delete(ctx, id)
}

func (s *baselineServiceImpl) FindBaselineForDevice(ctx context.Context, deviceID uint) (*model.Baseline, error) {
	deviceBaselines, err := s.baselineRepo.FindByDeviceID(ctx, deviceID)
	if err == nil && len(deviceBaselines) > 0 {
		return &deviceBaselines[0], nil
	}
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, nil
	}
	if device.GroupID != nil {
		groupBaselines, err := s.baselineRepo.FindByGroupID(ctx, *device.GroupID)
		if err == nil && len(groupBaselines) > 0 {
			return &groupBaselines[0], nil
		}
	}
	return nil, nil
}
