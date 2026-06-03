package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type deviationServiceImpl struct {
	deviationRepo repository.DeviationRepository
}

func NewDeviationService(deviationRepo repository.DeviationRepository) DeviationService {
	return &deviationServiceImpl{deviationRepo: deviationRepo}
}

func (s *deviationServiceImpl) Create(ctx context.Context, deviation *model.Deviation) error {
	return s.deviationRepo.Create(ctx, deviation)
}

func (s *deviationServiceImpl) GetByID(ctx context.Context, id uint) (*model.Deviation, error) {
	deviation, err := s.deviationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("偏差记录不存在")
	}
	return deviation, nil
}

func (s *deviationServiceImpl) List(ctx context.Context, baselineID *uint, deviceID *uint, offset, limit int) ([]model.Deviation, int64, error) {
	if deviceID != nil && baselineID != nil {
		return s.deviationRepo.FindByDeviceIDAndBaselineID(ctx, *deviceID, *baselineID, offset, limit)
	}
	if deviceID != nil {
		return s.deviationRepo.FindByDeviceID(ctx, *deviceID, offset, limit)
	}
	return s.deviationRepo.List(ctx, baselineID, offset, limit)
}

func (s *deviationServiceImpl) GetByDeviceID(ctx context.Context, deviceID uint) ([]model.Deviation, error) {
	deviations, _, err := s.deviationRepo.FindByDeviceID(ctx, deviceID, 0, 1000)
	if err != nil {
		return nil, err
	}
	return deviations, nil
}

func (s *deviationServiceImpl) ExportDeviations(ctx context.Context, deviceID, baselineID *uint) ([]byte, error) {
	deviations, _, err := s.deviationRepo.List(ctx, baselineID, 0, 10000)
	if err != nil {
		return nil, fmt.Errorf("查询偏差列表失败: %w", err)
	}
	if deviceID != nil {
		filtered := make([]model.Deviation, 0)
		for _, d := range deviations {
			if d.BackupTask.DeviceID == *deviceID {
				filtered = append(filtered, d)
			}
		}
		deviations = filtered
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"device_name", "baseline_scope", "diff_content", "detected_at"})
	for _, d := range deviations {
		deviceName := ""
		if d.BackupTask.Device.ID != 0 {
			deviceName = d.BackupTask.Device.Name
		}
		baselineScope := ""
		if d.Baseline.ID != 0 {
			baselineScope = string(d.Baseline.Scope)
		}
		writer.Write([]string{
			deviceName,
			baselineScope,
			d.DiffContent,
			d.DetectedAt.Format("2006-01-02 15:04:05"),
		})
	}
	writer.Flush()
	return buf.Bytes(), nil
}
