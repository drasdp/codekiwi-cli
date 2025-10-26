# 🥝 CodeKiwi

웹 브라우저 기반 터미널 개발 환경 - Docker + GUI 런처

## ✨ 특징

- 🖥️ **Split View UI**: 터미널 + 웹 프리뷰 동시 표시
- 🐳 **Docker 기반**: 격리된 개발 환경
- 🔧 **자동 개발 서버**: `npm install && npm run dev` 자동 실행
- 🌐 **실시간 프리뷰**: localhost:3000을 브라우저에서 바로 확인
- 🎨 **GUI 런처**: 디렉토리 선택만으로 간편하게 시작

## 📋 요구사항

- **Docker Desktop** (macOS, Linux, Windows)
- Python 3.9+ (개발 모드만 해당, 빌드된 실행 파일은 불필요)

## 🚀 빠른 시작 (런처 사용)

### 1. 런처 실행

#### 개발 모드 (Python 필요)
```bash
# 의존성 설치
uv sync

# 런처 실행
uv run python main.py
```

#### 빌드된 실행 파일 사용 (Python 불필요)
```bash
# macOS/Linux
./dist/CodeKiwi

# Windows
dist\CodeKiwi.exe
```

### 2. GUI에서 설정

1. **Browse 버튼** 클릭하여 작업할 프로젝트 디렉토리 선택
2. **▶ Start 버튼** 클릭
   - Docker 이미지 자동 빌드
   - 컨테이너 시작
   - 포트 충돌 자동 감지 및 처리
3. **🌐 Open Browser** 클릭하여 http://localhost:8080 열기

### 3. 개발 시작!

- **좌측 패널**: OpenCode AI 터미널
- **우상단 패널**: 웹 프리뷰 (localhost:3000)
- **우하단 패널**: DevServer 터미널 (자동으로 `npm install && npm run dev` 실행)

## 🔨 실행 파일 빌드

### macOS/Linux
```bash
./build.sh
```

### Windows
```batch
build.bat
```

빌드 결과물:
- macOS/Linux: `dist/CodeKiwi`
- Windows: `dist\CodeKiwi.exe`

## 🛠️ 구성

### 기술 스택
- **런처**: Python + tkinter + Docker SDK
- **컨테이너**: Ubuntu 22.04
- **도구**: ttyd, tmux, nginx, opencode-ai, claude-code, Node.js 20

### 포트
- `8080`: 메인 UI (nginx)
- `3000`: 개발 서버 (자동 프록시)
- `7681`: OpenCode 터미널 (내부)
- `7682`: DevServer 터미널 (내부)

### 볼륨 마운트
```
선택한_디렉토리  ←→  /workspace (컨테이너 내부)
```
양방향 실시간 동기화 ✅

## 📁 프로젝트 구조

```
codekiwi-web/
├── main.py              # GUI 런처 (tkinter + Docker SDK)
├── pyproject.toml       # uv 프로젝트 설정
├── build.sh             # macOS/Linux 빌드 스크립트
├── build.bat            # Windows 빌드 스크립트
├── .python-version      # Python 버전 명시
└── core/                # Docker 관련 파일들
    ├── Dockerfile       # 컨테이너 이미지 정의
    ├── docker-compose.yaml
    ├── nginx.conf       # 리버스 프록시 설정
    ├── index.html       # Split View UI
    ├── start.sh         # 컨테이너 시작 스크립트
    ├── .tmux.conf       # tmux 설정
    └── .bashrc          # bash 설정
```

## 🎯 주요 기능

### 자동 포트 충돌 처리
- 8080, 3000 포트 사용 중일 때 자동 감지
- 사용자 확인 후 기존 프로세스 종료

### 컨테이너 관리
- 컨테이너 이름: `codekiwi-{디렉토리명}`
- 상태 모니터링 (5초마다 갱신)
- 시작/중지/재시작 버튼
- 실시간 로그 출력

### 기존 컨테이너 처리
- 같은 이름 컨테이너 존재 시 경고
- 사용자 확인 후 기존 컨테이너 제거

## 💻 OS 호환

- ✅ macOS (Intel & Apple Silicon)
- ✅ Linux
- ✅ Windows (Docker Desktop 필요)

## 🐛 문제 해결

### "Docker error" 메시지가 나타날 때
- Docker Desktop이 설치되어 있는지 확인
- Docker Desktop이 실행 중인지 확인

### 포트 충돌 에러
- 8080 또는 3000 포트를 사용하는 프로세스 종료
- 런처가 자동으로 처리하도록 확인 다이얼로그에서 "Yes" 선택

### 빌드 실패
- `core/` 디렉토리에 모든 파일이 있는지 확인
- Docker Desktop이 실행 중인지 확인
- 로그에서 상세 에러 메시지 확인

## 📝 개발

### 의존성 추가
```bash
uv add package-name
```

### 개발 모드 실행
```bash
uv run python main.py
```

## 📄 라이선스

MIT License
