package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Router) listDeviations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var baselineID *uint
	if bid := c.Query("baseline_id"); bid != "" {
		v, err := strconv.ParseUint(bid, 10, 64)
		if err == nil {
			id := uint(v)
			baselineID = &id
		}
	}

	var deviceID *uint
	if did := c.Query("device_id"); did != "" {
		v, err := strconv.ParseUint(did, 10, 64)
		if err == nil {
			id := uint(v)
			deviceID = &id
		}
	}

	deviations, total, err := r.deviationService.List(c.Request.Context(), baselineID, deviceID, offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询偏差列表失败")
		return
	}
	PageResult(c, deviations, total, page, pageSize)
}

func (r *Router) getDeviation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	deviation, err := r.deviationService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "偏差记录不存在")
		return
	}
	Success(c, deviation)
}

func (r *Router) exportDeviations(c *gin.Context) {
	var deviceID *uint
	if did := c.Query("device_id"); did != "" {
		v, err := strconv.ParseUint(did, 10, 64)
		if err == nil {
			id := uint(v)
			deviceID = &id
		}
	}

	var baselineID *uint
	if bid := c.Query("baseline_id"); bid != "" {
		v, err := strconv.ParseUint(bid, 10, 64)
		if err == nil {
			id := uint(v)
			baselineID = &id
		}
	}

	data, err := r.deviationService.ExportDeviations(c.Request.Context(), deviceID, baselineID)
	if err != nil {
		DatabaseError(c, "导出偏差报告失败")
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=deviations.csv")
	c.Data(http.StatusOK, "text/csv", data)
}
