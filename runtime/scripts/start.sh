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

echo "Creating tmux session with opencode..."
tmux new-session -d -s opencode bash -c "cd \"$WORKSPACE\" && opencode; exec bash"

echo "Starting devserver in background..."
cd "$WORKSPACE" && npm install && npm run dev &

echo "Waiting for dev server to start..."
for i in {1..30}; do
    if curl -s http://localhost:$DEV_PORT > /dev/null 2>&1; then
        echo "Dev server is ready!"
        break
    fi
    echo "Waiting for dev server... ($i/30)"
    sleep 1
done

echo "Starting ttyd on port $TTYD_PORT (opencode)..."
ttyd -p "$TTYD_PORT" -t fontSize=14 -t 'theme={"background":"#1e1e1e"}' tmux attach-session -t opencode &

echo "Starting nginx..."
nginx -g 'daemon off;'
