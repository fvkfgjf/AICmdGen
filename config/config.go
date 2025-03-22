package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

// 获取配置文件路径
func getConfigPath() string {
	// 在Linux环境下使用/etc/AICmdGen目录
	if runtime.GOOS == "linux" {
		return "/etc/AICmdGen/config.toml"
	}
	// 其他环境使用当前目录
	return "config.toml"
}

// Load 从配置文件加载配置
func Load() (*Config, error) {
	configPath := getConfigPath()

	// 检查配置文件是否存在，不存在则返回错误
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在")
	}

	// 使用viper加载配置
	v := viper.New()
	v.SetConfigFile(filepath.Clean(configPath))
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 设置默认值
	if config.API.URL == "" {
		config.API.URL = "https://api.openai.com/v1"
	}

	return &config, nil
}

// Save 将配置保存到文件
func Save(config *Config) error {
	configPath := getConfigPath()

	// 使用viper保存配置
	v := viper.New()

	// 设置配置值
	v.Set("API.Key", config.API.Key)
	v.Set("API.Model", config.API.Model)
	v.Set("API.URL", config.API.URL)
	v.Set("App.DebugMode", config.App.DebugMode)

	// 设置配置文件路径和类型
	v.SetConfigFile(filepath.Clean(configPath))
	v.SetConfigType("toml")

	// 创建配置文件目录（如果不存在）
	dir := filepath.Dir(configPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %w", err)
		}
	}

	// 使用WriteConfigAs代替WriteConfig确保文件被创建
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
