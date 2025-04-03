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

// PrintCommandPanel 打印命令面板
func PrintCommandPanel(request, command string) {
	// width := 80

	// 打印顶部边框
	// fmt.Println(strings.Repeat("─", width))

	// requestColor := color.New(color.FgCyan, color.Bold)
	// requestColor.Print("请求: ")
	// fmt.Println(request)

	// 打印分隔线
	// fmt.Println(strings.Repeat("─", width))

	commandColor := color.New(color.FgGreen, color.Bold)
	commandColor.Print("命令: ")
	fmt.Println(command)

	// 打印底部边框
	// fmt.Println(strings.Repeat("─", width))
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

// PromptRetryChoice 提示用户选择操作
func PromptRetryChoice() int {
	fmt.Println("\n命令执行失败，请选择:")
	fmt.Println("1. 让AI提供替代方案")
	fmt.Println("2. 返回主菜单")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入选项 (1-2): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1", "2":
			fmt.Print("\033[5A\033[0J")
			return int(input[0] - '0')
		default:
			fmt.Println("无效选项，请重新输入")
		}
	}
}

// ExecuteCommand 执行命令
func ExecuteCommand(command string, debugMode bool) error {
	// 清空上五行内容
	fmt.Print("\033[7A\033[0J") // ANSI escape code 上移5行并清除
	cmdColor := color.New(color.FgMagenta, color.Bold)
	cmdColor.Println("执行命令:", command)

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

	return cmd.Run()
}
