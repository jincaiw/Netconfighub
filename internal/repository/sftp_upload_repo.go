package repository

import (
	"context"

	"github.com/netconfighub/netconfighub/internal/model"
	"gorm.io/gorm"
)

type sftpUploadRepoImpl struct {
	db *gorm.DB
}

func NewSFTPUploadRepository(db *gorm.DB) SFTPUploadRepository {
	return &sftpUploadRepoImpl{db: db}
}

func (r *sftpUploadRepoImpl) Create(ctx context.Context, upload *model.SFTPUpload) error {
	return r.db.WithContext(ctx).Create(upload).Error
}

func (r *sftpUploadRepoImpl) FindByID(ctx context.Context, id uint) (*model.SFTPUpload, error) {
	var upload model.SFTPUpload
	if err := r.db.WithContext(ctx).First(&upload, id).Error; err != nil {
		return nil, err
	}
	return &upload, nil
}

func (r *sftpUploadRepoImpl) ListByDevice(ctx context.Context, deviceID uint, offset, limit int) ([]model.SFTPUpload, int64, error) {
	var uploads []model.SFTPUpload
	var total int64

	countQuery := r.db.WithContext(ctx).Model(&model.SFTPUpload{}).Where("device_id = ?", deviceID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.db.WithContext(ctx).Model(&model.SFTPUpload{}).Where("device_id = ?", deviceID)
	if err := findQuery.Order("created_at DESC").Offset(offset).Limit(limit).Find(&uploads).Error; err != nil {
		return nil, 0, err
	}

	return uploads, total, nil
}

func (r *sftpUploadRepoImpl) ListByStage(ctx context.Context, stages []string, limit int) ([]model.SFTPUpload, error) {
	if len(stages) == 0 {
		return nil, nil
	}
	query := r.db.WithContext(ctx).Where("stage IN ?", stages).Order("created_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	var uploads []model.SFTPUpload
	if err := query.Find(&uploads).Error; err != nil {
		return nil, err
	}
	return uploads, nil
}

// UpdateStage applies updates only when the row is still in one of fromStages.
// The returned count lets callers detect that another archiver already moved
// the row, so the same upload is never archived twice.
func (r *sftpUploadRepoImpl) UpdateStage(ctx context.Context, id uint, fromStages []string, updates map[string]interface{}) (int64, error) {
	if len(fromStages) == 0 {
		return 0, nil
	}
	result := r.db.WithContext(ctx).Model(&model.SFTPUpload{}).
		Where("id = ? AND stage IN ?", id, fromStages).
		Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *sftpUploadRepoImpl) IncrementRetry(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.SFTPUpload{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}
