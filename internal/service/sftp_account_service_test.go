package service

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// memoryAccountRepo is an in-memory DeviceSFTPAccountRepository.
type memoryAccountRepo struct {
	accounts []*model.DeviceSFTPAccount
	nextID   uint
}

func (r *memoryAccountRepo) Create(_ context.Context, account *model.DeviceSFTPAccount) error {
	for _, existing := range r.accounts {
		if existing.Username == account.Username && existing.DeletedAt.Time.IsZero() {
			return fmt.Errorf("UNIQUE constraint failed: username")
		}
	}
	r.nextID++
	account.ID = r.nextID
	clone := *account
	r.accounts = append(r.accounts, &clone)
	return nil
}

func (r *memoryAccountRepo) FindByUsername(_ context.Context, username string) (*model.DeviceSFTPAccount, error) {
	for _, account := range r.accounts {
		if account.Username == username && account.DeletedAt.Time.IsZero() {
			clone := *account
			return &clone, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *memoryAccountRepo) FindByDeviceID(_ context.Context, deviceID uint) (*model.DeviceSFTPAccount, error) {
	for _, account := range r.accounts {
		if account.DeviceID == deviceID && account.DeletedAt.Time.IsZero() {
			clone := *account
			return &clone, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *memoryAccountRepo) Update(_ context.Context, account *model.DeviceSFTPAccount) error {
	for i, existing := range r.accounts {
		if existing.ID == account.ID {
			clone := *account
			r.accounts[i] = &clone
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func (r *memoryAccountRepo) Revoke(_ context.Context, deviceID uint) error {
	for i, account := range r.accounts {
		if account.DeviceID == deviceID && account.DeletedAt.Time.IsZero() {
			account.Status = model.SFTPAccountStatusDisabled
			account.DeletedAt.Valid = true
			_ = i
			return nil
		}
	}
	return nil
}

func (r *memoryAccountRepo) ListByDeviceID(_ context.Context, deviceID uint) ([]model.DeviceSFTPAccount, error) {
	var out []model.DeviceSFTPAccount
	for _, account := range r.accounts {
		if account.DeviceID == deviceID && account.DeletedAt.Time.IsZero() {
			out = append(out, *account)
		}
	}
	return out, nil
}

// memoryDeviceRepoForSFTP resolves devices by ID from a fixed set.
type memoryDeviceRepoForSFTP struct {
	devices map[uint]*model.Device
}

func (r *memoryDeviceRepoForSFTP) Create(context.Context, *model.Device) error { return nil }
func (r *memoryDeviceRepoForSFTP) FindByID(_ context.Context, id uint) (*model.Device, error) {
	if d, ok := r.devices[id]; ok {
		clone := *d
		return &clone, nil
	}
	return nil, fmt.Errorf("not found")
}
func (r *memoryDeviceRepoForSFTP) FindByName(context.Context, string) (*model.Device, error) {
	return nil, fmt.Errorf("not found")
}
func (r *memoryDeviceRepoForSFTP) FindByIP(context.Context, string) (*model.Device, error) {
	return nil, fmt.Errorf("not found")
}
func (r *memoryDeviceRepoForSFTP) List(context.Context, *uint, *string, *string, *string, int, int) ([]model.Device, int64, error) {
	return nil, 0, nil
}
func (r *memoryDeviceRepoForSFTP) Update(context.Context, *model.Device) error { return nil }
func (r *memoryDeviceRepoForSFTP) Delete(context.Context, uint) error          { return nil }
func (r *memoryDeviceRepoForSFTP) FindByGroupID(context.Context, uint) ([]model.Device, error) {
	return nil, nil
}
func (r *memoryDeviceRepoForSFTP) ListFailed(context.Context, int, int) ([]model.Device, int64, error) {
	return nil, 0, nil
}
func (r *memoryDeviceRepoForSFTP) UpdateBackupState(context.Context, uint, repository.BackupStateUpdate) error {
	return nil
}

// recordingAudit captures audit entries so tests can assert no password leaks.
type recordingAudit struct {
	details []string
}

func (a *recordingAudit) Log(_ context.Context, _ uint, _, action, _ string, _ uint, detail, _ string) error {
	a.details = append(a.details, action+"|"+detail)
	return nil
}

func (a *recordingAudit) List(context.Context, int, int) ([]model.AuditLog, int64, error) {
	return nil, 0, nil
}

func newSFTPAccountFixture(t *testing.T) (*sftpAccountServiceImpl, *memoryAccountRepo, *memoryDeviceRepoForSFTP, *recordingAudit) {
	t.Helper()
	accounts := &memoryAccountRepo{}
	devices := &memoryDeviceRepoForSFTP{devices: map[uint]*model.Device{
		1: {ID: 1, Name: "dev1", Enabled: true, CollectionMode: model.CollectionModeSFTPPush},
		2: {ID: 2, Name: "dev2", Enabled: true, CollectionMode: model.CollectionModePull},
	}}
	audit := &recordingAudit{}
	svc := NewSFTPAccountService(accounts, devices, audit).(*sftpAccountServiceImpl)
	return svc, accounts, devices, audit
}

func TestSFTPAccountCreateStoresOnlyHash(t *testing.T) {
	svc, accounts, _, audit := newSFTPAccountFixture(t)

	creds, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("创建账号失败: %v", err)
	}
	if creds.Password == "" {
		t.Fatalf("应返回一次性明文口令")
	}

	stored := accounts.accounts[0]
	if stored.PasswordHash == creds.Password {
		t.Fatalf("数据库中不得保存明文口令")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(creds.Password)); err != nil {
		t.Fatalf("哈希应能校验通过: %v", err)
	}

	// Metadata reads must never expose the hash.
	meta, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("读取账号失败: %v", err)
	}
	if meta.PasswordHash != "" {
		t.Fatalf("Get 不得返回口令哈希")
	}

	// Audit entries must not contain the password.
	for _, detail := range audit.details {
		if strings.Contains(detail, creds.Password) {
			t.Fatalf("审计详情泄露了口令: %s", detail)
		}
	}
}

func TestSFTPAccountResetInvalidatesOldPassword(t *testing.T) {
	svc, _, _, _ := newSFTPAccountFixture(t)

	creds, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("创建账号失败: %v", err)
	}

	reset, err := svc.Reset(context.Background(), 1)
	if err != nil {
		t.Fatalf("重置口令失败: %v", err)
	}
	if reset.Password == creds.Password {
		t.Fatalf("重置后口令必须变化")
	}

	if _, _, err := svc.Authenticate(context.Background(), creds.Account.Username, creds.Password); err == nil {
		t.Fatalf("旧口令必须失效")
	}
	if _, _, err := svc.Authenticate(context.Background(), reset.Account.Username, reset.Password); err != nil {
		t.Fatalf("新口令应通过认证: %v", err)
	}
}

func TestSFTPAccountDisabledAndRevokedCannotAuthenticate(t *testing.T) {
	svc, _, _, _ := newSFTPAccountFixture(t)

	creds, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("创建账号失败: %v", err)
	}

	if err := svc.Disable(context.Background(), 1); err != nil {
		t.Fatalf("禁用失败: %v", err)
	}
	if _, _, err := svc.Authenticate(context.Background(), creds.Account.Username, creds.Password); err == nil {
		t.Fatalf("禁用后不得认证")
	}

	if err := svc.Enable(context.Background(), 1); err != nil {
		t.Fatalf("启用失败: %v", err)
	}
	if _, _, err := svc.Authenticate(context.Background(), creds.Account.Username, creds.Password); err != nil {
		t.Fatalf("重新启用后应可认证: %v", err)
	}

	if err := svc.Revoke(context.Background(), 1); err != nil {
		t.Fatalf("撤销失败: %v", err)
	}
	if _, _, err := svc.Authenticate(context.Background(), creds.Account.Username, creds.Password); err == nil {
		t.Fatalf("撤销后不得认证")
	}
}

func TestSFTPAccountPullDeviceRejected(t *testing.T) {
	svc, _, _, _ := newSFTPAccountFixture(t)

	if _, err := svc.Create(context.Background(), 2); err == nil {
		t.Fatalf("pull 模式设备不得创建 SFTP 账号")
	}
}

func TestSFTPAccountAuthenticateRejectsCrossDevice(t *testing.T) {
	svc, accounts, devices, _ := newSFTPAccountFixture(t)

	credsA, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("创建设备1账号失败: %v", err)
	}

	// Add a second push device with its own account.
	devices.devices[3] = &model.Device{ID: 3, Name: "dev3", Enabled: true, CollectionMode: model.CollectionModeSFTPPush}
	credsB, err := svc.Create(context.Background(), 3)
	if err != nil {
		t.Fatalf("创建设备3账号失败: %v", err)
	}

	if credsA.Account.Username == credsB.Account.Username {
		t.Fatalf("不同设备应使用不同用户名")
	}

	device, _, err := svc.Authenticate(context.Background(), credsA.Account.Username, credsA.Password)
	if err != nil {
		t.Fatalf("设备1认证失败: %v", err)
	}
	if device.ID != 1 {
		t.Fatalf("账号必须绑定自身设备，实际 %d", device.ID)
	}

	// The wrong password for another account must not authenticate.
	if _, _, err := svc.Authenticate(context.Background(), credsB.Account.Username, credsA.Password); err == nil {
		t.Fatalf("跨设备口令不得认证成功")
	}
	_ = accounts
}

// TestSFTPAccountRecreateAfterRevoke guards the regression where the
// (device_id) unique index kept a soft-deleted row, so creating an account for
// a device that already had one destroyed the old account and failed.
func TestSFTPAccountRecreateAfterRevoke(t *testing.T) {
	svc, _, _, _ := newSFTPAccountFixture(t)

	first, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("首次创建账号失败: %v", err)
	}

	// Creating again while an account exists must succeed and rotate it.
	second, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("已有账号时再次创建应成功: %v", err)
	}
	if second.Password == "" {
		t.Fatalf("应返回新的一次性口令")
	}

	// The old password must no longer authenticate; the new one must.
	if _, _, err := svc.Authenticate(context.Background(), first.Account.Username, first.Password); err == nil {
		t.Fatalf("旧账号口令必须失效")
	}
	if _, _, err := svc.Authenticate(context.Background(), second.Account.Username, second.Password); err != nil {
		t.Fatalf("新账号应可认证: %v", err)
	}

	// Revoke then create again: this is the exact path that used to fail.
	if err := svc.Revoke(context.Background(), 1); err != nil {
		t.Fatalf("撤销账号失败: %v", err)
	}
	if _, err := svc.Create(context.Background(), 1); err != nil {
		t.Fatalf("撤销后应能重新创建账号: %v", err)
	}
}

func TestSFTPAccountDisabledDeviceCannotAuthenticate(t *testing.T) {
	svc, _, devices, _ := newSFTPAccountFixture(t)

	creds, err := svc.Create(context.Background(), 1)
	if err != nil {
		t.Fatalf("创建账号失败: %v", err)
	}

	devices.devices[1].Enabled = false
	if _, _, err := svc.Authenticate(context.Background(), creds.Account.Username, creds.Password); err == nil {
		t.Fatalf("设备禁用后不得认证")
	}

	// Switching the device back to pull must also stop SFTP authentication.
	devices.devices[1].Enabled = true
	devices.devices[1].CollectionMode = model.CollectionModePull
	if _, _, err := svc.Authenticate(context.Background(), creds.Account.Username, creds.Password); err == nil {
		t.Fatalf("设备切回 pull 后不得通过 SFTP 认证")
	}
}
