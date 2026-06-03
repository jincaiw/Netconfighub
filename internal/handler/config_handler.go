package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) getConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		BadRequest(c, "缺少 key 参数")
		return
	}
	value, err := r.configService.Get(c.Request.Context(), key)
	if err != nil {
		Error(c, http.StatusNotFound, 404, "配置项不存在")
		return
	}
	Success(c, gin.H{"key": key, "value": value})
}

func (r *Router) setConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数错误")
		return
	}
	if err := r.configService.Set(c.Request.Context(), req.Key, req.Value); err != nil {
		DatabaseError(c, "保存配置失败")
		return
	}
	Success(c, nil)
}

func (r *Router) listConfigs(c *gin.Context) {
	configs, err := r.configService.List(c.Request.Context())
	if err != nil {
		DatabaseError(c, "查询配置失败")
		return
	}
	Success(c, configs)
}
