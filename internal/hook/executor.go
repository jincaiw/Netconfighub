package hook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
)

type HookConfig struct {
	Command string `json:"command"`
	Remote  string `json:"remote"`
	Branch  string `json:"branch"`
	URL     string `json:"url"`
	Secret  string `json:"secret"`
}

func ParseConfig(configStr string) (*HookConfig, error) {
	var cfg HookConfig
	if err := json.Unmarshal([]byte(configStr), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func containsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// ExecuteExecHook runs a configured executable directly (no shell) and only
// when the operator has explicitly enabled exec hooks and allow-listed the
// absolute binary path.
func ExecuteExecHook(hook model.Hook, payload map[string]interface{}) {
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("NCH_ALLOW_EXEC_HOOKS")), "true") {
		log.Printf("Exec Hook 已禁用 [hookID=%d]（需显式设置 NCH_ALLOW_EXEC_HOOKS=true）", hook.ID)
		return
	}
	cfg, err := ParseConfig(hook.Config)
	if err != nil {
		log.Printf("解析 Exec Hook 配置失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	parts := strings.Fields(strings.TrimSpace(cfg.Command))
	if len(parts) == 0 {
		log.Printf("Exec Hook 命令为空 [hookID=%d]", hook.ID)
		return
	}
	exe, err := filepath.Abs(parts[0])
	if err != nil {
		log.Printf("解析 Exec Hook 可执行文件路径失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	allowlist := splitList(os.Getenv("NCH_EXEC_HOOK_ALLOWLIST"))
	if !containsString(allowlist, exe) {
		log.Printf("Exec Hook 可执行文件不在白名单 [hookID=%d]: %s", hook.ID, exe)
		return
	}

	payloadJSON, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, parts[1:]...)
	cmd.Stdin = bytes.NewReader(payloadJSON)
	cmd.Env = append(cmd.Environ(),
		fmt.Sprintf("NCH_HOOK_EVENT=%v", payload["event"]),
		fmt.Sprintf("NCH_HOOK_DEVICE=%v", payload["device_name"]),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Exec Hook 执行失败 [hookID=%d]: %v, output: %s", hook.ID, err, string(output))
		return
	}
	log.Printf("Exec Hook 执行成功 [hookID=%d]", hook.ID)
}

func ExecuteGitPushHook(hook model.Hook, payload map[string]interface{}) {
	cfg, err := ParseConfig(hook.Config)
	if err != nil {
		log.Printf("解析 Git Push Hook 配置失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	remote := cfg.Remote
	if remote == "" {
		remote = "origin"
	}
	branch := cfg.Branch
	if branch == "" {
		branch = "main"
	}
	cmd := exec.Command("git", "push", remote, branch)
	repoDir := os.Getenv("NCH_GIT_REPO_DIR")
	if repoDir != "" {
		cmd.Dir = repoDir
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Git Push Hook 执行失败 [hookID=%d]: %v, output: %s", hook.ID, err, string(output))
		return
	}
	log.Printf("Git Push Hook 执行成功 [hookID=%d]", hook.ID)
}

// lookupWebhookIP is replaceable in tests so DNS resolution (and therefore
// DNS rebinding protection) can be verified without real network access.
var lookupWebhookIP = net.LookupIP

// webhookNetDial is replaceable in tests to observe the exact dial target.
var webhookNetDial = (&net.Dialer{Timeout: 10 * time.Second}).DialContext

func isRestrictedWebhookIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified()
}

// validateWebhookTarget parses and validates a webhook URL and resolves its
// host to every A/AAAA record. Any resolved IP that is loopback/private/
// link-local/unspecified is rejected (fail closed) unless the host is
// explicitly allow-listed via NCH_WEBHOOK_ALLOW_PRIVATE. The returned IPs are
// the validated set that the request is pinned to, preventing DNS rebinding
// between validation and dialing.
func validateWebhookTarget(raw string) ([]net.IP, *url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("Webhook URL 解析失败: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, nil, fmt.Errorf("Webhook URL 仅允许 http/https")
	}
	host := u.Hostname()
	if host == "" {
		return nil, nil, fmt.Errorf("Webhook URL 缺少主机")
	}
	allowPrivate := splitList(os.Getenv("NCH_WEBHOOK_ALLOW_PRIVATE"))

	if strings.EqualFold(host, "localhost") {
		if !containsString(allowPrivate, host) {
			return nil, nil, fmt.Errorf("Webhook 禁止回环目标")
		}
		ips, err := lookupWebhookIP(host)
		if err != nil || len(ips) == 0 {
			return nil, nil, fmt.Errorf("Webhook 域名解析失败: %s", host)
		}
		return ips, u, nil
	}

	var ips []net.IP
	if ip := net.ParseIP(host); ip != nil {
		ips = []net.IP{ip}
	} else {
		resolved, err := lookupWebhookIP(host)
		if err != nil || len(resolved) == 0 {
			// Fail closed: an unresolvable domain must never be dialed.
			return nil, nil, fmt.Errorf("Webhook 域名解析失败: %s", host)
		}
		ips = resolved
	}
	for _, ip := range ips {
		if isRestrictedWebhookIP(ip) && !containsString(allowPrivate, host) {
			return nil, nil, fmt.Errorf("Webhook 禁止私网/回环/链路本地目标")
		}
	}
	return ips, u, nil
}

// validateWebhookURL is a thin compatibility wrapper around
// validateWebhookTarget.
func validateWebhookURL(raw string) error {
	_, _, err := validateWebhookTarget(raw)
	return err
}

// webhookDialer returns a DialContext that only connects to the validated IP
// set, trying each in order. The original hostname (and therefore the HTTP
// Host header and TLS SNI) is preserved because req.URL is left unchanged.
func webhookDialer(ips []net.IP, port string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var lastErr error
		for _, ip := range ips {
			target := net.JoinHostPort(ip.String(), port)
			conn, err := webhookNetDial(ctx, network, target)
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}
		if lastErr == nil {
			lastErr = fmt.Errorf("Webhook 无可用的已校验 IP")
		}
		return nil, lastErr
	}
}

func ExecuteWebhookHook(hook model.Hook, payload map[string]interface{}) {
	cfg, err := ParseConfig(hook.Config)
	if err != nil {
		log.Printf("解析 Webhook Hook 配置失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	ips, targetURL, err := validateWebhookTarget(cfg.URL)
	if err != nil {
		log.Printf("Webhook Hook 目标被拒绝 [hookID=%d]: %v", hook.ID, err)
		return
	}
	port := targetURL.Port()
	if port == "" {
		if targetURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("序列化 Webhook 载荷失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewReader(data))
	if err != nil {
		log.Printf("创建 Webhook 请求失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.Secret != "" {
		mac := hmac.New(sha256.New, []byte(cfg.Secret))
		_, _ = mac.Write(data)
		req.Header.Set("X-NCH-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)))
	}
	transport := &http.Transport{DialContext: webhookDialer(ips, port)}
	client := &http.Client{Timeout: 10 * time.Second, Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送 Webhook 请求失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Webhook 返回非成功状态码 [hookID=%d]: %d", hook.ID, resp.StatusCode)
		return
	}
	log.Printf("Webhook Hook 执行成功 [hookID=%d]", hook.ID)
}
