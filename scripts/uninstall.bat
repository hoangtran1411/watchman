@echo off
:: =============================================================================
:: WATCHMEN UNINSTALLER - Batch Wrapper
:: =============================================================================
:: Double-click friendly wrapper for uninstall.ps1
:: Automatically requests Administrator privileges
:: =============================================================================

cd /d "%~dp0"

:: Check for Administrator privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Requesting Administrator privileges...
    powershell -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

:: Run PowerShell uninstaller
PowerShell -ExecutionPolicy Bypass -File "%~dp0uninstall.ps1" %*

pause
