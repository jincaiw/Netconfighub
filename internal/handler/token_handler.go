package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type createTokenRequest struct {
	Name      string      `json:"name" binding:"required"`
	ExpiresAt interface{} `json:"expires_at"`
}

type createTokenResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func (r *Router) listTokens(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	tokens, err := r.tokenService.List(c.Request.Context(), adminID.(uint))
	if err != nil {
		DatabaseError(c, "查询 Token 列表失败")
		return
	}
	Success(c, tokens)
}

func (r *Router) createToken(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	var req createTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	token, rawToken, err := r.tokenService.Create(c.Request.Context(), adminID.(uint), req.Name, req.ExpiresAt)
	if err != nil {
		DatabaseError(c, "创建 Token 失败")
		return
	}

	resp := createTokenResponse{
		ID:    token.ID,
		Name:  token.Name,
		Token: rawToken,
	}
	if token.ExpiresAt != nil {
		resp.ExpiresAt = token.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
	}

	r.auditLog(c, "create", "token", token.ID, "创建 Token: "+req.Name)

	Success(c, resp)
}

func (r *Router) deleteToken(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.tokenService.Delete(c.Request.Context(), uint(id), adminID.(uint)); err != nil {
		DatabaseError(c, "删除 Token 失败")
		return
	}

	r.auditLog(c, "delete", "token", uint(id), "删除 Token")

	Success(c, nil)
}
