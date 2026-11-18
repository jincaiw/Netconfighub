package sftpserver

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// limits enforces the connection, concurrency and timeout bounds described in
// the SFTP configuration. It keeps a global connection semaphore plus a
// per-device upload lock so that one device cannot open parallel uploads.
type limits struct {
	maxConnections int
	idleTimeout    time.Duration
	sessionTimeout time.Duration

	conns chan struct{}

	mu         sync.Mutex
	deviceLock map[uint]struct{}
}

func newLimits(maxConnections int, idleTimeout, sessionTimeout time.Duration) *limits {
	if maxConnections <= 0 {
		maxConnections = 16
	}
	return &limits{
		maxConnections: maxConnections,
		idleTimeout:    idleTimeout,
		sessionTimeout: sessionTimeout,
		conns:          make(chan struct{}, maxConnections),
		deviceLock:     make(map[uint]struct{}),
	}
}

// acquireConn reserves a global connection slot. It reports false when the
// server is already at capacity instead of blocking an accept loop.
func (l *limits) acquireConn() bool {
	select {
	case l.conns <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *limits) releaseConn() {
	select {
	case <-l.conns:
	default:
	}
}

// tryLockDevice reserves the single upload slot for a device. A second
// concurrent upload for the same device is rejected.
func (l *limits) tryLockDevice(deviceID uint) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.deviceLock[deviceID]; exists {
		return false
	}
	l.deviceLock[deviceID] = struct{}{}
	return true
}

func (l *limits) unlockDevice(deviceID uint) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.deviceLock, deviceID)
}

// applyDeadlines bounds a connection by both an idle timeout (refreshed on
// activity) and an absolute session timeout.
func (l *limits) applyDeadlines(conn net.Conn) {
	now := time.Now()
	idle := now.Add(l.idleTimeout)
	abs := now.Add(l.sessionTimeout)
	if !abs.Before(idle) {
		conn.SetDeadline(idle)
		return
	}
	conn.SetDeadline(abs)
}

// touch refreshes the idle deadline after successful I/O while respecting the
// absolute session deadline established at accept time.
func (l *limits) touch(conn net.Conn, sessionDeadline time.Time) {
	next := time.Now().Add(l.idleTimeout)
	if !sessionDeadline.Before(next) {
		conn.SetDeadline(next)
		return
	}
	conn.SetDeadline(sessionDeadline)
}

func (l *limits) String() string {
	return fmt.Sprintf("max_connections=%d idle=%s session=%s", l.maxConnections, l.idleTimeout, l.sessionTimeout)
}
