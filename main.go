package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/28074/ai_go/config"
	"github.com/28074/ai_go/generator"
	"github.com/fatih/color"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		// 如果配置文件不存在，创建默认配置文件
		cfg = &config.Config{
			API: config.APIConfig{
				Key:   "your-api-key-here", // 请替换为实际的API密钥
				Model: "gpt-3.5-turbo",
			},
			App: config.AppConfig{
				DebugMode: false,
			},
		}

		// 保存默认配置到文件
		if err := config.Save(cfg); err != nil {
			log.Printf("警告: 无法创建配置文件: %v", err)
			log.Println("将使用内存中的默认配置继续运行")
		} else {
			log.Println("已创建默认配置文件 config.toml，请修改其中的API密钥后重新运行")
			return
		}
	} else {
		// 打印配置信息以便调试
		log.Printf("已加载配置: URL=%s, Model=%s", cfg.API.URL, cfg.API.Model)
	}

	// 创建命令生成器
	cmdGen := generator.New(cfg)

	// 获取用户请求
	request := strings.Join(os.Args[1:], " ")

	// 生成命令
	cmd, err := cmdGen.GenerateCommand(request)
	if err != nil {
		log.Fatalf("生成命令失败: %v", err)
	}

	// 输出结果
	printPanel(request, cmd)

	// 提供选项
	for {
		choice := promptChoice()
		switch choice {
		case 1: // 执行此条命令
			executeCommand(cmd)
			return
		case 2: // 换一条命令
			newRequest := request + " (请提供另一种实现方式)"
			cmd, err = cmdGen.GenerateCommand(newRequest)
			if err != nil {
				log.Fatalf("生成命令失败: %v", err)
			}
			printPanel(request, cmd)
		case 3: // 退出
			return
		}
	}
}

// 打印使用说明
func printUsage() {
	titleColor := color.New(color.FgGreen, color.Bold)
	titleColor.Println("AI CMD")

	descColor := color.New(color.FgYellow, color.Bold)
	descColor.Println("AI命令生成工具")

	fmt.Println()

	usageColor := color.New(color.FgWhite, color.Bold)
	usageColor.Print("用法: ")

	cmdColor := color.New(color.FgGreen)
	cmdColor.Print("AI.exe ")

	argColor := color.New(color.FgCyan)
	argColor.Println("<您需要执行的任务描述>")
}

// 打印命令面板
func printPanel(request, command string) {
	width := 80

	// 打印顶部边框
	fmt.Println("┏" + strings.Repeat("━", width-2) + "┓")

	// 打印标题
	titleColor := color.New(color.FgYellow, color.Bold)
	fmt.Print("┃ ")
	titleColor.Print(request)
	fmt.Println(strings.Repeat(" ", width-4-len(request)) + " ┃")

	// 打印分隔线
	fmt.Println("┣" + strings.Repeat("━", width-2) + "┫")

	// 打印命令内容
	cmdColor := color.New(color.FgWhite)
	lines := splitStringByWidth(command, width-4)
	for _, line := range lines {
		fmt.Print("┃ ")
		cmdColor.Print(line)
		fmt.Println(strings.Repeat(" ", width-4-len(line)) + " ┃")
	}

	// 打印底部边框
	fmt.Println("┗" + strings.Repeat("━", width-2) + "┛")
}

// 按宽度分割字符串
func splitStringByWidth(s string, width int) []string {
	var lines []string
	for len(s) > width {
		lines = append(lines, s[:width])
		s = s[width:]
	}
	if len(s) > 0 {
		lines = append(lines, s)
	}
	return lines
}

// 提示用户选择
func promptChoice() int {
	fmt.Println("\n选择操作:")
	fmt.Println("1. 执行此条命令")
	fmt.Println("2. 换一条命令")
	fmt.Println("3. 退出")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入选项 (1-3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1", "2", "3":
			choice := int(input[0] - '0')
			return choice
		default:
			fmt.Println("无效选项，请重新输入")
		}
	}
}

// 执行命令
func executeCommand(command string) {
	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("命令输出:")
	fmt.Println(strings.Repeat("─", 80))

	cmd := exec.Command("cmd.exe", "/c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n执行失败: %v\n", err)
	}

	fmt.Println(strings.Repeat("─", 80))
}
