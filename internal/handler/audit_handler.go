package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) listAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	logs, total, err := r.auditService.List(c.Request.Context(), offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询审计日志失败")
		return
	}
	PageResult(c, logs, total, page, pageSize)
}
