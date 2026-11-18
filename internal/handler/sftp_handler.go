package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/apperr"
	"github.com/netconfighub/netconfighub/internal/service"
)

// SFTPStatusProvider reports the runtime state of the inbound SFTP server so
// administrators can see whether it is listening and verify the host key
// fingerprint before configuring a device.
type SFTPStatusProvider interface {
	SFTPStatus() SFTPStatus
}

// SFTPStatus is the externally visible SFTP server state.
type SFTPStatus struct {
	Enabled     bool   `json:"enabled"`
	ListenAddr  string `json:"listen_addr"`
	Fingerprint string `json:"fingerprint"`
	ActiveConns int    `json:"active_conns"`
	MaxConns    int    `json:"max_conns"`
}

// sftpStatusHandler returns the SFTP server status.
func (r *Router) sftpStatusHandler(c *gin.Context) {
	if r.sftpStatus == nil {
		Success(c, SFTPStatus{Enabled: false})
		return
	}
	Success(c, r.sftpStatus.SFTPStatus())
}

// getSFTPAccount returns account metadata (never the password or hash).
func (r *Router) getSFTPAccount(c *gin.Context) {
	deviceID, ok := r.parseDeviceID(c)
	if !ok {
		return
	}
	account, err := r.sftpAccountService.Get(c.Request.Context(), deviceID)
	if err != nil {
		Error(c, http.StatusNotFound, 404, "设备尚未创建 SFTP 账号")
		return
	}
	Success(c, account)
}

// createSFTPAccount issues a new account and returns the plaintext password
// exactly once.
func (r *Router) createSFTPAccount(c *gin.Context) {
	deviceID, ok := r.parseDeviceID(c)
	if !ok {
		return
	}

	creds, err := r.sftpAccountService.Create(c.Request.Context(), deviceID)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, "sftp_account_create", "device", deviceID, "创建 SFTP 上传账号")
	// The password is only present in this response.
	Success(c, gin.H{
		"account":  creds.Account,
		"password": creds.Password,
	})
}

// resetSFTPAccount rotates the password and returns it once.
func (r *Router) resetSFTPAccount(c *gin.Context) {
	deviceID, ok := r.parseDeviceID(c)
	if !ok {
		return
	}

	creds, err := r.sftpAccountService.Reset(c.Request.Context(), deviceID)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, "sftp_account_reset", "device", deviceID, "重置 SFTP 上传口令")
	Success(c, gin.H{
		"account":  creds.Account,
		"password": creds.Password,
	})
}

func (r *Router) enableSFTPAccount(c *gin.Context) {
	r.setSFTPAccountStatus(c, true)
}

func (r *Router) disableSFTPAccount(c *gin.Context) {
	r.setSFTPAccountStatus(c, false)
}

func (r *Router) setSFTPAccountStatus(c *gin.Context, enable bool) {
	deviceID, ok := r.parseDeviceID(c)
	if !ok {
		return
	}

	var err error
	action := "sftp_account_disable"
	if enable {
		err = r.sftpAccountService.Enable(c.Request.Context(), deviceID)
		action = "sftp_account_enable"
	} else {
		err = r.sftpAccountService.Disable(c.Request.Context(), deviceID)
	}
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, action, "device", deviceID, "更新 SFTP 账号状态")
	Success(c, gin.H{"message": "账号状态已更新"})
}

// deleteSFTPAccount revokes the account so it can never authenticate again.
func (r *Router) deleteSFTPAccount(c *gin.Context) {
	deviceID, ok := r.parseDeviceID(c)
	if !ok {
		return
	}

	if err := r.sftpAccountService.Revoke(c.Request.Context(), deviceID); err != nil {
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, "sftp_account_revoke", "device", deviceID, "撤销 SFTP 上传账号")
	Success(c, gin.H{"message": "账号已撤销"})
}

// listSFTPUploads returns the upload history for one device.
func (r *Router) listSFTPUploads(c *gin.Context) {
	deviceID, ok := r.parseDeviceID(c)
	if !ok {
		return
	}
	if _, err := r.deviceService.GetByID(c.Request.Context(), deviceID); err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}

	page, pageSize := parsePagination(c)
	uploads, total, err := r.sftpUploadService.ListByDevice(c.Request.Context(), deviceID, (page-1)*pageSize, pageSize)
	if err != nil {
		DatabaseError(c, "查询上传记录失败")
		return
	}
	PageResult(c, uploads, total, page, pageSize)
}

// getSFTPUpload returns one upload record.
func (r *Router) getSFTPUpload(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	upload, err := r.sftpUploadService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "上传记录不存在")
		return
	}
	Success(c, upload)
}

// retrySFTPUpload re-queues a failed archive attempt.
func (r *Router) retrySFTPUpload(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.sftpUploadService.Retry(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, apperr.ErrUploadNotFound) {
			Error(c, http.StatusNotFound, 404, "上传记录不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}

	r.auditLog(c, "sftp_upload_retry", "sftp_upload", uint(id), "重试 SFTP 上传归档")
	Success(c, gin.H{"message": "已重新排队归档"})
}

func (r *Router) parseDeviceID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return 0, false
	}
	return uint(id), true
}

func parsePagination(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20
	if raw := c.Query("page"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			page = v
		}
	}
	if raw := c.Query("page_size"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 200 {
			pageSize = v
		}
	}
	return page, pageSize
}

var _ = service.SFTPAccountCredentials{}
