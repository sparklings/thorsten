package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// PageData는 템플릿에 전달할 데이터를 정의합니다.
type PageData struct {
	CurrentTime string
}

// HomeHandler는 메인 페이지 요청을 처리합니다.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// 템플릿 파일 경로 설정 (실행 위치 기준)
	// 실제 운영 환경에서는 절대 경로를 설정하거나 embed 패키지를 사용하는 것이 좋습니다.
	files := []string{
		filepath.Join("views", "layout.html"),
		filepath.Join("views", "home.html"),
	}

	// 템플릿 파싱
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("템플릿 파싱 오류: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		CurrentTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	// "layout" 템플릿 실행 (layout.html에서 define한 이름)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("템플릿 실행 오류: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}