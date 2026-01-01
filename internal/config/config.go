package config

import (
	"os"
	"path/filepath"

	"ccm/internal/provider"

	"gopkg.in/yaml.v3"
)

// Config 用户配置
type Config struct {
	Providers map[string]provider.Provider `yaml:"providers"`
}

// 配置文件路径
var configDir string
var configFile string

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, "claude-model", "configs")
	configFile = filepath.Join(configDir, "providers.yaml")
}

// GetConfigDir 获取配置目录
func GetConfigDir() string {
	return configDir
}

// GetConfigFile 获取配置文件路径
func GetConfigFile() string {
	return configFile
}

// Load 加载用户配置
func Load() (*Config, error) {
	cfg := &Config{
		Providers: make(map[string]provider.Provider),
	}

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// 如果文件不存在，返回空配置
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save 保存用户配置
func Save(cfg *Config) error {
	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// GetProvider 获取指定供应商配置（合并预置和用户配置）
func GetProvider(name string) (*provider.Provider, bool) {
	cfg, err := Load()
	if err != nil {
		return nil, false
	}

	// 优先查找用户配置
	if p, ok := cfg.Providers[name]; ok {
		return &p, true
	}

	// 查找预置配置（但需要有 API Key 才算配置完成）
	return nil, false
}

// AddProvider 添加或更新供应商
func AddProvider(p provider.Provider) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	cfg.Providers[p.Name] = p
	return Save(cfg)
}

// IsConfigured 检查供应商是否已配置（有 API Key）
func IsConfigured(name string) bool {
	cfg, err := Load()
	if err != nil {
		return false
	}

	p, ok := cfg.Providers[name]
	return ok && p.APIKey != ""
}
