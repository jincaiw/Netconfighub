package devicemodel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(path string) (*DeviceModelConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取模型文件失败: %w", err)
	}
	var cfg DeviceModelConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析模型文件失败 [%s]: %w", path, err)
	}
	if cfg.Model == "" || cfg.Vendor == "" {
		return nil, fmt.Errorf("模型文件缺少必要字段 (model, vendor): %s", path)
	}
	return &cfg, nil
}

func LoadFromDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取模型目录失败: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}
		path := filepath.Join(dir, name)
		cfg, err := LoadFromFile(path)
		if err != nil {
			return fmt.Errorf("加载模型文件失败 [%s]: %w", path, err)
		}
		if err := RegisterModel(cfg); err != nil {
			return fmt.Errorf("注册模型失败 [%s]: %w", path, err)
		}
	}
	return nil
}
