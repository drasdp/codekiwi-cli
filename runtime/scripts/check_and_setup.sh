#!/bin/bash

# UTF-8 설정
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

# Load configuration from environment (set by docker-compose)
WORKSPACE="${CODEKIWI_WORKSPACE_DIR:-/workspace}"

# 템플릿 저장소 URL (public repo)
TEMPLATE_REPO="${CODEKIWI_TEMPLATE_GITHUB_URL:-https://github.com/aardvark-team-dev/codekiwi-template.git}"

# 디렉토리가 비어있는지 확인
is_empty_dir() {
    [ -z "$(ls -A "$WORKSPACE" 2>/dev/null)" ]
}

# 템플릿 설치 함수
install_template() {
    echo "템플릿을 설치합니다..."

    # git clone으로 템플릿 가져오기
    git clone "$TEMPLATE_REPO" /tmp/template

    if [ $? -eq 0 ]; then
        # .git 디렉토리 제외하고 모든 파일 복사
        cp -r /tmp/template/* "$WORKSPACE/"
        cp -r /tmp/template/.[!.]* "$WORKSPACE/" 2>/dev/null
        rm -rf /tmp/template

        echo "✅ 템플릿 설치가 완료되었습니다!"

        # AUTH_SECRET 설정
        if [ -f "$WORKSPACE/.env.local" ]; then
            # .env.local이 있으면 AUTH_SECRET이 비어있는지 확인
            if grep -q '^AUTH_SECRET=""' "$WORKSPACE/.env.local" || grep -q '^AUTH_SECRET=$' "$WORKSPACE/.env.local"; then
                echo "🔑 AUTH_SECRET을 생성합니다..."
                # 안전한 암호학적 랜덤 문자열 생성 (base64, 32바이트)
                AUTH_SECRET=$(openssl rand -base64 32)
                # AUTH_SECRET 값 업데이트
                sed -i.bak "s|^AUTH_SECRET=.*|AUTH_SECRET=\"$AUTH_SECRET\"|" "$WORKSPACE/.env.local"
                rm -f "$WORKSPACE/.env.local.bak"
                echo "✅ AUTH_SECRET이 생성되었습니다!"
            fi
        else
            # .env.local이 없으면 생성
            echo "🔑 .env.local 파일을 생성합니다..."
            AUTH_SECRET=$(openssl rand -base64 32)
            echo "AUTH_SECRET=\"$AUTH_SECRET\" # Added by CodeKiwi" > "$WORKSPACE/.env.local"
            echo "✅ .env.local 파일이 생성되었습니다!"
        fi

        # package.json이 있으면 의존성 설치
        if [ -f "$WORKSPACE/package.json" ]; then
            echo "📦 의존성을 설치합니다..."
            cd "$WORKSPACE"
            # Use npm ci if package-lock.json exists (faster and more reliable)
            if [ -f "$WORKSPACE/package-lock.json" ]; then
                echo "Using npm ci for faster installation..."
                npm ci --prefer-offline
            else
                echo "Using npm install..."
                npm install
            fi
            echo "✅ 의존성 설치 완료!"
        fi

        return 0
    else
        echo "❌ 템플릿 설치에 실패했습니다."
        return 1
    fi
}

# 메인 로직
if is_empty_dir; then
    echo "======================================"
    echo "  CodeKiwi 초기 설정"
    echo "======================================"
    echo ""
    echo "📁 폴더가 비어있습니다."
    echo ""

    # 환경 변수로 템플릿 설치 여부 확인 (기본값: yes)
    INSTALL_TEMPLATE="${INSTALL_TEMPLATE:-${CODEKIWI_INSTALL_TEMPLATE_DEFAULT:-yes}}"

    if [ "$INSTALL_TEMPLATE" = "yes" ] || [ "$INSTALL_TEMPLATE" = "y" ]; then
        echo "CodeKiwi 템플릿을 설치합니다..."
        echo "템플릿에는 기본 개발 환경 설정과 예제 파일이 포함되어 있습니다."
        echo ""

        if install_template; then
            echo ""
            echo "✅ 템플릿 설치 완료!"
        else
            echo ""
            echo "❌ 템플릿 설치에 실패했습니다."
            echo "컨테이너를 종료합니다."
            echo ""
            echo "======================================"
            echo ""
            exit 1
        fi
    else
        echo "템플릿 설치를 건너뜁니다."
        echo "빈 디렉토리로 계속 진행합니다."
    fi

    echo ""
    echo "======================================"
    echo ""
else
    echo "✅ 작업 디렉토리에 파일이 존재합니다. 계속 진행합니다."
fi
