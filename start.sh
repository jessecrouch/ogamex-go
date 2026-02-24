#!/bin/bash

# OGameX Startup Script
# Prevents port conflicts and ensures clean startup

APP_NAME="ogamex"
PORT=8080

echo "=== OGameX Startup ==="

# Check if something is already running on the port
if lsof -i :$PORT > /dev/null 2>&1; then
    echo "Found existing process on port $PORT. Killing..."
    lsof -i :$PORT | grep -v COMMAND | awk '{print $2}' | xargs -r kill -9
    sleep 1
fi

# Check for any orphaned ogamex processes
if pgrep -x "$APP_NAME" > /dev/null; then
    echo "Found orphaned $APP_NAME processes. Killing..."
    pkill -9 -x "$APP_NAME"
    sleep 1
fi

echo "Starting $APP_NAME..."

# Set library path for Rust battle engine
export LD_LIBRARY_PATH=/home/bolt/Documents/ogamex-go/storage/rust-libs

# Run the server
cd /home/bolt/Documents/ogamex-go
./ogamex "$@"
