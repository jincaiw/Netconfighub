package hook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func ExecuteExecHook(hook model.Hook, payload map[string]interface{}) {
	cfg, err := ParseConfig(hook.Config)
	if err != nil {
		log.Printf("解析 Exec Hook 配置失败 [hookID=%d]: %v", hook.ID, err)
		return
	}
	payloadJSON, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", cfg.Command)
	cmd.Stdin = bytes.NewReader(payloadJSON)
	cmd.Env = append(cmd.Environ(),
		fmt.Sprintf("NCH_HOOK_EVENT=%s", payload["event"]),
		fmt.Sprintf("NCH_HOOK_DEVICE=%s", payload["device_name"]),
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

func ExecuteWebhookHook(hook model.Hook, payload map[string]interface{}) {
	cfg, err := ParseConfig(hook.Config)
	if err != nil {
		log.Printf("解析 Webhook Hook 配置失败 [hookID=%d]: %v", hook.ID, err)
		return
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
		req.Header.Set("X-NCH-Signature", cfg.Secret)
	}
	client := &http.Client{Timeout: 10 * time.Second}
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
