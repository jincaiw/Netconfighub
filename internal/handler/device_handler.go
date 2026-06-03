package handler

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/model"
)

type createDeviceRequest struct {
	Name           string `json:"name" binding:"required"`
	IP             string `json:"ip"`
	IPAddress      string `json:"ip_address"`
	Vendor         string `json:"vendor" binding:"required"`
	Model          string `json:"model" binding:"required"`
	Protocol       string `json:"protocol"`
	ConnProtocol   string `json:"conn_protocol"`
	Port           int    `json:"port" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	SSHKey         string `json:"ssh_key"`
	EnablePassword string `json:"enable_password"`
	BackupInterval string `json:"backup_interval"`
	GroupID        *uint  `json:"group_id"`
	Enabled        *bool  `json:"enabled"`
}

type updateDeviceRequest struct {
	Name           *string `json:"name"`
	IP             *string `json:"ip"`
	Vendor         *string `json:"vendor"`
	Model          *string `json:"model"`
	Protocol       *string `json:"protocol"`
	Port           *int    `json:"port"`
	Username       *string `json:"username"`
	Password       *string `json:"password"`
	SSHKey         *string `json:"ssh_key"`
	EnablePassword *string `json:"enable_password"`
	BackupInterval *string `json:"backup_interval"`
	GroupID        *uint   `json:"group_id"`
	Enabled        *bool   `json:"enabled"`
}

func (r *Router) listDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var groupID *uint
	if gid := c.Query("group_id"); gid != "" {
		v, err := strconv.ParseUint(gid, 10, 64)
		if err == nil {
			id := uint(v)
			groupID = &id
		}
	}

	var vendor *string
	if v := c.Query("vendor"); v != "" {
		vendor = &v
	}

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	var search *string
	if q := c.Query("search"); q != "" {
		search = &q
	}

	devices, total, err := r.deviceService.List(c.Request.Context(), groupID, vendor, status, search, offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询设备列表失败")
		return
	}
	PageResult(c, devices, total, page, pageSize)
}

func (r *Router) createDevice(c *gin.Context) {
	var req createDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	ip := strings.TrimSpace(req.IP)
	if ip == "" {
		ip = strings.TrimSpace(req.IPAddress)
	}
	if net.ParseIP(ip) == nil {
		InvalidIP(c, "IP 格式错误")
		return
	}
	protocol := strings.TrimSpace(req.Protocol)
	if protocol == "" {
		protocol = strings.TrimSpace(req.ConnProtocol)
	}
	if protocol == "" {
		BadRequest(c, "连接协议不能为空")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	device := &model.Device{
		Name:           req.Name,
		IP:             ip,
		Vendor:         model.Vendor(req.Vendor),
		Model:          model.DeviceModel(req.Model),
		Protocol:       model.ConnProtocol(protocol),
		Port:           req.Port,
		Username:       req.Username,
		Password:       req.Password,
		SSHKey:         req.SSHKey,
		EnablePassword: req.EnablePassword,
		BackupInterval: req.BackupInterval,
		GroupID:        req.GroupID,
		Enabled:        enabled,
	}

	result, err := r.deviceService.Create(c.Request.Context(), device)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "设备名称已存在") {
			DuplicateName(c, "设备名称已存在")
			return
		}
		if strings.Contains(errMsg, "设备 IP 已存在") {
			InvalidIP(c, "设备 IP 已存在")
			return
		}
		BadRequest(c, errMsg)
		return
	}

	r.auditLog(c, "create", "device", result.ID, "创建设备: "+req.Name)

	Success(c, result)
}

func (r *Router) getDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	device, err := r.deviceService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}
	Success(c, device)
}

func (r *Router) updateDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	var req updateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}

	device, err := r.deviceService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		BadRequest(c, "设备不存在")
		return
	}

	if req.Name != nil {
		device.Name = *req.Name
	}
	if req.IP != nil {
		if net.ParseIP(*req.IP) == nil {
			InvalidIP(c, "IP 格式错误")
			return
		}
		device.IP = *req.IP
	}
	if req.Vendor != nil {
		device.Vendor = model.Vendor(*req.Vendor)
	}
	if req.Model != nil {
		device.Model = model.DeviceModel(*req.Model)
	}
	if req.Protocol != nil {
		device.Protocol = model.ConnProtocol(*req.Protocol)
	}
	if req.Port != nil {
		device.Port = *req.Port
	}
	if req.Username != nil {
		device.Username = *req.Username
	}
	if req.Password != nil {
		device.Password = *req.Password
	}
	if req.SSHKey != nil {
		device.SSHKey = *req.SSHKey
	}
	if req.EnablePassword != nil {
		device.EnablePassword = *req.EnablePassword
	}
	if req.BackupInterval != nil {
		device.BackupInterval = *req.BackupInterval
	}
	if req.GroupID != nil {
		if *req.GroupID == 0 {
			device.GroupID = nil
		} else {
			device.GroupID = req.GroupID
		}
	}
	if req.Enabled != nil {
		device.Enabled = *req.Enabled
	}

	result, err := r.deviceService.Update(c.Request.Context(), uint(id), device)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "设备名称已存在") {
			DuplicateName(c, "设备名称已存在")
			return
		}
		if strings.Contains(errMsg, "设备 IP 已存在") {
			InvalidIP(c, "设备 IP 已存在")
			return
		}
		BadRequest(c, errMsg)
		return
	}

	r.auditLog(c, "update", "device", uint(id), "更新设备: "+device.Name)

	Success(c, result)
}

func (r *Router) deleteDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	device, err := r.deviceService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}

	if err := r.deviceService.Delete(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}

	r.auditLog(c, "delete", "device", uint(id), "删除设备: "+device.Name)

	Success(c, nil)
}

func (r *Router) enableDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.deviceService.EnableDevice(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}

	r.auditLog(c, "enable", "device", uint(id), "启用设备")

	Success(c, nil)
}

func (r *Router) disableDevice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if err := r.deviceService.DisableDevice(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}

	r.auditLog(c, "disable", "device", uint(id), "禁用设备")

	Success(c, nil)
}

func (r *Router) triggerBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	if _, err := r.deviceService.GetByID(c.Request.Context(), uint(id)); err != nil {
		Error(c, http.StatusNotFound, 404, "设备不存在")
		return
	}

	if r.scheduler == nil {
		Error(c, http.StatusInternalServerError, 500, "调度器未初始化")
		return
	}

	if err := r.scheduler.SubmitDevice(uint(id)); err != nil {
		BadRequest(c, "提交备份任务失败")
		return
	}

	r.auditLog(c, "trigger_backup", "device", uint(id), "触发设备备份")

	Success(c, gin.H{"message": "备份任务已提交"})
}

func (r *Router) getDeviceConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	content, err := r.deviceService.GetDeviceConfig(c.Request.Context(), uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "设备不存在") {
			Error(c, http.StatusNotFound, 404, "设备不存在")
			return
		}
		Error(c, http.StatusNotFound, 404, "暂无配置数据")
		return
	}
	Success(c, gin.H{"content": content})
}

func (r *Router) getDeviceVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	versions, err := r.deviceService.GetDeviceVersions(c.Request.Context(), uint(id), limit)
	if err != nil {
		if strings.Contains(err.Error(), "设备不存在") {
			Error(c, http.StatusNotFound, 404, "设备不存在")
			return
		}
		Error(c, http.StatusNotFound, 404, "暂无配置版本")
		return
	}
	Success(c, versions)
}

func (r *Router) getDeviceDeviations(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的 ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	deviations, total, err := r.deviceService.GetDeviceDeviations(c.Request.Context(), uint(id), offset, pageSize)
	if err != nil {
		Error(c, http.StatusNotFound, 404, "无法获取偏差信息")
		return
	}
	PageResult(c, deviations, total, page, pageSize)
}

func (r *Router) listFailedDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	devices, total, err := r.deviceService.ListFailedDevices(c.Request.Context(), offset, pageSize)
	if err != nil {
		DatabaseError(c, "查询失败设备列表失败")
		return
	}
	PageResult(c, devices, total, page, pageSize)
}

func (r *Router) getConfigDiff(c *gin.Context) {
	id := c.Param("id")
	fromHash := c.Query("from")
	toHash := c.Query("to")
	if fromHash == "" || toHash == "" {
		BadRequest(c, "缺少 from 或 to 参数")
		return
	}
	result, err := r.deviceService.GetConfigDiff(c.Request.Context(), id, fromHash, toHash)
	if err != nil {
		if strings.Contains(err.Error(), "设备不存在") {
			Error(c, http.StatusNotFound, 404, "设备不存在")
			return
		}
		Success(c, nil)
		return
	}
	Success(c, result)
}

func (r *Router) getLatestDiff(c *gin.Context) {
	id := c.Param("id")
	result, err := r.deviceService.GetLatestDiff(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "设备不存在") {
			Error(c, http.StatusNotFound, 404, "设备不存在")
			return
		}
		Success(c, nil)
		return
	}
	Success(c, result)
}

func (r *Router) downloadDiff(c *gin.Context) {
	id := c.Param("id")
	fromHash := c.Query("from")
	toHash := c.Query("to")
	if fromHash == "" || toHash == "" {
		BadRequest(c, "缺少 from 或 to 参数")
		return
	}
	result, err := r.deviceService.GetConfigDiff(c.Request.Context(), id, fromHash, toHash)
	if err != nil {
		if strings.Contains(err.Error(), "设备不存在") {
			Error(c, http.StatusNotFound, 404, "设备不存在")
			return
		}
		c.Header("Content-Disposition", "attachment; filename=config-diff.txt")
		c.String(http.StatusOK, "")
		return
	}
	diffContent := ""
	if result != nil {
		diffContent = result.DiffContent
	}
	c.Header("Content-Disposition", "attachment; filename=config-diff.txt")
	c.String(http.StatusOK, diffContent)
}

func (r *Router) importDevices(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		BadRequest(c, "请上传 CSV 文件")
		return
	}
	src, err := file.Open()
	if err != nil {
		DatabaseError(c, "读取上传文件失败")
		return
	}
	defer src.Close()
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		DatabaseError(c, "读取文件内容失败")
		return
	}
	result, err := r.deviceService.ImportDevices(c.Request.Context(), data)
	if err != nil {
		BadRequest(c, "导入设备失败")
		return
	}
	Success(c, result)
}

func (r *Router) exportDevices(c *gin.Context) {
	data, err := r.deviceService.ExportDevices(c.Request.Context())
	if err != nil {
		DatabaseError(c, "导出设备失败")
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=devices.csv")
	c.Data(http.StatusOK, "text/csv", data)
}

func (r *Router) getAuditUser(c *gin.Context) (uint, string) {
	userID := uint(0)
	username := ""
	if v, exists := c.Get("admin_id"); exists {
		if id, ok := v.(uint); ok {
			userID = id
		}
	}
	if v, exists := c.Get("username"); exists {
		if name, ok := v.(string); ok {
			username = name
		}
	}
	return userID, username
}

func (r *Router) auditLog(c *gin.Context, action, targetType string, targetID uint, detail string) {
	if r.auditService == nil {
		return
	}
	userID, username := r.getAuditUser(c)
	r.auditService.Log(c.Request.Context(), userID, username, action, targetType, targetID, detail, c.ClientIP())
}
