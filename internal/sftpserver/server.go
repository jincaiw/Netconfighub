package sftpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/netconfighub/netconfighub/internal/model"
	"golang.org/x/crypto/ssh"
)

// CredentialVerifier resolves the device identity behind a set of SFTP
// credentials. Implementations must only accept accounts whose device is
// enabled, not soft-deleted and configured for sftp_push collection.
type CredentialVerifier interface {
	// VerifySFTPCredentials returns the account and its device on success. The
	// error must not distinguish "unknown user" from "bad password" to avoid
	// account enumeration.
	VerifySFTPCredentials(ctx context.Context, username, password string) (*model.DeviceSFTPAccount, *model.Device, error)
}

// UploadSink receives a completed upload. It is invoked once the staging file
// has been atomically moved into the received directory, so a transfer
// failure never produces a record.
type UploadSink interface {
	OnReceived(ctx context.Context, upload ReceivedUpload) error
}

// ReceivedUpload describes a fully received staging file.
type ReceivedUpload struct {
	DeviceID       uint
	AccountID      uint
	StagingPath    string
	Bytes          int64
	Hash           string
	ClientFilename string
	RemoteAddr     string
	ReceivedAt     time.Time
}

// Server is the inbound SFTP listener. It is only started when the SFTP
// configuration is enabled and it accepts nothing but the sftp subsystem.
type Server struct {
	cfg       config.SftpConfig
	verifier  CredentialVerifier
	sink      UploadSink
	sshConfig *ssh.ServerConfig
	hostKey   ssh.Signer
	limits    *limits

	mu       sync.Mutex
	listener net.Listener
	stopping bool
	wg       sync.WaitGroup
	sessions sync.WaitGroup
}

// New builds a Server from configuration. It loads or creates the host key so
// callers can log the fingerprint before starting the listener.
func New(cfg config.SftpConfig, verifier CredentialVerifier, sink UploadSink) (*Server, error) {
	if verifier == nil {
		return nil, fmt.Errorf("缺少凭据校验实现")
	}
	if sink == nil {
		return nil, fmt.Errorf("缺少上传接收实现")
	}

	// A disabled listener must not create directories or keys: doing so would
	// make an unrelated misconfiguration fatal for the whole service.
	var signer ssh.Signer
	if cfg.Enabled {
		var err error
		signer, err = LoadOrCreateHostKey(cfg.HostKeyPath)
		if err != nil {
			return nil, err
		}
	}

	idle, err := parsePositiveDuration(cfg.IdleTimeout, 60*time.Second)
	if err != nil {
		return nil, fmt.Errorf("sftp.idle_timeout 无效: %w", err)
	}
	session, err := parsePositiveDuration(cfg.SessionTimeout, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("sftp.session_timeout 无效: %w", err)
	}

	s := &Server{
		cfg:      cfg,
		verifier: verifier,
		sink:     sink,
		hostKey:  signer,
		limits:   newLimits(cfg.MaxConnections, idle, session),
	}

	sshConfig := &ssh.ServerConfig{
		MaxAuthTries:  3,
		ServerVersion: "SSH-2.0-NetConfigHub-SFTP",
	}
	sshConfig.AddHostKey(signer)
	sshConfig.PasswordCallback = s.passwordCallback
	s.sshConfig = sshConfig

	return s, nil
}

// Fingerprint exposes the host key fingerprint so operators can verify the
// device is talking to the expected server.
func (s *Server) Fingerprint() string {
	return FingerprintSHA256(s.hostKey)
}

// Addr returns the configured listen address.
func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
}

func (s *Server) passwordCallback(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	account, device, err := s.verifier.VerifySFTPCredentials(context.Background(), conn.User(), string(password))
	if err != nil {
		log.Printf("SFTP 认证失败 [user=%s, remote=%s]: %v", conn.User(), conn.RemoteAddr(), err)
		return nil, err
	}

	// The device identity is fixed here and never derived from anything the
	// client sends later, so a filename cannot redirect an upload.
	return &ssh.Permissions{
		Extensions: map[string]string{
			extDeviceID:  fmt.Sprintf("%d", device.ID),
			extAccountID: fmt.Sprintf("%d", account.ID),
			extUsername:  conn.User(),
		},
	}, nil
}

const (
	extDeviceID  = "nch-device-id"
	extAccountID = "nch-account-id"
	extUsername  = "nch-username"
)

// Start begins listening. When SFTP is disabled it is a no-op so the main
// service keeps starting normally.
func (s *Server) Start() error {
	if !s.cfg.Enabled {
		return nil
	}

	ln, err := net.Listen("tcp", s.Addr())
	if err != nil {
		return fmt.Errorf("SFTP 监听失败 [%s]: %w", s.Addr(), err)
	}

	s.mu.Lock()
	if s.stopping {
		s.mu.Unlock()
		ln.Close()
		return fmt.Errorf("SFTP 服务正在停止")
	}
	s.listener = ln
	s.mu.Unlock()

	s.wg.Add(1)
	go s.acceptLoop(ln)

	log.Printf("SFTP 服务启动，监听地址: %s，主机密钥指纹: %s", s.Addr(), s.Fingerprint())
	return nil
}

func (s *Server) acceptLoop(ln net.Listener) {
	defer s.wg.Done()
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.mu.Lock()
			stopping := s.stopping
			s.mu.Unlock()
			if stopping {
				return
			}
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("SFTP 接受连接失败: %v", err)
			continue
		}

		if !s.limits.acquireConn() {
			log.Printf("SFTP 连接数已达上限，拒绝 %s", conn.RemoteAddr())
			conn.Close()
			continue
		}

		s.sessions.Add(1)
		go func() {
			defer s.sessions.Done()
			defer s.limits.releaseConn()
			s.handleConn(conn)
		}()
	}
}

// Stop closes the listener and waits for active sessions to finish or for the
// context deadline, whichever comes first.
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.stopping = true
	ln := s.listener
	s.listener = nil
	s.mu.Unlock()

	if ln != nil {
		if err := ln.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("关闭 SFTP 监听失败: %v", err)
		}
	}
	s.wg.Wait()

	done := make(chan struct{})
	go func() {
		s.sessions.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func parsePositiveDuration(raw string, fallback time.Duration) (time.Duration, error) {
	if raw == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, fmt.Errorf("必须为正时长")
	}
	return d, nil
}

// Status is the runtime state reported to the management API.
type Status struct {
	Enabled     bool   `json:"enabled"`
	ListenAddr  string `json:"listen_addr"`
	Fingerprint string `json:"fingerprint"`
	ActiveConns int    `json:"active_conns"`
	MaxConns    int    `json:"max_conns"`
}

// Status reports the current listener state. It is safe to call while the
// server is running or stopped.
func (s *Server) Status() Status {
	s.mu.Lock()
	addr := ""
	if s.listener != nil {
		addr = s.listener.Addr().String()
	}
	s.mu.Unlock()

	active := 0
	if s.limits != nil {
		active = len(s.limits.conns)
	}

	return Status{
		Enabled:     s.cfg.Enabled,
		ListenAddr:  addr,
		Fingerprint: s.Fingerprint(),
		ActiveConns: active,
		MaxConns:    s.cfg.MaxConnections,
	}
}
