# Thorsten 개발환경 구성 가이드

주간업무보고 애플리케이션 개발을 위한 환경 설정 및 사용 패키지 정보

## 목차
- [시스템 요구사항](#시스템-요구사항)
- [Go 설치](#go-설치)
- [프로젝트 설정](#프로젝트-설정)
- [사용 패키지](#사용-패키지)

## 시스템 요구사항
- Go 1.16 이상

## Go 설치

### Windows
1. Go 설치파일 다운로드 (64비트 x64 시스템용):
   ```bash
   curl -O https://go.dev/dl/go1.25.5.windows-amd64.msi
   ```
   > **참고**: Windows 10/11 x64의 경우 curl이 기본 제공됩니다. Git Bash 또는 WSL 사용도 가능합니다.

2. 설치파일 실행:
   ```bash
   go1.25.5.windows-amd64.msi
   ```

3. 설치 확인:
   ```bash
   go version
   ```

### Linux
1. Go 바이너리 다운로드:
   ```bash
   wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
   ```

2. 설치 및 환경 변수 설정:
   ```bash
   tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

3. 설치 확인:
   ```bash
   go version
   ```

## 프로젝트 설정

1. 프로젝트 디렉토리 생성:
   ```bash
   mkdir thorsten
   cd thorsten
   ```

2. Go 모듈 초기화:
   ```bash
   go mod init github.com/sparklings/thorsten
   ```

3. 필요한 패키지 설치:
   ```bash
   go get <package-name>@<version>
   ```

## 사용 패키지

### 표준 라이브러리
- `html/template` - HTML 안전 렌더링

### 데이터베이스
- `github.com/mattn/go-sqlite3` v1.14.14 - SQLite 데이터베이스 드라이버

### 웹 프레임워크 및 미들웨어
- `github.com/gin-gonic/gin` v1.7.4 - 웹 프레임워크 (라우팅 및 미들웨어)

### 인증 및 세션
- `github.com/dgrijalva/jwt-go` v3.2.0 - JWT 생성 및 검증
- `github.com/gorilla/sessions` v1.2.0 - 세션 관리

### HTTP 클라이언트
- `github.com/go-resty/resty/v2` v2.6.0 - RESTful API 클라이언트

### 테스트
- `github.com/stretchr/testify` v1.7.0 - 테스트 유틸리티
