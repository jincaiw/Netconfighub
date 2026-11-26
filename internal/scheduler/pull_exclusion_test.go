package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/netconfighub/netconfighub/internal/apperr"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"gorm.io/gorm"
)

func newDeviceRepoForTest(db *gorm.DB) repository.DeviceRepository {
	return repository.NewDeviceRepository(db)
}

func TestSubmitDeviceRejectsSFTPPushDevice(t *testing.T) {
	db := newIngestPathTestDB(t)
	deviceRepo := newDeviceRepoForTest(db)

	device := &model.Device{
		Name: "sftp-device", IP: "192.0.2.90", Vendor: model.VendorCisco,
		Protocol: model.ConnProtocolSSH, Port: 22, Username: "admin",
		Enabled: true, CollectionMode: model.CollectionModeSFTPPush,
	}
	if err := deviceRepo.Create(context.Background(), device); err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}

	sched := &Scheduler{deviceRepo: deviceRepo, stopCh: make(chan struct{})}
	err := sched.SubmitDevice(device.ID)
	if err == nil {
		t.Fatalf("sftp_push 设备不得被提交拉取")
	}
	if !errors.Is(err, apperr.ErrDeviceNotPullable) {
		t.Fatalf("应返回 ErrDeviceNotPullable，实际 %v", err)
	}
}

func TestSubmitDeviceAcceptsPullDevice(t *testing.T) {
	db := newIngestPathTestDB(t)
	deviceRepo := newDeviceRepoForTest(db)

	device := &model.Device{
		Name: "pull-device", IP: "192.0.2.91", Vendor: model.VendorCisco,
		Protocol: model.ConnProtocolSSH, Port: 22, Username: "admin",
		Enabled: true, CollectionMode: model.CollectionModePull,
	}
	if err := deviceRepo.Create(context.Background(), device); err != nil {
		t.Fatalf("创建设备失败: %v", err)
	}

	sched := &Scheduler{
		deviceRepo: deviceRepo,
		stopCh:     make(chan struct{}),
		taskQueue:  make(chan backupJob, 4),
	}
	if err := sched.SubmitDevice(device.ID); err != nil {
		t.Fatalf("pull 设备应可提交: %v", err)
	}
}
