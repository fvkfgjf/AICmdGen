package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fvkfgjf/AICmdGen/internal/config"
	"github.com/fvkfgjf/AICmdGen/internal/generator"
	"github.com/fvkfgjf/AICmdGen/internal/ui"
	"github.com/spf13/cobra"
)

var (
	debugFlag bool
	cfg       *config.Config
)

// rootCmd 表示没有调用子命令时的基础命令
// 修改 rootCmd 的参数校验逻辑
var rootCmd = &cobra.Command{
	Use:   "ai <命令描述>",
	Short: "AI命令生成工具",
	Long:  `AICmdGen 是一个基于AI的命令行工具生成器，可以将自然语言转换为命令行命令。`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			// 自动触发帮助信息显示
			cmd.Help()
			os.Exit(0)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		var err error
		cfg, err = loadConfig()
		if err != nil {
			handleConfigError(err)
			return
		}

		// 如果命令行指定了调试模式，覆盖配置
		if debugFlag {
			cfg.App.DebugMode = true
		}

		logConfigInfo()

		// 创建命令生成器
		cmdGen := generator.New(cfg)

		// 处理用户请求
		request := strings.Join(args, " ")
		processCommandRequest(cmdGen, request)
	},
}

func init() {
	// 初始化日志
	log.SetPrefix("[AICmdGen] ")
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	// 添加命令行标志
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "启用调试模式")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 加载配置文件
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		// 如果配置文件不存在，创建默认配置
		if strings.Contains(err.Error(), "配置文件不存在") {
			defaultCfg := config.GetDefaultConfig()
			if saveErr := config.Save(defaultCfg); saveErr != nil {
				return nil, fmt.Errorf("创建默认配置失败: %w", saveErr)
			}
			fmt.Println("已创建默认配置文件，请编辑配置并设置API密钥")
			return defaultCfg, nil
		}
		return nil, err
	}
	return cfg, nil
}

// 处理配置错误
func handleConfigError(err error) {
	fmt.Printf("配置加载失败: %v\n", err)
	fmt.Println("请确保配置文件存在且格式正确")
}

func logConfigInfo() {
	if cfg.App.DebugMode {
		log.Printf("[DEBUG] 已加载配置: URL=%s, Model=%s", cfg.API.URL, cfg.API.Model)
		log.Printf("[DEBUG] 调试模式已启用")
	}
}

// 处理命令请求
func processCommandRequest(cmdGen *generator.Generator, request string) {
	cmd, err := cmdGen.GenerateCommand(request)
	if err != nil {
		log.Fatalf("生成命令失败: %v", err)
	}

	logDebug("生成的命令: %s", cmd)

	ui.PrintCommandPanel(request, cmd)

	handleUserChoice(cmdGen, request, cmd)
}

// 处理用户选择
func handleUserChoice(cmdGen *generator.Generator, request, cmd string) {
    for {
        choice := ui.PromptChoice()
        switch choice {
        case 1:
            if err := ui.ExecuteCommand(cmd, cfg.App.DebugMode); err != nil {
                fmt.Printf("\n执行失败: %v\n", err)
                retryChoice := ui.PromptRetryChoice()
                if retryChoice == 1 {
                    // 让AI重新生成命令
                    newRequest := fmt.Sprintf("%s (执行失败: %v, 请提供替代方案)", request, err)
                    if newCmd, err := cmdGen.GenerateCommand(newRequest); err == nil {
                        cmd = newCmd
                        ui.PrintCommandPanel(request, cmd)
                        continue
                    }
                }
            }
            return
        case 2:
            newRequest := request + " (请提供另一种实现方式)"
            cmd, err := cmdGen.GenerateCommand(newRequest)
            if err != nil {
                log.Fatalf("生成命令失败: %v", err)
            }
            ui.PrintCommandPanel(request, cmd)
        case 3:
            return
        }
    }
}

// 记录调试信息
func logDebug(format string, v ...interface{}) {
	if cfg != nil && cfg.App.DebugMode {
		log.Printf("[DEBUG] "+format, v...)
	}
}
