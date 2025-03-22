//go:build !sysdir
// +build !sysdir

package config

// getConfigPath 返回当前目录中的配置文件路径
func getConfigPath() string {
	return "config.toml"
}
