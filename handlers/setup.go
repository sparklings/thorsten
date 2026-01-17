package handlers

import (
	"log"
	"net/http"
	"thorsten/handlers/api"
	"thorsten/handlers/web"
	"time"
)

// NewHandler는 전체 애플리케이션의 라우팅 설정을 통합하여 반환합니다.
func NewHandler() http.Handler {
	mux := http.NewServeMux()

	// 1. 정적 파일 서빙 (/assets/...)
	// assets 로직은 웹 뷰와 관련 있으므로 여기서 처리하거나 web 패키지 내부에 위임할 수 있습니다.
	// 여기서는 전역 설정으로 둡니다.
	fs := http.FileServer(http.Dir("./views/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// 2. WEB 라우팅 등록
	// "/" 경로에 대한 처리는 web 패키지의 핸들러를 사용
	webMux := web.NewWebHandler()
	
	// 3. API 라우팅 등록 
	// /api/ 경로로 들어오는 모든 요청은 api 패키지로 위임
	apiMux := api.NewGatewayHandler()
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// 메인 라우팅 로직
	// 간단한 구조에서는 mux.Handle("/", webMux) 만으로도 충분하지만,
	// URL 패턴 충돌을 방지하기 위해 Wrapper를 사용할 수도 있습니다.
	// 여기서는 표준 mux 패턴을 따릅니다.
	// "/"는 모든 경로와 매칭되므로 가장 범용적인 webMux를 연결합니다.
	// 주의: /api/ 요청은 위에서 먼저 잡아채므로 여기로 오지 않습니다.
	mux.Handle("/", webMux)

	return mux
}

// LoggingMiddleware 등 공통 미들웨어가 있다면 여기에 위치할 수 있습니다.
// 현재는 main.go에 구현되어 있지만, handlers 패키지로 옮겨오는 것도 고려할 수 있습니다.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

// contextKey는 컨텍스트 키 타입을 정의하여 충돌을 방지합니다.
type contextKey string

const (
	// UserKey는 로그인 사용자 정보를 컨텍스트에 저장할 때 사용하는 키
	UserKey contextKey = "user"
)
