package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/model"
)

func (r *Router) listBackups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var deviceID *uint
	if did := c.Query("device_id"); did != "" {
		v, err := strconv.ParseUint(did, 10, 64)
		if err == nil {
			id := uint(v)
			deviceID = &id
		}
	}

	var status *model.BackupStatus
	if s := c.Query("status"); s != "" {
		bs := model.BackupStatus(s)
		status = &bs
	}

	tasks, total, err := r.backupService.List(c.Request.Context(), deviceID, status, offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询备份列表失败")
		return
	}
	PageResult(c, tasks, total, page, pageSize)
}

func (r *Router) getBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	task, err := r.backupService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "备份任务不存在")
		return
	}
	Success(c, task)
}

func (r *Router) getBackupConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	content, err := r.backupService.GetConfig(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "无法获取备份配置")
		return
	}
	Success(c, gin.H{"content": content})
}
