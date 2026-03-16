<#
文件：scripts/build-release.ps1
作用：一键构建发布产物，负责先构建前端到 backend/frontend，再编译 Go 后端为 Linux/Windows 可执行文件，并组装为可直接部署的发布目录与压缩包。
职责边界：
1. 执行前端生产构建，保持产物输出到 backend/frontend。
2. 按目标参数交叉编译后端为 linux/amd64、windows/amd64 或同时编译两者。
3. 将二进制、config、frontend 复制到各自发布目录，并预创建 uploads 目录。
4. 生成 zip 压缩包，便于上传到服务器或分发到 Windows 环境。
对外接口：
- PowerShell CLI：pwsh -File .\scripts\build-release.ps1
- PowerShell CLI：.\scripts\build-release.ps1 linux
- 可选参数：-SkipFrontend、-OutputDir、-Target
#>

[CmdletBinding()]
param(
    [Parameter(Position = 1)]
    [switch]$SkipFrontend,
    [Parameter(Position = 2)]
    [string]$OutputDir = "release",
    [Parameter(Position = 0)]
    [ValidateSet('linux', 'windows', 'all')]
    [string]$Target = 'all'
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

$script:RootDir = Split-Path -Parent $PSScriptRoot
$script:VueDir = Join-Path $script:RootDir 'vue'
$script:BackendDir = Join-Path $script:RootDir 'backend'
$script:ReleaseDir = Join-Path $script:RootDir $OutputDir
$script:StagingDir = Join-Path $script:ReleaseDir '.build-staging'
$script:FrontendBuildDir = Join-Path $script:VueDir 'dist'
$script:ShouldCleanupFrontendOutput = $false

function Write-Step {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message
    )

    Write-Host ""
    Write-Host ">>> $Message" -ForegroundColor Cyan
}

function Resolve-CommandPath {
    param(
        [Parameter(Mandatory = $true)]
        [string[]]$CommandNames,
        [Parameter(Mandatory = $true)]
        [string]$NotFoundMessage
    )

    foreach ($commandName in $CommandNames) {
        $command = Get-Command $commandName -ErrorAction SilentlyContinue
        if ($null -ne $command) {
            return $command.Source
        }
    }

    throw $NotFoundMessage
}

function Get-NpmCommand {
    if ($IsWindows) {
        return Resolve-CommandPath -CommandNames @('npm.cmd', 'npm') -NotFoundMessage '未找到 npm，请先安装 Node.js 并确保 npm 在 PATH 中。'
    }

    return Resolve-CommandPath -CommandNames @('npm') -NotFoundMessage '未找到 npm，请先安装 Node.js 并确保 npm 在 PATH 中。'
}

function Get-GoCommand {
    return Resolve-CommandPath -CommandNames @('go') -NotFoundMessage '未找到 go，请先安装 Go 并确保 go 在 PATH 中。'
}

function Assert-DirectoryExists {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path,
        [Parameter(Mandatory = $true)]
        [string]$ErrorMessage
    )

    if (-not (Test-Path -LiteralPath $Path -PathType Container)) {
        throw $ErrorMessage
    }
}

function Remove-IfExists {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )

    if (Test-Path -LiteralPath $Path) {
        Remove-Item -LiteralPath $Path -Recurse -Force
    }
}

function Copy-Directory {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Source,
        [Parameter(Mandatory = $true)]
        [string]$Destination
    )

    if (-not (Test-Path -LiteralPath $Source)) {
        throw "目录不存在：$Source"
    }

    New-Item -ItemType Directory -Path $Destination -Force | Out-Null
    Copy-Item -Path (Join-Path $Source '*') -Destination $Destination -Recurse -Force
}

function Invoke-BuildCommand {
    param(
        [Parameter(Mandatory = $true)]
        [string]$FilePath,
        [Parameter(Mandatory = $true)]
        [string[]]$Arguments,
        [Parameter(Mandatory = $true)]
        [string]$WorkingDirectory,
        [hashtable]$EnvironmentVariables = @{}
    )

    $originalValues = @{}
    $locationPushed = $false

    try {
        foreach ($entry in $EnvironmentVariables.GetEnumerator()) {
            $originalValues[$entry.Key] = [Environment]::GetEnvironmentVariable($entry.Key, 'Process')
            [Environment]::SetEnvironmentVariable($entry.Key, $entry.Value, 'Process')
        }

        Push-Location -LiteralPath $WorkingDirectory
        $locationPushed = $true
        & $FilePath @Arguments
        $exitCode = $LASTEXITCODE
        if ($null -eq $exitCode) {
            $exitCode = 0
        }

        if ($exitCode -ne 0) {
            throw "命令执行失败：$FilePath $($Arguments -join ' ')"
        }
    }
    finally {
        if ($locationPushed) {
            Pop-Location
        }

        foreach ($entry in $EnvironmentVariables.GetEnumerator()) {
            [Environment]::SetEnvironmentVariable($entry.Key, $originalValues[$entry.Key], 'Process')
        }
    }
}

function Get-BuildTargets {
    $targets = @()

    if ($Target -eq 'linux' -or $Target -eq 'all') {
        $targets += @{
            Name = 'linux-amd64'
            GoOS = 'linux'
            GoArch = 'amd64'
            BinaryName = 'docplatform'
        }
    }

    if ($Target -eq 'windows' -or $Target -eq 'all') {
        $targets += @{
            Name = 'windows-amd64'
            GoOS = 'windows'
            GoArch = 'amd64'
            BinaryName = 'docplatform.exe'
        }
    }

    return $targets
}

function Build-Frontend {
    if ($SkipFrontend) {
        Write-Step '跳过前端构建'

        Assert-DirectoryExists -Path $script:FrontendBuildDir -ErrorMessage '已跳过前端构建，但 vue/dist 目录不存在。'

        if (-not (Test-Path -LiteralPath (Join-Path $script:FrontendBuildDir 'index.html') -PathType Leaf)) {
            throw '已跳过前端构建，但 vue/dist 缺少 index.html，无法继续发布。'
        }

        return
    }

    $npmCommand = Get-NpmCommand

    Write-Step '构建前端产物到 vue/dist'
    Invoke-BuildCommand -FilePath $npmCommand -Arguments @('run', 'build') -WorkingDirectory $script:VueDir

    Assert-DirectoryExists -Path $script:FrontendBuildDir -ErrorMessage '前端构建后未生成 vue/dist 目录。'

    if (-not (Test-Path -LiteralPath (Join-Path $script:FrontendBuildDir 'index.html') -PathType Leaf)) {
        throw '前端构建后缺少 index.html，发布产物不完整。'
    }

    $script:ShouldCleanupFrontendOutput = $true
}

function Build-GoBinary {
    param(
        [Parameter(Mandatory = $true)]
        [string]$TargetName,
        [Parameter(Mandatory = $true)]
        [string]$GoOS,
        [Parameter(Mandatory = $true)]
        [string]$GoArch,
        [Parameter(Mandatory = $true)]
        [string]$BinaryName
    )

    $goCommand = Get-GoCommand

    $targetDir = Join-Path $script:StagingDir $TargetName
    $binaryPath = Join-Path $targetDir $BinaryName
    $zipPath = Join-Path $script:StagingDir ("docplatform-$TargetName.zip")

    Write-Step "编译后端：$TargetName"
    New-Item -ItemType Directory -Path $targetDir -Force | Out-Null

    Invoke-BuildCommand -FilePath $goCommand -Arguments @('build', '-o', $binaryPath, 'cmd/server/main.go') -WorkingDirectory $script:BackendDir -EnvironmentVariables @{
        GOOS = $GoOS
        GOARCH = $GoArch
        CGO_ENABLED = '0'
    }

    if (-not (Test-Path -LiteralPath $binaryPath -PathType Leaf)) {
        throw "后端编译完成后未找到目标文件：$binaryPath"
    }

    Write-Step "组装发布目录：$TargetName"
    Copy-Directory -Source (Join-Path $script:BackendDir 'config') -Destination (Join-Path $targetDir 'config')
    Copy-Directory -Source $script:FrontendBuildDir -Destination (Join-Path $targetDir 'dist')
    New-Item -ItemType Directory -Path (Join-Path $targetDir 'uploads') -Force | Out-Null

    Write-Step "压缩发布包：$TargetName"
    Compress-Archive -Path (Join-Path $targetDir '*') -DestinationPath $zipPath -CompressionLevel Optimal
}

function Prepare-StagingDirectory {
    Remove-IfExists -Path $script:StagingDir
    New-Item -ItemType Directory -Path $script:StagingDir -Force | Out-Null
}

function Publish-ReleaseArtifacts {
    param(
        [Parameter(Mandatory = $true)]
        [hashtable[]]$BuildTargets
    )

    foreach ($buildTarget in $BuildTargets) {
        $finalTargetDir = Join-Path $script:ReleaseDir $buildTarget.Name
        $finalZipPath = Join-Path $script:ReleaseDir ("docplatform-$($buildTarget.Name).zip")
        $stagingTargetDir = Join-Path $script:StagingDir $buildTarget.Name
        $stagingZipPath = Join-Path $script:StagingDir ("docplatform-$($buildTarget.Name).zip")

        Remove-IfExists -Path $finalTargetDir
        Remove-IfExists -Path $finalZipPath

        Move-Item -LiteralPath $stagingTargetDir -Destination $finalTargetDir
        Move-Item -LiteralPath $stagingZipPath -Destination $finalZipPath
    }
}

try {
    Write-Step '准备发布目录'
    Assert-DirectoryExists -Path $script:VueDir -ErrorMessage "前端目录不存在：$script:VueDir"
    Assert-DirectoryExists -Path $script:BackendDir -ErrorMessage "后端目录不存在：$script:BackendDir"
    Assert-DirectoryExists -Path (Join-Path $script:BackendDir 'config') -ErrorMessage "后端配置目录不存在：$(Join-Path $script:BackendDir 'config')"
    New-Item -ItemType Directory -Path $script:ReleaseDir -Force | Out-Null

    $buildTargets = Get-BuildTargets
    Prepare-StagingDirectory

    Build-Frontend

    foreach ($buildTarget in $buildTargets) {
        Build-GoBinary -TargetName $buildTarget.Name -GoOS $buildTarget.GoOS -GoArch $buildTarget.GoArch -BinaryName $buildTarget.BinaryName
    }

    Write-Step '整理最终发布产物'
    Publish-ReleaseArtifacts -BuildTargets $buildTargets

    Write-Step '构建完成'
    Write-Host "发布目录：$script:ReleaseDir" -ForegroundColor Green

    foreach ($buildTarget in $buildTargets) {
        Write-Host "$($buildTarget.Name) 发布包：$(Join-Path $script:ReleaseDir ("docplatform-$($buildTarget.Name).zip"))" -ForegroundColor Green
    }

    Write-Host ''

    if ($Target -eq 'linux' -or $Target -eq 'all') {
        Write-Host '宝塔部署建议：上传 linux-amd64 目录内容或 docplatform-linux-amd64.zip 解压后的内容，并在服务器上将 config/config.yaml 改为生产配置。' -ForegroundColor Yellow
    }
}
finally {
    if ($script:ShouldCleanupFrontendOutput) {
        Write-Step '清理 vue/dist 前端中间产物'
        Remove-IfExists -Path $script:FrontendBuildDir
    }

    Remove-IfExists -Path $script:StagingDir
}
