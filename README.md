# 🥝 CodeKiwi

웹 브라우저에서 사용하는 AI 기반 개발 환경

Docker 기반의 통합 개발 환경으로, 좌측에는 AI 코드 에디터, 우측에는 실시간 웹 미리보기를 제공합니다.

## ✨ 특징

- 🚀 **원라이너 설치** - 한 줄로 설치 완료
- 💻 **opencode 스타일** - `cd my-project && codekiwi`로 즉시 실행
- 🤖 **AI 코드 에디터** - OpenCode AI 통합
- 👀 **실시간 미리보기** - 코드 변경 즉시 반영
- 🔢 **다중 인스턴스** - 여러 프로젝트 동시 실행 (자동 포트 할당)
- 🔄 **스마트 관리** - 프로젝트별 독립적인 컨테이너 및 포트
- 🌐 **크로스 플랫폼** - macOS, Linux, Windows 지원

## 🚀 빠른 시작

### 설치 (한 번만)

```bash
curl -fsSL https://raw.githubusercontent.com/aardvarkdev1/codekiwi-cli/main/cli/scripts/install.sh | bash
```

### 사용

```bash
# 프로젝트 디렉토리로 이동
cd ~/my-project

# CodeKiwi 실행 (현재 디렉토리 자동 마운트)
codekiwi

# 브라우저에서 http://localhost:8080 접속
```

**끝!** 😎

## 📖 사용법

### 기본 사용

```bash
# 현재 디렉토리로 시작 (포그라운드 실행, Ctrl+C로 종료)
cd ~/my-react-app
codekiwi

# 특정 디렉토리 지정
codekiwi ~/my-vue-app

# 다른 터미널에서 강제 종료
codekiwi kill ~/my-react-app

# 모든 인스턴스 목록 확인
codekiwi list
```

### 다중 프로젝트 동시 실행

```bash
# 터미널 1
cd ~/projects/react-app
codekiwi                          # 포트 8080, 3000 (포그라운드 실행)

# 터미널 2
cd ~/projects/vue-app
codekiwi                          # 포트 8081, 3001 (포그라운드 실행)

# 터미널 3 - 모든 인스턴스 확인
codekiwi list
# [인스턴스 #1]
#   📁 디렉토리: ~/projects/react-app
#   🌐 웹: http://localhost:8080
#   🔌 개발 서버: http://localhost:3000
#
# [인스턴스 #2]
#   📁 디렉토리: ~/projects/vue-app
#   🌐 웹: http://localhost:8081
#   🔌 개발 서버: http://localhost:3001

# 특정 프로젝트 강제 종료 (다른 터미널에서)
codekiwi kill ~/projects/react-app

# 모든 프로젝트 강제 종료
codekiwi kill --all
```

### 고급 명령어

```bash
# 모든 실행 중인 인스턴스 목록
codekiwi list

# 특정 디렉토리 강제 종료
codekiwi kill ~/my-project

# 최신 버전으로 업데이트
codekiwi update

# 완전 제거
codekiwi uninstall

# 도움말
codekiwi help
```

## 🎨 화면 구성

```
┌─────────────────────┬─────────────────────┐
│     좌측 50%        │     우측 50%        │
│                     │                     │
│  AI 코드 에디터     │  실시간 미리보기    │
│  (OpenCode AI)      │  (포트 3000)        │
│                     │                     │
│  - 코드 편집        │  - npm run dev      │
│  - AI 지원          │  - 자동 새로고침    │
│  - 터미널 접근      │  - 즉시 반영        │
│                     │                     │
└─────────────────────┴─────────────────────┘
```

**접속 주소:** http://localhost:8080

## 🛠️ 구성

- **ttyd** - 웹 터미널 에뮬레이터
- **tmux** - 터미널 세션 관리
- **nginx** - HTTP 역방향 프록시
- **opencode-ai** - AI 기반 코드 에디터
- **Node.js 20.x** - JavaScript 런타임
- **Docker** - 컨테이너 환경

## 💡 워크플로우 예시

### 일반적인 개발 흐름

```bash
# 1. 프로젝트로 이동
cd ~/my-awesome-app

# 2. CodeKiwi 시작 (포그라운드 실행)
codekiwi

# 3. 브라우저에서 http://localhost:8080 열기

# 4. 좌측 터미널에서 코드 편집
# 5. 우측에서 실시간으로 결과 확인

# 6. 작업 완료 후 Ctrl+C로 종료
```

### 여러 프로젝트 작업

```bash
# 오전: React 프로젝트
cd ~/work/client-dashboard
codekiwi

# 오후: API 서버 작업
codekiwi ~/work/api-server

# 저녁: 개인 프로젝트
codekiwi ~/personal/blog
```

## 📁 디렉토리 마운트

CodeKiwi는 지정한 디렉토리를 컨테이너의 `/workspace`에 마운트합니다.

- **자동 마운트**: 인자 없이 실행하면 현재 디렉토리(`$PWD`) 사용
- **명시적 지정**: `codekiwi /path/to/project`로 특정 경로 지정
- **실시간 동기화**: 호스트와 컨테이너 간 파일 변경사항 즉시 반영
- **자유로운 전환**: 언제든 다른 디렉토리로 전환 가능

## 📌 포트 사용

- **8080** - 메인 웹 인터페이스
- **3000** - 개발 서버 (npm run dev)

컨테이너 실행 중에는 해당 포트를 사용할 수 없습니다.

## ⚙️ 시스템 요구사항

- **Docker** 20.10 이상
- **Docker Compose** 2.0 이상
- **메모리** 최소 2GB 권장
- **디스크** 약 1GB (이미지 + 의존성)

### 지원 플랫폼

✅ macOS (Intel/Apple Silicon)
✅ Linux (Ubuntu, Debian, Fedora 등)
✅ Windows 10/11 (Docker Desktop 필요)

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

CodeKiwi는 자동으로 사용 가능한 포트를 찾습니다:
- 웹 포트: 8080부터 시작하여 사용 가능한 포트 자동 할당
- 개발 서버 포트: 3000부터 시작하여 자동 할당

수동으로 포트 확인이 필요한 경우:
```bash
# 8080 포트를 사용 중인 프로세스 확인
lsof -i :8080

# 실행 중인 모든 CodeKiwi 인스턴스 확인
codekiwi list
```

### 완전히 초기화

```bash
# 모든 인스턴스 강제 종료
codekiwi kill --all

# Docker 정리
docker system prune -a

# 제거 및 재설치
codekiwi uninstall
curl -fsSL https://raw.githubusercontent.com/aardvarkdev1/codekiwi-cli/main/cli/scripts/install.sh | bash
```

## 📚 명령어 레퍼런스

| 명령어 | 설명 |
|--------|------|
| `codekiwi` | 현재 디렉토리로 시작 (포그라운드 실행) |
| `codekiwi <path>` | 지정한 디렉토리로 시작 |
| `codekiwi list` | 실행 중인 모든 인스턴스 목록 표시 |
| `codekiwi kill [path]` | 지정한 디렉토리의 인스턴스 강제 종료 |
| `codekiwi kill --all` | 모든 인스턴스 강제 종료 |
| `codekiwi update` | 최신 버전으로 업데이트 |
| `codekiwi uninstall` | 완전히 제거 |
| `codekiwi help` | 도움말 표시 |

## 🤝 기여

이슈와 PR은 언제나 환영합니다!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 라이선스

MIT License

## 🙏 감사

- [ttyd](https://github.com/tsl0922/ttyd) - Web terminal emulator
- [OpenCode AI](https://github.com/opencodeiiit/opencode-ai) - AI code editor
- [tmux](https://github.com/tmux/tmux) - Terminal multiplexer
- [nginx](https://nginx.org/) - Web server

---

**Made with ❤️ for developers**
