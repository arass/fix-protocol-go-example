@echo off
echo -----------------------------------------------------------------------------
echo PROJECT SETUP SCRIPT (Go FIX Client)
echo -----------------------------------------------------------------------------

:: 1. Initialize config.cfg if it's missing
if not exist "config.cfg" (
    if exist "config.cfg.example" (
        echo Initializing config.cfg from template...
        copy config.cfg.example config.cfg
        echo Go edit config.cfg for your destination
    ) else (
        echo Error: Neither config.cfg nor config.cfg.example found.
        pause
        exit /b 1
    )
)

:: 2. Create data directories
if not exist "log" mkdir log
if not exist "store" mkdir store

:: 3. Fix/Download dependencies
echo Downloading dependencies (go mod tidy)...
go mod tidy

echo.
echo Setup Complete! You can now run the app with:
echo    go run cmd/thisgofix/main.go
echo.
pause
