//go:build sysdir
// +build sysdir

package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// getConfigPath 返回系统配置目录中的配置文件路径
func getConfigPath() string {
	// 在Linux环境下使用/etc/AICmdGen目录
	if runtime.GOOS == "linux" {
		return "/etc/AICmdGen/config.toml"
	}

	// 在Windows环境下使用AppData目录
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "AICmdGen", "config.toml")
		}
	}

	// 如果无法确定系统配置目录，回退到当前目录
	return "config.toml"
}
