package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/model"
)

func (r *Router) listHooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	hooks, total, err := r.hookService.List(c.Request.Context(), offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询 Hook 列表失败")
		return
	}
	PageResult(c, hooks, total, page, pageSize)
}

func (r *Router) createHook(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Config  string `json:"config" binding:"required"`
		Events  string `json:"events" binding:"required"`
		Enabled *bool  `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	hook := &model.Hook{
		Name:    req.Name,
		Type:    req.Type,
		Config:  req.Config,
		Events:  req.Events,
		Enabled: enabled,
	}
	result, err := r.hookService.Create(c.Request.Context(), hook)
	if err != nil {
		BadRequest(c, "创建 Hook 失败")
		return
	}

	r.auditLog(c, "create", "hook", result.ID, "创建 Hook: "+req.Name)

	Success(c, result)
}

func (r *Router) getHook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}
	hook, err := r.hookService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "Hook 不存在")
		return
	}
	Success(c, hook)
}

func (r *Router) updateHook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}
	var req struct {
		Name    string `json:"name" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Config  string `json:"config" binding:"required"`
		Events  string `json:"events" binding:"required"`
		Enabled *bool  `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	hook := &model.Hook{
		Name:    req.Name,
		Type:    req.Type,
		Config:  req.Config,
		Events:  req.Events,
		Enabled: enabled,
	}
	result, err := r.hookService.Update(c.Request.Context(), uint(id), hook)
	if err != nil {
		BadRequest(c, "更新 Hook 失败")
		return
	}

	r.auditLog(c, "update", "hook", uint(id), "更新 Hook: "+req.Name)

	Success(c, result)
}

func (r *Router) deleteHook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}
	if err := r.hookService.Delete(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "Hook 不存在")
		return
	}

	r.auditLog(c, "delete", "hook", uint(id), "删除 Hook")

	Success(c, nil)
}
