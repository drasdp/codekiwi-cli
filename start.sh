#!/bin/bash

# UTF-8 설정
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

echo "Creating tmux session with opencode..."
tmux new-session -d -s opencode bash -c 'opencode; exec bash'

echo "Starting devserver in background..."
cd /workspace && npm install && npm run dev &

echo "Starting ttyd on port 7681 (opencode)..."
ttyd -p 7681 -t fontSize=14 -t 'theme={"background":"#1e1e1e"}' tmux attach-session -t opencode &

echo "Starting nginx..."
nginx -g 'daemon off;'
