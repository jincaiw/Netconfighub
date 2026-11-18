package sftpserver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/pkg/sftp"
)

const (
	virtualRoot = "/upload"
)

type uploadContext struct {
	deviceID   uint
	accountID  uint
	remoteAddr string
	limits     *limits
	conn       net.Conn
	sessionEnd time.Time
	sink       UploadSink
}

// uploadHandlers implements the three SFTP server handler interfaces. It
// deliberately does not implement sftp.FileReader, so download requests are
// rejected by the request server before reaching any handler.
type uploadHandlers struct {
	cfg config.SftpConfig
	ctx uploadContext

	mu       sync.Mutex
	inflight map[string]*restrictedFile
}

func newUploadHandlers(cfg config.SftpConfig, ctx uploadContext) *uploadHandlers {
	return &uploadHandlers{
		cfg:      cfg,
		ctx:      ctx,
		inflight: make(map[string]*restrictedFile),
	}
}

var errPermissionDenied = sftp.ErrSSHFxPermissionDenied

// resolvePath maps a client supplied path onto the fixed virtual root and
// rejects anything escaping it. The client filename is only a label; it never
// determines the owning device.
func resolvePath(clientPath string) (string, error) {
	raw := strings.TrimSpace(clientPath)
	if raw == "" || strings.ContainsRune(raw, 0) {
		return "", errPermissionDenied
	}
	// Reject traversal before cleaning: path.Clean would otherwise collapse
	// "../x" into something that looks harmless.
	for _, part := range strings.Split(raw, "/") {
		if part == ".." {
			return "", errPermissionDenied
		}
	}

	// Absolute paths must already live under the virtual root; relative paths
	// are resolved against it, which is how clients behave after Realpath.
	var cleaned string
	if strings.HasPrefix(raw, "/") {
		cleaned = path.Clean(raw)
	} else {
		cleaned = path.Clean(virtualRoot + "/" + raw)
	}
	if cleaned != virtualRoot && !strings.HasPrefix(cleaned, virtualRoot+"/") {
		return "", errPermissionDenied
	}

	rel := strings.TrimPrefix(strings.TrimPrefix(cleaned, virtualRoot), "/")
	if rel == "" || rel == "." {
		return virtualRoot, nil
	}
	return rel, nil
}

// Filewrite handles Put/Open. Note the pkg/sftp interface name is Filewrite.
func (h *uploadHandlers) Filewrite(request *sftp.Request) (io.WriterAt, error) {
	rel, err := resolvePath(request.Filepath)
	if err != nil {
		return nil, err
	}
	if rel == virtualRoot {
		return nil, errPermissionDenied
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if _, busy := h.inflight[rel]; busy {
		return nil, sftp.ErrSSHFxFailure
	}

	file, err := newRestrictedFile(h.cfg, h.ctx, rel)
	if err != nil {
		return nil, err
	}
	h.inflight[rel] = file
	return file, nil
}

// Filecmd handles rename-based completion and refuses every mutation command.
func (h *uploadHandlers) Filecmd(request *sftp.Request) error {
	switch request.Method {
	case "Rename", "PosixRename":
		return h.finish(request)
	case "Setstat", "Remove", "Rmdir", "Mkdir", "Link", "Symlink", "StatVFS":
		return errPermissionDenied
	default:
		return errPermissionDenied
	}
}

// Filelist exposes only minimal metadata. Directory listing is empty so a
// client cannot enumerate other uploads.
func (h *uploadHandlers) Filelist(request *sftp.Request) (sftp.ListerAt, error) {
	switch request.Method {
	case "Stat", "Lstat":
		if _, err := resolvePath(request.Filepath); err != nil {
			return nil, err
		}
		return listerAt{&statInfo{
			name:  path.Base(request.Filepath),
			mode:  os.FileMode(0o600),
			mtime: time.Now(),
		}}, nil
	case "List":
		if _, err := resolvePath(request.Filepath); err != nil {
			return nil, err
		}
		return listerAt{}, nil
	default:
		return nil, errPermissionDenied
	}
}

// finish completes the rename-then-commit path used by clients that upload to
// a temporary name first.
func (h *uploadHandlers) finish(request *sftp.Request) error {
	fromRel, err := resolvePath(request.Filepath)
	if err != nil {
		return err
	}
	if _, err := resolvePath(request.Target); err != nil {
		return err
	}

	h.mu.Lock()
	file, ok := h.inflight[fromRel]
	if ok {
		delete(h.inflight, fromRel)
	}
	h.mu.Unlock()

	if !ok {
		return errPermissionDenied
	}
	return file.Close()
}

// listerAt adapts a fixed slice of file info to sftp.ListerAt.
type listerAt []os.FileInfo

func (l listerAt) ListAt(f []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l)) {
		return 0, io.EOF
	}
	n := copy(f, l[offset:])
	if n < len(f) {
		return n, io.EOF
	}
	return n, nil
}

// statInfo is a read-only os.FileInfo used for Stat/Lstat responses.
type statInfo struct {
	name  string
	mode  os.FileMode
	mtime time.Time
}

func (s *statInfo) Name() string       { return s.name }
func (s *statInfo) Size() int64        { return 0 }
func (s *statInfo) Mode() os.FileMode  { return s.mode }
func (s *statInfo) ModTime() time.Time { return s.mtime }
func (s *statInfo) IsDir() bool        { return s.mode.IsDir() }
func (s *statInfo) Sys() interface{}   { return nil }

// restrictedFile streams one upload into a 0600 staging file, enforcing the
// configured size bound and offset sanity checks before the file is moved into
// the received directory.
type restrictedFile struct {
	cfg config.SftpConfig
	ctx uploadContext
	rel string

	file *os.File
	tmp  string

	mu      sync.Mutex
	written int64
	closed  bool
}

func newRestrictedFile(cfg config.SftpConfig, ctx uploadContext, rel string) (*restrictedFile, error) {
	tmpDir := filepath.Join(cfg.StagingDir, "tmp")
	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		return nil, fmt.Errorf("创建暂存目录失败: %w", err)
	}

	tmp, err := os.CreateTemp(tmpDir, "upload-*.part")
	if err != nil {
		return nil, fmt.Errorf("创建暂存文件失败: %w", err)
	}
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("设置暂存文件权限失败: %w", err)
	}

	return &restrictedFile{
		cfg:  cfg,
		ctx:  ctx,
		rel:  rel,
		file: tmp,
		tmp:  tmp.Name(),
	}, nil
}

// WriteAt enforces the size bound and rejects sparse writes far beyond it
// before touching the disk.
func (f *restrictedFile) WriteAt(p []byte, off int64) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return 0, sftp.ErrSSHFxFailure
	}

	maxSize := f.cfg.MaxFileSizeBytes
	if maxSize <= 0 {
		maxSize = 10 << 20
	}
	if off < 0 || off > maxSize {
		return 0, sftp.ErrSSHFxFailure
	}
	if off+int64(len(p)) > maxSize {
		return 0, sftp.ErrSSHFxFailure
	}

	if f.ctx.limits != nil && f.ctx.conn != nil {
		f.ctx.limits.touch(f.ctx.conn, f.ctx.sessionEnd)
	}

	n, err := f.file.WriteAt(p, off)
	if n > 0 {
		f.written += int64(n)
	}
	return n, err
}

// Close validates the received content and moves it into the received
// directory. Any validation failure discards the staging file.
func (f *restrictedFile) Close() error {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return nil
	}
	f.closed = true
	written := f.written
	f.mu.Unlock()

	defer func() {
		if f.file != nil {
			f.file.Close()
		}
	}()

	if written == 0 {
		os.Remove(f.tmp)
		return sftp.ErrSSHFxFailure
	}

	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		os.Remove(f.tmp)
		return sftp.ErrSSHFxFailure
	}

	// A device configuration is text; reject binary payloads.
	buf := make([]byte, 8192)
	for {
		n, err := f.file.Read(buf)
		if n > 0 && bytesContainNUL(buf[:n]) {
			os.Remove(f.tmp)
			return sftp.ErrSSHFxFailure
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(f.tmp)
			return sftp.ErrSSHFxFailure
		}
	}

	if err := f.file.Sync(); err != nil {
		os.Remove(f.tmp)
		return sftp.ErrSSHFxFailure
	}

	receivedDir := filepath.Join(f.cfg.StagingDir, "received", fmt.Sprintf("%d", f.ctx.deviceID))
	if err := os.MkdirAll(receivedDir, 0o700); err != nil {
		os.Remove(f.tmp)
		return sftp.ErrSSHFxFailure
	}

	finalPath := filepath.Join(receivedDir, fmt.Sprintf("%d-%d.txt", time.Now().UnixNano(), f.ctx.accountID))
	if err := os.Rename(f.tmp, finalPath); err != nil {
		os.Remove(f.tmp)
		return sftp.ErrSSHFxFailure
	}
	if err := os.Chmod(finalPath, 0o600); err != nil {
		log.Printf("设置接收文件权限失败 [path=%s]: %v", finalPath, err)
	}

	// The hash must be computed from the completed file: pkg/sftp's request
	// server dispatches WRITE packets from several workers, so incremental
	// hashing on WriteAt would depend on arrival order and break verification.
	contentHash, hashErr := hashFile(finalPath)
	if hashErr != nil {
		os.Remove(finalPath)
		return sftp.ErrSSHFxFailure
	}

	received := ReceivedUpload{
		DeviceID:       f.ctx.deviceID,
		AccountID:      f.ctx.accountID,
		StagingPath:    finalPath,
		Bytes:          written,
		Hash:           contentHash,
		ClientFilename: path.Base(f.rel),
		RemoteAddr:     f.ctx.remoteAddr,
		ReceivedAt:     time.Now(),
	}

	if f.ctx.sink != nil {
		if err := f.ctx.sink.OnReceived(context.Background(), received); err != nil {
			log.Printf("记录 SFTP 上传失败 [deviceID=%d]: %v", f.ctx.deviceID, err)
			return sftp.ErrSSHFxFailure
		}
	}
	return nil
}

// TransferError cleans up a partially written file when the connection breaks.
func (f *restrictedFile) TransferError(err error) {
	f.mu.Lock()
	f.closed = true
	f.mu.Unlock()
	if f.file != nil {
		f.file.Close()
	}
	if f.tmp != "" {
		os.Remove(f.tmp)
	}
}

// hashFile returns the hex SHA256 of a file's contents.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func bytesContainNUL(b []byte) bool {
	for _, c := range b {
		if c == 0 {
			return true
		}
	}
	return false
}
