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
- 💻 **간단한 사용법** - `cd my-project && codekiwi start`로 즉시 실행
- 🤖 **AI 코드 에디터** - OpenCode AI 통합
- 👀 **실시간 미리보기** - 코드 변경 즉시 반영
- 🔢 **다중 인스턴스** - 여러 프로젝트 동시 실행 (자동 포트 할당)
- 🔄 **스마트 관리** - 프로젝트별 독립적인 컨테이너 및 포트
- 🎯 **포그라운드 실행** - 실시간 로그, Ctrl+C로 자동 정리 (v0.2.0+)
- 🏃 **백그라운드 옵션** - `-d` 플래그로 데몬 모드 실행 가능
- 🌐 **크로스 플랫폼** - macOS, Linux, Windows 지원

## 🚀 빠른 시작

### 1️⃣ 설치 (한 번만)

```bash
curl -fsSL https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/cli/scripts/install.sh | bash
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
4. **Docker 이미지 다운로드**: `drasdp/codekiwi-runtime:latest` 이미지 pull

### 2️⃣ 사용

```bash
# 프로젝트 디렉토리로 이동
cd ~/my-project

# CodeKiwi 실행 (포그라운드 모드 - 기본)
codekiwi start

# 브라우저가 자동으로 http://localhost:8080 열림
# 로그가 실시간으로 표시됨
# Ctrl+C로 자동 종료 및 정리

# 백그라운드 모드로 실행 (선택사항)
codekiwi start -d
# 또는
codekiwi start --detached
```

#### 실행 과정 상세

`codekiwi start` 명령 실행 시:

1. **설정 로드**: `~/.codekiwi/config.env`에서 설정 값 로드
2. **포트 할당**: WEB_PORT 8080부터 사용 가능한 포트 자동 찾기
3. **컨테이너 시작**: Docker Compose로 컨테이너 실행
   - 작업 디렉토리를 `/workspace`로 마운트
   - 환경 변수 전달 (템플릿 설치 여부 등)
4. **헬스체크 대기**: 서비스가 준비될 때까지 대기 (최대 60초)
5. **브라우저 열기**: 서비스 준비 완료 후 자동으로 브라우저 열기 (`-n` 플래그로 비활성화 가능)
6. **실행 모드**:
   - **포그라운드 (기본)**: 로그 스트리밍, Ctrl+C로 자동 정리
   - **백그라운드 (`-d`)**: 시작 후 프롬프트 복귀, `codekiwi kill`로 수동 종료

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
     - `/` → 웹 UI (index.html)
     - `/terminal/` → 웹 터미널 (ttyd)
     - `/preview/` → 개발 서버

## 📖 사용법

### 기본 명령어

```bash
# 현재 디렉토리로 시작 (포그라운드 모드)
codekiwi start

# 백그라운드 모드로 시작
codekiwi start -d

# 특정 디렉토리 지정
codekiwi start ~/my-project

# 특정 포트로 시작
codekiwi start -p 8090

# 브라우저 열지 않고 시작
codekiwi start -n

# 실행 중인 모든 인스턴스 보기
codekiwi list
# 또는 별칭 사용
codekiwi ls

# 중지된 인스턴스 포함하여 보기
codekiwi list -a

# 컨테이너 이름만 표시 (quiet 모드)
codekiwi list -q

# 특정 프로젝트 종료
codekiwi kill ~/my-project

# 확인 없이 강제 종료
codekiwi kill -f ~/my-project

# 모든 인스턴스 종료
codekiwi kill -a
# 또는
codekiwi kill --all
```

### 다중 프로젝트 동시 실행

```bash
# 백그라운드 모드로 여러 프로젝트 실행
cd ~/project-a
codekiwi start -d  # localhost:8080

cd ~/project-b
codekiwi start -d  # localhost:8081 (자동 할당)

cd ~/project-c
codekiwi start -d  # localhost:8082 (자동 할당)

# 실행 중인 모든 인스턴스 확인
codekiwi list

# 포그라운드 모드는 터미널별로 실행
# 터미널 1: cd ~/project-a && codekiwi start
# 터미널 2: cd ~/project-b && codekiwi start
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

## 🌐 포트 구조 및 네트워크

CodeKiwi는 **nginx 기반 경로 라우팅**을 통해 단일 포트로 여러 서비스를 제공합니다.

### 포트 매핑 구조

```
[브라우저] → [호스트 머신] → [Docker 컨테이너 내부]

🌐 WEB_PORT (기본: 8080) - 단일 진입점
   http://localhost:8080
   ↓ Docker 포트 매핑 (8080:80)
   ↓ nginx :80 (컨테이너 내부)
   ├─ /                        → 웹 UI (index.html)
   ├─ /terminal/               → ttyd :7681 (웹 터미널)
   └─ /preview/                → dev server :3000
      ├─ /preview/             → 웹 애플리케이션 루트
      └─ /preview/codekiwi-dashboard → 대시보드

🔒 내부 포트 (호스트에 노출 안 됨)
   - dev server :3000 (컨테이너 내부)
   - ttyd :7681 (컨테이너 내부)
   → 모두 nginx 리버스 프록시를 통해서만 접근 가능
   → 호스트의 localhost:3000은 사용 불가
```

### 멀티 인스턴스 포트 할당

여러 프로젝트를 동시에 실행할 때, 각 인스턴스는 자동으로 다른 WEB_PORT를 할당받습니다:

```bash
# 첫 번째 인스턴스
$ cd ~/project-a && codekiwi
→ WEB: localhost:8080

# 두 번째 인스턴스
$ cd ~/project-b && codekiwi
→ WEB: localhost:8081

# 세 번째 인스턴스
$ cd ~/project-c && codekiwi
→ WEB: localhost:8082
```

각 인스턴스는 **하나의 포트만 사용**하며, 내부 서비스들은 nginx를 통해 경로로 구분됩니다.

### 웹 UI의 버튼 동작

우측 패널의 버튼들은 **현재 접속한 호스트를 기준**으로 새 창을 엽니다:

- **웹페이지 보기**: `현재주소/preview/`
  - 예: `localhost:8081` 접속 중 → `localhost:8081/preview/` 새 창
- **대시보드 보기**: `현재주소/preview/codekiwi-dashboard`
  - 예: `localhost:8081` 접속 중 → `localhost:8081/preview/codekiwi-dashboard` 새 창

### 단일 포트 구조의 장점

- ✅ **단순함**: 포트 하나만 기억하면 됨 (localhost:8080)
- ✅ **보안**: 최소한의 포트만 노출 (공격 표면 최소화)
  - dev server(3000), ttyd(7681)는 호스트에 직접 노출 안 됨
  - nginx가 유일한 진입점 역할
- ✅ **멀티 인스턴스**: 포트 충돌 없이 여러 프로젝트 동시 실행
  - 인스턴스당 포트 1개만 필요 (WEB_PORT)
- ✅ **일관성**: 모든 서비스가 동일한 경로 패턴 사용

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
curl -fsSL https://raw.githubusercontent.com/drasdp/codekiwi-cli/main/cli/scripts/install.sh | bash
```

## 📚 명령어 레퍼런스

### 명령어

| 명령어 | 별칭 | 설명 | 주요 플래그 |
|--------|------|------|------------|
| `codekiwi start [path]` | - | CodeKiwi 시작 | `-d` 백그라운드 모드<br>`-n` 브라우저 안 열기<br>`-p` 웹 포트 지정<br>`-t` 템플릿 지정 |
| `codekiwi list` | `ls` | 실행 중인 인스턴스 목록 | `-a` 모든 인스턴스 표시<br>`-q` 컨테이너 이름만 표시 |
| `codekiwi kill [path\|container]` | - | 인스턴스 종료 | `-f` 확인 없이 강제<br>`-a` 모든 인스턴스 종료 |
| `codekiwi update` | - | CLI와 Docker 이미지 업데이트 | `--skip-cli` CLI 업데이트 건너뛰기<br>`--skip-image` 이미지 업데이트 건너뛰기 |
| `codekiwi uninstall` | - | CodeKiwi 제거 | `-f` 확인 없이 강제<br>`--keep-data` 설정 파일 유지<br>`--keep-images` Docker 이미지 유지 |

### 전역 플래그

| 플래그 | 설명 |
|--------|------|
| `--help`, `-h` | 도움말 표시 |
| `--version`, `-v` | 버전 정보 표시 |

### Start 명령어 상세 플래그

| 플래그 | 단축 | 설명 | 기본값 |
|--------|------|------|--------|
| `--detached` | `-d` | 백그라운드 모드로 실행 | false (포그라운드) |
| `--no-open` | `-n` | 브라우저 자동 열기 비활성화 | false |
| `--web-port` | `-p` | 웹 인터페이스 포트 지정 | 자동 (8080부터) |
| `--dev-port` | | 개발 서버 포트 지정 | 자동 (3000부터) |
| `--template` | `-t` | 프로젝트 템플릿 지정 | 자동 감지 |

---

# 👩‍💻 개발자 가이드

CodeKiwi 프로젝트를 개발하고 기여하는 방법입니다.

## 🏗️ 프로젝트 구조

```
codekiwi-web/
├── cli-go/                   # Go 기반 CLI (v0.2.0+)
│   ├── cmd/codekiwi/
│   │   └── main.go          # 메인 엔트리포인트
│   ├── internal/
│   │   ├── commands/        # CLI 명령어 구현
│   │   │   ├── start.go     # start 명령
│   │   │   ├── list.go      # list 명령
│   │   │   ├── kill.go      # kill 명령
│   │   │   ├── update.go    # update 명령
│   │   │   └── uninstall.go # uninstall 명령
│   │   ├── config/          # 설정 관리
│   │   ├── docker/          # Docker 작업
│   │   ├── platform/        # 플랫폼별 기능
│   │   └── state/           # 인스턴스 상태 관리
│   ├── go.mod               # Go 모듈 정의
│   └── go.sum               # 의존성 lock 파일
├── cli/                      # 레거시 Bash CLI (deprecated)
│   └── scripts/
│       └── install.sh       # 설치 스크립트
├── runtime/                 # Docker 컨테이너 런타임
│   ├── Dockerfile          # 런타임 이미지 정의
│   ├── config/
│   │   └── nginx.conf      # Nginx 프록시 설정 (health 엔드포인트 포함)
│   ├── scripts/
│   │   ├── check_and_setup.sh  # 템플릿 설치 스크립트
│   │   └── start.sh            # 컨테이너 시작 스크립트
│   └── web/
│       └── index.html      # 웹 인터페이스 (iframe 구조)
├── config.env              # 중앙 설정 파일 (SSOT)
├── docker-compose.yaml     # 프로덕션 Compose 파일 (healthcheck 포함)
└── docker-compose.dev.yaml # 개발용 Compose 파일 (healthcheck 포함)
```

## 🚀 개발 환경 설정

### 1. 저장소 클론

```bash
git clone https://github.com/drasdp/codekiwi-web.git
cd codekiwi-web
```

### 2. 개발 모드 실행

#### Go CLI 빌드 및 실행

```bash
# Go CLI 빌드
cd cli-go
go build -o codekiwi cmd/codekiwi/main.go

# 개발 모드로 실행
./codekiwi start ~/test-project
```

#### 개발 모드 감지

- `docker-compose.dev.yaml` 파일이 있으면 개발 모드로 인식
- 로컬에서 `runtime/Dockerfile`을 빌드
- Docker Hub 이미지 대신 로컬 빌드 이미지 사용

#### Go CLI 개발 워크플로우

```bash
# 코드 수정 후
cd cli-go
go build -o codekiwi cmd/codekiwi/main.go

# 테스트
./codekiwi start              # 포그라운드 모드
./codekiwi start -d           # 백그라운드 모드
./codekiwi list               # 인스턴스 목록
./codekiwi kill .             # 현재 프로젝트 종료
```

### 3. 설정 관리 (SSOT)

모든 설정은 `config.env`에서 중앙 관리됩니다:

```bash
# 포트 설정
CODEKIWI_WEB_PORT_DEFAULT=8080      # 웹 인터페이스 포트 (호스트에 노출)
CODEKIWI_DEV_PORT_DEFAULT=3000      # 개발 서버 내부 포트 (컨테이너 내부 전용, 호스트에 노출 안 됨)
CODEKIWI_TTYD_PORT=7681             # 웹 터미널 내부 포트 (컨테이너 내부 전용, 호스트에 노출 안 됨)
CODEKIWI_NGINX_PORT=80              # Nginx 내부 포트 (컨테이너 내부 전용)

# 경로 설정
CODEKIWI_WORKSPACE_DIR=/workspace   # 컨테이너 내 작업 디렉토리
```

**포트 노출 정책:**
- **호스트에 노출**: `WEB_PORT` (8080) 하나만
- **내부 전용**: dev server(3000), ttyd(7681) - nginx를 통해서만 접근 가능

설정 로드 방식:
- **Go CLI**: `internal/config/config.go`에서 godotenv로 파싱
- **Runtime**: Docker Compose의 `env_file`로 환경 변수 주입
- **레거시 Bash CLI**: `lib/config-loader.sh`를 source (deprecated)

## 📦 빌드 및 배포

### 로컬 이미지 빌드

```bash
# 개발용 이미지 빌드
docker-compose -f docker-compose.dev.yaml build
```

### 프로덕션 배포

1. Docker Hub에 이미지 푸시:
```bash
docker build -t drasdp/codekiwi-runtime:latest ./runtime
docker push drasdp/codekiwi-runtime:latest
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
# Go CLI 빌드
cd cli-go
go build -o codekiwi cmd/codekiwi/main.go

# 다양한 시나리오 테스트
./codekiwi start                    # 현재 디렉토리 (포그라운드)
./codekiwi start -d ~/test-dir      # 특정 디렉토리 (백그라운드)
./codekiwi start -p 9000             # 커스텀 포트
./codekiwi list                      # 실행 중인 인스턴스
./codekiwi list -a                   # 모든 인스턴스
./codekiwi kill ~/test-dir           # 특정 프로젝트 종료
./codekiwi kill -a                   # 모든 인스턴스 종료

# 도움말
./codekiwi --help
./codekiwi start --help
```

## 📄 라이선스
- 상업적 사용 시 team@aardvark.co.kr 에 연락 후 협의. For commercial use, please contact team@aardvark.co.kr to discuss terms
