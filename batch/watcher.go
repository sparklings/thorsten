package batch

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RunWatcher는 파일 변경을 감지하고 서버를 재시작하는 관리자 프로세스입니다.
func RunWatcher() {
	log.Println("파일 감시 모드(Watch Mode)가 시작되었습니다...")

	for {
		log.Println(">> 애플리케이션 시작 중...")
		
		// 1. 실행 인자 구성 (--watch 제외)
		args := []string{"run", "main.go"}
		for _, arg := range os.Args[1:] {
			if !strings.Contains(arg, "watch") {
				args = append(args, arg)
			}
		}

		// 2. go run 명령 실행
		cmd := exec.Command("go", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Printf("애플리케이션 시작 실패: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		// 파일 변경 감지 채널
		done := make(chan error, 1)
		change := make(chan bool, 1)

		// 프로세스 종료 감시
		go func() {
			done <- cmd.Wait()
		}()

		// 파일 시스템 감시 시작
		stopWatch := make(chan bool)
		go watchFileSystem(".", stopWatch, change)

		select {
		case err := <-done:
			if err != nil {
				// 의도치 않은 종료 (예: 런타임 에러)
				log.Printf("애플리케이션이 에러와 함께 종료되었습니다: %v", err)
			} else {
				log.Println("애플리케이션이 정상 종료되었습니다.")
			}
			
			// 감시 고루틴 종료
			close(stopWatch)
			log.Println(">> 소스 코드 변경을 대기합니다...")
			
			// 3. 종료 후 변경 대기
			waitForChange()
			log.Println(">> 변경 감지됨. 재시작합니다.")

		case <-change:
			log.Println("\n>> 파일 변경 감지! 애플리케이션을 재시작합니다...")
			
			// 프로세스 트리 강제 종료 (Windows의 taskkill 활용)
			// go run은 자식 프로세스로 실제 앱을 실행하므로 트리 전체를 죽여야 함
			killProcessTree(cmd.Process.Pid)
			
			// 이전 프로세스가 완전히 종료될 때까지 대기
			<-done 
		}
	}
}

// killProcessTree는 지정된 PID와 그 자식 프로세스들을 강제 종료합니다.
func killProcessTree(pid int) {
	if runtime.GOOS == "windows" {
		// /F: 강제 종료, /T: 트리 종료(자식 포함)
		// strconv.Itoa를 사용하여 int를 string으로 변환
		exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid)).Run()
	} else {
		// Unix 계열 (Linux/macOS)
		proc, err := os.FindProcess(pid)
		if err == nil {
			proc.Kill()
		}
	}
}

// waitForChange는 파일 변경이 발생할 때까지 대기합니다.


// waitForChange는 파일 변경이 발생할 때까지 대기합니다.
func waitForChange() {
	change := make(chan bool, 1)
	go watchFileSystem(".", nil, change)
	<-change
}

// watchFileSystem은 지정된 경로의 파일 변경을 주기적으로 확인합니다.
func watchFileSystem(root string, stop <-chan bool, change chan<- bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastModTimes := make(map[string]time.Time)
	var mu sync.Mutex

	// 초기 상태 기록
	scanFiles(root, lastModTimes, &mu)

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if hasChanges(root, lastModTimes, &mu) {
				change <- true
				return
			}
		}
	}
}

// scanFiles는 디렉토리를 순회하며 마지막 수정 시간을 기록합니다.
func scanFiles(root string, modTimes map[string]time.Time, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// 감시 제외 대상
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "assets" || info.Name() == "tmp" || info.Name() == "batch") {
			return filepath.SkipDir
		}
		// 감시 대상 확장자
		ext := filepath.Ext(path)
		if ext == ".go" || ext == ".html" || ext == ".css" || ext == ".js" {
			modTimes[path] = info.ModTime()
		}
		return nil
	})
}

// hasChanges는 파일 변경 여부를 확인합니다.
func hasChanges(root string, lastModTimes map[string]time.Time, mu *sync.Mutex) bool {
	changed := false
	mu.Lock()
	defer mu.Unlock()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil 
		}
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "assets" || info.Name() == "tmp" || info.Name() == "batch") {
			return filepath.SkipDir
		}
		
		ext := filepath.Ext(path)
		if ext == ".go" || ext == ".html" || ext == ".css" || ext == ".js" {
			lastTime, exists := lastModTimes[path]
			if !exists || info.ModTime().After(lastTime) {
				changed = true
				return filepath.SkipDir // 변경 발견 시 즉시 중단
			}
		}
		return nil
	})

	return changed
}
