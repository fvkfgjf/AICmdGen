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

### Windows安装步骤
```bash
# 克隆仓库
git clone https://github.com/fvkfgjf/AICmdGen.git
cd AICmdGen

# 使用安装脚本（需要管理员权限）
powershell -ExecutionPolicy ByPass -File .\scripts\install.ps1

# 或手动编译
# 注意：需要先安装make工具
make install && make build
```

### Linux安装步骤
```bash
# 克隆仓库
git clone https://github.com/fvkfgjf/AICmdGen.git
cd AICmdGen

# 使用安装脚本（需要root权限）
sudo ./scripts/install.sh
```

 ## 使用示例
```bash
# 基本使用
ai "查找当前目录下所有的go文件"

# 调试模式（查看详细日志）
ai -debug "将demo.txt重命名为test.txt"
```

## 配置指南

1. 首次运行会自动生成配置文件模板`config.toml`
2. 修改配置文件内容：
```toml
[API]
key = "your-api-key-here"  # 替换为实际API密钥
model = "gpt-3.5-turbo"    # 可选模型
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
AICmdGen/
├── cmd/
│   └── ai/         
│       └── main.go       # 主程序入口
├── internal/
│   ├── config/
│   │   ├── config.go     # 配置结构和方法
│   │   ├── localdir.go   # 本地目录配置实现
│   │   └── sysdir.go     # 系统目录配置实现
│   ├── generator/
│   │   └── generator.go  # 命令生成实现
│   └── ui/
│       └── terminal.go   # 终端UI实现
├── scripts/
│   ├── install.ps1       # Windows安装脚本
│   └── install.sh        # Linux安装脚本
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