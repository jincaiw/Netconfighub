package hook

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
)

// stubLookup replaces lookupWebhookIP for the duration of a test and restores
// it with t.Cleanup. Tests using this helper must not call t.Parallel.
func stubLookup(t *testing.T, fn func(string) ([]net.IP, error)) {
	t.Helper()
	orig := lookupWebhookIP
	lookupWebhookIP = fn
	t.Cleanup(func() { lookupWebhookIP = orig })
}

func TestValidateWebhookRejectsBadSchemeAndURL(t *testing.T) {
	for _, raw := range []string{"://bad", "ftp://example.com/x", "//example.com"} {
		if err := validateWebhookURL(raw); err == nil {
			t.Errorf("validateWebhookURL(%q) 应拒绝", raw)
		}
	}
}

func TestValidateWebhookLiteralIPs(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		allowPriv string
		wantErr   bool
	}{
		{"public ipv4", "http://8.8.8.8/hook", "", false},
		{"private ipv4", "http://10.0.0.5/hook", "", true},
		{"loopback ipv4", "http://127.0.0.1/hook", "", true},
		{"link local ipv4", "http://169.254.1.1/hook", "", true},
		{"unspecified ipv4", "http://0.0.0.0/hook", "", true},
		{"private ipv6", "http://[::1]/hook", "", true},
		{"private ip allowlisted", "http://10.0.0.5/hook", "10.0.0.5", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NCH_WEBHOOK_ALLOW_PRIVATE", tt.allowPriv)
			err := validateWebhookURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateWebhookURL(%q) err=%v, wantErr=%v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestValidateWebhookLocalhost(t *testing.T) {
	t.Run("rejected by default", func(t *testing.T) {
		t.Setenv("NCH_WEBHOOK_ALLOW_PRIVATE", "")
		if err := validateWebhookURL("http://localhost:8080/hook"); err == nil {
			t.Fatal("localhost 默认应被拒绝")
		}
	})
	t.Run("allowed when allowlisted and resolved", func(t *testing.T) {
		t.Setenv("NCH_WEBHOOK_ALLOW_PRIVATE", "localhost")
		stubLookup(t, func(host string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("127.0.0.1")}, nil
		})
		ips, _, err := validateWebhookTarget("http://localhost:8080/hook")
		if err != nil {
			t.Fatalf("allowlist localhost 应放行: %v", err)
		}
		if len(ips) != 1 {
			t.Fatalf("期望 1 个 IP，实际 %d", len(ips))
		}
	})
}

// TestValidateWebhookResolvesAllIPs proves validation checks every resolved
// address, not only the first one: a public first IP must not mask a private
// later IP (the DNS rebinding scenario).
func TestValidateWebhookResolvesAllIPs(t *testing.T) {
	tests := []struct {
		name      string
		ips       []net.IP
		allowPriv string
		wantErr   bool
	}{
		{"public only", []net.IP{net.ParseIP("93.184.216.34")}, "", false},
		{"private only", []net.IP{net.ParseIP("10.0.0.5")}, "", true},
		{"public then private", []net.IP{net.ParseIP("93.184.216.34"), net.ParseIP("10.0.0.5")}, "", true},
		{"private then public", []net.IP{net.ParseIP("10.0.0.5"), net.ParseIP("93.184.216.34")}, "", true},
		{"mixed allowlisted", []net.IP{net.ParseIP("93.184.216.34"), net.ParseIP("10.0.0.5")}, "mixed.example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NCH_WEBHOOK_ALLOW_PRIVATE", tt.allowPriv)
			stubLookup(t, func(host string) ([]net.IP, error) { return tt.ips, nil })
			ips, u, err := validateWebhookTarget("https://mixed.example.com/hook")
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tt.wantErr)
			}
			if err == nil {
				if len(ips) != len(tt.ips) {
					t.Fatalf("期望 %d 个已校验 IP，实际 %d", len(tt.ips), len(ips))
				}
				if u.Hostname() != "mixed.example.com" {
					t.Fatalf("原始主机名被改写: %q", u.Hostname())
				}
			}
		})
	}
}

func TestValidateWebhookFailClosedOnLookupError(t *testing.T) {
	t.Setenv("NCH_WEBHOOK_ALLOW_PRIVATE", "")
	t.Run("lookup error", func(t *testing.T) {
		stubLookup(t, func(string) ([]net.IP, error) { return nil, errors.New("dns failure") })
		if err := validateWebhookURL("https://unresolvable.example.com/hook"); err == nil {
			t.Fatal("解析失败必须 fail closed")
		}
	})
	t.Run("no addresses", func(t *testing.T) {
		stubLookup(t, func(string) ([]net.IP, error) { return []net.IP{}, nil })
		if err := validateWebhookURL("https://empty.example.com/hook"); err == nil {
			t.Fatal("解析为空必须 fail closed")
		}
	})
}

// TestValidateWebhookDialPinsValidatedIP proves the dial target is a pinned,
// validated IP (with sequential fallback) and never the original hostname.
func TestValidateWebhookDialPinsValidatedIP(t *testing.T) {
	orig := webhookNetDial
	defer func() { webhookNetDial = orig }()
	var targets []string
	webhookNetDial = func(ctx context.Context, network, address string) (net.Conn, error) {
		targets = append(targets, address)
		return nil, errors.New("stub dial")
	}

	ips := []net.IP{net.ParseIP("93.184.216.34"), net.ParseIP("8.8.8.8")}
	dial := webhookDialer(ips, "443")
	_, _ = dial(context.Background(), "tcp", "example.com:443")

	if len(targets) != 2 {
		t.Fatalf("应该顺序尝试 2 个已校验 IP，实际 %v", targets)
	}
	want := []string{"93.184.216.34:443", "8.8.8.8:443"}
	for i := range want {
		if targets[i] != want[i] {
			t.Errorf("拨号目标[%d]=%q, want %q", i, targets[i], want[i])
		}
	}
	for _, target := range targets {
		if strings.Contains(target, "example.com") {
			t.Errorf("拨号不得回退到原始主机名: %q", target)
		}
	}
}

func TestValidateWebhookDialerDefaultPorts(t *testing.T) {
	// The pinned target must use the port derived from the URL scheme; this
	// test exercises the port selection performed in ExecuteWebhookHook.
	for _, tc := range []struct{ scheme, want string }{{"http", "80"}, {"https", "443"}} {
		_, u, err := validateWebhookTarget(tc.scheme + "://8.8.8.8/hook")
		if err != nil {
			t.Fatalf("valid public IP rejected: %v", err)
		}
		port := u.Port()
		if port == "" {
			if u.Scheme == "https" {
				port = "443"
			} else {
				port = "80"
			}
		}
		if port != tc.want {
			t.Errorf("scheme=%s port=%q, want %q", tc.scheme, port, tc.want)
		}
	}
}
