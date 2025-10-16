#!/bin/bash

# Simple watch script for hot reload
# Watches Go files and restarts the server when changes are detected

BINARY="./zipcodereader"
BUILD_CMD="go build -o $BINARY ."
RUN_ARGS="${@:-}"  # Pass any arguments to the binary

echo "🔄 Starting hot reload watcher..."
echo "📝 Watching for changes in *.go files"
echo "🚀 Server will restart automatically on file changes"
echo "Press Ctrl+C to stop"
echo ""

# Kill any existing zipcodereader processes
pkill -f zipcodereader 2>/dev/null

# Build initial binary
echo "🔨 Building..."
if ! $BUILD_CMD; then
    echo "❌ Build failed"
    exit 1
fi

# Start the server in background
echo "✅ Starting server..."
$BINARY $RUN_ARGS &
SERVER_PID=$!
echo "📍 Server running (PID: $SERVER_PID)"
echo ""

# Watch for file changes
watch_files() {
    # Get initial timestamp
    LAST_CHANGE=$(find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./.git/*" -exec stat -f "%m" {} \; | sort -n | tail -1)

    while true; do
        sleep 1

        # Get current timestamp
        CURRENT_CHANGE=$(find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./.git/*" -exec stat -f "%m" {} \; | sort -n | tail -1)

        # Check if files changed
        if [ "$CURRENT_CHANGE" != "$LAST_CHANGE" ]; then
            echo ""
            echo "🔄 Change detected, rebuilding..."

            # Kill old server
            kill $SERVER_PID 2>/dev/null
            wait $SERVER_PID 2>/dev/null

            # Rebuild
            if $BUILD_CMD; then
                echo "✅ Build successful, restarting server..."
                $BINARY $RUN_ARGS &
                SERVER_PID=$!
                echo "📍 Server restarted (PID: $SERVER_PID)"
                echo ""
            else
                echo "❌ Build failed, keeping old server"
                $BINARY $RUN_ARGS &
                SERVER_PID=$!
            fi

            LAST_CHANGE=$CURRENT_CHANGE
        fi
    done
}

# Trap Ctrl+C and cleanup
trap 'echo ""; echo "🛑 Stopping..."; kill $SERVER_PID 2>/dev/null; pkill -f zipcodereader 2>/dev/null; exit 0' INT TERM

# Start watching
watch_files
