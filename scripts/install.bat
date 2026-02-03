@echo off
:: =============================================================================
:: WATCHMEN INSTALLER - Batch Wrapper
:: =============================================================================
:: Double-click friendly wrapper for install.ps1
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

:: Run PowerShell installer
PowerShell -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*

pause
