package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/model"
)

type createBaselineRequest struct {
	DeviceID *uint               `json:"device_id"`
	GroupID  *uint               `json:"group_id"`
	Scope    model.BaselineScope `json:"scope" binding:"required"`
	Content  string              `json:"content" binding:"required"`
}

type updateBaselineRequest struct {
	DeviceID *uint                `json:"device_id"`
	GroupID  *uint                `json:"group_id"`
	Scope    *model.BaselineScope `json:"scope"`
	Content  *string              `json:"content"`
}

func (r *Router) listBaselines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var scope *model.BaselineScope
	if s := c.Query("scope"); s != "" {
		bs := model.BaselineScope(s)
		scope = &bs
	}

	baselines, total, err := r.baselineService.List(c.Request.Context(), scope, offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询基线列表失败")
		return
	}
	PageResult(c, baselines, total, page, pageSize)
}

func (r *Router) createBaseline(c *gin.Context) {
	var req createBaselineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	baseline := &model.Baseline{
		DeviceID: req.DeviceID,
		GroupID:  req.GroupID,
		Scope:    req.Scope,
		Content:  req.Content,
	}

	result, err := r.baselineService.Create(c.Request.Context(), baseline)
	if err != nil {
		BadRequest(c, "创建基线失败")
		return
	}

	r.auditLog(c, "create", "baseline", result.ID, "创建基线")

	Success(c, result)
}

func (r *Router) getBaseline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	baseline, err := r.baselineService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "基线不存在")
		return
	}
	Success(c, baseline)
}

func (r *Router) updateBaseline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	var req updateBaselineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	baseline := &model.Baseline{}
	if req.DeviceID != nil {
		baseline.DeviceID = req.DeviceID
	}
	if req.GroupID != nil {
		baseline.GroupID = req.GroupID
	}
	if req.Scope != nil {
		baseline.Scope = *req.Scope
	}
	if req.Content != nil {
		baseline.Content = *req.Content
	}

	result, err := r.baselineService.Update(c.Request.Context(), uint(id), baseline)
	if err != nil {
		BadRequest(c, "更新基线失败")
		return
	}

	r.auditLog(c, "update", "baseline", uint(id), "更新基线")

	Success(c, result)
}

func (r *Router) deleteBaseline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.baselineService.Delete(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "基线不存在" {
			Error(c, http.StatusNotFound, 404, "基线不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, "delete", "baseline", uint(id), "删除基线")

	Success(c, nil)
}
