package scheduler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/collector"
	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/netconfighub/netconfighub/internal/crypto"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/notifier"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/service"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type RetryConfig struct {
	MaxRetries int           `yaml:"max_retries"`
	RetryDelay time.Duration `yaml:"retry_delay"`
}

type Scheduler struct {
	cron           *cron.Cron
	db             *gorm.DB
	cfg            *config.SchedulerConfig
	deviceRepo     repository.DeviceRepository
	backupTaskRepo repository.BackupTaskRepository
	baselineRepo   repository.BaselineRepository
	deviationRepo  repository.DeviationRepository
	collector      *collector.UnifiedCollector
	gitStore       *gitstore.Store
	baselineEngine baseline.Engine
	alertService   service.AlertService
	auditService   service.AuditService
	notifier       notifier.Notifier
	hookService    service.HookService
	workerPoolSize int
	taskQueue      chan uint
	wg             sync.WaitGroup
	cancel         context.CancelFunc
	encryptionKey  string
	maxRetries     int
	retryDelay     time.Duration
	taskTimeout    time.Duration
}

func NewScheduler(
	db *gorm.DB,
	cfg *config.SchedulerConfig,
	deviceRepo repository.DeviceRepository,
	backupTaskRepo repository.BackupTaskRepository,
	baselineRepo repository.BaselineRepository,
	deviationRepo repository.DeviationRepository,
	c *collector.UnifiedCollector,
	gs *gitstore.Store,
	be baseline.Engine,
	alertService service.AlertService,
	ntf notifier.Notifier,
	encryptionKey string,
	hookService service.HookService,
	auditService service.AuditService,
) *Scheduler {
	maxRetries := cfg.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}
	retryDelay, err := time.ParseDuration(cfg.RetryDelay)
	if err != nil || retryDelay == 0 {
		retryDelay = 30 * time.Second
	}
	taskTimeout, err := time.ParseDuration(cfg.TaskTimeout)
	if err != nil || taskTimeout == 0 {
		taskTimeout = 5 * time.Minute
	}

	return &Scheduler{
		cron:           cron.New(cron.WithSeconds()),
		db:             db,
		cfg:            cfg,
		deviceRepo:     deviceRepo,
		backupTaskRepo: backupTaskRepo,
		baselineRepo:   baselineRepo,
		deviationRepo:  deviationRepo,
		collector:      c,
		gitStore:       gs,
		baselineEngine: be,
		alertService:   alertService,
		auditService:   auditService,
		notifier:       ntf,
		hookService:    hookService,
		workerPoolSize: cfg.WorkerPoolSize,
		taskQueue:      make(chan uint, cfg.WorkerPoolSize*2),
		encryptionKey:  encryptionKey,
		maxRetries:     maxRetries,
		retryDelay:     retryDelay,
		taskTimeout:    taskTimeout,
	}
}

func (s *Scheduler) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	for i := 0; i < s.workerPoolSize; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}

	go s.processRetryingTasks(ctx)

	interval := s.cfg.DefaultInterval
	_, err := s.cron.AddFunc(fmt.Sprintf("@every %s", interval), func() {
		s.scheduleAllDevices(ctx)
	})
	if err != nil {
		cancel()
		return fmt.Errorf("添加定时任务失败: %w", err)
	}

	s.cron.Start()
	return nil
}

func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.cron.Stop()
	close(s.taskQueue)
	s.wg.Wait()
}

func (s *Scheduler) SubmitDevice(deviceID uint) error {
	select {
	case s.taskQueue <- deviceID:
		return nil
	default:
		return fmt.Errorf("备份队列已满，请稍后重试")
	}
}

func (s *Scheduler) SubmitGroup(groupID uint) error {
	if s.deviceRepo == nil {
		return fmt.Errorf("设备仓库未初始化")
	}
	devices, err := s.deviceRepo.FindByGroupID(context.Background(), groupID)
	if err != nil {
		return fmt.Errorf("查询分组设备失败: %w", err)
	}
	if len(devices) == 0 {
		return fmt.Errorf("分组下没有设备")
	}
	for _, device := range devices {
		if !device.Enabled {
			continue
		}
		select {
		case s.taskQueue <- device.ID:
		default:
			log.Printf("队列已满，跳过设备 %d", device.ID)
		}
	}
	return nil
}

func (s *Scheduler) worker(ctx context.Context, id int) {
	defer s.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case deviceID, ok := <-s.taskQueue:
			if !ok {
				return
			}
			s.executeBackup(ctx, deviceID)
		}
	}
}

func (s *Scheduler) executeBackup(ctx context.Context, deviceID uint) {
	if s.deviceRepo == nil {
		log.Printf("deviceRepo 未初始化，跳过设备 %d", deviceID)
		return
	}

	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		log.Printf("查询设备失败 [deviceID=%d]: %v", deviceID, err)
		return
	}

	if !device.Enabled {
		now := time.Now()
		task := &model.BackupTask{
			DeviceID:  deviceID,
			Status:    model.BackupStatusSkipped,
			StartedAt: &now,
			MaxRetries: s.maxRetries,
		}
		task.FinishedAt = &now
		task.FailureReason = "设备已禁用"
		task.ErrorType = "device_disabled"
		if s.backupTaskRepo != nil {
			if err := s.backupTaskRepo.Create(ctx, task); err != nil {
				log.Printf("创建跳过任务失败 [deviceID=%d]: %v", deviceID, err)
			}
		}
		log.Printf("设备已禁用，跳过备份 [deviceID=%d]", deviceID)
		return
	}

	s.decryptDeviceCredentials(device)

	now := time.Now()
	task := &model.BackupTask{
		DeviceID:   deviceID,
		Status:     model.BackupStatusRunning,
		StartedAt:  &now,
		MaxRetries: s.maxRetries,
	}
	if s.backupTaskRepo == nil {
		log.Printf("backupTaskRepo 未初始化，跳过设备 %d", deviceID)
		return
	}
	if err := s.backupTaskRepo.Create(ctx, task); err != nil {
		log.Printf("创建备份任务失败 [deviceID=%d]: %v", deviceID, err)
		return
	}

	if s.collector == nil {
		s.handleBackupFailure(ctx, task, device, fmt.Errorf("采集器未初始化"))
		return
	}

	backupCtx, cancel := context.WithTimeout(ctx, s.taskTimeout)
	defer cancel()

	type collectResult struct {
		content string
		err     error
	}
	resultCh := make(chan collectResult, 1)
	go func() {
		content, err := s.collector.Collect(device)
		resultCh <- collectResult{content: content, err: err}
	}()

	select {
	case <-backupCtx.Done():
		if backupCtx.Err() == context.DeadlineExceeded {
			s.handleBackupFailure(ctx, task, device, fmt.Errorf("采集超时，已超过 %v", s.taskTimeout))
			s.handleTimeout(ctx, task, device)
		} else {
			s.handleBackupFailure(ctx, task, device, fmt.Errorf("采集被取消"))
		}
	case result := <-resultCh:
		if result.err != nil {
			s.handleBackupFailure(ctx, task, device, result.err)
		} else {
			s.handleBackupSuccess(ctx, task, device, result.content)
		}
	}
}

func (s *Scheduler) handleBackupSuccess(ctx context.Context, task *model.BackupTask, device *model.Device, configContent string) {
	now := time.Now()

	if s.gitStore != nil {
		if _, err := s.gitStore.WriteConfigAsync(string(device.Vendor), device.Name, configContent); err != nil {
			log.Printf("写入 Git 存储失败 [device=%s]: %v", device.Name, err)
		}
	}

	if s.baselineEngine != nil && s.baselineRepo != nil {
		b := s.findBaselineForDevice(ctx, device)
		if b != nil {
			result, err := s.baselineEngine.Compare(ctx, b.Content, configContent, b.ID)
			if err != nil {
				log.Printf("基线比对失败 [deviceID=%d]: %v", device.ID, err)
			} else if result != nil && result.HasDeviation {
				log.Printf("检测到配置偏差 [deviceID=%d]: added=%d, removed=%d", device.ID, result.AddedLines, result.RemovedLines)
				if s.deviationRepo != nil {
					deviation := &model.Deviation{
						BackupTaskID: task.ID,
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
					s.auditService.Log(ctx, 0, "system", "baseline_deviation", "device", device.ID,
						fmt.Sprintf("设备 %s 检测到配置偏差: 新增 %d 行，删除 %d 行", device.Name, result.AddedLines, result.RemovedLines), "")
				}
			}
		}
	}

	hash := sha256.Sum256([]byte(configContent))
	configHash := hex.EncodeToString(hash[:])

	task.Status = model.BackupStatusSuccess
	task.ConfigHash = configHash
	finishedAt := time.Now()
	task.FinishedAt = &finishedAt
	if err := s.backupTaskRepo.Update(ctx, task); err != nil {
		log.Printf("更新备份任务失败 [taskID=%d]: %v", task.ID, err)
	}

	device.LastBackupStatus = string(model.BackupStatusSuccess)
	device.LastBackupAt = &now
	device.LastFailureReason = ""
	device.RetryCount = 0
	device.Group = nil
	s.encryptDeviceCredentials(device)
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		log.Printf("更新设备状态失败 [deviceID=%d]: %v", device.ID, err)
	}

	if s.hookService != nil {
		s.hookService.ExecuteHooks(ctx, "backup_success", map[string]interface{}{
			"event":       "backup_success",
			"device_id":   device.ID,
			"device_name": device.Name,
			"timestamp":   now.Format(time.RFC3339),
		})
	}

	if s.auditService != nil {
		s.auditService.Log(ctx, 0, "system", "backup_success", "device", device.ID,
			fmt.Sprintf("设备 %s 备份成功", device.Name), "")
	}
}

func (s *Scheduler) handleBackupFailure(ctx context.Context, task *model.BackupTask, device *model.Device, backupErr error) {
	now := time.Now()

	errorType := classifyError(backupErr)

	task.RetryCount = task.RetryCount + 1
	task.ErrorType = errorType
	task.FailureReason = backupErr.Error()

	if task.RetryCount < task.MaxRetries {
		task.Status = model.BackupStatusRetrying
	} else {
		task.Status = model.BackupStatusFailed
	}

	finishedAt := time.Now()
	task.FinishedAt = &finishedAt
	if err := s.backupTaskRepo.Update(ctx, task); err != nil {
		log.Printf("更新备份任务失败 [taskID=%d]: %v", task.ID, err)
	}

	device.LastBackupStatus = string(task.Status)
	device.LastBackupAt = &now
	device.LastFailureReason = backupErr.Error()
	device.RetryCount = task.RetryCount
	device.Group = nil
	s.encryptDeviceCredentials(device)
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		log.Printf("更新设备状态失败 [deviceID=%d]: %v", device.ID, err)
	}

	if task.Status == model.BackupStatusFailed {
		if s.alertService != nil {
			alert := &model.Alert{
				Type:       "backup_failed",
				DeviceID:   &device.ID,
				DeviceName: device.Name,
				Title:      fmt.Sprintf("设备 %s 备份失败", device.Name),
				Message:    fmt.Sprintf("错误类型: %s, 原因: %s", errorType, backupErr.Error()),
				Severity:   "error",
			}
			if err := s.alertService.Create(ctx, alert); err != nil {
				log.Printf("创建告警失败 [deviceID=%d]: %v", device.ID, err)
			}
		}

		if s.notifier != nil {
			subject := fmt.Sprintf("NetConfigHub 告警: 设备 %s 备份失败", device.Name)
			body := fmt.Sprintf("设备 %s 备份失败: %s", device.Name, backupErr.Error())
			if err := s.notifier.Send(ctx, subject, body); err != nil {
				log.Printf("发送通知失败 [deviceID=%d]: %v", device.ID, err)
			}
		}

		if s.hookService != nil {
			s.hookService.ExecuteHooks(ctx, "backup_failed", map[string]interface{}{
				"event":       "backup_failed",
				"device_id":   device.ID,
				"device_name": device.Name,
				"error":       backupErr.Error(),
				"error_type":  errorType,
				"timestamp":   now.Format(time.RFC3339),
			})
		}

		if s.auditService != nil {
			s.auditService.Log(ctx, 0, "system", "backup_failed", "device", device.ID,
				fmt.Sprintf("设备 %s 备份失败: [%s] %s", device.Name, errorType, backupErr.Error()), "")
		}
	}
}

func (s *Scheduler) handleTimeout(ctx context.Context, task *model.BackupTask, device *model.Device) {
	task.Status = model.BackupStatusTimeout
	task.ErrorType = "command_timeout"
	if err := s.backupTaskRepo.Update(ctx, task); err != nil {
		log.Printf("更新超时任务失败 [taskID=%d]: %v", task.ID, err)
	}

	device.LastBackupStatus = string(model.BackupStatusTimeout)
	device.Group = nil
	s.encryptDeviceCredentials(device)
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		log.Printf("更新设备状态失败 [deviceID=%d]: %v", device.ID, err)
	}

	if s.alertService != nil {
		alert := &model.Alert{
			Type:       "backup_timeout",
			DeviceID:   &device.ID,
			DeviceName: device.Name,
			Title:      fmt.Sprintf("设备 %s 备份超时", device.Name),
			Message:    fmt.Sprintf("采集超时，已超过 %v", s.taskTimeout),
			Severity:   "warning",
		}
		if err := s.alertService.Create(ctx, alert); err != nil {
			log.Printf("创建超时告警失败 [deviceID=%d]: %v", device.ID, err)
		}
	}
}

func classifyError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "连接超时") || strings.Contains(msg, "连接失败"):
		return "connect_timeout"
	case strings.Contains(msg, "认证失败") || strings.Contains(msg, "password"):
		return "auth_failed"
	case strings.Contains(msg, "提示符超时") || strings.Contains(msg, "prompt"):
		return "prompt_timeout"
	case strings.Contains(msg, "命令执行超时") || strings.Contains(msg, "超时"):
		return "command_timeout"
	case strings.Contains(msg, "模型不存在") || strings.Contains(msg, "厂商"):
		return "model_not_found"
	case strings.Contains(msg, "Git"):
		return "git_commit_failed"
	case strings.Contains(msg, "配置为空"):
		return "config_empty"
	default:
		return "unknown_error"
	}
}

func (s *Scheduler) processRetryingTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.retryDelay):
			if s.backupTaskRepo == nil {
				continue
			}
			tasks, err := s.backupTaskRepo.FindByStatus(ctx, model.BackupStatusRetrying)
			if err != nil {
				log.Printf("查询重试任务失败: %v", err)
				continue
			}
			for _, task := range tasks {
				select {
				case <-ctx.Done():
					return
				case s.taskQueue <- task.DeviceID:
					log.Printf("重新提交重试任务 [deviceID=%d, retryCount=%d]", task.DeviceID, task.RetryCount)
				default:
					log.Printf("队列已满，跳过重试任务 [deviceID=%d]", task.DeviceID)
				}
			}
		}
	}
}

func (s *Scheduler) findBaselineForDevice(ctx context.Context, device *model.Device) *model.Baseline {
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

func (s *Scheduler) scheduleAllDevices(ctx context.Context) {
	if s.deviceRepo == nil {
		return
	}
	devices, _, err := s.deviceRepo.List(ctx, nil, nil, nil, nil, 0, 10000)
	if err != nil {
		log.Printf("查询所有设备失败: %v", err)
		return
	}
	now := time.Now()
	for _, device := range devices {
		if !device.Enabled {
			continue
		}
		interval := s.cfg.DefaultInterval
		if device.BackupInterval != "" {
			interval = device.BackupInterval
		}
		if device.LastBackupAt != nil {
			dur, err := time.ParseDuration(interval)
			if err != nil {
				dur, _ = time.ParseDuration(s.cfg.DefaultInterval)
			}
			if now.Sub(*device.LastBackupAt) < dur {
				continue
			}
		}
		select {
		case <-ctx.Done():
			return
		case s.taskQueue <- device.ID:
		default:
		}
	}
}

func (s *Scheduler) AddCronJob(spec string, cmd func()) error {
	_, err := s.cron.AddFunc(spec, cmd)
	if err != nil {
		return fmt.Errorf("添加 cron 任务失败 [%s]: %w", spec, err)
	}
	return nil
}

func (s *Scheduler) decryptDeviceCredentials(device *model.Device) {
	if device.Password != "" {
		decrypted, err := crypto.Decrypt(device.Password, s.encryptionKey)
		if err != nil {
			log.Printf("解密设备密码失败 [deviceID=%d]: %v", device.ID, err)
		} else {
			device.Password = decrypted
		}
	}
	if device.SSHKey != "" {
		decrypted, err := crypto.Decrypt(device.SSHKey, s.encryptionKey)
		if err != nil {
			log.Printf("解密设备SSH密钥失败 [deviceID=%d]: %v", device.ID, err)
		} else {
			device.SSHKey = decrypted
		}
	}
	if device.EnablePassword != "" {
		decrypted, err := crypto.Decrypt(device.EnablePassword, s.encryptionKey)
		if err != nil {
			log.Printf("解密设备Enable密码失败 [deviceID=%d]: %v", device.ID, err)
		} else {
			device.EnablePassword = decrypted
		}
	}
}

func (s *Scheduler) encryptDeviceCredentials(device *model.Device) {
	if device.Password != "" {
		enc, err := crypto.Encrypt(device.Password, s.encryptionKey)
		if err == nil {
			device.Password = enc
		}
	}
	if device.SSHKey != "" {
		enc, err := crypto.Encrypt(device.SSHKey, s.encryptionKey)
		if err == nil {
			device.SSHKey = enc
		}
	}
	if device.EnablePassword != "" {
		enc, err := crypto.Encrypt(device.EnablePassword, s.encryptionKey)
		if err == nil {
			device.EnablePassword = enc
		}
	}
}
