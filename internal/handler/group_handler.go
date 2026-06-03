package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/model"
)

type createGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type updateGroupRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (r *Router) listGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	groups, total, err := r.groupService.List(c.Request.Context(), offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询分组列表失败")
		return
	}
	PageResult(c, groups, total, page, pageSize)
}

func (r *Router) createGroup(c *gin.Context) {
	var req createGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	group := &model.Group{
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := r.groupService.Create(c.Request.Context(), group)
	if err != nil {
		BadRequest(c, "创建分组失败")
		return
	}
	Success(c, result)
}

func (r *Router) getGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	group, err := r.groupService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "分组不存在")
		return
	}
	Success(c, group)
}

func (r *Router) updateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	var req updateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	group := &model.Group{}
	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}

	result, err := r.groupService.Update(c.Request.Context(), uint(id), group)
	if err != nil {
		BadRequest(c, "更新分组失败")
		return
	}
	Success(c, result)
}

func (r *Router) deleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.groupService.Delete(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "分组不存在" {
			Error(c, http.StatusNotFound, 404, "分组不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}
	Success(c, nil)
}

func (r *Router) triggerGroupBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if _, err := r.groupService.GetByID(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "分组不存在")
		return
	}

	if r.scheduler == nil {
		Error(c, http.StatusInternalServerError, 500, "调度器未初始化")
		return
	}

	if err := r.scheduler.SubmitGroup(uint(id)); err != nil {
		if strings.Contains(err.Error(), "分组下没有设备") {
			BadRequest(c, "分组下没有设备")
			return
		}
		BadRequest(c, "提交分组备份任务失败")
		return
	}

	r.auditLog(c, "trigger_backup", "group", uint(id), "触发分组备份")

	Success(c, gin.H{"message": "分组备份任务已提交"})
}
