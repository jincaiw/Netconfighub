package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) listAlerts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var alertType *string
	if t := c.Query("type"); t != "" {
		alertType = &t
	}

	var severity *string
	if s := c.Query("severity"); s != "" {
		severity = &s
	}

	var deviceID *uint
	if d := c.Query("device_id"); d != "" {
		if id, err := strconv.ParseUint(d, 10, 64); err == nil {
			uid := uint(id)
			deviceID = &uid
		}
	}

	alerts, total, err := r.alertService.List(c.Request.Context(), alertType, severity, deviceID, offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询告警列表失败")
		return
	}
	PageResult(c, alerts, total, page, pageSize)
}

func (r *Router) countUnreadAlerts(c *gin.Context) {
	count, err := r.alertService.CountUnread(c.Request.Context())
	if err != nil {
		DatabaseError(c, "查询未读告警数失败")
		return
	}
	Success(c, gin.H{"count": count})
}

func (r *Router) markAlertAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.alertService.MarkAsRead(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "告警不存在")
		return
	}
	Success(c, nil)
}

func (r *Router) markAllAlertsAsRead(c *gin.Context) {
	if err := r.alertService.MarkAllAsRead(c.Request.Context()); err != nil {
		DatabaseError(c, "标记全部已读失败")
		return
	}
	Success(c, nil)
}

func (r *Router) deleteAlert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}
	if err := r.alertService.Delete(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "告警不存在")
		return
	}
	Success(c, nil)
}
