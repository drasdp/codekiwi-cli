# 🥝 CodeKiwi

웹 브라우저에서 사용하는 터미널 개발 환경

## 🚀 빠른 시작

```bash
docker-compose up -d --build
```

브라우저에서 http://localhost:8080 열기

**끝!** `~/CodeKiwi` 디렉토리가 자동으로 생성됩니다.

## 📁 프로젝트 사용법

### 1. 프로젝트 추가
```bash
# 호스트에서 프로젝트 폴더 추가
cp -r my-project ~/CodeKiwi/
```

### 2. 작업하기
- 좌측 터미널에서 `/workspace`로 이동
- opencode로 코딩 시작!

### 3. 양방향 실시간 동기화 ✅
```
호스트: ~/CodeKiwi  ←→  컨테이너: /workspace
```

## 🛠️ 구성

- **ttyd** - 웹 터미널
- **tmux** - 세션 유지
- **nginx** - 프록시
- **opencode-ai** - AI 코드 에디터
- **Docker** - 컨테이너화

## 💻 OS 호환

macOS · Linux · Windows (Docker Desktop)
