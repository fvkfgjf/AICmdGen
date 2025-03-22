# AICmdGen Windows��װ�ű�
# �˽ű�������Windowsϵͳ�ϰ�װAICmdGen�����й���

# ֹͣ�ű�ִ��ʱ��ʾ����
$ErrorActionPreference = "Stop"

# ��װĿ¼
$InstallDir = "$env:LOCALAPPDATA\AICmdGen"
$ConfigDir = "$env:APPDATA\AICmdGen"
$ConfigFile = "$ConfigDir\config.toml"
$Executable = "ai.exe"
$FinalPath = "$InstallDir\$Executable"

# ��ɫ����
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

# ������ԱȨ��
function Check-Admin {
    $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    
    if (-not $isAdmin) {
        Print-Error "���Թ���ԱȨ�����д˽ű�"
    }
}

# ���Go����
function Check-Go {
    Print-Info "���Go����..."
    try {
        $goVersion = (go version) -replace "go version go" -replace " windows.*"
        Print-Info "��⵽Go�汾: $goVersion"
        
        # �򵥰汾�Ƚϣ�ȷ��������1.20
        $versionParts = $goVersion.Split(".")
        $major = [int]$versionParts[0]
        $minor = [int]$versionParts[1]
        
        if ($major -lt 1 -or ($major -eq 1 -and $minor -lt 20)) {
            Print-Error "Go�汾���ͣ���Ҫ1.20+�汾"
        }
    }
    catch {
        Print-Error "δ�ҵ�Go���������Ȱ�װGo 1.20+"
    }
}

# ������Ŀ
function Build-Project {
    Print-Info "������Ŀ..."
    
    # ִ��go mod tidy
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Print-Error "go mod tidy ʧ��"
    }
    
    # ������Ŀ
    $env:CGO_ENABLED = 0
    go build -tags sysdir -ldflags="-s -w" -trimpath -o $Executable .\main.go
    if ($LASTEXITCODE -ne 0) {
        Print-Error "����ʧ��"
    }
    
    if (-not (Test-Path $Executable)) {
        Print-Error "�����Ŀ�ִ���ļ�������"
    }
    
    Print-Info "����ɹ�"
}

# ��������Ŀ¼���ļ�
function Setup-Config {
    Print-Info "���������ļ�..."
    
    # ��������Ŀ¼
    if (-not (Test-Path $ConfigDir)) {
        New-Item -ItemType Directory -Path $ConfigDir -Force | Out-Null
    }
    
    # ��������ļ������ڣ�����Ĭ������
    if (-not (Test-Path $ConfigFile)) {
        $configContent = @"
[api]
key = 'your-api-key-here'
model = 'gpt-3.5-turbo'
url = 'https://api.openai.com/v1'

[app]
debugmode = false
"@
        $configContent | Out-File -FilePath $ConfigFile -Encoding utf8
        Print-Info "�Ѵ���Ĭ�������ļ�: $ConfigFile"
        Print-Warn "��༭�����ļ�����������API��Կ"
    }
    else {
        Print-Info "�����ļ��Ѵ���: $ConfigFile"
    }
}

# ��װ��ִ���ļ�
function Install-Executable {
    Print-Info "��װ��ִ���ļ��� $InstallDir..."
    
    # ������װĿ¼
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # ���ƿ�ִ���ļ�����װĿ¼
    Copy-Item -Path $Executable -Destination $FinalPath -Force
    
    Print-Info "��װ���"
}

# ��ӵ�PATH��������
function Add-ToPath {
    Print-Info "��������ӵ�PATH��������..."
    
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        Print-Info "����ӵ��û�PATH��������"
    }
    else {
        Print-Info "����·������PATH����������"
    }
}

# ������ʱ�ļ�
function Cleanup {
    Print-Info "������ʱ�ļ�..."
    Remove-Item -Path $Executable -Force
}

# ������
function Main {
    Write-Output "======================================"
    Write-Output "      AICmdGen Windows ��װ����"
    Write-Output "======================================"
    
    Check-Admin
    Check-Go
    Build-Project
    Setup-Config
    Install-Executable
    Add-ToPath
    Cleanup
    
    Write-Output "======================================"
    Write-Output "      ��װ�ɹ�!"
    Write-Output "======================================"
    Write-Output "ʹ�÷���: ai ""������������"""
    Write-Output "�����ļ�: $ConfigFile"
    Write-Output "��ȷ���༭�����ļ�����������API��Կ"
    Write-Output ""
    Write-Output "ע��: ��������Ҫ��������������ʾ����PowerShell��ʹPATH����������Ч"
}

# ִ��������
Main