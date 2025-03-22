# AICmdGen Windows安装脚本
# 此脚本用于在Windows系统上安装AICmdGen命令行工具

# 停止脚本执行时显示错误
$ErrorActionPreference = "Stop"

# 安装目录
$InstallDir = "$env:LOCALAPPDATA\AICmdGen"
$ConfigDir = "$env:APPDATA\AICmdGen"
$ConfigFile = "$ConfigDir\config.toml"
$Executable = "ai.exe"
$FinalPath = "$InstallDir\$Executable"

# 颜色定义
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    else {
        $input | Write-Output
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Print-Info($message) {
    Write-Output "[INFO] $message"
}

function Print-Warn($message) {
    Write-ColorOutput Yellow "[WARN] $message"
}

function Print-Error($message) {
    Write-ColorOutput Red "[ERROR] $message"
    exit 1
}

# 检查管理员权限
function Check-Admin {
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    
    if (-not $isAdmin) {
        Print-Error "请以管理员权限运行此脚本"
    }
}

# 检查Go环境
function Check-Go {
    Print-Info "检查Go环境..."
    try {
        $goVersion = (go version) -replace "go version go" -replace " windows.*"
        Print-Info "检测到Go版本: $goVersion"
        
        # 简单版本比较，确保至少是1.20
        $versionParts = $goVersion.Split(".")
        $major = [int]$versionParts[0]
        $minor = [int]$versionParts[1]
        
        if ($major -lt 1 -or ($major -eq 1 -and $minor -lt 20)) {
            Print-Error "Go版本过低，需要1.20+版本"
        }
    }
    catch {
        Print-Error "未找到Go环境，请先安装Go 1.20+"
    }
}

# 编译项目
function Build-Project {
    Print-Info "编译项目..."
    
    # 执行go mod tidy
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Print-Error "go mod tidy 失败"
    }
    
    # 编译项目
    $env:CGO_ENABLED = 0
    go build -ldflags="-s -w" -trimpath -o $Executable .\main.go
    if ($LASTEXITCODE -ne 0) {
        Print-Error "编译失败"
    }
    
    if (-not (Test-Path $Executable)) {
        Print-Error "编译后的可执行文件不存在"
    }
    
    Print-Info "编译成功"
}

# 创建配置目录和文件
function Setup-Config {
    Print-Info "设置配置文件..."
    
    # 创建配置目录
    if (-not (Test-Path $ConfigDir)) {
        New-Item -ItemType Directory -Path $ConfigDir -Force | Out-Null
    }
    
    # 如果配置文件不存在，创建默认配置
    if (-not (Test-Path $ConfigFile)) {
        $configContent = @"
[API]
URL = "https://api.openai.com/v1"
Key = "your-api-key-here"
Model = "gpt-3.5-turbo"

[App]
DebugMode = false
"@
        $configContent | Out-File -FilePath $ConfigFile -Encoding utf8
        Print-Info "已创建默认配置文件: $ConfigFile"
        Print-Warn "请编辑配置文件并设置您的API密钥"
    }
    else {
        Print-Info "配置文件已存在: $ConfigFile"
    }
}

# 安装可执行文件
function Install-Executable {
    Print-Info "安装可执行文件到 $InstallDir..."
    
    # 创建安装目录
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # 复制可执行文件到安装目录
    Copy-Item -Path $Executable -Destination $FinalPath -Force
    
    Print-Info "安装完成"
}

# 添加到PATH环境变量
function Add-ToPath {
    Print-Info "将程序添加到PATH环境变量..."
    
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        Print-Info "已添加到用户PATH环境变量"
    }
    else {
        Print-Info "程序路径已在PATH环境变量中"
    }
}

# 清理临时文件
function Cleanup {
    Print-Info "清理临时文件..."
    Remove-Item -Path $Executable -Force
}

# 主函数
function Main {
    Write-Output "======================================"
    Write-Output "      AICmdGen Windows 安装程序"
    Write-Output "======================================"
    
    Check-Admin
    Check-Go
    Build-Project
    Setup-Config
    Install-Executable
    Add-ToPath
    Cleanup
    
    Write-Output "======================================"
    Write-Output "      安装成功!"
    Write-Output "======================================"
    Write-Output "使用方法: ai ""您的命令描述"""
    Write-Output "配置文件: $ConfigFile"
    Write-Output "请确保编辑配置文件并设置您的API密钥"
    Write-Output ""
    Write-Output "注意: 您可能需要重新启动命令提示符或PowerShell以使PATH环境变量生效"
}

# 执行主函数
Main