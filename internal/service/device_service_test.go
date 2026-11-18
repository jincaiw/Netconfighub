package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/netconfighub/netconfighub/internal/model"
)

type countingGroupRepo struct {
	listAllNamesCalls int
	groups            []model.Group
}

func (r *countingGroupRepo) Create(ctx context.Context, group *model.Group) error { return nil }
func (r *countingGroupRepo) FindByID(ctx context.Context, id uint) (*model.Group, error) {
	return nil, fmt.Errorf("not found")
}
func (r *countingGroupRepo) List(ctx context.Context, offset, limit int) ([]model.Group, int64, error) {
	return r.groups, int64(len(r.groups)), nil
}
func (r *countingGroupRepo) ListAllNames(ctx context.Context, limit int) ([]model.Group, error) {
	r.listAllNamesCalls++
	return r.groups, nil
}
func (r *countingGroupRepo) Update(ctx context.Context, group *model.Group) error { return nil }
func (r *countingGroupRepo) Delete(ctx context.Context, id uint) error            { return nil }

type recordingDeviceRepo struct {
	created           []*model.Device
	listCalls         int
	unboundedListCall bool
	nextID            uint
}

func (r *recordingDeviceRepo) Create(ctx context.Context, device *model.Device) error {
	r.nextID++
	device.ID = r.nextID
	r.created = append(r.created, device)
	return nil
}
func (r *recordingDeviceRepo) FindByID(ctx context.Context, id uint) (*model.Device, error) {
	return nil, fmt.Errorf("not found")
}
func (r *recordingDeviceRepo) FindByName(ctx context.Context, name string) (*model.Device, error) {
	return nil, fmt.Errorf("not found")
}
func (r *recordingDeviceRepo) FindByIP(ctx context.Context, ip string) (*model.Device, error) {
	return nil, fmt.Errorf("not found")
}
func (r *recordingDeviceRepo) List(ctx context.Context, groupID *uint, vendor *string, status *string, search *string, offset, limit int) ([]model.Device, int64, error) {
	r.listCalls++
	if limit == 0 || limit > 10000 {
		r.unboundedListCall = true
	}
	return nil, 0, nil
}
func (r *recordingDeviceRepo) Update(ctx context.Context, device *model.Device) error { return nil }
func (r *recordingDeviceRepo) Delete(ctx context.Context, id uint) error              { return nil }
func (r *recordingDeviceRepo) FindByGroupID(ctx context.Context, groupID uint) ([]model.Device, error) {
	return nil, nil
}
func (r *recordingDeviceRepo) ListFailed(ctx context.Context, offset, limit int) ([]model.Device, int64, error) {
	return nil, 0, nil
}

func TestImportDevicesLoadsGroupsOnce(t *testing.T) {
	groups := &countingGroupRepo{groups: []model.Group{{ID: 7, Name: "core"}}}
	devices := &recordingDeviceRepo{}
	svc := NewDeviceService(devices, groups, nil, nil, nil, nil, nil, "test-encryption-key-32bytes!!")

	csvData := []byte("name,ip,vendor,model,protocol,port,username,password,group_name\n" +
		"sw1,10.0.0.1,cisco,ios,ssh,22,admin,pass,core\n" +
		"sw2,10.0.0.2,cisco,ios,ssh,22,admin,pass,core\n" +
		"sw3,10.0.0.3,cisco,ios,ssh,22,admin,pass,core\n")

	result, err := svc.ImportDevices(context.Background(), csvData)
	if err != nil {
		t.Fatalf("ImportDevices failed: %v", err)
	}
	if result.Success != 3 || result.Failed != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if groups.listAllNamesCalls != 1 {
		t.Fatalf("group lookup called %d times, want 1 (O(rows) N+1 regression)", groups.listAllNamesCalls)
	}
	if len(devices.created) != 3 {
		t.Fatalf("created %d devices, want 3", len(devices.created))
	}
	id := uint(7)
	for i, d := range devices.created {
		if d.GroupID == nil || *d.GroupID != id {
			t.Fatalf("device %d groupID = %v, want 7", i, d.GroupID)
		}
	}
	if devices.unboundedListCall {
		t.Fatal("unexpected unbounded device List call during import")
	}
}
