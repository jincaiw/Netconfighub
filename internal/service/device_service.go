package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/crypto"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/sanitizer"
)

type deviceServiceImpl struct {
	deviceRepo     repository.DeviceRepository
	groupRepo      repository.GroupRepository
	backupTaskRepo repository.BackupTaskRepository
	gitStore       *gitstore.Store
	deviationRepo  repository.DeviationRepository
	baselineEngine baseline.Engine
	sanitizer      *sanitizer.Sanitizer
	encryptionKey  string
}

func NewDeviceService(
	deviceRepo repository.DeviceRepository,
	groupRepo repository.GroupRepository,
	backupTaskRepo repository.BackupTaskRepository,
	gitStore *gitstore.Store,
	deviationRepo repository.DeviationRepository,
	baselineEngine baseline.Engine,
	san *sanitizer.Sanitizer,
	encryptionKey string,
) DeviceService {
	return &deviceServiceImpl{
		deviceRepo:     deviceRepo,
		groupRepo:      groupRepo,
		backupTaskRepo: backupTaskRepo,
		gitStore:       gitStore,
		deviationRepo:  deviationRepo,
		baselineEngine: baselineEngine,
		sanitizer:      san,
		encryptionKey:  encryptionKey,
	}
}

func (s *deviceServiceImpl) Create(ctx context.Context, device *model.Device) (*model.Device, error) {
	if err := validateStorageName("设备名称", device.Name); err != nil {
		return nil, err
	}
	if !device.Vendor.IsValid() {
		return nil, fmt.Errorf("无效的厂商类型，支持: cisco, h3c, huawei, ruijie")
	}
	if !device.Protocol.IsValid() {
		return nil, fmt.Errorf("无效的连接协议，支持: ssh, telnet")
	}
	if device.Port < 1 || device.Port > 65535 {
		return nil, fmt.Errorf("无效的端口号，范围: 1-65535")
	}
	if existing, err := s.deviceRepo.FindByName(ctx, device.Name); err == nil && existing != nil {
		return nil, fmt.Errorf("设备名称已存在")
	}
	if existing, err := s.deviceRepo.FindByIP(ctx, device.IP); err == nil && existing != nil {
		return nil, fmt.Errorf("设备 IP 已存在")
	}
	if device.GroupID != nil {
		if _, err := s.groupRepo.FindByID(ctx, *device.GroupID); err != nil {
			return nil, fmt.Errorf("指定的分组不存在")
		}
	}
	device.LastBackupStatus = "never"
	if device.Password != "" {
		encrypted, err := crypto.Encrypt(device.Password, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("加密密码失败: %w", err)
		}
		device.Password = encrypted
	}
	if device.SSHKey != "" {
		encrypted, err := crypto.Encrypt(device.SSHKey, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("加密SSH密钥失败: %w", err)
		}
		device.SSHKey = encrypted
	}
	if device.EnablePassword != "" {
		encrypted, err := crypto.Encrypt(device.EnablePassword, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("加密Enable密码失败: %w", err)
		}
		device.EnablePassword = encrypted
	}
	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, err
	}
	result, err := s.deviceRepo.FindByID(ctx, device.ID)
	if err != nil {
		return nil, err
	}
	s.maskDeviceCredentials(result)
	populateDeviceGroupName(result)
	return result, nil
}

func (s *deviceServiceImpl) GetByID(ctx context.Context, id uint) (*model.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("设备不存在")
	}
	s.maskDeviceCredentials(device)
	populateDeviceGroupName(device)
	return device, nil
}

func (s *deviceServiceImpl) List(ctx context.Context, groupID *uint, vendor *string, status *string, search *string, offset, limit int) ([]model.Device, int64, error) {
	devices, total, err := s.deviceRepo.List(ctx, groupID, vendor, status, search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for i := range devices {
		s.maskDeviceCredentials(&devices[i])
		populateDeviceGroupName(&devices[i])
	}
	return devices, total, nil
}

func (s *deviceServiceImpl) Update(ctx context.Context, id uint, device *model.Device) (*model.Device, error) {
	if err := validateStorageName("设备名称", device.Name); err != nil {
		return nil, err
	}
	existing, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("设备不存在")
	}
	if device.Name != existing.Name {
		if dup, err := s.deviceRepo.FindByName(ctx, device.Name); err == nil && dup != nil {
			return nil, fmt.Errorf("设备名称已存在")
		}
	}
	if device.IP != existing.IP {
		if dup, err := s.deviceRepo.FindByIP(ctx, device.IP); err == nil && dup != nil {
			return nil, fmt.Errorf("设备 IP 已存在")
		}
	}
	existing.Name = device.Name
	existing.IP = device.IP
	existing.Vendor = device.Vendor
	existing.Model = device.Model
	existing.Protocol = device.Protocol
	existing.Port = device.Port
	existing.Username = device.Username
	if device.Password != "" {
		encrypted, err := crypto.Encrypt(device.Password, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("加密密码失败: %w", err)
		}
		existing.Password = encrypted
	}
	if device.SSHKey != "" {
		encrypted, err := crypto.Encrypt(device.SSHKey, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("加密SSH密钥失败: %w", err)
		}
		existing.SSHKey = encrypted
	}
	if device.EnablePassword != "" {
		encrypted, err := crypto.Encrypt(device.EnablePassword, s.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("加密Enable密码失败: %w", err)
		}
		existing.EnablePassword = encrypted
	}
	existing.GroupID = device.GroupID
	existing.BackupInterval = device.BackupInterval
	existing.Enabled = device.Enabled

	if !existing.Vendor.IsValid() {
		return nil, fmt.Errorf("无效的厂商类型，支持: cisco, h3c, huawei, ruijie")
	}
	if !existing.Protocol.IsValid() {
		return nil, fmt.Errorf("无效的连接协议，支持: ssh, telnet")
	}
	if existing.Port < 1 || existing.Port > 65535 {
		return nil, fmt.Errorf("无效的端口号，范围: 1-65535")
	}
	if existing.GroupID != nil {
		if _, err := s.groupRepo.FindByID(ctx, *existing.GroupID); err != nil {
			return nil, fmt.Errorf("指定的分组不存在")
		}
	}

	existing.Group = nil
	if err := s.deviceRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	result, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.maskDeviceCredentials(result)
	populateDeviceGroupName(result)
	return result, nil
}

func (s *deviceServiceImpl) Delete(ctx context.Context, id uint) error {
	if _, err := s.deviceRepo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("设备不存在")
	}
	return s.deviceRepo.Delete(ctx, id)
}

func (s *deviceServiceImpl) TriggerBackup(ctx context.Context, deviceID uint) (*model.BackupTask, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("设备不存在")
	}
	if !device.Enabled {
		return nil, fmt.Errorf("设备已禁用，无法触发备份")
	}
	task := &model.BackupTask{
		DeviceID:   deviceID,
		Status:     model.BackupStatusPending,
		MaxRetries: 3,
	}
	if err := s.backupTaskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	device.LastBackupStatus = string(model.BackupStatusPending)
	device.Group = nil
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *deviceServiceImpl) GetDeviceConfig(ctx context.Context, deviceID uint) (string, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return "", fmt.Errorf("设备不存在")
	}
	if s.gitStore == nil {
		return "", fmt.Errorf("Git 存储未初始化")
	}
	content, err := s.gitStore.ReadLatestConfig(string(device.Vendor), device.Name)
	if err != nil {
		return "", fmt.Errorf("读取设备配置失败: %w", err)
	}
	if s.sanitizer != nil {
		content = s.sanitizer.Sanitize(content)
	}
	return content, nil
}

func (s *deviceServiceImpl) GetDeviceVersions(ctx context.Context, deviceID uint, limit int) ([]model.ConfigVersion, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("设备不存在")
	}
	if s.gitStore == nil {
		return nil, fmt.Errorf("Git 存储未初始化")
	}
	commits, err := s.gitStore.ConfigHistory(string(device.Vendor), device.Name, limit)
	if err != nil {
		return nil, fmt.Errorf("查询配置版本失败: %w", err)
	}
	versions := make([]model.ConfigVersion, 0, len(commits))
	for _, c := range commits {
		versions = append(versions, model.ConfigVersion{
			Hash:    c.Hash.String(),
			Message: c.Message,
			Author:  c.Author.Name,
			Date:    c.Author.When,
		})
	}
	return versions, nil
}

func (s *deviceServiceImpl) GetDeviceDeviations(ctx context.Context, deviceID uint, offset, limit int) ([]model.Deviation, int64, error) {
	if _, err := s.deviceRepo.FindByID(ctx, deviceID); err != nil {
		return nil, 0, fmt.Errorf("设备不存在")
	}
	return s.deviationRepo.FindByDeviceID(ctx, deviceID, offset, limit)
}

func (s *deviceServiceImpl) ListFailedDevices(ctx context.Context, offset, limit int) ([]model.Device, int64, error) {
	devices, total, err := s.deviceRepo.ListFailed(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for i := range devices {
		s.maskDeviceCredentials(&devices[i])
		populateDeviceGroupName(&devices[i])
	}
	return devices, total, nil
}

func populateDeviceGroupName(device *model.Device) {
	if device.Group != nil {
		device.GroupName = device.Group.Name
	}
}

func (s *deviceServiceImpl) maskDeviceCredentials(device *model.Device) {
	if device.Password != "" {
		decrypted, err := crypto.Decrypt(device.Password, s.encryptionKey)
		if err == nil {
			_ = decrypted
		}
		device.Password = "******"
	}
	if device.SSHKey != "" {
		device.SSHKey = "******"
	}
	if device.EnablePassword != "" {
		device.EnablePassword = "******"
	}
}

func (s *deviceServiceImpl) DecryptDeviceCredentials(device *model.Device) {
	if device.Password != "" && device.Password != "******" {
		decrypted, err := crypto.Decrypt(device.Password, s.encryptionKey)
		if err == nil {
			device.Password = decrypted
		}
	}
	if device.SSHKey != "" && device.SSHKey != "******" {
		decrypted, err := crypto.Decrypt(device.SSHKey, s.encryptionKey)
		if err == nil {
			device.SSHKey = decrypted
		}
	}
	if device.EnablePassword != "" && device.EnablePassword != "******" {
		decrypted, err := crypto.Decrypt(device.EnablePassword, s.encryptionKey)
		if err == nil {
			device.EnablePassword = decrypted
		}
	}
}

func (s *deviceServiceImpl) GetConfigDiff(ctx context.Context, deviceID, fromHash, toHash string) (*model.DiffResult, error) {
	deviceIDInt, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的设备 ID")
	}
	device, err := s.deviceRepo.FindByID(ctx, uint(deviceIDInt))
	if err != nil {
		return nil, fmt.Errorf("设备不存在")
	}
	if s.gitStore == nil {
		return nil, fmt.Errorf("Git 存储未初始化")
	}
	fromContent, err := s.gitStore.ReadConfigAtRevision(string(device.Vendor), device.Name, fromHash)
	if err != nil {
		return nil, fmt.Errorf("读取源版本配置失败: %w", err)
	}
	toContent, err := s.gitStore.ReadConfigAtRevision(string(device.Vendor), device.Name, toHash)
	if err != nil {
		return nil, fmt.Errorf("读取目标版本配置失败: %w", err)
	}
	compareResult, err := s.baselineEngine.Compare(ctx, fromContent, toContent, 0)
	if err != nil {
		return nil, fmt.Errorf("比对失败: %w", err)
	}
	return &model.DiffResult{
		FromHash:     fromHash,
		ToHash:       toHash,
		DiffContent:  compareResult.DiffContent,
		AddedLines:   compareResult.AddedLines,
		RemovedLines: compareResult.RemovedLines,
	}, nil
}

func (s *deviceServiceImpl) GetLatestDiff(ctx context.Context, deviceID string) (*model.DiffResult, error) {
	deviceIDInt, err := strconv.ParseUint(deviceID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的设备 ID")
	}
	device, err := s.deviceRepo.FindByID(ctx, uint(deviceIDInt))
	if err != nil {
		return nil, fmt.Errorf("设备不存在")
	}
	if s.gitStore == nil {
		return nil, fmt.Errorf("Git 存储未初始化")
	}
	commits, err := s.gitStore.ConfigHistory(string(device.Vendor), device.Name, 2)
	if err != nil {
		return nil, fmt.Errorf("查询配置历史失败: %w", err)
	}
	if len(commits) < 2 {
		return nil, fmt.Errorf("配置历史不足两个版本，无法对比")
	}
	fromHash := commits[1].Hash.String()
	toHash := commits[0].Hash.String()
	return s.GetConfigDiff(ctx, deviceID, fromHash, toHash)
}

func (s *deviceServiceImpl) ImportDevices(ctx context.Context, csvData []byte) (*model.ImportResult, error) {
	reader := csv.NewReader(bytes.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("解析 CSV 失败: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV 文件为空或缺少数据行")
	}
	result := &model.ImportResult{Total: len(records) - 1}
	for i, record := range records[1:] {
		row := i + 2
		if len(record) < 9 {
			result.Failed++
			result.Errors = append(result.Errors, model.ImportError{
				Row:    row,
				Reason: "字段不足，需要 9 列",
			})
			continue
		}
		name := strings.TrimSpace(record[0])
		ip := strings.TrimSpace(record[1])
		vendor := strings.TrimSpace(record[2])
		modelVal := strings.TrimSpace(record[3])
		protocol := strings.TrimSpace(record[4])
		portStr := strings.TrimSpace(record[5])
		username := strings.TrimSpace(record[6])
		password := strings.TrimSpace(record[7])
		groupName := strings.TrimSpace(record[8])

		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			result.Failed++
			result.Errors = append(result.Errors, model.ImportError{
				Row:    row,
				Name:   name,
				Reason: fmt.Sprintf("无效端口号: %s", portStr),
			})
			continue
		}
		if !model.Vendor(vendor).IsValid() {
			result.Failed++
			result.Errors = append(result.Errors, model.ImportError{
				Row:    row,
				Name:   name,
				Reason: fmt.Sprintf("无效厂商: %s", vendor),
			})
			continue
		}
		if !model.ConnProtocol(protocol).IsValid() {
			result.Failed++
			result.Errors = append(result.Errors, model.ImportError{
				Row:    row,
				Name:   name,
				Reason: fmt.Sprintf("无效连接协议: %s", protocol),
			})
			continue
		}
		var groupID *uint
		if groupName != "" {
			groups, _, err := s.groupRepo.List(ctx, 0, 1000)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, model.ImportError{
					Row:    row,
					Name:   name,
					Reason: "查询分组失败",
				})
				continue
			}
			found := false
			for _, g := range groups {
				if g.Name == groupName {
					groupID = &g.ID
					found = true
					break
				}
			}
			if !found {
				result.Failed++
				result.Errors = append(result.Errors, model.ImportError{
					Row:    row,
					Name:   name,
					Reason: fmt.Sprintf("分组不存在: %s", groupName),
				})
				continue
			}
		}
		device := &model.Device{
			Name:             name,
			IP:               ip,
			Vendor:           model.Vendor(vendor),
			Model:            model.DeviceModel(modelVal),
			Protocol:         model.ConnProtocol(protocol),
			Port:             port,
			Username:         username,
			Password:         password,
			GroupID:          groupID,
			LastBackupStatus: "never",
		}
		if device.Password != "" {
			encrypted, err := crypto.Encrypt(device.Password, s.encryptionKey)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, model.ImportError{
					Row:    row,
					Name:   name,
					Reason: fmt.Sprintf("加密密码失败: %v", err),
				})
				continue
			}
			device.Password = encrypted
		}
		if device.EnablePassword != "" {
			encrypted, err := crypto.Encrypt(device.EnablePassword, s.encryptionKey)
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, model.ImportError{
					Row:    row,
					Name:   name,
					Reason: fmt.Sprintf("加密Enable密码失败: %v", err),
				})
				continue
			}
			device.EnablePassword = encrypted
		}
		if err := s.deviceRepo.Create(ctx, device); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, model.ImportError{
				Row:    row,
				Name:   name,
				Reason: err.Error(),
			})
			continue
		}
		result.Success++
	}
	return result, nil
}

func (s *deviceServiceImpl) ExportDevices(ctx context.Context) ([]byte, error) {
	devices, _, err := s.deviceRepo.List(ctx, nil, nil, nil, nil, 0, 10000)
	if err != nil {
		return nil, fmt.Errorf("查询设备列表失败: %w", err)
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"name", "ip", "vendor", "model", "protocol", "port", "username", "password", "group_name"})
	for _, d := range devices {
		groupName := ""
		if d.Group != nil {
			groupName = d.Group.Name
		}
		writer.Write([]string{
			d.Name,
			d.IP,
			string(d.Vendor),
			string(d.Model),
			string(d.Protocol),
			strconv.Itoa(d.Port),
			d.Username,
			"",
			groupName,
		})
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *deviceServiceImpl) EnableDevice(ctx context.Context, id uint) error {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("设备不存在")
	}
	device.Enabled = true
	device.Group = nil
	return s.deviceRepo.Update(ctx, device)
}

func (s *deviceServiceImpl) DisableDevice(ctx context.Context, id uint) error {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("设备不存在")
	}
	device.Enabled = false
	device.Group = nil
	return s.deviceRepo.Update(ctx, device)
}
