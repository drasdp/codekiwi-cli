FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ttyd \
    nginx \
    curl \
    tmux \
    git \
    vim \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g opencode-ai \
    && npm install -g @anthropic-ai/claude-code \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY nginx.conf /etc/nginx/sites-enabled/default
COPY index.html /var/www/html/index.html
COPY .tmux.conf /root/.tmux.conf
COPY start.sh /start.sh
RUN chmod +x /start.sh

# workspace 디렉토리 생성 (마운트 포인트)
RUN mkdir -p /workspace

EXPOSE 80
EXPOSE 3000

CMD ["/start.sh"]
