#!/bin/bash

# -----------------------------------------------------------------------------
# PROJECT STARTUP SCRIPT
# -----------------------------------------------------------------------------
# This script prepares the environment and starts the Go FIX Client.
#
# Usage:
#   chmod +x run.sh
#   ./run.sh
# -----------------------------------------------------------------------------

# 1. Stop if any command fails
set -e

echo "System: Preparing to start FIX Client..."

# 2. Check if 'go' is installed
if ! command -v go &> /dev/null; then
    echo "Error: 'go' command not found. Please install Go (1.20+) first."
    exit 1
fi

# 3. Check for config.cfg, copy from example if absent
#    This allows developers to have their own local settings without
#    worrying about overwriting them on 'git pull'.
if [ ! -f "config.cfg" ]; then
    if [ -f "config.cfg.example" ]; then
        echo "System: Initializing config.cfg from example template..."
        cp config.cfg.example config.cfg
        echo "WARNING: YOU MUST EDIT config.cfg to configure destination settings"
    else
        echo "Error: Neither config.cfg nor config.cfg.example found."
        exit 1
    fi
fi

# 4. Create necessary directories for the app
#    The app will write its logs and session data here.
mkdir -p log store

# 5. Fix/Download dependencies
#    This ensures go.sum is correct and the quickfix library is downloaded.
echo "System: Downloading dependencies (go mod tidy)..."
go mod tidy

# 6. Run the Application
#    We point to the entry point in cmd/thisgofix/
echo "System: Starting application..."
echo "-----------------------------------------------------------------------------"
go run cmd/thisgofix/main.go
