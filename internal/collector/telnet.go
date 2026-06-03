package collector

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
)

type TelnetCollector struct {
	connectTimeout time.Duration
	readTimeout    time.Duration
}

func NewTelnetCollector(timeout time.Duration) *TelnetCollector {
	return &TelnetCollector{
		connectTimeout: timeout,
		readTimeout:    30 * time.Second,
	}
}

func (c *TelnetCollector) Collect(device model.Device) (string, error) {
	adapter, err := GetAdapter(device.Vendor)
	if err != nil {
		return "", fmt.Errorf("获取厂商适配器失败: %w", err)
	}

	addr := net.JoinHostPort(device.IP, fmt.Sprintf("%d", device.Port))
	conn, err := net.DialTimeout("tcp", addr, c.connectTimeout)
	if err != nil {
		return "", fmt.Errorf("Telnet 连接失败 [%s]: %w", addr, err)
	}
	defer conn.Close()

	promptRe := regexp.MustCompile(adapter.PromptPattern())
	usernameRe := regexp.MustCompile(`(?i)(username|login|user)[:]\s*`)
	passwordRe := regexp.MustCompile(`(?i)password[:]\s*`)

	if _, err := c.readUntil(conn, usernameRe, 10*time.Second); err != nil {
		return "", fmt.Errorf("等待登录提示失败: %w", err)
	}
	c.sendLine(conn, device.Username)

	if _, err := c.readUntil(conn, passwordRe, 10*time.Second); err != nil {
		return "", fmt.Errorf("等待密码提示失败: %w", err)
	}
	c.sendLine(conn, device.Password)

	if _, err := c.readUntil(conn, promptRe, 10*time.Second); err != nil {
		return "", fmt.Errorf("等待命令提示符失败: %w", err)
	}

	if device.EnablePassword != "" {
		enableCmd := adapter.EnableCommand()
		if enableCmd != "" {
			c.sendLine(conn, enableCmd)
			passwordRe := regexp.MustCompile(`(?i)password[:]\s*`)
			c.readUntil(conn, passwordRe, 5*time.Second)
			c.sendLine(conn, device.EnablePassword)
			c.readUntil(conn, promptRe, 10*time.Second)
		}
	}

	pagingCmds := getDisablePagingCommands(device.Vendor)
	for _, cmd := range pagingCmds {
		c.sendLine(conn, cmd)
		c.readUntil(conn, promptRe, 10*time.Second)
	}

	command := adapter.ShowRunningConfig()
	c.sendLine(conn, command)

	output, err := c.readUntil(conn, promptRe, c.readTimeout)
	if err != nil {
		return "", fmt.Errorf("读取配置输出失败: %w", err)
	}

	c.sendLine(conn, "exit")

	config, err := adapter.ParseOutput(output)
	if err != nil {
		return "", fmt.Errorf("解析输出失败: %w", err)
	}

	return config, nil
}

func (c *TelnetCollector) sendLine(conn net.Conn, line string) {
	fmt.Fprintf(conn, "%s\r\n", line)
	time.Sleep(100 * time.Millisecond)
}

func (c *TelnetCollector) readUntil(conn net.Conn, pattern *regexp.Regexp, timeout time.Duration) (string, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	var output bytes.Buffer
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if n > 0 {
			data := buf[:n]
			clean, responses := processTelnetData(data)
			if len(responses) > 0 {
				conn.Write(responses)
			}
			output.Write(clean)
			if pattern.MatchString(output.String()) {
				return output.String(), nil
			}
		}
		if err != nil {
			if err == io.EOF {
				return output.String(), fmt.Errorf("连接关闭")
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return output.String(), fmt.Errorf("读取超时")
			}
			return output.String(), err
		}
	}
}

func processTelnetData(data []byte) (clean []byte, responses []byte) {
	var cleanBuf bytes.Buffer
	var respBuf bytes.Buffer
	i := 0
	for i < len(data) {
		if data[i] == 0xFF && i+1 < len(data) {
			if data[i+1] == 0xFF {
				cleanBuf.WriteByte(0xFF)
				i += 2
				continue
			}
			if i+2 < len(data) {
				cmd := data[i+1]
				opt := data[i+2]
				switch cmd {
				case 0xFD:
					respBuf.Write([]byte{0xFF, 0xFC, opt})
				case 0xFB:
					respBuf.Write([]byte{0xFF, 0xFE, opt})
				}
				i += 3
			} else {
				i += 2
			}
		} else {
			cleanBuf.WriteByte(data[i])
			i++
		}
	}
	return cleanBuf.Bytes(), respBuf.Bytes()
}
