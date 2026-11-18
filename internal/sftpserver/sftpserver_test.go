package sftpserver

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

// fakeVerifier accepts exactly one username/password pair and returns the
// device behind it.
type fakeVerifier struct {
	username string
	password string
	device   *model.Device
	account  *model.DeviceSFTPAccount
	err      error
}

func (f *fakeVerifier) VerifySFTPCredentials(_ context.Context, username, password string) (*model.DeviceSFTPAccount, *model.Device, error) {
	if f.err != nil {
		return nil, nil, f.err
	}
	if username != f.username || password != f.password {
		return nil, nil, errors.New("认证失败")
	}
	return f.account, f.device, nil
}

// recordingSink captures OnReceived invocations.
type recordingSink struct {
	uploads []ReceivedUpload
	err     error
}

func (r *recordingSink) OnReceived(_ context.Context, upload ReceivedUpload) error {
	if r.err != nil {
		return r.err
	}
	r.uploads = append(r.uploads, upload)
	return nil
}

func testConfig(t *testing.T) config.SftpConfig {
	t.Helper()
	dir := t.TempDir()
	return config.SftpConfig{
		Enabled:          true,
		Host:             "127.0.0.1",
		Port:             0,
		HostKeyPath:      filepath.Join(dir, "sftp", "host_key"),
		StagingDir:       filepath.Join(dir, "sftp"),
		MaxFileSizeBytes: 4096,
		MaxConnections:   4,
		IdleTimeout:      "30s",
		SessionTimeout:   "2m",
	}
}

func newTestServer(t *testing.T, cfg config.SftpConfig, sink UploadSink) (*Server, *fakeVerifier) {
	t.Helper()
	verifier := &fakeVerifier{
		username: "device-1",
		password: "secret-password",
		device:   &model.Device{ID: 1, Name: "dev1", Enabled: true, CollectionMode: model.CollectionModeSFTPPush},
		account:  &model.DeviceSFTPAccount{ID: 10, DeviceID: 1, Username: "device-1", Status: model.SFTPAccountStatusActive},
	}
	server, err := New(cfg, verifier, sink)
	if err != nil {
		t.Fatalf("创建 SFTP 服务失败: %v", err)
	}
	return server, verifier
}

// startTestListener binds an ephemeral port and returns the address plus a
// stop function, so tests do not depend on the configured port.
func startTestListener(t *testing.T, server *Server) (string, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("监听失败: %v", err)
	}
	server.mu.Lock()
	server.listener = ln
	server.mu.Unlock()

	server.wg.Add(1)
	go server.acceptLoop(ln)

	return ln.Addr().String(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		server.Stop(ctx)
	}
}

func dialTest(t *testing.T, addr, user, pass string) *ssh.Client {
	t.Helper()
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	})
	if err != nil {
		t.Fatalf("SSH 连接失败: %v", err)
	}
	return client
}

func TestHostKeyPersistsAndIsRestricted(t *testing.T) {
	cfg := testConfig(t)

	first, err := LoadOrCreateHostKey(cfg.HostKeyPath)
	if err != nil {
		t.Fatalf("首次加载主机密钥失败: %v", err)
	}
	firstFP := FingerprintSHA256(first)

	second, err := LoadOrCreateHostKey(cfg.HostKeyPath)
	if err != nil {
		t.Fatalf("再次加载主机密钥失败: %v", err)
	}
	if got := FingerprintSHA256(second); got != firstFP {
		t.Fatalf("主机密钥指纹不稳定: %s != %s", got, firstFP)
	}

	info, err := os.Stat(cfg.HostKeyPath)
	if err != nil {
		t.Fatalf("主机密钥文件不存在: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("主机密钥权限应为 0600，实际 %o", perm)
	}
	if !strings.HasPrefix(firstFP, "SHA256:") {
		t.Fatalf("指纹格式错误: %s", firstFP)
	}
}

func TestStartDisabledDoesNotListen(t *testing.T) {
	cfg := testConfig(t)
	cfg.Enabled = false

	server, _ := newTestServer(t, cfg, &recordingSink{})
	if err := server.Start(); err != nil {
		t.Fatalf("禁用时 Start 不应报错: %v", err)
	}
	server.mu.Lock()
	ln := server.listener
	server.mu.Unlock()
	if ln != nil {
		t.Fatalf("禁用时不应创建监听")
	}
}

func TestShellAndExecAreRejected(t *testing.T) {
	cfg := testConfig(t)
	sink := &recordingSink{}
	server, _ := newTestServer(t, cfg, sink)
	addr, stop := startTestListener(t, server)
	defer stop()

	client := dialTest(t, addr, "device-1", "secret-password")
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}
	defer session.Close()

	if err := session.Shell(); err == nil {
		t.Fatalf("shell 请求必须被拒绝")
	}

	session2, err := client.NewSession()
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}
	defer session2.Close()
	if err := session2.Run("id"); err == nil {
		t.Fatalf("exec 请求必须被拒绝")
	}
}

func TestPortForwardingIsRejected(t *testing.T) {
	cfg := testConfig(t)
	server, _ := newTestServer(t, cfg, &recordingSink{})
	addr, stop := startTestListener(t, server)
	defer stop()

	client := dialTest(t, addr, "device-1", "secret-password")
	defer client.Close()

	if _, err := client.Dial("tcp", "127.0.0.1:22"); err == nil {
		t.Fatalf("direct-tcpip 必须被拒绝")
	}
}

func TestBadCredentialsRejected(t *testing.T) {
	cfg := testConfig(t)
	server, _ := newTestServer(t, cfg, &recordingSink{})
	addr, stop := startTestListener(t, server)
	defer stop()

	if _, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            "device-1",
		Auth:            []ssh.AuthMethod{ssh.Password("wrong")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}); err == nil {
		t.Fatalf("错误口令必须被拒绝")
	}
}

func TestUploadSucceedsAndRecordsSink(t *testing.T) {
	cfg := testConfig(t)
	sink := &recordingSink{}
	server, _ := newTestServer(t, cfg, sink)
	addr, stop := startTestListener(t, server)
	defer stop()

	client := dialTest(t, addr, "device-1", "secret-password")
	defer client.Close()

	sshSession, err := client.NewSession()
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}
	defer sshSession.Close()

	// The session must be usable as an sftp client channel.
	pipe, err := sshSession.StdinPipe()
	if err != nil {
		t.Fatalf("stdin 失败: %v", err)
	}
	stdout, err := sshSession.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout 失败: %v", err)
	}
	if err := sshSession.RequestSubsystem("sftp"); err != nil {
		t.Fatalf("请求 sftp 子系统失败: %v", err)
	}

	// pkg/sftp client needs a ReadWriteCloser; use the low level newClientPipe.
	sc, err := sftp.NewClientPipe(stdout, pipe)
	if err != nil {
		t.Fatalf("创建 SFTP 客户端失败: %v", err)
	}
	defer sc.Close()

	content := "hostname dev1\ninterface GigabitEthernet0/1\n"
	f, err := sc.Create("/upload/running-config.txt")
	if err != nil {
		t.Fatalf("创建远端文件失败: %v", err)
	}
	if _, err := f.Write([]byte(content)); err != nil {
		t.Fatalf("写入失败: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("关闭失败: %v", err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if len(sink.uploads) > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if len(sink.uploads) != 1 {
		t.Fatalf("应记录 1 次上传，实际 %d", len(sink.uploads))
	}
	got := sink.uploads[0]
	if got.DeviceID != 1 {
		t.Fatalf("设备归属错误: %d", got.DeviceID)
	}
	if got.Bytes != int64(len(content)) {
		t.Fatalf("接收字节数错误: %d", got.Bytes)
	}

	receivedDir := filepath.Join(cfg.StagingDir, "received", "1")
	entries, err := os.ReadDir(receivedDir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("接收目录应有 1 个文件: %v", err)
	}
	info, err := entries[0].Info()
	if err != nil {
		t.Fatalf("读取文件信息失败: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("接收文件权限应为 0600，实际 %o", perm)
	}
}

func TestDownloadIsRejected(t *testing.T) {
	cfg := testConfig(t)
	server, _ := newTestServer(t, cfg, &recordingSink{})
	addr, stop := startTestListener(t, server)
	defer stop()

	client := dialTest(t, addr, "device-1", "secret-password")
	defer client.Close()

	sshSession, err := client.NewSession()
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}
	defer sshSession.Close()

	pipe, err := sshSession.StdinPipe()
	if err != nil {
		t.Fatalf("stdin 失败: %v", err)
	}
	stdout, err := sshSession.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout 失败: %v", err)
	}
	if err := sshSession.RequestSubsystem("sftp"); err != nil {
		t.Fatalf("请求 sftp 子系统失败: %v", err)
	}

	sc, err := sftp.NewClientPipe(stdout, pipe)
	if err != nil {
		t.Fatalf("创建 SFTP 客户端失败: %v", err)
	}
	defer sc.Close()

	if _, err := sc.Open("/upload/running-config.txt"); err == nil {
		t.Fatalf("下载（读）必须被拒绝")
	}
}

func TestPathEscapeRejected(t *testing.T) {
	for _, bad := range []string{"../etc/passwd", "/etc/passwd", "/upload/../../etc/passwd", "/other/file"} {
		if _, err := resolvePath(bad); err == nil {
			t.Fatalf("越界路径必须被拒绝: %s", bad)
		}
	}
	for _, ok := range []string{"/upload/a.txt", "a.txt", ".", "/upload/sub/b.txt"} {
		if _, err := resolvePath(ok); err != nil {
			t.Fatalf("合法路径被误拒: %s (%v)", ok, err)
		}
	}
}

func TestSparseAndOversizeWritesRejected(t *testing.T) {
	cfg := testConfig(t)
	cfg.MaxFileSizeBytes = 128
	ctx := uploadContext{deviceID: 1, accountID: 1}

	file, err := newRestrictedFile(cfg, ctx, "a.txt")
	if err != nil {
		t.Fatalf("创建受限文件失败: %v", err)
	}
	defer file.TransferError(errors.New("cleanup"))

	if _, err := file.WriteAt([]byte("x"), 1<<20); err == nil {
		t.Fatalf("稀疏超大偏移必须被拒绝")
	}
	if _, err := file.WriteAt(make([]byte, 200), 0); err == nil {
		t.Fatalf("超过大小上限的写入必须被拒绝")
	}
}

func TestEmptyUploadRejected(t *testing.T) {
	cfg := testConfig(t)
	sink := &recordingSink{}
	ctx := uploadContext{deviceID: 1, accountID: 1, sink: sink}

	file, err := newRestrictedFile(cfg, ctx, "empty.txt")
	if err != nil {
		t.Fatalf("创建受限文件失败: %v", err)
	}
	if err := file.Close(); err == nil {
		t.Fatalf("空文件必须被拒绝")
	}
	if len(sink.uploads) != 0 {
		t.Fatalf("被拒绝的上传不应产生记录")
	}
}

func TestBinaryUploadRejected(t *testing.T) {
	cfg := testConfig(t)
	sink := &recordingSink{}
	ctx := uploadContext{deviceID: 1, accountID: 1, sink: sink}

	file, err := newRestrictedFile(cfg, ctx, "bin.txt")
	if err != nil {
		t.Fatalf("创建受限文件失败: %v", err)
	}
	if _, err := file.WriteAt([]byte{0x00, 0x01, 0x02}, 0); err != nil {
		t.Fatalf("写入失败: %v", err)
	}
	if err := file.Close(); err == nil {
		t.Fatalf("含 NUL 的内容必须被拒绝")
	}
}

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("生成哈希失败: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte("secret-password")); err != nil {
		t.Fatalf("哈希校验失败: %v", err)
	}
}
