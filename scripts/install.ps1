# =============================================================================
# WATCHMEN INSTALLER
# =============================================================================
# PowerShell script to install Watchmen as Windows Service
#
# Usage:
#   .\install.ps1                           # Interactive mode
#   .\install.ps1 -Silent                   # Silent mode (no prompts)
#   .\install.ps1 -ConfigPath "D:\config.yaml"  # Custom config path
#
# Requires: Administrator privileges
# =============================================================================

#Requires -Version 5.1
#Requires -RunAsAdministrator

[CmdletBinding()]
param(
    [switch]$Silent,
    [string]$ConfigPath = "",
    [switch]$Help
)

# =============================================================================
# Configuration
# =============================================================================
$ServiceName = "Watchmen"
$DisplayName = "Watchmen - SQL Agent Monitor"
$Description = "Monitors SQL Server Agent jobs and sends Windows Toast notifications when jobs fail."
$InstallDir = "$env:ProgramData\Watchmen"
$LogDir = "$InstallDir\logs"
$ExeName = "watchmen.exe"
$ConfigFileName = "config.yaml"

# Colors
$ColorSuccess = "Green"
$ColorError = "Red"
$ColorWarning = "Yellow"
$ColorInfo = "Cyan"

# =============================================================================
# Functions
# =============================================================================

function Write-Banner {
    $version = "1.0.0"
    Write-Host ""
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    Write-Host "   🔧 WATCHMEN INSTALLER v$version" -ForegroundColor $ColorInfo
    Write-Host "   SQL Server Agent Job Monitor" -ForegroundColor $ColorInfo
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    Write-Host ""
}

function Write-Step {
    param([string]$Message, [string]$Status = "...")
    
    switch ($Status) {
        "OK"      { Write-Host "[✓] $Message" -ForegroundColor $ColorSuccess }
        "FAIL"    { Write-Host "[✗] $Message" -ForegroundColor $ColorError }
        "WARN"    { Write-Host "[!] $Message" -ForegroundColor $ColorWarning }
        "INFO"    { Write-Host "[i] $Message" -ForegroundColor $ColorInfo }
        default   { Write-Host "[...] $Message" -ForegroundColor White }
    }
}

function Show-Help {
    Write-Host @"

Watchmen Installer

Usage:
  .\install.ps1 [options]

Options:
  -Silent         Run without prompts (for automation)
  -ConfigPath     Path to custom config file
  -Help           Show this help message

Examples:
  .\install.ps1                              # Interactive install
  .\install.ps1 -Silent                      # Silent install
  .\install.ps1 -ConfigPath "D:\config.yaml" # Custom config

"@
}

function Test-IsUpgrade {
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    return $null -ne $service
}

function Stop-ExistingService {
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($service -and $service.Status -eq "Running") {
        Write-Step "Stopping existing service..."
        Stop-Service -Name $ServiceName -Force
        Start-Sleep -Seconds 2
        Write-Step "Service stopped" "OK"
    }
}

function Backup-ExistingConfig {
    $configFile = "$InstallDir\$ConfigFileName"
    if (Test-Path $configFile) {
        $backupFile = "$configFile.backup"
        Copy-Item -Path $configFile -Destination $backupFile -Force
        Write-Step "Config backed up to config.yaml.backup" "OK"
    }
}

function Install-Files {
    param([string]$SourceDir)
    
    # Create directories
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-Step "Created installation directory" "OK"
    }
    
    if (-not (Test-Path $LogDir)) {
        New-Item -ItemType Directory -Path $LogDir -Force | Out-Null
        Write-Step "Created logs directory" "OK"
    }
    
    # Copy executable
    $sourceExe = Join-Path $SourceDir $ExeName
    if (Test-Path $sourceExe) {
        Copy-Item -Path $sourceExe -Destination "$InstallDir\$ExeName" -Force
        Write-Step "Copied $ExeName" "OK"
    } else {
        Write-Step "Executable not found: $sourceExe" "FAIL"
        return $false
    }
    
    # Copy config (only if not exists or custom path provided)
    $destConfig = "$InstallDir\$ConfigFileName"
    if ($ConfigPath -and (Test-Path $ConfigPath)) {
        Copy-Item -Path $ConfigPath -Destination $destConfig -Force
        Write-Step "Copied custom config" "OK"
    } elseif (-not (Test-Path $destConfig)) {
        $defaultConfig = Join-Path $SourceDir "config.example.yaml"
        if (Test-Path $defaultConfig) {
            Copy-Item -Path $defaultConfig -Destination $destConfig -Force
            Write-Step "Created default config" "OK"
        } else {
            Write-Step "No config template found" "WARN"
        }
    } else {
        Write-Step "Keeping existing config" "INFO"
    }
    
    return $true
}

function Register-WindowsService {
    # Check if service already exists
    $existingService = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    
    if ($existingService) {
        # Update existing service
        Write-Step "Updating existing service..." 
        sc.exe config $ServiceName binPath= "`"$InstallDir\$ExeName`" service" start= delayed-auto | Out-Null
    } else {
        # Create new service
        Write-Step "Registering Windows Service..."
        $binPath = "`"$InstallDir\$ExeName`" service"
        
        sc.exe create $ServiceName `
            binPath= $binPath `
            DisplayName= $DisplayName `
            start= delayed-auto | Out-Null
        
        # Set description
        sc.exe description $ServiceName $Description | Out-Null
    }
    
    # Verify service registration
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($service) {
        Write-Step "Windows Service registered" "OK"
        return $true
    } else {
        Write-Step "Failed to register service" "FAIL"
        return $false
    }
}

function Start-WatchmenService {
    Write-Step "Starting service..."
    
    try {
        Start-Service -Name $ServiceName -ErrorAction Stop
        Start-Sleep -Seconds 2
        
        $service = Get-Service -Name $ServiceName
        if ($service.Status -eq "Running") {
            Write-Step "Service started" "OK"
            return $true
        } else {
            Write-Step "Service failed to start (Status: $($service.Status))" "WARN"
            return $false
        }
    } catch {
        Write-Step "Failed to start service: $($_.Exception.Message)" "FAIL"
        return $false
    }
}

function Show-Summary {
    param([bool]$IsUpgrade)
    
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    $status = if ($service) { $service.Status } else { "Unknown" }
    
    Write-Host ""
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    if ($IsUpgrade) {
        Write-Host "   ✅ UPGRADE COMPLETE!" -ForegroundColor $ColorSuccess
    } else {
        Write-Host "   ✅ INSTALLATION COMPLETE!" -ForegroundColor $ColorSuccess
    }
    Write-Host ""
    Write-Host "   Service Name:  $ServiceName" -ForegroundColor White
    Write-Host "   Status:        $status" -ForegroundColor White
    Write-Host "   Startup:       Automatic (Delayed Start)" -ForegroundColor White
    Write-Host "   Config:        $InstallDir\$ConfigFileName" -ForegroundColor White
    Write-Host "   Logs:          $LogDir\" -ForegroundColor White
    Write-Host ""
    Write-Host "   Next steps:" -ForegroundColor $ColorWarning
    Write-Host "   1. Edit config.yaml to add your SQL Server(s)" -ForegroundColor White
    Write-Host "   2. Run: watchmen reload (to apply changes)" -ForegroundColor White
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    Write-Host ""
}

# =============================================================================
# Main
# =============================================================================

if ($Help) {
    Show-Help
    exit 0
}

# Get script directory (where install.ps1 is located)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Banner

# Check if this is an upgrade
$isUpgrade = Test-IsUpgrade
if ($isUpgrade) {
    Write-Step "Existing installation detected - Upgrading..." "INFO"
    Stop-ExistingService
    Backup-ExistingConfig
}

# Check for admin privileges (already required via #Requires)
Write-Step "Administrator privileges confirmed" "OK"

# Install files
if (-not (Install-Files -SourceDir $ScriptDir)) {
    Write-Step "Installation failed" "FAIL"
    exit 1
}

# Register service
if (-not (Register-WindowsService)) {
    exit 1
}

# Start service
$started = Start-WatchmenService

# Show summary
Show-Summary -IsUpgrade $isUpgrade

if (-not $Silent) {
    Write-Host "Press any key to exit..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
}

exit 0
