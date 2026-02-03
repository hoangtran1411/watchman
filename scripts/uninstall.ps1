# =============================================================================
# WATCHMEN UNINSTALLER
# =============================================================================
# PowerShell script to uninstall Watchmen Windows Service
#
# Usage:
#   .\uninstall.ps1                    # Interactive mode
#   .\uninstall.ps1 -KeepConfig        # Keep config and logs
#   .\uninstall.ps1 -RemoveAll         # Remove everything
#   .\uninstall.ps1 -Yes               # No confirmation prompt
#
# Requires: Administrator privileges
# =============================================================================

#Requires -Version 5.1
#Requires -RunAsAdministrator

[CmdletBinding()]
param(
    [switch]$KeepConfig,
    [switch]$RemoveAll,
    [switch]$Yes,
    [switch]$Help
)

# =============================================================================
# Configuration
# =============================================================================
$ServiceName = "Watchmen"
$InstallDir = "$env:ProgramData\Watchmen"

# Colors
$ColorSuccess = "Green"
$ColorError = "Red"
$ColorWarning = "Yellow"
$ColorInfo = "Cyan"

# =============================================================================
# Functions
# =============================================================================

function Write-Banner {
    Write-Host ""
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    Write-Host "   🗑️ WATCHMEN UNINSTALLER" -ForegroundColor $ColorInfo
    Write-Host "   SQL Server Agent Job Monitor" -ForegroundColor $ColorInfo
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    Write-Host ""
}

function Write-Step {
    param([string]$Message, [string]$Status = "...")
    
    switch ($Status) {
        "OK" { Write-Host "[✓] $Message" -ForegroundColor $ColorSuccess }
        "FAIL" { Write-Host "[✗] $Message" -ForegroundColor $ColorError }
        "WARN" { Write-Host "[!] $Message" -ForegroundColor $ColorWarning }
        "INFO" { Write-Host "[i] $Message" -ForegroundColor $ColorInfo }
        default { Write-Host "[...] $Message" -ForegroundColor White }
    }
}

function Show-Help {
    Write-Host @"

Watchmen Uninstaller

Usage:
  .\uninstall.ps1 [options]

Options:
  -KeepConfig     Keep configuration and log files
  -RemoveAll      Remove everything including config and logs
  -Yes            Skip confirmation prompt
  -Help           Show this help message

Examples:
  .\uninstall.ps1                    # Interactive uninstall
  .\uninstall.ps1 -KeepConfig        # Keep config for reinstall
  .\uninstall.ps1 -RemoveAll -Yes    # Remove all, no prompt

"@
}

function Test-ServiceExists {
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    return $null -ne $service
}

function Stop-WatchmenService {
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    
    if ($service -and $service.Status -eq "Running") {
        Write-Step "Stopping service..."
        try {
            Stop-Service -Name $ServiceName -Force -ErrorAction Stop
            Start-Sleep -Seconds 2
            Write-Step "Service stopped" "OK"
        }
        catch {
            Write-Step "Failed to stop service: $($_.Exception.Message)" "WARN"
        }
    }
    elseif ($service) {
        Write-Step "Service is not running" "INFO"
    }
}

function Remove-WatchmenService {
    Write-Step "Removing Windows Service..."
    
    try {
        sc.exe delete $ServiceName | Out-Null
        Start-Sleep -Seconds 1
        
        # Verify removal
        $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
        if ($null -eq $service) {
            Write-Step "Windows Service removed" "OK"
            return $true
        }
        else {
            Write-Step "Service still exists (may require reboot)" "WARN"
            return $true
        }
    }
    catch {
        Write-Step "Failed to remove service: $($_.Exception.Message)" "FAIL"
        return $false
    }
}

function Remove-InstallDirectory {
    param([bool]$KeepConfig)
    
    if (-not (Test-Path $InstallDir)) {
        Write-Step "Installation directory not found" "INFO"
        return
    }
    
    if ($KeepConfig) {
        # Remove only executable, keep config and logs
        $exePath = "$InstallDir\watchmen.exe"
        if (Test-Path $exePath) {
            Remove-Item -Path $exePath -Force
            Write-Step "Removed executable" "OK"
        }
        Write-Step "Kept config and logs in $InstallDir" "INFO"
    }
    else {
        # Remove everything
        try {
            Remove-Item -Path $InstallDir -Recurse -Force
            Write-Step "Removed installation directory" "OK"
        }
        catch {
            Write-Step "Failed to remove directory: $($_.Exception.Message)" "WARN"
        }
    }
}

function Get-UserChoice {
    if ($Yes -or $RemoveAll -or $KeepConfig) {
        return $KeepConfig
    }
    
    Write-Host ""
    Write-Host "Do you want to keep configuration and log files?" -ForegroundColor $ColorWarning
    Write-Host "  [Y] Yes - Keep config for future reinstall"
    Write-Host "  [N] No  - Remove everything"
    Write-Host ""
    
    $choice = Read-Host "Enter choice [Y/N]"
    return $choice -match "^[Yy]"
}

function Show-Summary {
    param([bool]$KeptConfig)
    
    Write-Host ""
    Write-Host ("=" * 63) -ForegroundColor $ColorInfo
    Write-Host "   ✅ UNINSTALLATION COMPLETE!" -ForegroundColor $ColorSuccess
    Write-Host ""
    
    if ($KeptConfig -and (Test-Path $InstallDir)) {
        Write-Host "   Config preserved at: $InstallDir" -ForegroundColor White
        Write-Host "   To reinstall, run install.ps1" -ForegroundColor White
    }
    else {
        Write-Host "   All files have been removed" -ForegroundColor White
    }
    
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

Write-Banner

# Check if service exists
if (-not (Test-ServiceExists)) {
    Write-Step "Watchmen service not found" "WARN"
    
    if (Test-Path $InstallDir) {
        Write-Step "Found installation directory, cleaning up..." "INFO"
    }
    else {
        Write-Host "Nothing to uninstall." -ForegroundColor $ColorInfo
        exit 0
    }
}

# Confirmation (unless -Yes flag)
if (-not $Yes -and -not $RemoveAll -and -not $KeepConfig) {
    Write-Host "This will uninstall Watchmen from your system." -ForegroundColor $ColorWarning
    $confirm = Read-Host "Are you sure you want to continue? [Y/N]"
    
    if ($confirm -notmatch "^[Yy]") {
        Write-Host "Uninstallation cancelled." -ForegroundColor $ColorInfo
        exit 0
    }
}

# Check for admin privileges (already required via #Requires)
Write-Step "Administrator privileges confirmed" "OK"

# Stop service
Stop-WatchmenService

# Remove service
Remove-WatchmenService

# Ask about config
$keepConfigChoice = Get-UserChoice

# Remove files
Remove-InstallDirectory -KeepConfig $keepConfigChoice

# Show summary
Show-Summary -KeptConfig $keepConfigChoice

if (-not $Yes) {
    Write-Host "Press any key to exit..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
}

exit 0
