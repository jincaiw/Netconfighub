package collector

import (
	"bytes"
	"fmt"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
	"golang.org/x/crypto/ssh"
)

type SSHCollector struct {
	connectTimeout time.Duration
	commandTimeout time.Duration
}

func NewSSHCollector(timeout time.Duration) *SSHCollector {
	return &SSHCollector{
		connectTimeout: timeout,
		commandTimeout: 30 * time.Second,
	}
}

func (c *SSHCollector) Collect(device model.Device) (string, error) {
	adapter, err := GetAdapter(device.Vendor)
	if err != nil {
		return "", fmt.Errorf("获取厂商适配器失败: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		Timeout:         c.connectTimeout,
		User:            device.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if device.SSHKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(device.SSHKey))
		if err != nil {
			return "", fmt.Errorf("解析 SSH 私钥失败: %w", err)
		}
		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(device.Password)}
	}

	addr := fmt.Sprintf("%s:%d", device.IP, device.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return "", fmt.Errorf("SSH 连接失败 [%s]: %w", addr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建 SSH 会话失败: %w", err)
	}
	defer session.Close()

	command := adapter.ShowRunningConfig()
	needEnable := device.EnablePassword != ""

	if needEnable {
		stdin, err := session.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("获取 stdin 管道失败: %w", err)
		}

		var buf bytes.Buffer
		session.Stdout = &buf
		session.Stderr = &buf

		if err := session.Shell(); err != nil {
			return "", fmt.Errorf("启动 shell 失败: %w", err)
		}

		time.Sleep(300 * time.Millisecond)

		enableCmd := adapter.EnableCommand()
		if enableCmd != "" {
			fmt.Fprintf(stdin, "%s\n", enableCmd)
			time.Sleep(300 * time.Millisecond)
			fmt.Fprintf(stdin, "%s\n", device.EnablePassword)
			time.Sleep(500 * time.Millisecond)
		}

		pagingCmds := getDisablePagingCommands(device.Vendor)
		for _, cmd := range pagingCmds {
			fmt.Fprintf(stdin, "%s\n", cmd)
			time.Sleep(300 * time.Millisecond)
		}

		fmt.Fprintf(stdin, "%s\n", command)

		stdin.Close()

		done := make(chan struct{})
		go func() {
			session.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(c.commandTimeout):
			session.Close()
			return "", fmt.Errorf("命令执行超时")
		}

		config, err := adapter.ParseOutput(buf.String())
		if err != nil {
			return "", fmt.Errorf("解析输出失败: %w", err)
		}
		return config, nil
	}

	type cmdResult struct {
		output []byte
		err    error
	}
	done := make(chan cmdResult, 1)
	go func() {
		output, err := session.CombinedOutput(command)
		done <- cmdResult{output: output, err: err}
	}()

	select {
	case r := <-done:
		if r.err != nil {
			return "", fmt.Errorf("执行命令失败 [%s]: %w", command, r.err)
		}
		config, err := adapter.ParseOutput(string(r.output))
		if err != nil {
			return "", fmt.Errorf("解析输出失败: %w", err)
		}
		return config, nil
	case <-time.After(c.commandTimeout):
		session.Close()
		return "", fmt.Errorf("命令执行超时 [%s]", command)
	}
}
