FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ttyd \
    nginx \
    curl \
    tmux \
    git \
    vim \
    lsof \
    dos2unix \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g opencode-ai \
    && npm install -g @anthropic-ai/claude-code \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY config/nginx.conf /etc/nginx/sites-enabled/default
COPY web/index.html /var/www/html/index.html
COPY config/.tmux.conf /root/.tmux.conf
COPY scripts/start.sh /start.sh
COPY scripts/check_and_setup.sh /check_and_setup.sh
RUN dos2unix /start.sh /check_and_setup.sh && chmod +x /start.sh /check_and_setup.sh

# workspace 디렉토리 생성 (마운트 포인트)
RUN mkdir -p /workspace

EXPOSE 80
EXPOSE 3000

CMD ["/start.sh"]