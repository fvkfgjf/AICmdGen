package generator

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fvkfgjf/AICmdGen/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Generator 负责生成命令
type Generator struct {
	client   openai.Client
	messages []openai.ChatCompletionMessageParamUnion
	config   *config.Config
}

// New 创建新的命令生成器
func New(cfg *config.Config) *Generator {
	return &Generator{
		client: openai.NewClient(
			option.WithBaseURL(cfg.API.URL),
			option.WithAPIKey(cfg.API.Key),
		),
		messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(genSystemPrompt()),
		},
		config: cfg,
	}
}

// GenerateCommand 根据用户请求生成命令
func (g *Generator) GenerateCommand(request string) (string, error) {
	// 添加用户消息
	g.messages = append(g.messages, openai.UserMessage(request))

	// 创建API请求
	params := openai.ChatCompletionNewParams{
		Model:    g.config.API.Model,
		Messages: g.messages,
	}

	g.logDebug(func() {
		log.Printf("[DEBUG] 发送请求到API: 模型=%s, 消息数=%d", g.config.API.Model, len(g.messages))
	})

	// 发送API请求并获取响应
	command, err := g.streamAPIRequest(params)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}

	// 保留最后两条消息（系统提示和用户请求）
	if len(g.messages) > 2 {
		g.messages = []openai.ChatCompletionMessageParamUnion{
			g.messages[0],
			g.messages[len(g.messages)-1],
		}
	}

	g.messages = append(g.messages, openai.AssistantMessage(command))

	return command, nil
}

// streamAPIRequest 发送流式API请求并返回完整响应
func (g *Generator) streamAPIRequest(params openai.ChatCompletionNewParams) (string, error) {
	ctx := context.Background()
	stream := g.client.Chat.Completions.NewStreaming(ctx, params)
	defer stream.Close()

	var builder strings.Builder

	g.logDebug(func() {
		log.Printf("[DEBUG] 开始接收API响应...")
	})

	for stream.Next() {
		chunk := stream.Current()

		g.logDebug(func() {
			log.Printf("[DEBUG] 收到响应块: %+v", chunk)
		})

		// 提取当前块的内容
		if chunk.Choices != nil && len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta

			g.logDebug(func() {
				log.Printf("[DEBUG] Delta内容: %+v", delta)
			})

			content := delta.Content
			if content != "" {
				builder.WriteString(content)

				// 在调试模式下实时输出内容
				g.logDebug(func() {
					fmt.Print(content)
				})
			}
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("流式请求错误: %w", err)
	}

	// 获取完整的生成内容
	result := builder.String()

	result = cleanCommandResult(result)

	g.logDebug(func() {
		fmt.Println()
		log.Printf("[DEBUG] 完整响应内容: %s", result)
	})

	return result, nil
}

// cleanCommandResult 清理命令结果，移除可能的代码块标记和多余空白
func cleanCommandResult(result string) string {
	// 移除可能的代码块标记
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimPrefix(result, "bash")
	result = strings.TrimPrefix(result, "cmd")
	result = strings.TrimPrefix(result, "powershell")
	result = strings.TrimSuffix(result, "```")

	// 移除开头和结尾的空白
	result = strings.TrimSpace(result)

	return result
}

// logDebug 在调试模式下执行日志函数
func (g *Generator) logDebug(logFunc func()) {
	if g.config != nil && g.config.App.DebugMode {
		logFunc()
	}
}

// genSystemPrompt 生成系统提示
func genSystemPrompt() string {
	osHint := runtime.GOOS

	// 如果是Linux系统，尝试获取发行版信息
	if osHint == "linux" {
		if distro, err := getLinuxDistro(); err == nil && distro != "" {
			osHint = fmt.Sprintf("Linux %s", distro)
		}
	}

	return fmt.Sprintf(`作为%s命令行专家，严格遵循：
1. 仅返回可直接在终端执行的纯命令文本
2. 禁止包含任何解释性文字、代码块标记或注释
3. 命令必须适用于%s系统
4. 使用最简洁有效的命令实现用户需求
5. 确保命令安全，避免危险操作
6. 如无法确定命令，返回最可能正确的版本`, osHint, osHint)
}

// getLinuxDistro 获取Linux发行版信息
func getLinuxDistro() (string, error) {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}

	// 解析文件内容
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "NAME=") {
			// 提取NAME的值
			value := strings.TrimPrefix(line, "NAME=")
			value = strings.Trim(value, "\"")
			return value, nil
		}
	}

	return "", fmt.Errorf("未找到发行版信息")
}
