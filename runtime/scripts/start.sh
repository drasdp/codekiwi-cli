#!/bin/bash

# UTF-8 설정
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

# Load configuration from environment (set by docker-compose)
WORKSPACE="${CODEKIWI_WORKSPACE_DIR:-/workspace}"
TTYD_PORT="${CODEKIWI_TTYD_PORT:-7681}"
DEV_PORT="${CODEKIWI_DEV_PORT_DEFAULT:-3000}"

# 디렉토리 체크 및 템플릿 설정
/check_and_setup.sh

# OpenCode auth.json 확인 및 설정
AUTH_DIR="/root/.local/share/opencode"
if [ -f "$AUTH_DIR/auth.json" ]; then
    echo "OpenCode auth.json found at $AUTH_DIR/auth.json"
else
    echo "Warning: OpenCode auth.json not found at $AUTH_DIR/auth.json"
    echo "OpenCode may require API key configuration."
fi

echo "Checking for package.json..."
if [ -f "$WORKSPACE/package.json" ]; then
    echo "package.json found, starting dev server..."
    cd "$WORKSPACE"

    # Install dependencies if node_modules is empty or doesn't exist
    # Named volume creates the directory but it may be empty
    if [ ! -d "$WORKSPACE/node_modules/.bin" ] || [ -z "$(ls -A $WORKSPACE/node_modules 2>/dev/null)" ]; then
        echo "Installing dependencies..."
        # Use npm ci if package-lock.json exists (faster and more reliable)
        if [ -f "$WORKSPACE/package-lock.json" ]; then
            echo "Using npm ci for faster installation..."
            npm ci --prefer-offline
        else
            echo "Using npm install..."
            npm install
        fi
    else
        echo "Dependencies already installed, skipping..."
    fi

    echo "✅ npm install completed successfully"

    echo "Starting devserver in background..."
    npm run dev &

    echo "Waiting for dev server to start..."
    for i in {1..60}; do
        if curl -s http://localhost:$DEV_PORT > /dev/null 2>&1; then
            echo "Dev server is ready!"
            # Create ready flag for health check
            touch /tmp/services_ready
            break
        fi
        echo "Waiting for dev server... ($i/60)"
        sleep 1
    done

    if [ ! -f /tmp/services_ready ]; then
        echo "Warning: Dev server did not start within 60 seconds"
        echo "Continuing anyway..."
        touch /tmp/services_ready
    fi
else
    echo "No package.json found. Skipping dev server."
    echo "Dev server will be available after template installation."
    # Create ready flag anyway for health check
    touch /tmp/services_ready
fi

echo "Creating tmux session with opencode..."
tmux new-session -d -s opencode bash -c "cd \"$WORKSPACE\" && opencode; exec bash"

echo "Starting ttyd on port $TTYD_PORT (opencode)..."
ttyd -p "$TTYD_PORT" -t fontSize=14 -t 'theme={"background":"#1e1e1e"}' tmux attach-session -t opencode &

echo "Starting nginx..."
nginx -g 'daemon off;'
