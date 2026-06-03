package notifier

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type Notifier interface {
	Send(ctx context.Context, subject, body string) error
}

type NotifyConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	Type       string `yaml:"type" json:"type"`
	WebhookURL string `yaml:"webhook_url" json:"webhook_url"`
	SMTPHost   string `yaml:"smtp_host" json:"smtp_host"`
	SMTPPort   int    `yaml:"smtp_port" json:"smtp_port"`
	SMTPUser   string `yaml:"smtp_user" json:"smtp_user"`
	SMTPPass   string `yaml:"smtp_pass" json:"smtp_pass"`
	FromAddr   string `yaml:"from_addr" json:"from_addr"`
	ToAddrs    string `yaml:"to_addrs" json:"to_addrs"`
}

type WebhookNotifier struct {
	url string
}

func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{url: url}
}

func (n *WebhookNotifier) Send(ctx context.Context, subject, body string) error {
	payload := map[string]string{
		"subject":   subject,
		"body":      body,
		"timestamp": time.Now().Format(time.RFC3339),
		"source":    "NetConfigHub",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化 webhook 载荷失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("创建 webhook 请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送 webhook 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook 返回非成功状态码: %d", resp.StatusCode)
	}
	return nil
}

type EmailNotifier struct {
	host     string
	port     int
	user     string
	password string
	fromAddr string
	toAddrs  []string
}

func NewEmailNotifier(config NotifyConfig) *EmailNotifier {
	return &EmailNotifier{
		host:     config.SMTPHost,
		port:     config.SMTPPort,
		user:     config.SMTPUser,
		password: config.SMTPPass,
		fromAddr: config.FromAddr,
		toAddrs:  strings.Split(config.ToAddrs, ","),
	}
}

func (n *EmailNotifier) Send(ctx context.Context, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", n.host, n.port)

	msg := fmt.Sprintf("From: %s\r\n", n.fromAddr)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(n.toAddrs, ","))
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	msg += "\r\n"
	msg += body

	if n.port == 465 {
		tlsConfig := &tls.Config{
			ServerName: n.host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("TLS 连接失败: %w", err)
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, n.host)
		if err != nil {
			return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
		}
		defer c.Close()

		if err := c.Auth(smtp.PlainAuth("", n.user, n.password, n.host)); err != nil {
			return fmt.Errorf("SMTP 认证失败: %w", err)
		}
		if err := c.Mail(n.fromAddr); err != nil {
			return fmt.Errorf("设置发件人失败: %w", err)
		}
		for _, to := range n.toAddrs {
			if err := c.Rcpt(to); err != nil {
				return fmt.Errorf("设置收件人失败: %w", err)
			}
		}
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("获取数据写入器失败: %w", err)
		}
		if _, err := w.Write([]byte(msg)); err != nil {
			return fmt.Errorf("写入邮件内容失败: %w", err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("关闭写入器失败: %w", err)
		}
		return c.Quit()
	}

	auth := smtp.PlainAuth("", n.user, n.password, n.host)
	if err := smtp.SendMail(addr, auth, n.fromAddr, n.toAddrs, []byte(msg)); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}

func NewNotifier(config NotifyConfig) Notifier {
	if !config.Enabled {
		return nil
	}
	switch config.Type {
	case "webhook":
		return NewWebhookNotifier(config.WebhookURL)
	case "email":
		return NewEmailNotifier(config)
	default:
		return nil
	}
}
