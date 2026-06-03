package collector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/netconfighub/netconfighub/internal/devicemodel"
	"github.com/netconfighub/netconfighub/internal/model"
)

type Collector interface {
	Collect(device model.Device) (string, error)
}

type VendorAdapter interface {
	Vendor() model.Vendor
	ShowRunningConfig() string
	ShowStartupConfig() string
	SaveConfig() string
	EnableCommand() string
	ParseOutput(raw string) (string, error)
	PromptPattern() string
}

var adapterRegistry = map[model.Vendor]VendorAdapter{}

func RegisterAdapter(adapter VendorAdapter) {
	adapterRegistry[adapter.Vendor()] = adapter
}

func GetAdapter(vendor model.Vendor) (VendorAdapter, error) {
	adapter, ok := adapterRegistry[vendor]
	if !ok {
		return nil, fmt.Errorf("不支持的厂商类型: %s", vendor)
	}
	return adapter, nil
}

type DynamicAdapter struct {
	config *devicemodel.DeviceModelConfig
}

func NewDynamicAdapter(cfg *devicemodel.DeviceModelConfig) *DynamicAdapter {
	return &DynamicAdapter{config: cfg}
}

func (a *DynamicAdapter) Vendor() model.Vendor { return model.Vendor(a.config.Vendor) }

func (a *DynamicAdapter) ShowRunningConfig() string {
	if len(a.config.Commands) > 0 {
		return a.config.Commands[0]
	}
	return ""
}

func (a *DynamicAdapter) ShowStartupConfig() string { return "" }

func (a *DynamicAdapter) SaveConfig() string { return "" }

func (a *DynamicAdapter) EnableCommand() string { return a.config.EnableCmd }

func (a *DynamicAdapter) PromptPattern() string { return a.config.PromptRegex }

func (a *DynamicAdapter) ParseOutput(raw string) (string, error) {
	return parseOutputGeneric(raw, a.ShowRunningConfig(), a.config.PromptRegex), nil
}

func getDisablePagingCommands(vendor model.Vendor) []string {
	cfg, err := devicemodel.GetByVendor(string(vendor))
	if err != nil {
		return nil
	}
	return cfg.DisablePaging
}

func parseOutputGeneric(raw, command, promptPattern string) string {
	re := regexp.MustCompile(promptPattern)
	lines := strings.Split(raw, "\n")
	var result []string
	for i, line := range lines {
		trimmed := strings.TrimRight(line, "\r ")
		if strings.Contains(trimmed, "--More--") || strings.Contains(trimmed, "---- More ----") {
			continue
		}
		if strings.Contains(trimmed, "\x1b[") {
			trimmed = stripANSI(trimmed)
		}
		if strings.TrimSpace(trimmed) == strings.TrimSpace(command) {
			continue
		}
		if i == len(lines)-1 && re.MatchString(strings.TrimSpace(trimmed)) {
			continue
		}
		if i == 0 && re.MatchString(strings.TrimSpace(trimmed)) {
			continue
		}
		result = append(result, trimmed)
	}
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}
	for len(result) > 0 && strings.TrimSpace(result[0]) == "" {
		result = result[1:]
	}
	return strings.Join(result, "\n")
}

func stripANSI(s string) string {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansi.ReplaceAllString(s, "")
}

func init() {
	models := devicemodel.ListModels()
	for _, cfg := range models {
		RegisterAdapter(NewDynamicAdapter(cfg))
	}
}
