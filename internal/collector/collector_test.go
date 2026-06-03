package collector

import (
	"testing"

	"github.com/netconfighub/netconfighub/internal/model"
)

func TestVendorAdapterCommands(t *testing.T) {
	tests := []struct {
		vendor      model.Vendor
		showRunning string
	}{
		{model.VendorCisco, "show running-config"},
		{model.VendorH3C, "display current-configuration"},
		{model.VendorHuawei, "display current-configuration"},
		{model.VendorRuijie, "show running-config"},
	}

	for _, tt := range tests {
		t.Run(string(tt.vendor), func(t *testing.T) {
			adapter, err := GetAdapter(tt.vendor)
			if err != nil {
				t.Fatalf("获取 %s 适配器失败: %v", tt.vendor, err)
			}
			if adapter.ShowRunningConfig() != tt.showRunning {
				t.Errorf("ShowRunningConfig() = %q, want %q", adapter.ShowRunningConfig(), tt.showRunning)
			}
		})
	}
}

func TestGetAdapterUnsupported(t *testing.T) {
	_, err := GetAdapter("unsupported")
	if err == nil {
		t.Error("不支持的厂商应返回错误")
	}
}

func TestNewCollector(t *testing.T) {
	c := NewCollector()
	if c == nil {
		t.Fatal("NewCollector() returned nil")
	}
	if c.sshClient == nil {
		t.Error("sshClient not initialized")
	}
	if c.telnetClient == nil {
		t.Error("telnetClient not initialized")
	}
}

func TestCiscoAdapterParseOutput(t *testing.T) {
	adapter, err := GetAdapter(model.VendorCisco)
	if err != nil {
		t.Fatalf("获取适配器失败: %v", err)
	}
	raw := "show running-config\nhostname Router\ninterface Gig0/1\n ip address 10.0.0.1\nRouter#"
	output, err := adapter.ParseOutput(raw)
	if err != nil {
		t.Fatalf("ParseOutput failed: %v", err)
	}
	if output == "" {
		t.Error("ParseOutput returned empty string")
	}
}

func TestH3CAdapterParseOutput(t *testing.T) {
	adapter, err := GetAdapter(model.VendorH3C)
	if err != nil {
		t.Fatalf("获取适配器失败: %v", err)
	}
	raw := "display current-configuration\nhostname Switch\n[Switch]"
	output, err := adapter.ParseOutput(raw)
	if err != nil {
		t.Fatalf("ParseOutput failed: %v", err)
	}
	if output == "" {
		t.Error("ParseOutput returned empty string")
	}
}

func TestDisablePagingCommands(t *testing.T) {
	tests := []struct {
		vendor   model.Vendor
		expected []string
	}{
		{model.VendorCisco, []string{"terminal length 0", "terminal width 0"}},
		{model.VendorRuijie, []string{"terminal length 0"}},
		{model.VendorH3C, []string{"screen-length disable"}},
		{model.VendorHuawei, []string{"screen-length 0 temporary"}},
	}

	for _, tt := range tests {
		result := getDisablePagingCommands(tt.vendor)
		if len(result) != len(tt.expected) {
			t.Errorf("getDisablePagingCommands(%s) = %v, want %v", tt.vendor, result, tt.expected)
			continue
		}
		for i, cmd := range result {
			if cmd != tt.expected[i] {
				t.Errorf("getDisablePagingCommands(%s)[%d] = %q, want %q", tt.vendor, i, cmd, tt.expected[i])
			}
		}
	}
}

func TestEnableCommand(t *testing.T) {
	tests := []struct {
		vendor   model.Vendor
		expected string
	}{
		{model.VendorCisco, "enable"},
		{model.VendorRuijie, "enable"},
		{model.VendorH3C, "super 0"},
		{model.VendorHuawei, ""},
	}

	for _, tt := range tests {
		adapter, err := GetAdapter(tt.vendor)
		if err != nil {
			t.Fatalf("获取 %s 适配器失败: %v", tt.vendor, err)
		}
		if adapter.EnableCommand() != tt.expected {
			t.Errorf("EnableCommand(%s) = %q, want %q", tt.vendor, adapter.EnableCommand(), tt.expected)
		}
	}
}
