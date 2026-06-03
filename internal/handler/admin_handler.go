package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func (r *Router) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "用户名和密码不能为空")
		return
	}

	token, expiresAt, err := r.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	Success(c, loginResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (r *Router) logout(c *gin.Context) {
	_ = r.authService.Logout(c.Request.Context())
	Success(c, nil)
}

func (r *Router) refreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		Unauthorized(c, "缺少 Authorization 头部")
		return
	}

	token, expiresAt, err := r.authService.RefreshToken(c.Request.Context(), extractBearerToken(authHeader))
	if err != nil {
		Unauthorized(c, err.Error())
		return
	}

	Success(c, loginResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func extractBearerToken(authHeader string) string {
	parts := []byte(authHeader)
	for i, b := range parts {
		if b == ' ' {
			return string(parts[i+1:])
		}
	}
	return authHeader
}

type createAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type updateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

func (r *Router) listAdmins(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	admins, total, err := r.adminService.List(c.Request.Context(), offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询管理员列表失败")
		return
	}
	PageResult(c, admins, total, page, pageSize)
}

func (r *Router) createAdmin(c *gin.Context) {
	var req createAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	admin, err := r.adminService.Create(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		DatabaseError(c, "创建管理员失败")
		return
	}

	r.auditLog(c, "create", "admin", admin.ID, "创建管理员: "+req.Username)

	Success(c, admin)
}

func (r *Router) getAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	admin, err := r.adminService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "管理员不存在")
		return
	}
	Success(c, admin)
}

func (r *Router) updateAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	var req updateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	var password *string
	if req.Password != "" {
		password = &req.Password
	}

	admin, err := r.adminService.Update(c.Request.Context(), uint(id), req.Username, password)
	if err != nil {
		DatabaseError(c, "更新管理员失败")
		return
	}

	r.auditLog(c, "update", "admin", uint(id), "更新管理员: "+req.Username)

	Success(c, admin)
}

func (r *Router) deleteAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.adminService.Delete(c.Request.Context(), uint(id)); err != nil {
		DatabaseError(c, "删除管理员失败")
		return
	}

	r.auditLog(c, "delete", "admin", uint(id), "删除管理员")

	Success(c, nil)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (r *Router) changePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效，新密码至少6位")
		return
	}

	adminID, exists := c.Get("admin_id")
	if !exists {
		Unauthorized(c, "未认证")
		return
	}

	if err := r.authService.ChangePassword(c.Request.Context(), adminID.(uint), req.OldPassword, req.NewPassword); err != nil {
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, "change_password", "admin", adminID.(uint), "修改密码")

	Success(c, nil)
}
