package api

import (
	"encoding/json"
	"net/http"
)

// NewGatewayHandler는 API 관련 라우팅을 설정하여 반환합니다.
func NewGatewayHandler() http.Handler {
	mux := http.NewServeMux()

	// API 핸들러 등록 (예시)
	// /api/health 처럼 호출됨 (상위 라우터에서 StripPrefix 여부에 따라 다름)
	mux.HandleFunc("/health", HealthCheckHandler)

	return mux
}

// HealthCheckHandler는 API 상태를 확인하는 샘플 핸들러입니다.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"message": "Thorsten API is running",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
