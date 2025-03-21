package generator

import (
	"context"
	"fmt"
	"log"
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

	// 调试模式下打印请求信息
	if g.config.App.DebugMode {
		log.Printf("[DEBUG] 发送请求到: %s", g.config.API.URL)
		log.Printf("[DEBUG] 使用模型: %s", g.config.API.Model)
		log.Printf("[DEBUG] 请求内容: %s", request)
	}

	ctx := context.Background()
	stream := g.client.Chat.Completions.NewStreaming(ctx, params)
	defer stream.Close()

	var builder strings.Builder
	acc := openai.ChatCompletionAccumulator{}

	// 调试模式下记录响应片段
	var debugChunks []string

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		// 调试模式下记录每个响应片段
		if g.config.App.DebugMode {
			debugChunks = append(debugChunks, fmt.Sprintf("%+v", chunk))
		}

		if content, ok := acc.JustFinishedContent(); ok {
			builder.WriteString(content)
			// 调试模式下打印每个完成的内容片段
			if g.config.App.DebugMode {
				log.Printf("[DEBUG] 收到内容片段: %s", content)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}

	// 从响应对象中获取完整内容
	result := ""
	if len(acc.Choices) > 0 && acc.Choices[0].Message.Content != "" {
		result = acc.Choices[0].Message.Content
	} else {
		// 如果响应对象中没有内容，则使用累积的内容
		result = strings.TrimSpace(builder.String())
	}

	// 检查结果是否为空
	if result == "" {
		return "", fmt.Errorf("API返回了空命令")
	}

	// 调试模式下打印完整响应信息
	if g.config.App.DebugMode {
		log.Printf("[DEBUG] 响应完成，总共收到 %d 个片段", len(debugChunks))
		log.Printf("[DEBUG] 完整响应内容: %s", result)

		// 打印响应对象信息
		log.Printf("[DEBUG] 响应对象: %+v", acc)
	}

	g.saveResponse(acc)
	return result, nil
}

// 生成系统提示
func genSystemPrompt() string {
	osHint := "Linux Bash"
	if runtime.GOOS == "windows" {
		osHint = "Windows CMD"
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

// 保存响应到消息历史
func (g *CommandGenerator) saveResponse(acc openai.ChatCompletionAccumulator) {
	if len(acc.Choices) > 0 && acc.Choices[0].Message.Content != "" {
		g.messages = append(g.messages, acc.Choices[0].Message.ToParam())
	}
}
