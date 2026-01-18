package web

import (
	"net/http"
)

// NewWebHandler는 웹 페이지 관련 라우팅을 설정하여 반환합니다.
func NewWebHandler() http.Handler {
	mux := http.NewServeMux()

	// 웹 페이지 핸들러 등록
	// 주의: 메인 핸들러("/")는 모든 경로를 잡아먹으므로 가장 마지막에 등록하거나 주의해야 함
	// 하지만 ServeMux는 가장 긴 패턴을 먼저 매칭하므로, 
	// /login, /about 등을 먼저 등록하고 / 는 마지막에 매칭됨.

	mux.HandleFunc("/", HomeHandler)
	
	// 추후 추가될 페이지들
	// mux.HandleFunc("/login", LoginHandler)
	// mux.HandleFunc("/about", AboutHandler)

	return mux
}
