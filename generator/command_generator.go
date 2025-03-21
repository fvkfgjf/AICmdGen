package generator

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/28074/ai_go/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// CommandGenerator 负责生成命令
type CommandGenerator struct {
	client   openai.Client // 使用值类型
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
	g.messages = append(g.messages, openai.UserMessage(request))

	params := openai.ChatCompletionNewParams{
		Model:       openai.ChatModel(g.config.API.Model), // 模型类型转换
		Messages:    g.messages,
		Temperature: openai.Float(0.2),
	}

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
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}

	result := strings.TrimSpace(builder.String())
	g.saveResponse(acc)
	return result, nil
}

// 生成系统提示
func genSystemPrompt() string {
	osHint := "Linux Bash"
	if runtime.GOOS == "windows" {
		osHint = "Windows CMD"
	}
	return fmt.Sprintf(`作为%s专家，严格遵循：
1. 只返回可直接执行的命令
2. 路径空格自动加引号
3. 多步骤用&&连接
4. 使用相对路径`, osHint)
}

// 保存响应到消息历史
func (g *CommandGenerator) saveResponse(acc openai.ChatCompletionAccumulator) {
	if len(acc.Choices) > 0 && acc.Choices[0].Message.Content != "" {
		g.messages = append(g.messages, acc.Choices[0].Message.ToParam())
	}
}
