package sftpserver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/sftp"
)

// TestConcurrentWriteHashMatchesFileContent guards the regression where the
// upload hash was accumulated in WriteAt order. pkg/sftp dispatches WRITE
// packets from several workers, so hashing on arrival produced a digest that
// did not match the finished file and every large upload failed to archive.
func TestConcurrentWriteHashMatchesFileContent(t *testing.T) {
	cfg := testConfig(t)
	cfg.MaxFileSizeBytes = 1 << 20

	sink := &recordingSink{}
	ctx := uploadContext{deviceID: 1, accountID: 1, sink: sink}

	file, err := newRestrictedFile(cfg, ctx, "concurrent.txt")
	if err != nil {
		t.Fatalf("创建受限文件失败: %v", err)
	}

	// Build content and write it out of order, as concurrent workers would.
	content := make([]byte, 64*1024)
	for i := range content {
		content[i] = byte('a' + i%26)
	}
	const chunk = 8 * 1024
	order := []int{5, 0, 3, 7, 1, 6, 2, 4}
	for _, idx := range order {
		start := idx * chunk
		if _, err := file.WriteAt(content[start:start+chunk], int64(start)); err != nil {
			t.Fatalf("写入分片 %d 失败: %v", idx, err)
		}
	}

	if err := file.Close(); err != nil {
		t.Fatalf("关闭失败: %v", err)
	}
	if len(sink.uploads) != 1 {
		t.Fatalf("应记录 1 次上传，实际 %d", len(sink.uploads))
	}

	want := sha256.Sum256(content)
	wantHex := hex.EncodeToString(want[:])
	if sink.uploads[0].Hash != wantHex {
		t.Fatalf("乱序写入后哈希必须与最终文件一致\nwant=%s\ngot =%s", wantHex, sink.uploads[0].Hash)
	}

	// The staged file on disk must match too.
	onDisk, err := os.ReadFile(sink.uploads[0].StagingPath)
	if err != nil {
		t.Fatalf("读取暂存文件失败: %v", err)
	}
	if len(onDisk) != len(content) {
		t.Fatalf("暂存文件大小 = %d, want %d", len(onDisk), len(content))
	}
	diskSum := sha256.Sum256(onDisk)
	if hex.EncodeToString(diskSum[:]) != wantHex {
		t.Fatalf("暂存文件内容与写入不一致")
	}
}

// TestDisabledServerDoesNotCreateHostKey guards the regression where a disabled
// listener still created directories and keys, making a bad path fatal for the
// whole service.
func TestDisabledServerDoesNotCreateHostKey(t *testing.T) {
	cfg := testConfig(t)
	cfg.Enabled = false
	// Point at a path whose parent cannot be created: if New() tried, it would
	// fail or leave artifacts behind.
	cfg.HostKeyPath = filepath.Join(cfg.HostKeyPath, "nested", "host_key")

	server, err := New(cfg, &fakeVerifier{}, &recordingSink{})
	if err != nil {
		t.Fatalf("禁用状态下创建服务不应失败: %v", err)
	}
	if server == nil {
		t.Fatalf("应返回可用的服务实例")
	}
	if _, statErr := os.Stat(cfg.HostKeyPath); !os.IsNotExist(statErr) {
		t.Fatalf("禁用状态下不得创建主机密钥")
	}
	if err := server.Start(); err != nil {
		t.Fatalf("禁用状态下 Start 不应报错: %v", err)
	}
}

// TestDownloadIsDeniedByReader documents that a nil FileGet is not acceptable:
// it panics inside pkg/sftp, so a rejecting reader must be installed.
func TestDeniedReaderRejectsRead(t *testing.T) {
	var reader deniedReader
	if _, err := reader.Fileread(&sftp.Request{Method: "Get", Filepath: "/upload/x.txt"}); err == nil {
		t.Fatalf("下载请求必须被拒绝")
	}
}

var _ = context.Background
