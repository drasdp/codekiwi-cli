#!/bin/bash

set -e

# 색상 코드
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

echo "======================================"
echo "   CodeKiwi 설치 스크립트"
echo "======================================"
echo ""

# Docker 확인
print_info "Docker 설치 확인 중..."
if ! command -v docker &> /dev/null; then
    print_error "Docker가 설치되어 있지 않습니다."
    echo ""
    echo "다음 링크에서 Docker를 설치해주세요:"
    echo "  macOS/Windows: https://www.docker.com/products/docker-desktop"
    echo "  Linux: https://docs.docker.com/engine/install/"
    echo ""
    exit 1
fi

if ! docker info &> /dev/null; then
    print_error "Docker가 실행되고 있지 않습니다."
    echo "Docker Desktop을 실행한 후 다시 시도해주세요."
    exit 1
fi

print_success "Docker 확인 완료"

# 설치 디렉토리 생성
INSTALL_DIR="$HOME/.codekiwi"
print_info "설치 디렉토리: $INSTALL_DIR"

if [ -d "$INSTALL_DIR" ]; then
    print_info "기존 설치를 발견했습니다. 업데이트합니다..."
fi

mkdir -p "$INSTALL_DIR"

# GitHub 리포지토리 URL
REPO_URL="https://raw.githubusercontent.com/your-username/codekiwi-web/main"

# 필요한 파일 다운로드
print_info "필요한 파일을 다운로드합니다..."

files=(
    "codekiwi"
    "Dockerfile"
    "docker-compose.yaml"
    "start.sh"
    "nginx.conf"
    "index.html"
    ".tmux.conf"
    ".bashrc"
)

for file in "${files[@]}"; do
    print_info "다운로드 중: $file"
    if ! curl -fsSL "$REPO_URL/$file" -o "$INSTALL_DIR/$file"; then
        print_error "$file 다운로드 실패"
        exit 1
    fi
done

# 실행 권한 부여
chmod +x "$INSTALL_DIR/codekiwi"
chmod +x "$INSTALL_DIR/start.sh"

print_success "파일 다운로드 완료"

# 심볼릭 링크 생성
print_info "글로벌 명령어 설치 중..."

BIN_DIR="/usr/local/bin"
if [ ! -d "$BIN_DIR" ]; then
    sudo mkdir -p "$BIN_DIR"
fi

# 기존 심볼릭 링크 제거
if [ -L "$BIN_DIR/codekiwi" ]; then
    sudo rm "$BIN_DIR/codekiwi"
fi

# 새 심볼릭 링크 생성
sudo ln -s "$INSTALL_DIR/codekiwi" "$BIN_DIR/codekiwi"

print_success "글로벌 명령어 설치 완료"

# Docker 이미지 빌드
print_info "Docker 이미지를 빌드합니다... (1-2분 소요)"
cd "$INSTALL_DIR"
docker build -t codekiwi:latest . > /dev/null 2>&1

print_success "Docker 이미지 빌드 완료"

echo ""
echo "======================================"
print_success "CodeKiwi 설치가 완료되었습니다!"
echo "======================================"
echo ""
echo "🚀 사용 방법:"
echo ""
echo "  1. 프로젝트 디렉토리로 이동:"
echo "     ${BLUE}cd ~/my-project${NC}"
echo ""
echo "  2. CodeKiwi 실행:"
echo "     ${BLUE}codekiwi${NC}"
echo ""
echo "  3. 브라우저에서 접속:"
echo "     ${BLUE}http://localhost:8080${NC}"
echo ""
echo "💡 기타 명령어:"
echo "  codekiwi stop        - 중지"
echo "  codekiwi restart     - 재시작"
echo "  codekiwi status      - 상태 확인"
echo "  codekiwi help        - 도움말"
echo ""
echo "📚 자세한 정보:"
echo "  https://github.com/your-username/codekiwi-web"
echo ""
