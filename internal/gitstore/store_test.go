package gitstore

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	dir, err := os.MkdirTemp("", "gitstore-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	store, err := NewStore(dir, "test-author", "test@test.com", "", false, "")
	if err != nil {
		os.RemoveAll(dir)
		t.Fatalf("NewStore 失败: %v", err)
	}
	return store, dir
}

func TestNewStore(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitstore-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(dir)

	store, err := NewStore(dir, "test-author", "test@test.com", "", false, "")
	if err != nil {
		t.Fatalf("NewStore 失败: %v", err)
	}
	if store == nil {
		t.Fatal("store 不应为 nil")
	}
	if store.repo == nil {
		t.Fatal("repo 不应为 nil")
	}

	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Fatal(".git 目录应存在")
	}
}

func TestNewStoreExistingRepo(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitstore-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(dir)

	store1, err := NewStore(dir, "test-author", "test@test.com", "", false, "")
	if err != nil {
		t.Fatalf("第一次 NewStore 失败: %v", err)
	}
	_ = store1

	store2, err := NewStore(dir, "test-author", "test@test.com", "", false, "")
	if err != nil {
		t.Fatalf("第二次 NewStore 失败: %v", err)
	}
	if store2 == nil {
		t.Fatal("store2 不应为 nil")
	}
	if store2.repo == nil {
		t.Fatal("repo 不应为 nil")
	}
}

func TestNewStoreNonexistentDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitstore-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(dir)

	repoPath := filepath.Join(dir, "subdir", "repo")
	store, err := NewStore(repoPath, "test-author", "test@test.com", "", false, "")
	if err != nil {
		t.Fatalf("NewStore 失败: %v", err)
	}
	if store == nil {
		t.Fatal("store 不应为 nil")
	}
}

func TestWriteConfig(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	hash, err := store.WriteConfig("group1", "device1", "config content v1")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}
	if hash == "" {
		t.Fatal("hash 不应为空")
	}

	configFile := filepath.Join(dir, "group1", "device1", "running-config.txt")
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}
	if string(content) != "config content v1" {
		t.Fatalf("配置内容不匹配: got %q, want %q", string(content), "config content v1")
	}
}

func TestWriteConfigUngrouped(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	hash, err := store.WriteConfig("", "device1", "config content")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}
	if hash == "" {
		t.Fatal("hash 不应为空")
	}

	configFile := filepath.Join(dir, "ungrouped", "device1", "running-config.txt")
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}
	if string(content) != "config content" {
		t.Fatalf("配置内容不匹配: got %q, want %q", string(content), "config content")
	}
}

func TestWriteConfigDuplicate(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	hash1, err := store.WriteConfig("group1", "device1", "config content")
	if err != nil {
		t.Fatalf("第一次 WriteConfig 失败: %v", err)
	}
	if hash1 == "" {
		t.Fatal("第一次 hash 不应为空")
	}

	hash2, err := store.WriteConfig("group1", "device1", "config content")
	if err != nil {
		t.Fatalf("第二次 WriteConfig 失败: %v", err)
	}
	if hash2 != "" {
		t.Fatal("重复写入相同配置应返回空哈希")
	}
}

func TestWriteConfigDifferent(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	hash1, err := store.WriteConfig("group1", "device1", "config content v1")
	if err != nil {
		t.Fatalf("第一次 WriteConfig 失败: %v", err)
	}
	if hash1 == "" {
		t.Fatal("第一次 hash 不应为空")
	}

	hash2, err := store.WriteConfig("group1", "device1", "config content v2")
	if err != nil {
		t.Fatalf("第二次 WriteConfig 失败: %v", err)
	}
	if hash2 == "" {
		t.Fatal("写入不同配置应返回非空哈希")
	}
	if hash1 == hash2 {
		t.Fatal("不同配置的哈希应不同")
	}
}

func TestConfigHistory(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	_, err := store.WriteConfig("group1", "device1", "config content v1")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	_, err = store.WriteConfig("group1", "device1", "config content v2")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	history, err := store.ConfigHistory("group1", "device1", 0)
	if err != nil {
		t.Fatalf("ConfigHistory 失败: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("历史版本数量不匹配: got %d, want 2", len(history))
	}
}

func TestConfigHistoryWithLimit(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	_, err := store.WriteConfig("group1", "device1", "config content v1")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	_, err = store.WriteConfig("group1", "device1", "config content v2")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	_, err = store.WriteConfig("group1", "device1", "config content v3")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	history, err := store.ConfigHistory("group1", "device1", 2)
	if err != nil {
		t.Fatalf("ConfigHistory 失败: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("历史版本数量不匹配: got %d, want 2", len(history))
	}
}

func TestReadConfigAtRevision(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	_, err := store.WriteConfig("group1", "device1", "config content v1")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	_, err = store.WriteConfig("group1", "device1", "config content v2")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	history, err := store.ConfigHistory("group1", "device1", 0)
	if err != nil {
		t.Fatalf("ConfigHistory 失败: %v", err)
	}
	if len(history) < 2 {
		t.Fatalf("历史版本数量不足: %d", len(history))
	}

	oldHash := history[1].Hash.String()
	content, err := store.ReadConfigAtRevision("group1", "device1", oldHash)
	if err != nil {
		t.Fatalf("ReadConfigAtRevision 失败: %v", err)
	}
	if content != "config content v1" {
		t.Fatalf("旧版本配置内容不匹配: got %q, want %q", content, "config content v1")
	}

	latestHash := history[0].Hash.String()
	latestContent, err := store.ReadConfigAtRevision("group1", "device1", latestHash)
	if err != nil {
		t.Fatalf("ReadConfigAtRevision 失败: %v", err)
	}
	if latestContent != "config content v2" {
		t.Fatalf("最新版本配置内容不匹配: got %q, want %q", latestContent, "config content v2")
	}
}

func TestReadLatestConfig(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	_, err := store.WriteConfig("group1", "device1", "config content v1")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	content, err := store.ReadLatestConfig("group1", "device1")
	if err != nil {
		t.Fatalf("ReadLatestConfig 失败: %v", err)
	}
	if content != "config content v1" {
		t.Fatalf("配置内容不匹配: got %q, want %q", content, "config content v1")
	}

	_, err = store.WriteConfig("group1", "device1", "config content v2")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	content, err = store.ReadLatestConfig("group1", "device1")
	if err != nil {
		t.Fatalf("ReadLatestConfig 失败: %v", err)
	}
	if content != "config content v2" {
		t.Fatalf("配置内容不匹配: got %q, want %q", content, "config content v2")
	}
}

func TestReadLatestConfigNotExist(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	_, err := store.ReadLatestConfig("group1", "device1")
	if err == nil {
		t.Fatal("读取不存在的配置应返回错误")
	}
}

func TestListDeviceConfigs(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	_, err := store.WriteConfig("group1", "device1", "config content 1")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	_, err = store.WriteConfig("group1", "device2", "config content 2")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	_, err = store.WriteConfig("group2", "device3", "config content 3")
	if err != nil {
		t.Fatalf("WriteConfig 失败: %v", err)
	}

	devices, err := store.ListDeviceConfigs()
	if err != nil {
		t.Fatalf("ListDeviceConfigs 失败: %v", err)
	}
	if len(devices) != 3 {
		t.Fatalf("设备数量不匹配: got %d, want 3", len(devices))
	}

	found := make(map[string]bool)
	for _, d := range devices {
		key := d.GroupName + "/" + d.DeviceName
		found[key] = true
	}
	if !found["group1/device1"] {
		t.Fatal("缺少 group1/device1")
	}
	if !found["group1/device2"] {
		t.Fatal("缺少 group1/device2")
	}
	if !found["group2/device3"] {
		t.Fatal("缺少 group2/device3")
	}
}

func TestListDeviceConfigsEmpty(t *testing.T) {
	store, dir := newTestStore(t)
	defer os.RemoveAll(dir)

	devices, err := store.ListDeviceConfigs()
	if err != nil {
		t.Fatalf("ListDeviceConfigs 失败: %v", err)
	}
	if len(devices) != 0 {
		t.Fatalf("空仓库设备数量应为 0, got %d", len(devices))
	}
}
