#!/bin/bash

# -----------------------------------------------------------------------------
# PROJECT STARTUP SCRIPT
# -----------------------------------------------------------------------------
# This script prepares the environment and starts the Go FIX Client.
#
# Usage:
#   chmod +x run-fix-cert.sh
#   ./run-fix-cert.sh [arguments]
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
mkdir -p log store

# 5. Fix/Download dependencies
echo "System: Downloading dependencies (go mod tidy)..."
go mod tidy

# 6. Run the Application
#    "$@" passes all arguments from this script directly to the go run command.
echo "System: Starting application with args: $@"
echo "-----------------------------------------------------------------------------"
go run cmd/fixtest/main.go "$@"
