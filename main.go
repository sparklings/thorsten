package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"thorsten/batch"
	"thorsten/handlers"
)

func main() {
	// 플래그 정의
	watch := flag.Bool("watch", false, "소스 코드 변경 감지 및 자동 재시작")
	debug := flag.Bool("debug", false, "디버그 모드 활성화 (시스템 로그 표시)")
	addr := flag.String("addr", "127.0.0.1", "서버 IP 주소 (기본값: 127.0.0.1)")
	port := flag.Int("port", 8080, "서버 포트 번호 (기본값: 8080)")

	flag.Parse()

	// Watch 모드일 경우 감시자 로직 실행
	if *watch {
		batch.RunWatcher()
		return
	}

	// 라우터 초기화 (handlers 패키지에서 통합 관리)
	mux := handlers.NewHandler()

	// 핸들러 설정 (디버그 모드일 경우 로깅 미들웨어 추가)
	var handler http.Handler = mux
	if *debug {
		handler = handlers.LoggingMiddleware(mux)
	}

	serverAddr := fmt.Sprintf("%s:%d", *addr, *port)
	log.Printf("서버가 시작되었습니다. http://%s", serverAddr)

	// 서버 시작
	err := http.ListenAndServe(serverAddr, handler)
	if err != nil {
		log.Fatal("서버 시작 실패: ", err)
	}
}