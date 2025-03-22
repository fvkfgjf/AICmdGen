#!/bin/bash

# AICmdGen Linux安装脚本
# 此脚本用于在Linux系统上安装AICmdGen命令行工具

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # 无颜色

# 安装目录
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/AICmdGen"
CONFIG_FILE="$CONFIG_DIR/config.toml"
EXECUTABLE="ai"

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否以root权限运行
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "请以root权限运行此脚本"
        exit 1
    fi
}

# 检查Go环境
check_go() {
    print_info "检查Go环境..."
    if ! command -v go &> /dev/null; then
        print_error "未找到Go环境，请先安装Go 1.20+"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "检测到Go版本: $GO_VERSION"
    
    # 简单版本比较，确保至少是1.20
    MAJOR=$(echo $GO_VERSION | cut -d. -f1)
    MINOR=$(echo $GO_VERSION | cut -d. -f2)
    
    if [ "$MAJOR" -lt 1 ] || ([ "$MAJOR" -eq 1 ] && [ "$MINOR" -lt 20 ]); then
        print_error "Go版本过低，需要1.20+版本"
        exit 1
    fi
}

# 编译项目
build_project() {
    print_info "编译项目..."
    go mod tidy
    CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $EXECUTABLE ./main.go
    
    if [ ! -f "$EXECUTABLE" ]; then
        print_error "编译失败"
        exit 1
    fi
    
    print_info "编译成功"
}

# 创建配置目录和文件
setup_config() {
    print_info "设置配置文件..."
    
    # 创建配置目录
    mkdir -p "$CONFIG_DIR"
    
    # 如果配置文件不存在，创建默认配置
    if [ ! -f "$CONFIG_FILE" ]; then
        cat > "$CONFIG_FILE" << EOF
[API]
URL = "https://api.openai.com/v1"
Key = "your-api-key-here"
Model = "gpt-3.5-turbo"

[App]
DebugMode = false
EOF
        print_info "已创建默认配置文件: $CONFIG_FILE"
        print_warn "请编辑配置文件并设置您的API密钥"
    else
        print_info "配置文件已存在: $CONFIG_FILE"
    fi
    
    # 设置配置文件权限
    chmod 644 "$CONFIG_FILE"
}

# 安装可执行文件
install_executable() {
    print_info "安装可执行文件到 $INSTALL_DIR..."
    
    # 复制可执行文件到安装目录
    cp "$EXECUTABLE" "$INSTALL_DIR/"
    chmod 755 "$INSTALL_DIR/$EXECUTABLE"
    
    print_info "安装完成"
}

# 清理临时文件
cleanup() {
    print_info "清理临时文件..."
    rm -f "$EXECUTABLE"
}

# 主函数
main() {
    echo "======================================"
    echo "      AICmdGen Linux 安装程序"
    echo "======================================"
    
    check_root
    check_go
    build_project
    setup_config
    install_executable
    cleanup
    
    echo "======================================"
    echo "      安装成功!"
    echo "======================================"
    echo "使用方法: ai \"您的命令描述\""
    echo "配置文件: $CONFIG_FILE"
    echo "请确保编辑配置文件并设置您的API密钥"
}

# 执行主函数
main