package service

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"gorm.io/gorm"
)

func setupBaselineTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.Admin{}, &model.APIToken{}, &model.Group{}, &model.Device{}, &model.BackupTask{}, &model.Baseline{}, &model.Deviation{}); err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}
	return db
}

func setupBaselineTestService(t *testing.T, db *gorm.DB) BaselineService {
	t.Helper()
	baselineRepo := repository.NewBaselineRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	deviationRepo := repository.NewDeviationRepository(db)
	return NewBaselineService(baselineRepo, deviceRepo, deviationRepo)
}

func TestCreateDeviceBaseline(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)

	deviceID := createTestDevice(t, deviceRepo)

	deviceIDCopy := deviceID
	baseline := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "hostname Router1\ninterface Gig0/1\n ip address 10.0.0.1\n",
	}

	result, err := svc.Create(context.Background(), baseline)
	if err != nil {
		t.Fatalf("创建设备级基线失败: %v", err)
	}
	if result.ID == 0 {
		t.Error("基线 ID 不应为 0")
	}
	if result.Scope != model.BaselineScopeDevice {
		t.Errorf("期望 scope=device，实际为 %s", result.Scope)
	}
	if *result.DeviceID != deviceID {
		t.Errorf("期望 device_id=%d，实际为 %d", deviceID, *result.DeviceID)
	}
}

func TestCreateGroupBaseline(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	groupRepo := repository.NewGroupRepository(db)

	groupID := createTestGroup(t, groupRepo)

	groupIDCopy := groupID
	baseline := &model.Baseline{
		Scope:   model.BaselineScopeGroup,
		GroupID: &groupIDCopy,
		Content: "hostname Switch\nlogging host 10.0.0.100\n",
	}

	result, err := svc.Create(context.Background(), baseline)
	if err != nil {
		t.Fatalf("创建分组级基线失败: %v", err)
	}
	if result.ID == 0 {
		t.Error("基线 ID 不应为 0")
	}
	if result.Scope != model.BaselineScopeGroup {
		t.Errorf("期望 scope=group，实际为 %s", result.Scope)
	}
	if *result.GroupID != groupID {
		t.Errorf("期望 group_id=%d，实际为 %d", groupID, *result.GroupID)
	}
}

func TestCreateDuplicateDeviceBaseline(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)

	deviceID := createTestDevice(t, deviceRepo)

	deviceIDCopy := deviceID
	first := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "baseline content 1",
	}
	if _, err := svc.Create(context.Background(), first); err != nil {
		t.Fatalf("创建第一个基线失败: %v", err)
	}

	second := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "baseline content 2",
	}
	_, err := svc.Create(context.Background(), second)
	if err == nil {
		t.Fatal("重复创建设备级基线应返回错误")
	}
}

func TestCreateInvalidScope(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)

	baseline := &model.Baseline{
		Scope:   BaselineScopeInvalid,
		Content: "test",
	}
	_, err := svc.Create(context.Background(), baseline)
	if err == nil {
		t.Fatal("无效作用域应返回错误")
	}
}

func TestCreateEmptyContent(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)

	baseline := &model.Baseline{
		Scope:   model.BaselineScopeDevice,
		Content: "",
	}
	_, err := svc.Create(context.Background(), baseline)
	if err == nil {
		t.Fatal("空内容应返回错误")
	}
}

func TestCreateDeviceScopeWithoutDeviceID(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)

	baseline := &model.Baseline{
		Scope:   model.BaselineScopeDevice,
		Content: "test content",
	}
	_, err := svc.Create(context.Background(), baseline)
	if err == nil {
		t.Fatal("设备级基线缺少 device_id 应返回错误")
	}
}

func TestCreateGroupScopeWithoutGroupID(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)

	baseline := &model.Baseline{
		Scope:   model.BaselineScopeGroup,
		Content: "test content",
	}
	_, err := svc.Create(context.Background(), baseline)
	if err == nil {
		t.Fatal("分组级基线缺少 group_id 应返回错误")
	}
}

func TestGetByID(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)

	deviceID := createTestDevice(t, deviceRepo)
	deviceIDCopy := deviceID

	baseline := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "test content for get",
	}
	created, _ := svc.Create(context.Background(), baseline)

	found, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("查询基线失败: %v", err)
	}
	if found.Content != "test content for get" {
		t.Errorf("期望内容=test content for get，实际为 %s", found.Content)
	}
}

func TestGetByNotFound(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)

	_, err := svc.GetByID(context.Background(), 99999)
	if err == nil {
		t.Fatal("不存在的基线应返回错误")
	}
}

func TestUpdateBaseline(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)

	deviceID := createTestDevice(t, deviceRepo)
	deviceIDCopy := deviceID

	baseline := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "original content",
	}
	created, _ := svc.Create(context.Background(), baseline)

	updated := &model.Baseline{Content: "updated content"}
	result, err := svc.Update(context.Background(), created.ID, updated)
	if err != nil {
		t.Fatalf("更新基线失败: %v", err)
	}
	if result.Content != "updated content" {
		t.Errorf("期望内容=updated content，实际为 %s", result.Content)
	}
}

func TestDeleteBaseline(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)

	deviceID := createTestDevice(t, deviceRepo)
	deviceIDCopy := deviceID

	baseline := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "to be deleted",
	}
	created, _ := svc.Create(context.Background(), baseline)

	err := svc.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("删除基线失败: %v", err)
	}

	_, err = svc.GetByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("删除后查询应返回错误")
	}
}

func TestFindBaselineForDeviceDeviceLevel(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)
	groupRepo := repository.NewGroupRepository(db)

	groupID := createTestGroup(t, groupRepo)
	deviceID := createTestDeviceInGroup(t, deviceRepo, groupID)

	deviceIDCopy := deviceID
	deviceBaseline := &model.Baseline{
		Scope:    model.BaselineScopeDevice,
		DeviceID: &deviceIDCopy,
		Content:  "device-level baseline",
	}
	svc.Create(context.Background(), deviceBaseline)

	groupIDCopy := groupID
	groupBaseline := &model.Baseline{
		Scope:   model.BaselineScopeGroup,
		GroupID: &groupIDCopy,
		Content: "group-level baseline",
	}
	svc.Create(context.Background(), groupBaseline)

	found, err := svc.FindBaselineForDevice(context.Background(), deviceID)
	if err != nil {
		t.Fatalf("查找设备基线失败: %v", err)
	}
	if found == nil {
		t.Fatal("应找到设备级基线")
	}
	if found.Content != "device-level baseline" {
		t.Errorf("优先返回设备级基线，期望内容=device-level baseline，实际为 %s", found.Content)
	}
}

func TestFindBaselineForDeviceGroupFallback(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)
	groupRepo := repository.NewGroupRepository(db)

	groupID := createTestGroup(t, groupRepo)
	deviceID := createTestDeviceInGroup(t, deviceRepo, groupID)

	groupIDCopy := groupID
	groupBaseline := &model.Baseline{
		Scope:   model.BaselineScopeGroup,
		GroupID: &groupIDCopy,
		Content: "group-level baseline only",
	}
	svc.Create(context.Background(), groupBaseline)

	found, err := svc.FindBaselineForDevice(context.Background(), deviceID)
	if err != nil {
		t.Fatalf("查找设备基线失败: %v", err)
	}
	if found == nil {
		t.Fatal("无设备级时应回退到分组级基线")
	}
	if found.Content != "group-level baseline only" {
		t.Errorf("期望内容=group-level baseline only，实际为 %s", found.Content)
	}
}

func TestFindBaselineForDeviceNoBaseline(t *testing.T) {
	db := setupBaselineTestDB(t)
	svc := setupBaselineTestService(t, db)
	deviceRepo := repository.NewDeviceRepository(db)

	deviceID := createTestDevice(t, deviceRepo)

	found, err := svc.FindBaselineForDevice(context.Background(), deviceID)
	if err != nil {
		t.Fatalf("查找设备基线失败: %v", err)
	}
	if found != nil {
		t.Error("无基线时应返回 nil")
	}
}

var BaselineScopeInvalid model.BaselineScope = "invalid"

func createTestDevice(t *testing.T, repo repository.DeviceRepository) uint {
	t.Helper()
	device := &model.Device{
		Name:     "test-device",
		IP:       "10.0.0.1",
		Vendor:   model.VendorCisco,
		Model:    model.DeviceModelIOS,
		Protocol: model.ConnProtocolSSH,
		Port:     22,
		Username: "admin",
		Password: "secret",
	}
	if err := repo.Create(nil, device); err != nil {
		t.Fatalf("创建测试设备失败: %v", err)
	}
	return device.ID
}

func createTestDeviceInGroup(t *testing.T, repo repository.DeviceRepository, groupID uint) uint {
	t.Helper()
	gid := groupID
	device := &model.Device{
		Name:     "test-device-in-group",
		IP:       "10.0.0.2",
		Vendor:   model.VendorCisco,
		Model:    model.DeviceModelIOS,
		Protocol: model.ConnProtocolSSH,
		Port:     22,
		Username: "admin",
		Password: "secret",
		GroupID:  &gid,
	}
	if err := repo.Create(nil, device); err != nil {
		t.Fatalf("创建测试设备失败: %v", err)
	}
	return device.ID
}

func createTestGroup(t *testing.T, repo repository.GroupRepository) uint {
	t.Helper()
	group := &model.Group{Name: "test-group", Description: "for testing"}
	if err := repo.Create(nil, group); err != nil {
		t.Fatalf("创建测试分组失败: %v", err)
	}
	return group.ID
}
