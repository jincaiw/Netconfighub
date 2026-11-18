package sftpserver

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// LoadOrCreateHostKey loads the server host key from path, generating and
// persisting a new ed25519 key when the file does not exist yet. The key is
// created with 0600 permissions inside a 0700 directory so that a device can
// verify a stable fingerprint across restarts.
func LoadOrCreateHostKey(path string) (ssh.Signer, error) {
	if path == "" {
		return nil, fmt.Errorf("主机密钥路径不能为空")
	}

	if data, err := os.ReadFile(path); err == nil {
		signer, err := ssh.ParsePrivateKey(data)
		if err != nil {
			return nil, fmt.Errorf("解析主机密钥失败: %w", err)
		}
		return signer, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("读取主机密钥失败: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("创建主机密钥目录失败: %w", err)
	}

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成主机密钥失败: %w", err)
	}

	block, err := ssh.MarshalPrivateKey(privateKey, "netconfighub-sftp")
	if err != nil {
		return nil, fmt.Errorf("编码主机密钥失败: %w", err)
	}
	encoded := pem.EncodeToMemory(block)

	// Write through a temporary file so a crash cannot leave a truncated key.
	tmp, err := os.CreateTemp(dir, "host_key-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("创建主机密钥临时文件失败: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return nil, fmt.Errorf("设置主机密钥权限失败: %w", err)
	}
	if _, err := tmp.Write(encoded); err != nil {
		tmp.Close()
		return nil, fmt.Errorf("写入主机密钥失败: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return nil, fmt.Errorf("同步主机密钥失败: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("关闭主机密钥临时文件失败: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return nil, fmt.Errorf("保存主机密钥失败: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("解析新生成的主机密钥失败: %w", err)
	}
	return signer, nil
}

// FingerprintSHA256 returns the SHA256 fingerprint of the public key in the
// same "SHA256:<base64>" form printed by ssh-keygen -lf.
func FingerprintSHA256(signer ssh.Signer) string {
	if signer == nil {
		return ""
	}
	sum := sha256.Sum256(signer.PublicKey().Marshal())
	return "SHA256:" + base64.RawStdEncoding.EncodeToString(sum[:])
}
