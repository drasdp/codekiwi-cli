# 🥝 CodeKiwi

웹 브라우저에서 사용하는 AI 기반 개발 환경

Docker 기반의 통합 개발 환경으로, 좌측에는 AI 코드 에디터, 우측에는 실시간 웹 미리보기를 제공합니다.

---

## 📖 목차

- [사용자 가이드](#-사용자-가이드) - CodeKiwi를 설치하고 사용하는 방법
- [개발자 가이드](#-개발자-가이드) - CodeKiwi 프로젝트 자체를 개발하는 방법

---

# 👤 사용자 가이드

## ✨ 특징

- 🚀 **원라이너 설치** - 한 줄로 설치 완료
- 💻 **opencode 스타일** - `cd my-project && codekiwi`로 즉시 실행
- 🤖 **AI 코드 에디터** - OpenCode AI 통합
- 👀 **실시간 미리보기** - 코드 변경 즉시 반영
- 🔢 **다중 인스턴스** - 여러 프로젝트 동시 실행 (자동 포트 할당)
- 🔄 **스마트 관리** - 프로젝트별 독립적인 컨테이너 및 포트
- 🌐 **크로스 플랫폼** - macOS, Linux, Windows 지원

## 🚀 빠른 시작

### 1️⃣ 설치 (한 번만)

```bash
curl -fsSL https://raw.githubusercontent.com/aardvarkdev1/codekiwi-cli/main/cli/scripts/install.sh | bash
```

#### 설치 과정 상세

설치 스크립트는 다음 작업을 수행합니다:

1. **Docker 확인**: Docker가 설치되어 있고 실행 중인지 확인
2. **파일 다운로드**: `~/.codekiwi/` 디렉토리에 필요한 파일들 설치
   ```
   ~/.codekiwi/
   ├── codekiwi                 # CLI 실행 파일
   ├── config.env               # 중앙 설정 파일 (포트, 경로 등)
   ├── docker-compose.yaml      # Docker Compose 설정
   └── lib/
       └── config-loader.sh     # 설정 로더 스크립트
   ```
3. **심볼릭 링크 생성**: `/usr/local/bin/codekiwi` → `~/.codekiwi/codekiwi`
4. **Docker 이미지 다운로드**: `aardvarkdev1/codekiwi-runtime:latest` 이미지 pull

### 2️⃣ 사용

```bash
# 프로젝트 디렉토리로 이동
cd ~/my-project

# CodeKiwi 실행
codekiwi

# 브라우저가 자동으로 http://localhost:8080 열림
# Ctrl+C로 종료
```

#### 실행 과정 상세

`codekiwi` 명령 실행 시:

1. **설정 로드**: `~/.codekiwi/config.env`에서 설정 값 로드
2. **포트 할당**:
   - 웹 인터페이스: 8080부터 사용 가능한 포트 자동 찾기
   - 개발 서버: 3000부터 사용 가능한 포트 자동 찾기
3. **컨테이너 시작**: Docker Compose로 컨테이너 실행
   - 작업 디렉토리를 `/workspace`로 마운트
   - 환경 변수 전달 (템플릿 설치 여부 등)
4. **브라우저 열기**: 자동으로 브라우저에서 웹 인터페이스 열기
5. **로그 스트리밍**: 포그라운드로 실행되며 실시간 로그 출력

#### 컨테이너 내부 동작

컨테이너가 시작되면:

1. **디렉토리 체크** (`check_and_setup.sh`):
   - 작업 디렉토리가 비어있으면 템플릿 설치 여부 물음
   - 템플릿에는 기본 React 앱과 개발 설정 포함
2. **개발 환경 시작** (`start.sh`):
   - **tmux 세션**: `opencode` AI 에디터를 tmux 세션에서 실행
   - **개발 서버**: `npm install && npm run dev` 백그라운드 실행
   - **웹 터미널**: ttyd가 7681 포트에서 터미널 제공
   - **프록시 서버**: nginx가 80 포트에서 다음을 프록시:
     - `/` → 웹 터미널 (ttyd)
     - `/devserver/` → 개발 서버 (포트 3000)

## 📖 사용법

### 기본 명령어

```bash
# 현재 디렉토리로 시작
codekiwi

# 특정 디렉토리 지정
codekiwi ~/my-project

# 실행 중인 모든 인스턴스 보기
codekiwi list

# 특정 프로젝트 강제 종료
codekiwi kill ~/my-project

# 모든 인스턴스 강제 종료
codekiwi kill --all
```

### 다중 프로젝트 동시 실행

```bash
# 터미널 1
cd ~/project-a
codekiwi  # localhost:8080, localhost:3000

# 터미널 2
cd ~/project-b
codekiwi  # localhost:8081, localhost:3001 (자동 할당)

# 터미널 3
codekiwi list  # 모든 실행 중인 인스턴스 확인
```

## 🎨 화면 구성

```
브라우저: http://localhost:8080
┌─────────────────────┬─────────────────────┐
│     좌측 50%        │     우측 50%        │
│                     │                     │
│  AI 코드 에디터     │  실시간 미리보기    │
│  (OpenCode AI)      │  (npm run dev)      │
│                     │                     │
│  터미널에서:        │  자동 새로고침:     │
│  - 코드 편집        │  - React/Vue/etc    │
│  - AI 지원          │  - 파일 저장 시     │
│  - 터미널 명령      │  - 즉시 반영        │
│                     │                     │
└─────────────────────┴─────────────────────┘
```

## ⚙️ 시스템 요구사항

- **Docker** 20.10 이상 및 Docker Compose 2.0 이상
- **메모리** 최소 2GB 권장
- **디스크** 약 1GB (Docker 이미지)

## 🔧 문제 해결

### Docker가 실행되지 않을 때

```bash
# macOS
open -a Docker

# Linux
sudo systemctl start docker

# Windows
# Docker Desktop 실행
```

### 포트 충돌

CodeKiwi는 자동으로 사용 가능한 포트를 찾습니다. 수동 확인이 필요한 경우:

```bash
# 포트 사용 확인
lsof -i :8080

# 실행 중인 인스턴스 확인
codekiwi list
```

### 완전히 초기화

```bash
# 모든 인스턴스 종료
codekiwi kill --all

# 제거 및 재설치
codekiwi uninstall
curl -fsSL https://raw.githubusercontent.com/aardvarkdev1/codekiwi-cli/main/cli/scripts/install.sh | bash
```

## 📚 명령어 레퍼런스

| 명령어 | 설명 |
|--------|------|
| `codekiwi` | 현재 디렉토리로 시작 |
| `codekiwi <path>` | 지정한 디렉토리로 시작 |
| `codekiwi list` | 실행 중인 모든 인스턴스 목록 |
| `codekiwi kill [path]` | 특정 인스턴스 종료 |
| `codekiwi kill --all` | 모든 인스턴스 종료 |
| `codekiwi update` | 최신 버전으로 업데이트 |
| `codekiwi uninstall` | 완전히 제거 |
| `codekiwi help` | 도움말 표시 |

---

# 👩‍💻 개발자 가이드

CodeKiwi 프로젝트를 개발하고 기여하는 방법입니다.

## 🏗️ 프로젝트 구조

```
codekiwi-web/
├── cli/                      # CLI 도구
│   ├── bin/
│   │   └── codekiwi         # 메인 CLI 스크립트
│   └── scripts/
│       └── install.sh       # 설치 스크립트
├── runtime/                 # Docker 컨테이너 런타임
│   ├── Dockerfile          # 런타임 이미지 정의
│   ├── config/
│   │   └── nginx.conf      # Nginx 프록시 설정
│   ├── scripts/
│   │   ├── check_and_setup.sh  # 템플릿 설치 스크립트
│   │   └── start.sh            # 컨테이너 시작 스크립트
│   └── web/
│       └── index.html      # 웹 인터페이스 (iframe 구조)
├── lib/
│   └── config-loader.sh    # 설정 로더 헬퍼
├── config.env              # 중앙 설정 파일 (SSOT)
├── docker-compose.yaml     # 프로덕션 Compose 파일
└── docker-compose.dev.yaml # 개발용 Compose 파일
```

## 🚀 개발 환경 설정

### 1. 저장소 클론

```bash
git clone https://github.com/aardvarkdev1/codekiwi-web.git
cd codekiwi-web
```

### 2. 개발 모드 실행

개발 모드에서는 로컬에서 Docker 이미지를 빌드하고 실행합니다:

```bash
# 프로젝트 루트에서
./cli/bin/codekiwi ~/test-project
```

개발 모드 감지:
- `docker-compose.dev.yaml` 파일이 있으면 개발 모드로 인식
- 로컬에서 `runtime/Dockerfile`을 빌드
- Docker Hub 이미지 대신 로컬 빌드 이미지 사용

### 3. 설정 관리 (SSOT)

모든 설정은 `config.env`에서 중앙 관리됩니다:

```bash
# 주요 설정 값
CODEKIWI_WEB_PORT_DEFAULT=8080      # 웹 인터페이스 기본 포트
CODEKIWI_DEV_PORT_DEFAULT=3000      # 개발 서버 기본 포트
CODEKIWI_TTYD_PORT=7681             # 웹 터미널 포트
CODEKIWI_WORKSPACE_DIR=/workspace   # 컨테이너 내 작업 디렉토리
```

설정 로드 방식:
- **CLI**: `lib/config-loader.sh`를 source하여 Bash 변수로 로드
- **Runtime**: Docker Compose의 `env_file`로 환경 변수 주입

## 📦 빌드 및 배포

### 로컬 이미지 빌드

```bash
# 개발용 이미지 빌드
docker-compose -f docker-compose.dev.yaml build
```

### 프로덕션 배포

1. Docker Hub에 이미지 푸시:
```bash
docker build -t aardvarkdev1/codekiwi-runtime:latest ./runtime
docker push aardvarkdev1/codekiwi-runtime:latest
```

2. GitHub에 코드 푸시:
```bash
git add .
git commit -m "Update version"
git push origin main
```

사용자는 `codekiwi update` 명령으로 최신 버전을 받을 수 있습니다.

## 🧪 테스트

### 설치 테스트

```bash
# 로컬 install.sh 테스트
./cli/scripts/install.sh
```

### CLI 테스트

```bash
# 다양한 시나리오 테스트
./cli/bin/codekiwi                    # 현재 디렉토리
./cli/bin/codekiwi ~/test-dir         # 특정 디렉토리
./cli/bin/codekiwi list               # 목록 확인
./cli/bin/codekiwi kill ~/test-dir    # 종료
```

## 🤝 기여 방법

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Test thoroughly in development mode
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### 코드 스타일

- Bash 스크립트: ShellCheck 권장사항 준수
- Docker: 멀티스테이지 빌드로 이미지 크기 최소화
- 설정: 모든 하드코딩 값은 `config.env`에 정의

## 📄 라이선스

MIT License

## 🙏 감사

- [ttyd](https://github.com/tsl0922/ttyd) - Web terminal emulator
- [OpenCode AI](https://github.com/opencodeiiit/opencode-ai) - AI code editor
- [tmux](https://github.com/tmux/tmux) - Terminal multiplexer
- [nginx](https://nginx.org/) - Web server

---

**Made with ❤️ for developers**