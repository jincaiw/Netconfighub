package sftpserver

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// handleConn runs the SSH handshake and then serves session channels. Shell,
// exec, PTY and any forwarding request are refused; only the sftp subsystem is
// served.
func (s *Server) handleConn(raw net.Conn) {
	defer raw.Close()

	sessionDeadline := time.Now().Add(s.limits.sessionTimeout)
	s.limits.applyDeadlines(raw)

	serverConn, chans, reqs, err := ssh.NewServerConn(raw, s.sshConfig)
	if err != nil {
		log.Printf("SFTP SSH 握手失败 [remote=%s]: %v", raw.RemoteAddr(), err)
		return
	}
	defer serverConn.Close()

	deviceID, accountID, ok := identityFromPermissions(serverConn.Permissions)
	if !ok {
		log.Printf("SFTP 会话缺少设备身份 [remote=%s]", raw.RemoteAddr())
		return
	}

	// Refuse every global request (tcpip-forward, streamlocal-forward, ...).
	go ssh.DiscardRequests(reqs)

	var sessions sync.WaitGroup
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			log.Printf("SFTP 拒绝非 session 通道 [deviceID=%d, type=%s]", deviceID, newChannel.ChannelType())
			newChannel.Reject(ssh.UnknownChannelType, "仅支持 sftp 会话")
			continue
		}

		// Each session runs independently so one rejected request cannot tear
		// down the whole SSH connection.
		sessions.Add(1)
		go func(ch ssh.NewChannel) {
			defer sessions.Done()
			s.serveSession(ch, deviceID, accountID, raw.RemoteAddr().String(), raw, sessionDeadline)
		}(newChannel)
	}
	sessions.Wait()
}

func (s *Server) serveSession(newChannel ssh.NewChannel, deviceID, accountID uint, remoteAddr string, raw net.Conn, sessionDeadline time.Time) {
	channel, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("SFTP 接受会话通道失败 [deviceID=%d]: %v", deviceID, err)
		return
	}
	defer channel.Close()

	if !awaitSFTPSubsystem(channel, requests) {
		return
	}

	if !s.limits.tryLockDevice(deviceID) {
		log.Printf("SFTP 设备已有上传进行中，拒绝并发上传 [deviceID=%d]", deviceID)
		return
	}
	defer s.limits.unlockDevice(deviceID)

	handler := newUploadHandlers(s.cfg, uploadContext{
		deviceID:   deviceID,
		accountID:  accountID,
		remoteAddr: remoteAddr,
		limits:     s.limits,
		conn:       raw,
		sessionEnd: sessionDeadline,
		sink:       s.sink,
	})

	requestServer := sftp.NewRequestServer(channel, sftp.Handlers{
		// Downloads are unsupported: the reader returns a permission error
		// instead of leaving the handler nil (which would panic).
		FileGet:  deniedReader{},
		FilePut:  handler,
		FileCmd:  handler,
		FileList: handler,
	})

	if err := requestServer.Serve(); err != nil && err != io.EOF {
		log.Printf("SFTP 会话结束 [deviceID=%d]: %v", deviceID, err)
	}
	requestServer.Close()
}

// deniedReader refuses every download attempt.
type deniedReader struct{}

func (deniedReader) Fileread(*sftp.Request) (io.ReaderAt, error) {
	return nil, sftp.ErrSSHFxPermissionDenied
}

// awaitSFTPSubsystem accepts only the subsystem request with payload "sftp"
// and rejects everything else, including shell, exec and pty-req.
func awaitSFTPSubsystem(channel ssh.Channel, requests <-chan *ssh.Request) bool {
	for req := range requests {
		if req.Type != "subsystem" {
			log.Printf("SFTP 拒绝会话请求 [type=%s]", req.Type)
			req.Reply(false, nil)
			continue
		}

		var payload struct {
			Subsystem string
		}
		if err := ssh.Unmarshal(req.Payload, &payload); err != nil || payload.Subsystem != "sftp" {
			log.Printf("SFTP 拒绝非 sftp 子系统请求")
			req.Reply(false, nil)
			continue
		}

		req.Reply(true, nil)
		return true
	}
	return false
}

func identityFromPermissions(perms *ssh.Permissions) (deviceID, accountID uint, ok bool) {
	if perms == nil || perms.Extensions == nil {
		return 0, 0, false
	}
	rawDevice, hasDevice := perms.Extensions[extDeviceID]
	rawAccount, hasAccount := perms.Extensions[extAccountID]
	if !hasDevice || !hasAccount {
		return 0, 0, false
	}
	d, err := strconv.ParseUint(rawDevice, 10, 64)
	if err != nil || d == 0 {
		return 0, 0, false
	}
	a, err := strconv.ParseUint(rawAccount, 10, 64)
	if err != nil {
		return 0, 0, false
	}
	return uint(d), uint(a), true
}
