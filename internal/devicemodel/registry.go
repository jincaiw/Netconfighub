package devicemodel

import (
	"fmt"
	"sync"
)

var (
	registry     = make(map[string]*DeviceModelConfig)
	registryOnce sync.Once
)

func init() {
	registryOnce.Do(func() {
		registerBuiltinModels()
	})
}

func registerBuiltinModels() {
	ciscoIOS := &DeviceModelConfig{
		Model:     "cisco_ios",
		Vendor:    "cisco",
		Protocols: []string{"ssh", "telnet"},
		DisablePaging: []string{
			"terminal length 0",
			"terminal width 0",
		},
		Commands:    []string{"show running-config"},
		Exit:        []string{"exit"},
		PromptRegex: `[>#]\s*$`,
		EnableCmd:   "enable",
		SecretRules: []SecretRule{
			{Pattern: `^enable secret .+$`, Replace: `enable secret <removed>`},
			{Pattern: `^username (\S+) password .+$`, Replace: `username $1 password <removed>`},
		},
	}

	huaweiVRP := &DeviceModelConfig{
		Model:     "huawei_vrp",
		Vendor:    "huawei",
		Protocols: []string{"ssh", "telnet"},
		DisablePaging: []string{
			"screen-length 0 temporary",
		},
		Commands:    []string{"display current-configuration"},
		Exit:        []string{"quit"},
		PromptRegex: `[<\[]\S+[>\]]`,
		EnableCmd:   "",
		SecretRules: []SecretRule{
			{Pattern: `^local-user .+ password .+$`, Replace: `local-user <removed> password <removed>`},
		},
	}

	h3cComware := &DeviceModelConfig{
		Model:     "h3c_comware",
		Vendor:    "h3c",
		Protocols: []string{"ssh", "telnet"},
		DisablePaging: []string{
			"screen-length disable",
		},
		Commands:    []string{"display current-configuration"},
		Exit:        []string{"quit"},
		PromptRegex: `[<\[]\S+[>\]]`,
		EnableCmd:   "super 0",
		SecretRules: []SecretRule{
			{Pattern: `^password .+$`, Replace: `password <removed>`},
		},
	}

	ruijieRGOS := &DeviceModelConfig{
		Model:     "ruijie_rg-os",
		Vendor:    "ruijie",
		Protocols: []string{"ssh", "telnet"},
		DisablePaging: []string{
			"terminal length 0",
		},
		Commands:    []string{"show running-config"},
		Exit:        []string{"exit"},
		PromptRegex: `[>#]\s*$`,
		EnableCmd:   "enable",
		SecretRules: []SecretRule{
			{Pattern: `^enable secret .+$`, Replace: `enable secret <removed>`},
		},
	}

	registry[ciscoIOS.Model] = ciscoIOS
	registry[huaweiVRP.Model] = huaweiVRP
	registry[h3cComware.Model] = h3cComware
	registry[ruijieRGOS.Model] = ruijieRGOS
}

func key(vendor, model string) string {
	return fmt.Sprintf("%s_%s", vendor, model)
}

func GetModel(vendor, model string) (*DeviceModelConfig, error) {
	k := key(vendor, model)
	if cfg, ok := registry[k]; ok {
		return cfg, nil
	}
	return nil, fmt.Errorf("未找到设备模型配置: vendor=%s, model=%s", vendor, model)
}

func GetByVendor(vendor string) (*DeviceModelConfig, error) {
	for _, cfg := range registry {
		if cfg.Vendor == vendor {
			return cfg, nil
		}
	}
	return nil, fmt.Errorf("未找到厂商模型配置: vendor=%s", vendor)
}

func RegisterModel(cfg *DeviceModelConfig) error {
	if cfg == nil {
		return fmt.Errorf("模型配置不能为空")
	}
	if cfg.Model == "" || cfg.Vendor == "" {
		return fmt.Errorf("模型名称和厂商不能为空")
	}
	k := key(cfg.Vendor, cfg.Model)
	registry[k] = cfg
	return nil
}

func ListModels() []*DeviceModelConfig {
	result := make([]*DeviceModelConfig, 0, len(registry))
	for _, cfg := range registry {
		result = append(result, cfg)
	}
	return result
}
