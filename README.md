# AI CMD 命令行工具

基于AI技术的命令行生成工具，能够根据自然语言描述自动生成并执行对应的系统命令

## 功能特性

- 🚀 智能命令生成：通过自然语言描述自动生成可执行命令
- ⚙️ 多平台支持：自动适配Windows/Linux系统命令格式
- 🔧 配置管理：支持配置文件自动生成与持久化存储
- 🐞 调试模式：提供详细的运行日志和API交互信息

## 安装说明

### 前置要求
- Go 1.20+ 开发环境
- 有效的OpenAI API密钥

### 安装步骤
```bash
# 克隆仓库
git clone https://github.com/28074/ai_go.git
cd ai_go

# 安装依赖
make install

# 编译项目
make build
```

## 配置指南

1. 首次运行会自动生成配置文件模板`config.toml`
2. 修改配置文件内容：
```toml
[API]
key = "your-api-key-here"  # 替换为实际API密钥
model = "gpt-3.5-turbo"    # 可选模型类型
url = "https://api.openai.com/v1"  # API端点

[App]
debug_mode = false  # 调试模式开关
```

## 使用示例
```bash
# 基本使用
./ai "查找当前目录下所有的go文件"

# 调试模式（查看详细日志）
./ai -debug "将demo.txt重命名为test.txt"
```

## 项目结构
```
ai_go/
├── config/              # 配置管理模块
│   └── config.go        # 配置文件加载/保存实现
├── generator/           # 命令生成模块 
│   └── command_generator.go  # OpenAI API交互实现
├── main.go              # 主程序入口
├── go.mod               # 依赖管理
└── README.md            # 本文档
```

## 贡献指南
欢迎通过Issue提交问题或PR贡献代码，请遵循以下规范：
1. 新功能开发请创建feature分支
2. Bug修复请创建hotfix分支
3. 提交前请执行格式检查：
```bash
go fmt ./...
go vet ./...
```

## 许可证
[MIT License](LICENSE)