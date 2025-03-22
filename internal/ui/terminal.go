package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// PrintUsage 打印使用说明
func PrintUsage() {
	titleColor := color.New(color.FgGreen, color.Bold)
	titleColor.Println("AICmdGen")

	descColor := color.New(color.FgYellow, color.Bold)
	descColor.Println("AI命令生成工具")

	fmt.Println()

	usageColor := color.New(color.FgWhite, color.Bold)
	usageColor.Println("使用方法:")
	fmt.Println("  ai [选项] \"<命令描述>\"")

	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -debug    启用调试模式")
	fmt.Println("  -help     显示帮助信息")

	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  ai \"查找当前目录下所有的go文件\"")
	fmt.Println("  ai -debug \"将demo.txt重命名为test.txt\"")
}

// PrintCommandPanel 打印命令面板
func PrintCommandPanel(request, command string) {
	width := 80

	// 打印顶部边框
	fmt.Println(strings.Repeat("─", width))

	requestColor := color.New(color.FgCyan, color.Bold)
	requestColor.Print("请求: ")
	fmt.Println(request)

	// 打印分隔线
	fmt.Println(strings.Repeat("─", width))

	commandColor := color.New(color.FgGreen, color.Bold)
	commandColor.Print("命令: ")
	fmt.Println(command)

	// 打印底部边框
	fmt.Println(strings.Repeat("─", width))
}

// PromptChoice 提示用户选择操作
func PromptChoice() int {
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

// ExecuteCommand 执行命令
func ExecuteCommand(command string, debugMode bool) {
	fmt.Println("\n" + strings.Repeat("─", 80))
	fmt.Println("命令输出:")
	fmt.Println(strings.Repeat("─", 80))

	if debugMode {
		log.Printf("[DEBUG] 执行命令: %s", command)
	}

	var cmd *exec.Cmd
	if os.PathSeparator == '\\' {
		cmd = exec.Command("cmd.exe", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n执行失败: %v\n", err)
	}

	fmt.Println(strings.Repeat("─", 80))
}
