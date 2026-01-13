#Requires -Version 5.1
<#
.SYNOPSIS
    CCM Quick Install Script for Windows
.DESCRIPTION
    Downloads and installs the latest version of CCM (Claude Code Manager)
.EXAMPLE
    iwr -useb https://raw.githubusercontent.com/taliove/go-claude-model/main/scripts/install.ps1 | iex
#>

$ErrorActionPreference = "Stop"

$Repo = "taliove/go-claude-model"
$BinaryName = "ccm"
$InstallDir = if ($env:CCM_INSTALL_DIR) { $env:CCM_INSTALL_DIR } else { "$env:USERPROFILE\.local\bin" }

# Colors
function Write-Info { param($msg) Write-Host "[INFO] $msg" -ForegroundColor Blue }
function Write-Success { param($msg) Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-Warn { param($msg) Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Write-Err { param($msg) Write-Host "[ERROR] $msg" -ForegroundColor Red; exit 1 }

function Get-Architecture {
    $arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { Write-Err "Unsupported architecture: $arch" }
    }
}

function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        return $response.tag_name
    } catch {
        Write-Err "Failed to get latest version: $_"
    }
}

function Check-Dependencies {
    Write-Host ""
    Write-Info "Checking dependencies..."

    # Check npm
    if (Get-Command npm -ErrorAction SilentlyContinue) {
        Write-Success "npm is installed"
    } else {
        Write-Warn "npm is not installed"
        Write-Host "  Install Node.js: https://nodejs.org/"
    }

    # Check claude
    if (Get-Command claude -ErrorAction SilentlyContinue) {
        Write-Success "claude-code is installed"
    } else {
        Write-Warn "claude-code is not installed"
        Write-Host "  Install: npm install -g @anthropic-ai/claude-code"
    }
}

function Install-CCM {
    Write-Info "Installing CCM..."

    $arch = Get-Architecture
    $version = Get-LatestVersion

    Write-Info "Detected platform: windows_$arch"
    Write-Info "Latest version: $version"

    $versionNum = $version.TrimStart('v')
    $downloadUrl = "https://github.com/$Repo/releases/download/$version/${BinaryName}_${versionNum}_windows_${arch}.zip"

    Write-Info "Downloading from: $downloadUrl"

    # Create temp directory
    $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    $zipPath = Join-Path $tempDir "ccm.zip"

    try {
        # Download
        [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing

        # Extract
        Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force

        # Install
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        $exePath = Get-ChildItem -Path $tempDir -Filter "$BinaryName.exe" -Recurse | Select-Object -First 1
        if ($exePath) {
            Move-Item -Path $exePath.FullName -Destination (Join-Path $InstallDir "$BinaryName.exe") -Force
        } else {
            Write-Err "Binary not found in archive"
        }

        Write-Success "CCM installed to $InstallDir\$BinaryName.exe"

        # Check PATH
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($userPath -notlike "*$InstallDir*") {
            Write-Host ""
            Write-Info "Add to your PATH (run in PowerShell as Administrator):"
            Write-Host ""
            Write-Host "  `$env:Path += `";$InstallDir`""
            Write-Host ""
            Write-Host "Or add permanently:"
            Write-Host "  [Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$InstallDir', 'User')"
        }

        Write-Host ""
        Write-Host "==========================================" -ForegroundColor Cyan
        Write-Host "  Getting Started" -ForegroundColor Cyan
        Write-Host "==========================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "1. Run initial setup:"
        Write-Host "   ccm init" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "2. Configure a provider:"
        Write-Host "   ccm add doubao --key `"your-api-key`"" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "3. Start Claude Code:"
        Write-Host "   ccm run doubao" -ForegroundColor Yellow
        Write-Host ""

        # Check dependencies
        Check-Dependencies

        Write-Host ""
        Write-Success "Installation complete!"

    } finally {
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

Install-CCM
