package generator

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/fvkfgjf/AICmdGen/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// CommandGenerator 负责生成命令
type CommandGenerator struct {
	client   openai.Client
	messages []openai.ChatCompletionMessageParamUnion
	config   *config.Config
}

// New 创建新生成器
func New(cfg *config.Config) *CommandGenerator {
	return &CommandGenerator{
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

// GenerateCommand 生成命令
func (g *CommandGenerator) GenerateCommand(request string) (string, error) {
	// 添加用户消息
	g.messages = append(g.messages, openai.UserMessage(request))

	params := openai.ChatCompletionNewParams{
		Model:       openai.ChatModel(g.config.API.Model),
		Messages:    g.messages,
		Temperature: openai.Float(0.2),
	}

	// 调试日志
	g.logDebug(func() {
		log.Printf("[DEBUG] 发送请求到: %s", g.config.API.URL)
		log.Printf("[DEBUG] 使用模型: %s", g.config.API.Model)
		log.Printf("[DEBUG] 请求内容: %s", request)
	})

	// 处理API请求
	result, err := g.streamAPIRequest(params)
	if err != nil {
		return "", err
	}

	return result, nil
}

// 流式处理API请求并返回结果
func (g *CommandGenerator) streamAPIRequest(params openai.ChatCompletionNewParams) (string, error) {
	ctx := context.Background()
	stream := g.client.Chat.Completions.NewStreaming(ctx, params)
	defer stream.Close()

	var builder strings.Builder
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if content, ok := acc.JustFinishedContent(); ok {
			builder.WriteString(content)
			g.logDebug(func() {
				log.Printf("[DEBUG] 收到内容片段: %s", content)
			})
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}

	// 获取结果
	result := strings.TrimSpace(builder.String())
	if result == "" && len(acc.Choices) > 0 {
		result = acc.Choices[0].Message.Content
	}

	// 检查结果
	if result == "" {
		return "", fmt.Errorf("API返回了空命令")
	}

	// 保存响应到消息历史
	if len(acc.Choices) > 0 {
		g.messages = append(g.messages, acc.Choices[0].Message.ToParam())
	}

	return result, nil
}

// 调试日志辅助函数
func (g *CommandGenerator) logDebug(logFunc func()) {
	if g.config != nil && g.config.App.DebugMode {
		logFunc()
	}
}

// 生成系统提示
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
3. 路径含空格时自动添加双引号
4. 多步骤操作使用&&连接
5. 只使用相对路径
6. 保持命令简洁高效

正确格式示例：
ren "old file.txt" "new file.txt"

错误格式示例：
"请使用以下命令："  # 解释性文字
ren old.txt new.txt  # 缺少必要引号
move file1.txt file2.txt && echo "完成"  # 多余的解释性步骤`, osHint)
}

// getLinuxDistro 从/etc/os-release获取Linux发行版信息
func getLinuxDistro() (string, error) {
	// 读取/etc/os-release文件
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}

	// 解析文件内容
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// 查找NAME字段
		if strings.HasPrefix(line, "NAME=") {
			// 提取NAME的值
			name := strings.TrimPrefix(line, "NAME=")
			// 去除引号
			name = strings.Trim(name, "\"")
			return name, nil
		}
	}

	return "", fmt.Errorf("无法从/etc/os-release获取发行版信息")
}
