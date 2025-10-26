#!/bin/bash

set -e

# Load configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

if [ -f "$PROJECT_ROOT/lib/config-loader.sh" ]; then
    source "$PROJECT_ROOT/lib/config-loader.sh"
fi

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
INSTALL_DIR="${CODEKIWI_INSTALL_DIR:-$HOME/.codekiwi}"
print_info "설치 디렉토리: $INSTALL_DIR"

if [ -d "$INSTALL_DIR" ]; then
    print_info "기존 설치를 발견했습니다. 업데이트합니다..."
fi

mkdir -p "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR/lib"

# GitHub 리포지토리 URL
REPO_URL="${RAW_CODEKIWI_GITHUB_URL:-https://raw.githubusercontent.com/drasdp/codekiwi-cli/main}"

# 필요한 파일 다운로드
print_info "필요한 파일을 다운로드합니다..."

# CLI 파일 다운로드
print_info "다운로드 중: codekiwi"
if ! curl -fsSL "$REPO_URL/cli/bin/codekiwi" -o "$INSTALL_DIR/codekiwi"; then
    print_error "codekiwi 다운로드 실패"
    exit 1
fi

# 설정 파일 다운로드
print_info "다운로드 중: config.env"
if ! curl -fsSL "$REPO_URL/config.env" -o "$INSTALL_DIR/config.env"; then
    print_error "config.env 다운로드 실패"
    exit 1
fi

# config-loader 다운로드
print_info "다운로드 중: lib/config-loader.sh"
if ! curl -fsSL "$REPO_URL/lib/config-loader.sh" -o "$INSTALL_DIR/lib/config-loader.sh"; then
    print_error "config-loader.sh 다운로드 실패"
    exit 1
fi

# docker-compose.yaml 다운로드 (루트에서)
print_info "다운로드 중: docker-compose.yaml"
if ! curl -fsSL "$REPO_URL/docker-compose.yaml" -o "$INSTALL_DIR/docker-compose.yaml"; then
    print_error "docker-compose.yaml 다운로드 실패"
    exit 1
fi

# 실행 권한 부여
chmod +x "$INSTALL_DIR/codekiwi"
chmod +x "$INSTALL_DIR/lib/config-loader.sh"

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

# Docker 이미지 pull
print_info "Docker 이미지를 다운로드합니다..."
IMAGE_NAME="${CODEKIWI_FULL_IMAGE_NAME:-drasdp/codekiwi-runtime}"
IMAGE_TAG="${CODEKIWI_IMAGE_TAG_DEFAULT:-latest}"
docker pull "$IMAGE_NAME:$IMAGE_TAG"

print_success "Docker 이미지 다운로드 완료"

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
echo "  codekiwi list        - 실행 중인 인스턴스 목록"
echo "  codekiwi kill        - 종료"
echo "  codekiwi update      - 업데이트"
echo "  codekiwi help        - 도움말"
echo ""
echo "📚 자세한 정보:"
echo "  https://github.com/drasdp/codekiwi-cli"
echo ""
