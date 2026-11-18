package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// SFTPAccountService manages the per-device credentials used by devices that
// push their configuration over SFTP. Plaintext passwords are generated,
// returned exactly once and never persisted or logged.
type SFTPAccountService interface {
	// Create issues a new account for a device and returns the plaintext
	// password once. Existing accounts for the device are revoked first.
	Create(ctx context.Context, deviceID uint) (*SFTPAccountCredentials, error)
	// Reset replaces the password of an existing account and returns the new
	// plaintext once.
	Reset(ctx context.Context, deviceID uint) (*SFTPAccountCredentials, error)
	// Enable and Disable toggle whether an account may authenticate.
	Enable(ctx context.Context, deviceID uint) error
	Disable(ctx context.Context, deviceID uint) error
	// Revoke soft-deletes the account so it can never authenticate again.
	Revoke(ctx context.Context, deviceID uint) error
	// Get returns account metadata without any credential material.
	Get(ctx context.Context, deviceID uint) (*model.DeviceSFTPAccount, error)
	// Authenticate verifies SFTP credentials and returns the bound device.
	Authenticate(ctx context.Context, username, password string) (*model.Device, *model.DeviceSFTPAccount, error)
}

// SFTPAccountCredentials is returned only when a password is created or reset.
type SFTPAccountCredentials struct {
	Account  *model.DeviceSFTPAccount `json:"account"`
	Password string                   `json:"password"`
}

type sftpAccountServiceImpl struct {
	accountRepo repository.DeviceSFTPAccountRepository
	deviceRepo  repository.DeviceRepository
	audit       AuditService
}

func NewSFTPAccountService(
	accountRepo repository.DeviceSFTPAccountRepository,
	deviceRepo repository.DeviceRepository,
	audit AuditService,
) SFTPAccountService {
	return &sftpAccountServiceImpl{
		accountRepo: accountRepo,
		deviceRepo:  deviceRepo,
		audit:       audit,
	}
}

const (
	sftpPasswordBytes = 18 // 36 hex characters
	sftpUsernameTries = 3
)

func (s *sftpAccountServiceImpl) Create(ctx context.Context, deviceID uint) (*SFTPAccountCredentials, error) {
	device, err := s.requirePushDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	password, err := randomPassword()
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成口令哈希失败: %w", err)
	}

	// Revoking first keeps "one active account per device" true even when an
	// account already exists.
	if existing, findErr := s.accountRepo.FindByDeviceID(ctx, device.ID); findErr == nil && existing != nil {
		if revokeErr := s.accountRepo.Revoke(ctx, device.ID); revokeErr != nil {
			return nil, fmt.Errorf("撤销旧账号失败: %w", revokeErr)
		}
	}

	var created *model.DeviceSFTPAccount
	for attempt := 0; attempt < sftpUsernameTries; attempt++ {
		username, genErr := randomUsername(device.ID)
		if genErr != nil {
			return nil, genErr
		}
		account := &model.DeviceSFTPAccount{
			DeviceID:     device.ID,
			Username:     username,
			PasswordHash: string(hash),
			Status:       model.SFTPAccountStatusActive,
		}
		if createErr := s.accountRepo.Create(ctx, account); createErr != nil {
			if isUniqueViolation(createErr) {
				continue
			}
			return nil, fmt.Errorf("创建 SFTP 账号失败: %w", createErr)
		}
		created = account
		break
	}
	if created == nil {
		return nil, fmt.Errorf("创建 SFTP 账号失败: 用户名冲突")
	}

	s.auditLog(ctx, "sftp_account_create", device.ID, fmt.Sprintf("为设备 %s 创建 SFTP 上传账号 %s", device.Name, created.Username))

	return &SFTPAccountCredentials{Account: created, Password: password}, nil
}

func (s *sftpAccountServiceImpl) Reset(ctx context.Context, deviceID uint) (*SFTPAccountCredentials, error) {
	device, err := s.requirePushDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	account, err := s.accountRepo.FindByDeviceID(ctx, device.ID)
	if err != nil || account == nil {
		return nil, fmt.Errorf("设备尚未创建 SFTP 账号")
	}

	password, err := randomPassword()
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成口令哈希失败: %w", err)
	}

	account.PasswordHash = string(hash)
	account.Status = model.SFTPAccountStatusActive
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("重置 SFTP 口令失败: %w", err)
	}

	s.auditLog(ctx, "sftp_account_reset", device.ID, fmt.Sprintf("重置设备 %s 的 SFTP 上传口令", device.Name))
	return &SFTPAccountCredentials{Account: account, Password: password}, nil
}

func (s *sftpAccountServiceImpl) Enable(ctx context.Context, deviceID uint) error {
	return s.setStatus(ctx, deviceID, model.SFTPAccountStatusActive, "sftp_account_enable")
}

func (s *sftpAccountServiceImpl) Disable(ctx context.Context, deviceID uint) error {
	return s.setStatus(ctx, deviceID, model.SFTPAccountStatusDisabled, "sftp_account_disable")
}

func (s *sftpAccountServiceImpl) setStatus(ctx context.Context, deviceID uint, status, action string) error {
	account, err := s.accountRepo.FindByDeviceID(ctx, deviceID)
	if err != nil || account == nil {
		return fmt.Errorf("设备尚未创建 SFTP 账号")
	}
	account.Status = status
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return fmt.Errorf("更新账号状态失败: %w", err)
	}
	s.auditLog(ctx, action, deviceID, fmt.Sprintf("账号 %s 状态改为 %s", account.Username, status))
	return nil
}

func (s *sftpAccountServiceImpl) Revoke(ctx context.Context, deviceID uint) error {
	account, err := s.accountRepo.FindByDeviceID(ctx, deviceID)
	if err != nil || account == nil {
		return fmt.Errorf("设备尚未创建 SFTP 账号")
	}
	if err := s.accountRepo.Revoke(ctx, deviceID); err != nil {
		return fmt.Errorf("撤销账号失败: %w", err)
	}
	s.auditLog(ctx, "sftp_account_revoke", deviceID, fmt.Sprintf("撤销账号 %s", account.Username))
	return nil
}

func (s *sftpAccountServiceImpl) Get(ctx context.Context, deviceID uint) (*model.DeviceSFTPAccount, error) {
	account, err := s.accountRepo.FindByDeviceID(ctx, deviceID)
	if err != nil || account == nil {
		return nil, fmt.Errorf("设备尚未创建 SFTP 账号")
	}
	// Never leak the hash to callers outside the service.
	account.PasswordHash = ""
	return account, nil
}

func (s *sftpAccountServiceImpl) Authenticate(ctx context.Context, username, password string) (*model.Device, *model.DeviceSFTPAccount, error) {
	// A single generic error avoids revealing whether the username exists.
	denied := errors.New("认证失败")

	if strings.TrimSpace(username) == "" || password == "" {
		return nil, nil, denied
	}

	account, err := s.accountRepo.FindByUsername(ctx, username)
	if err != nil || account == nil {
		return nil, nil, denied
	}
	if !account.IsActive() {
		return nil, nil, denied
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		return nil, nil, denied
	}

	device, err := s.deviceRepo.FindByID(ctx, account.DeviceID)
	if err != nil || device == nil {
		return nil, nil, denied
	}
	if !device.Enabled {
		return nil, nil, denied
	}
	if device.CollectionMode != model.CollectionModeSFTPPush {
		return nil, nil, denied
	}

	return device, account, nil
}

func (s *sftpAccountServiceImpl) requirePushDevice(ctx context.Context, deviceID uint) (*model.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil || device == nil {
		return nil, fmt.Errorf("设备不存在")
	}
	if device.CollectionMode != model.CollectionModeSFTPPush {
		return nil, fmt.Errorf("设备采集方式不是 sftp_push，无法创建上传账号")
	}
	return device, nil
}

func (s *sftpAccountServiceImpl) auditLog(ctx context.Context, action string, deviceID uint, detail string) {
	if s.audit == nil {
		return
	}
	// detail never contains the password: only usernames and state changes.
	if err := s.audit.Log(ctx, 0, "system", action, "device", deviceID, detail, ""); err != nil {
		_ = err
	}
}

func randomPassword() (string, error) {
	buf := make([]byte, sftpPasswordBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("生成随机口令失败: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func randomUsername(deviceID uint) (string, error) {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("生成随机用户名失败: %w", err)
	}
	return fmt.Sprintf("sftp_%d_%s", deviceID, hex.EncodeToString(buf)), nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
