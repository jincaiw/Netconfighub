package repository

import (
	"context"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type backupTaskRepoImpl struct {
	db *gorm.DB
}

func NewBackupTaskRepository(db *gorm.DB) BackupTaskRepository {
	return &backupTaskRepoImpl{db: db}
}

func (r *backupTaskRepoImpl) Create(ctx context.Context, task *model.BackupTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *backupTaskRepoImpl) FindByID(ctx context.Context, id uint) (*model.BackupTask, error) {
	var task model.BackupTask
	if err := r.db.WithContext(ctx).Preload("Device").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *backupTaskRepoImpl) List(ctx context.Context, deviceID *uint, status *model.BackupStatus, offset, limit int) ([]model.BackupTask, int64, error) {
	var tasks []model.BackupTask
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.BackupTask{})
	if deviceID != nil {
		countQuery = countQuery.Where("device_id = ?", *deviceID)
	}
	if status != nil {
		countQuery = countQuery.Where("status = ?", *status)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.BackupTask{})
	if deviceID != nil {
		findQuery = findQuery.Where("device_id = ?", *deviceID)
	}
	if status != nil {
		findQuery = findQuery.Where("status = ?", *status)
	}
	if err := findQuery.Preload("Device").Offset(offset).Limit(limit).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	fillDeviceName(tasks)

	return tasks, total, nil
}

func (r *backupTaskRepoImpl) Update(ctx context.Context, task *model.BackupTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *backupTaskRepoImpl) UpdateStatus(ctx context.Context, id uint, status model.BackupStatus, failureReason string, errorType string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}
	if errorType != "" {
		updates["error_type"] = errorType
	}
	return r.db.WithContext(ctx).Model(&model.BackupTask{}).Where("id = ?", id).Updates(updates).Error
}

func (r *backupTaskRepoImpl) FindByStatus(ctx context.Context, status model.BackupStatus) ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	if err := r.db.WithContext(ctx).Where("status = ?", status).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func fillDeviceName(tasks []model.BackupTask) {
	for i := range tasks {
		if tasks[i].Device.Name != "" {
			tasks[i].DeviceName = tasks[i].Device.Name
		}
	}
}

// FindRetryableTasks returns up to limit retrying tasks that still have retries
// left, oldest first. Bounding the query prevents the retry loop from loading
// the entire history into memory on every tick.
func (r *backupTaskRepoImpl) FindRetryableTasks(ctx context.Context, limit int) ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	q := r.db.WithContext(ctx).Where("status = ? AND retry_count < max_retries", model.BackupStatusRetrying).Order("id ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// ClaimRetryingTask atomically transitions a retrying task to running. It returns
// true only when this caller won the claim, making retry execution idempotent and
// preventing duplicate concurrent retries of the same task.
func (r *backupTaskRepoImpl) ClaimRetryingTask(ctx context.Context, id uint) (bool, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.BackupTask{}).
		Where("id = ? AND status = ?", id, model.BackupStatusRetrying).
		Updates(map[string]interface{}{
			"status":      model.BackupStatusRunning,
			"started_at":  now,
			"finished_at": nil,
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// ListSince returns tasks created at/after since (newest first), capped by limit.
// Used by the dashboard to aggregate over a bounded, complete range instead of a
// small sample.
func (r *backupTaskRepoImpl) ListSince(ctx context.Context, since time.Time, limit int) ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	q := r.db.WithContext(ctx).Where("created_at >= ?", since).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Preload("Device").Find(&tasks).Error; err != nil {
		return nil, err
	}
	fillDeviceName(tasks)
	return tasks, nil
}

// PruneHistory bounds backup history per device: each device keeps its most
// recent keepSuccessPerDevice successful tasks and keepFailedPerDevice other
// tasks; everything else is removed. This clears runaway retry/failure storms
// (which may all be recent) that an age-only window would not remove. keepDays
// is retained as an absolute upper age cap: rows older than keepDays that are
// not in a device's kept sets are also removed. Returns rows removed.
func (r *backupTaskRepoImpl) PruneHistory(ctx context.Context, keepDays int, keepSuccessPerDevice int, keepFailedPerDevice int) (int64, error) {
	if keepSuccessPerDevice < 0 {
		keepSuccessPerDevice = 0
	}
	if keepFailedPerDevice < 0 {
		keepFailedPerDevice = 0
	}
	var devIDs []uint
	if err := r.db.WithContext(ctx).Model(&model.BackupTask{}).Distinct().Pluck("device_id", &devIDs).Error; err != nil {
		return 0, err
	}
	kept := make([]uint, 0, len(devIDs)*(keepSuccessPerDevice+keepFailedPerDevice))
	for _, d := range devIDs {
		if keepSuccessPerDevice > 0 {
			var ids []uint
			if err := r.db.WithContext(ctx).Model(&model.BackupTask{}).
				Where("device_id = ? AND status = ?", d, model.BackupStatusSuccess).
				Order("created_at DESC").Limit(keepSuccessPerDevice).Pluck("id", &ids).Error; err != nil {
				return 0, err
			}
			kept = append(kept, ids...)
		}
		if keepFailedPerDevice > 0 {
			var ids []uint
			if err := r.db.WithContext(ctx).Model(&model.BackupTask{}).
				Where("device_id = ? AND status <> ?", d, model.BackupStatusSuccess).
				Order("created_at DESC").Limit(keepFailedPerDevice).Pluck("id", &ids).Error; err != nil {
				return 0, err
			}
			kept = append(kept, ids...)
		}
	}
	var res *gorm.DB
	if len(kept) == 0 {
		res = r.db.WithContext(ctx).Where("1 = 1").Delete(&model.BackupTask{})
	} else {
		res = r.db.WithContext(ctx).Where("id NOT IN ?", kept).Delete(&model.BackupTask{})
	}
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// FinalizeStaleTasks marks leftover non-terminal tasks (pending/running/retrying)
// as failed. Called at startup before workers begin so stale retry artifacts from
// a previous run do not linger as "进行中" or get re-processed indefinitely.
func (r *backupTaskRepoImpl) FinalizeStaleTasks(ctx context.Context) (int64, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.BackupTask{}).
		Where("status IN ?", []model.BackupStatus{model.BackupStatusPending, model.BackupStatusRunning, model.BackupStatusRetrying}).
		Updates(map[string]interface{}{
			"status":         model.BackupStatusFailed,
			"error_type":     "stale_interrupted",
			"failure_reason": "系统重启前未完成的历史任务，已归并为失败",
			"finished_at":    now,
		})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
