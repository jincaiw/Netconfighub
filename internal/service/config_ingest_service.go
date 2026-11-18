package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/netconfighub/netconfighub/internal/apperr"
	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/notifier"
	"github.com/netconfighub/netconfighub/internal/repository"
)

// IngestDeps groups the dependencies of the shared configuration ingest
// pipeline. Both the pull scheduler and the SFTP archive workers feed it, so
// the Git write, baseline comparison, alerting, hooks and audit trail behave
// identically for every source.
type IngestDeps struct {
	DeviceRepo     repository.DeviceRepository
	BackupTaskRepo repository.BackupTaskRepository
	BaselineRepo   repository.BaselineRepository
	DeviationRepo  repository.DeviationRepository
	GitStore       *gitstore.Store
	BaselineEngine baseline.Engine
	AlertService   AlertService
	Notifier       notifier.Notifier
	HookService    HookService
	AuditService   AuditService
}

type configIngestServiceImpl struct {
	deviceRepo     repository.DeviceRepository
	backupTaskRepo repository.BackupTaskRepository
	baselineRepo   repository.BaselineRepository
	deviationRepo  repository.DeviationRepository
	gitStore       *gitstore.Store
	baselineEngine baseline.Engine
	alertService   AlertService
	notifier       notifier.Notifier
	hookService    HookService
	auditService   AuditService
}

func NewConfigIngestService(deps IngestDeps) ConfigIngestService {
	return &configIngestServiceImpl{
		deviceRepo:     deps.DeviceRepo,
		backupTaskRepo: deps.BackupTaskRepo,
		baselineRepo:   deps.BaselineRepo,
		deviationRepo:  deps.DeviationRepo,
		gitStore:       deps.GitStore,
		baselineEngine: deps.BaselineEngine,
		alertService:   deps.AlertService,
		notifier:       deps.Notifier,
		hookService:    deps.HookService,
		auditService:   deps.AuditService,
	}
}

// Ingest persists one configuration and applies every downstream side effect in
// a fixed order: Git write, content hash, baseline comparison, device state,
// backup task, success hook and audit.
//
// The Git write is synchronous and any failure is returned to the caller: a
// configuration that never reached the Git store must never be reported as a
// successful backup.
func (s *configIngestServiceImpl) Ingest(ctx context.Context, req IngestRequest) (*IngestResult, error) {
	if req.Device == nil {
		return nil, fmt.Errorf("入库失败: 设备不能为空")
	}
	if s.gitStore == nil {
		return nil, fmt.Errorf("入库失败: Git 存储未初始化")
	}
	device := req.Device

	// 1) Write to the Git store synchronously. A duplicate content write is a
	// no-op that returns an empty hash.
	commitHash, err := s.gitStore.WriteConfig(string(device.Vendor), device.Name, req.Content)
	if err != nil {
		return nil, apperr.WrapCause(apperr.ErrGitCommit, err, "写入 Git 存储失败: %v", err)
	}

	// 2) Content hash. This is the SHA256 of the configuration text, not the Git
	// commit hash, so identical content always maps to the same value.
	sum := sha256.Sum256([]byte(req.Content))
	configHash := hex.EncodeToString(sum[:])
	changed := commitHash != ""

	now := time.Now()

	// 3) Baseline comparison: device baseline first, group baseline as fallback.
	s.applyBaselineComparison(ctx, device, req.Content, req.TaskID, now)

	// 4) Partial device state update. Never Save(device): the caller's copy holds
	// decrypted credentials, so a full-row write would re-encrypt them.
	if s.deviceRepo != nil {
		if err := s.deviceRepo.UpdateBackupState(ctx, device.ID, repository.BackupStateUpdate{
			Status:        string(model.BackupStatusSuccess),
			At:            &now,
			FailureReason: "",
			RetryCount:    0,
		}); err != nil {
			return nil, fmt.Errorf("更新设备状态失败: %w", err)
		}
	}

	// 5) Backup task result.
	if s.backupTaskRepo != nil && req.TaskID != 0 {
		if err := s.markTaskSuccess(ctx, req.TaskID, configHash, now); err != nil {
			return nil, err
		}
	}

	// 6) Success hook and audit.
	if s.hookService != nil {
		s.hookService.ExecuteHooks(ctx, "backup_success", map[string]interface{}{
			"event":       "backup_success",
			"device_id":   device.ID,
			"device_name": device.Name,
			"source":      req.Source,
			"timestamp":   now.Format(time.RFC3339),
		})
	}
	if s.auditService != nil {
		if err := s.auditService.Log(ctx, 0, "system", "backup_success", "device", device.ID,
			fmt.Sprintf("设备 %s 备份成功 (来源: %s)", device.Name, req.Source), ""); err != nil {
			log.Printf("记录审计日志失败 [deviceID=%d]: %v", device.ID, err)
		}
	}

	return &IngestResult{ConfigHash: configHash, Changed: changed}, nil
}

func (s *configIngestServiceImpl) markTaskSuccess(ctx context.Context, taskID uint, configHash string, finishedAt time.Time) error {
	task, err := s.backupTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("查询备份任务失败: %w", err)
	}
	task.Status = model.BackupStatusSuccess
	task.ConfigHash = configHash
	task.FinishedAt = &finishedAt
	task.FailureReason = ""
	task.ErrorType = ""
	if err := s.backupTaskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("更新备份任务失败: %w", err)
	}
	return nil
}

// applyBaselineComparison mirrors the previous scheduler behaviour: a device
// baseline wins, a group baseline is the fallback. Comparison and its side
// effects are best effort and never fail the ingest.
func (s *configIngestServiceImpl) applyBaselineComparison(ctx context.Context, device *model.Device, content string, taskID uint, now time.Time) {
	if s.baselineEngine == nil || s.baselineRepo == nil {
		return
	}
	b := s.findBaselineForDevice(ctx, device)
	if b == nil {
		return
	}
	result, err := s.baselineEngine.Compare(ctx, b.Content, content, b.ID)
	if err != nil {
		log.Printf("基线比对失败 [deviceID=%d]: %v", device.ID, err)
		return
	}
	if result == nil || !result.HasDeviation {
		return
	}

	log.Printf("检测到配置偏差 [deviceID=%d]: added=%d, removed=%d", device.ID, result.AddedLines, result.RemovedLines)

	if s.deviationRepo != nil {
		deviation := &model.Deviation{
			BackupTaskID: taskID,
			BaselineID:   b.ID,
			DiffContent:  result.DiffContent,
			DetectedAt:   now,
		}
		if err := s.deviationRepo.Create(ctx, deviation); err != nil {
			log.Printf("创建偏差记录失败 [deviceID=%d]: %v", device.ID, err)
		}
	}
	if s.alertService != nil {
		alert := &model.Alert{
			Type:       "baseline_deviation",
			DeviceID:   &device.ID,
			DeviceName: device.Name,
			Title:      fmt.Sprintf("设备 %s 检测到配置偏差", device.Name),
			Message:    fmt.Sprintf("新增 %d 行，删除 %d 行", result.AddedLines, result.RemovedLines),
			Severity:   "warning",
		}
		if err := s.alertService.Create(ctx, alert); err != nil {
			log.Printf("创建告警失败 [deviceID=%d]: %v", device.ID, err)
		}
	}
	if s.notifier != nil {
		subject := fmt.Sprintf("NetConfigHub 告警: 设备 %s 配置偏差", device.Name)
		body := fmt.Sprintf("设备 %s 检测到配置偏差: 新增 %d 行，删除 %d 行", device.Name, result.AddedLines, result.RemovedLines)
		if err := s.notifier.Send(ctx, subject, body); err != nil {
			log.Printf("发送通知失败 [deviceID=%d]: %v", device.ID, err)
		}
	}
	if s.hookService != nil {
		s.hookService.ExecuteHooks(ctx, "baseline_deviation", map[string]interface{}{
			"event":         "baseline_deviation",
			"device_id":     device.ID,
			"device_name":   device.Name,
			"added_lines":   result.AddedLines,
			"removed_lines": result.RemovedLines,
			"timestamp":     now.Format(time.RFC3339),
		})
	}
	if s.auditService != nil {
		if err := s.auditService.Log(ctx, 0, "system", "baseline_deviation", "device", device.ID,
			fmt.Sprintf("设备 %s 检测到配置偏差: 新增 %d 行，删除 %d 行", device.Name, result.AddedLines, result.RemovedLines), ""); err != nil {
			log.Printf("记录审计日志失败 [deviceID=%d]: %v", device.ID, err)
		}
	}
}

func (s *configIngestServiceImpl) findBaselineForDevice(ctx context.Context, device *model.Device) *model.Baseline {
	baselines, err := s.baselineRepo.FindByDeviceID(ctx, device.ID)
	if err == nil && len(baselines) > 0 {
		return &baselines[0]
	}
	if device.GroupID != nil {
		groupBaselines, err := s.baselineRepo.FindByGroupID(ctx, *device.GroupID)
		if err == nil && len(groupBaselines) > 0 {
			return &groupBaselines[0]
		}
	}
	return nil
}
