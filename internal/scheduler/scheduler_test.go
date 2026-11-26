package scheduler

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/netconfighub/netconfighub/internal/apperr"
	"github.com/netconfighub/netconfighub/internal/config"
)

func TestSchedulerStartStop(t *testing.T) {
	cfg := &config.SchedulerConfig{
		DefaultInterval: "1h",
		WorkerPoolSize:  2,
		MaxRetries:      3,
		RetryDelay:      "30s",
		TaskTimeout:     "5m",
	}
	s := NewScheduler(SchedulerDeps{Cfg: cfg, EncryptionKey: "test-encryption-key-32bytes!!"})

	if err := s.Start(); err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}

	time.Sleep(200 * time.Millisecond)
	s.Stop()
}

func TestSubmitDevice(t *testing.T) {
	cfg := &config.SchedulerConfig{
		DefaultInterval: "1h",
		WorkerPoolSize:  2,
		MaxRetries:      3,
		RetryDelay:      "30s",
		TaskTimeout:     "5m",
	}
	s := NewScheduler(SchedulerDeps{Cfg: cfg, EncryptionKey: "test-encryption-key-32bytes!!"})

	if err := s.Start(); err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	defer s.Stop()

	if err := s.SubmitDevice(1); err != nil {
		t.Fatalf("提交设备失败: %v", err)
	}

	if err := s.SubmitDevice(2); err != nil {
		t.Fatalf("提交设备失败: %v", err)
	}
}

func TestSubmitDeviceQueueFull(t *testing.T) {
	cfg := &config.SchedulerConfig{
		DefaultInterval: "1h",
		WorkerPoolSize:  1,
		MaxRetries:      3,
		RetryDelay:      "30s",
		TaskTimeout:     "5m",
	}
	s := NewScheduler(SchedulerDeps{Cfg: cfg, EncryptionKey: "test-encryption-key-32bytes!!"})

	for i := 0; i < cap(s.taskQueue); i++ {
		s.taskQueue <- backupJob{deviceID: uint(i)}
	}

	err := s.SubmitDevice(999)
	if err == nil {
		t.Error("队列已满时应返回错误")
	}
}

func TestSchedulerStopIsIdempotent(t *testing.T) {
	cfg := &config.SchedulerConfig{
		DefaultInterval: "1h",
		WorkerPoolSize:  2,
		MaxRetries:      3,
		RetryDelay:      "30s",
		TaskTimeout:     "5m",
	}
	s := NewScheduler(SchedulerDeps{Cfg: cfg, EncryptionKey: "test-encryption-key-32bytes!!"})
	if err := s.Start(); err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	s.Stop()
	s.Stop() // 第二次 Stop 不得 panic
}

func TestConcurrentSubmitDuringStopDoesNotPanic(t *testing.T) {
	cfg := &config.SchedulerConfig{
		DefaultInterval: "1h",
		WorkerPoolSize:  4,
		MaxRetries:      3,
		RetryDelay:      "30s",
		TaskTimeout:     "5m",
	}
	s := NewScheduler(SchedulerDeps{Cfg: cfg, EncryptionKey: "test-encryption-key-32bytes!!"})
	if err := s.Start(); err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			_ = s.SubmitDevice(id)
		}(uint(i))
	}
	s.Stop()
	wg.Wait()
}

func TestSubmitDeviceAfterStopReturnsError(t *testing.T) {
	cfg := &config.SchedulerConfig{
		DefaultInterval: "1h",
		WorkerPoolSize:  2,
		MaxRetries:      3,
		RetryDelay:      "30s",
		TaskTimeout:     "5m",
	}
	s := NewScheduler(SchedulerDeps{Cfg: cfg, EncryptionKey: "test-encryption-key-32bytes!!"})
	if err := s.Start(); err != nil {
		t.Fatalf("启动调度器失败: %v", err)
	}
	s.Stop()
	if err := s.SubmitDevice(1); err == nil {
		t.Error("停止后提交应返回错误")
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		errMsg   string
		wantType string
	}{
		{"SSH 连接超时 [192.168.1.1:22]", "connect_timeout"},
		{"Telnet 连接失败 [192.168.1.1:23]", "connect_timeout"},
		{"认证失败", "auth_failed"},
		{"password incorrect", "auth_failed"},
		{"提示符超时", "prompt_timeout"},
		{"prompt not found", "prompt_timeout"},
		{"命令执行超时", "command_timeout"},
		{"采集超时，已超过 5m", "command_timeout"},
		{"模型不存在", "model_not_found"},
		{"不支持的厂商类型", "model_not_found"},
		{"Git 提交失败", "git_commit_failed"},
		{"配置为空", "config_empty"},
		{"其他错误", "unknown_error"},
	}
	for _, tt := range tests {
		got := classifyError(fmt.Errorf("%s", tt.errMsg))
		if got != tt.wantType {
			t.Errorf("classifyError(%q) = %q, want %q", tt.errMsg, got, tt.wantType)
		}
	}
}

func TestClassifyErrorTyped(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"connect timeout", apperr.Wrap(apperr.ErrConnectTimeout, "SSH 连接失败 [x]"), "connect_timeout"},
		{"auth failed", apperr.Wrap(apperr.ErrAuthFailed, "认证失败"), "auth_failed"},
		{"prompt timeout", apperr.Wrap(apperr.ErrPromptTimeout, "提示符超时"), "prompt_timeout"},
		{"command timeout", apperr.Wrap(apperr.ErrCommandTimeout, "命令执行超时"), "command_timeout"},
		{"model not found", apperr.Wrap(apperr.ErrModelNotFound, "不支持的厂商类型"), "model_not_found"},
		{"git commit", apperr.Wrap(apperr.ErrGitCommit, "提交变更失败"), "git_commit_failed"},
		{"config empty", apperr.Wrap(apperr.ErrConfigEmpty, "配置为空"), "config_empty"},
		{"double wrapped", fmt.Errorf("collect: %w", apperr.Wrap(apperr.ErrConnectTimeout, "连接失败")), "connect_timeout"},
		{"unknown", errors.New("something else"), "unknown_error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyError(tt.err); got != tt.want {
				t.Errorf("classifyError(%v) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}
