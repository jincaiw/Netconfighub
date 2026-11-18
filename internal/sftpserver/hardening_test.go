package sftpserver

import (
	"net"
	"testing"
	"time"

	"github.com/pkg/sftp"
)

// fakeDeadlineConn records the deadline requested by the limits helper without
// requiring a real socket or any sleeping.
type fakeDeadlineConn struct {
	net.Conn
	deadline time.Time
}

func (c *fakeDeadlineConn) SetDeadline(t time.Time) error {
	c.deadline = t
	return nil
}

// TestSymlinkAndUnsafeMutationRejected covers 验收 4「链接」：client issued
// SSH_FXP_SYMLINK / rename targets must be denied so a device cannot create
// links to escape the virtual root or move other uploads.
func TestSymlinkAndUnsafeMutationRejected(t *testing.T) {
	cfg := testConfig(t)
	server, _ := newTestServer(t, cfg, &recordingSink{})
	addr, stop := startTestListener(t, server)
	defer stop()

	client := dialTest(t, addr, "device-1", "secret-password")
	defer client.Close()

	sc, err := sftp.NewClient(client)
	if err != nil {
		t.Fatalf("创建 SFTP 客户端失败: %v", err)
	}
	defer sc.Close()

	if err := sc.Symlink("/upload/a.txt", "/upload/b.txt"); err == nil {
		t.Fatalf("symlink 必须被拒绝")
	}
	// Rename of a path with no in-flight upload must be denied, otherwise a
	// client could commit an arbitrary file into the received directory.
	if err := sc.Rename("/upload/notexist.txt", "/upload/x.txt"); err == nil {
		t.Fatalf("无 inflight 的 rename 必须被拒绝")
	}
	// Setstat (chmod) is a mutation and must be denied.
	if err := sc.Chmod("/upload/a.txt", 0o644); err == nil {
		t.Fatalf("chmod/setstat 必须被拒绝")
	}
}

// TestLimitsDeadlineAndConcurrency covers 验收 4「超时」与「并发」的控制边界：
// the connection cap, the per-device single-upload lock, and the idle/session
// deadline selection must all behave as documented. Deterministic, no sleeps.
func TestLimitsDeadlineAndConcurrency(t *testing.T) {
	l := newLimits(2, 10*time.Second, time.Minute)

	// Connection semaphore: exactly maxConnections slots.
	if !l.acquireConn() || !l.acquireConn() {
		t.Fatalf("前 %d 个连接应成功", 2)
	}
	if l.acquireConn() {
		t.Fatalf("超过 maxConnections 的连接必须被拒绝")
	}
	l.releaseConn()
	if !l.acquireConn() {
		t.Fatalf("释放后应能再次获取连接槽")
	}

	// Per-device lock: only one concurrent upload per device.
	if !l.tryLockDevice(1) {
		t.Fatalf("首次设备锁应成功")
	}
	if l.tryLockDevice(1) {
		t.Fatalf("同一设备第二把锁必须被拒绝")
	}
	l.unlockDevice(1)
	if !l.tryLockDevice(1) {
		t.Fatalf("解锁后应能再次加锁")
	}

	// idle < session => idle deadline is the effective bound.
	c1 := &fakeDeadlineConn{}
	before := time.Now()
	l.applyDeadlines(c1)
	if c1.deadline.IsZero() {
		t.Fatalf("应设置 deadline")
	}
	if d := c1.deadline.Sub(before); d < 9*time.Second || d > 11*time.Second {
		t.Fatalf("idle 较短时应使用 idle 截止时间，实际 %v", d)
	}

	// session < idle => absolute session deadline wins.
	l2 := newLimits(1, time.Minute, 5*time.Second)
	c2 := &fakeDeadlineConn{}
	before2 := time.Now()
	l2.applyDeadlines(c2)
	if d := c2.deadline.Sub(before2); d < 4*time.Second || d > 6*time.Second {
		t.Fatalf("session 较短时应使用 session 截止时间，实际 %v", d)
	}

	// touch must never extend past the absolute session deadline.
	c3 := &fakeDeadlineConn{}
	sessionDeadline := time.Now().Add(2 * time.Second)
	l.touch(c3, sessionDeadline)
	if !c3.deadline.Equal(sessionDeadline) {
		t.Fatalf("touch 不应超过绝对会话截止时间，got=%v want=%v", c3.deadline, sessionDeadline)
	}
}
