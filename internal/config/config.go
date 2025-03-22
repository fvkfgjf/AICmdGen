package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 结构体保存应用程序配置
type Config struct {
	API APIConfig
	App AppConfig
}

// APIConfig 保存API相关配置
type APIConfig struct {
	URL   string
	Key   string
	Model string
}

// AppConfig 保存应用程序相关配置
type AppConfig struct {
	DebugMode bool
}

// GetDefaultConfig 返回默认配置
func GetDefaultConfig() *Config {
	return &Config{
		API: APIConfig{
			URL:   "https://api.openai.com/v1",
			Key:   "your-api-key-here",
			Model: "gpt-3.5-turbo",
		},
		App: AppConfig{
			DebugMode: false,
		},
	}
}

// Load 从配置文件加载配置
func Load() (*Config, error) {
	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在")
	}

	// 加载配置
	v := viper.New()
	v.SetConfigFile(filepath.Clean(configPath))
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	if config.API.URL == "" {
		config.API.URL = "https://api.openai.com/v1"
	}

	return &config, nil
}

// Save 将配置保存到文件
func Save(config *Config) error {
	configPath := getConfigPath()

	v := viper.New()
	v.Set("API.Key", config.API.Key)
	v.Set("API.Model", config.API.Model)
	v.Set("API.URL", config.API.URL)
	v.Set("App.DebugMode", config.App.DebugMode)
	v.SetConfigFile(filepath.Clean(configPath))
	v.SetConfigType("toml")

	// 创建配置文件目录（如果不存在）
	dir := filepath.Dir(configPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %w", err)
		}
	}

	// 确保文件被创建
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
