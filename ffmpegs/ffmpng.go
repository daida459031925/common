package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

//生成的sh或者其他linux下的内容记得需要使用chmod +x 文件   给上执行权限

// FFmpegParams 包含FFmpeg命令的参数
type FFmpegParams struct {
	Name        string
	InputURL    string
	HlsTime     int
	HlsListSize int
	HlsFlags    int
	OutputFile  string
}

// DefaultParams 返回默认的FFmpeg参数
func DefaultParams() FFmpegParams {
	return FFmpegParams{
		Name:        "./ffmpeg",
		InputURL:    "rtsp://admin:hH87793891@@192.168.20.125:554/Streaming/Channels/101",
		HlsTime:     10,
		HlsListSize: 6,
		HlsFlags:    10,
		OutputFile:  "output.m3u8",
	}
}

// FFmpegProcess 结构体用于管理FFmpeg进程以及相关的定时器等信息
type FFmpegProcess struct {
	cmd    *exec.Cmd
	timer  *time.Timer
	active bool
}

// NewFFmpegProcess 创建一个新的FFmpegProcess实例并启动初始定时器
func NewFFmpegProcess(params FFmpegParams) *FFmpegProcess {
	ffmpegProcess := &FFmpegProcess{
		active: false,
	}
	ffmpegProcess.Start(params)
	ffmpegProcess.timer = time.AfterFunc(60*time.Second, func() {
		if ffmpegProcess.active {
			ffmpegProcess.Stop()
		}
	})
	return ffmpegProcess
}

// Start 启动FFmpeg命令
func (fp *FFmpegProcess) Start(params FFmpegParams) {
	cmd := exec.Command(params.Name,
		"-i", params.InputURL,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", params.HlsTime),
		"-hls_list_size", fmt.Sprintf("%d", params.HlsListSize),
		"-hls_flags", fmt.Sprintf("%d", params.HlsFlags),
		"-y", params.OutputFile,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start FFmpeg command: %v\n", err)
		return
	}
	fp.cmd = cmd
	fp.active = true
}

// Stop 停止FFmpeg进程
func (fp *FFmpegProcess) Stop() {
	if fp.cmd != nil && fp.cmd.Process != nil {
		err := fp.cmd.Process.Kill()
		if err != nil {
			fmt.Printf("Failed to kill FFmpeg process: %v\n", err)
			return
		}
	}
	fp.active = false
}

// ResetTimer 重置定时器，每次有新请求时调用该方法来重新计时60秒
func (fp *FFmpegProcess) ResetTimer() {
	if fp.timer != nil {
		fp.timer.Reset(60 * time.Second)
	}
}

func handleFFmpegRequest(w http.ResponseWriter, r *http.Request, ffmpegProcess *FFmpegProcess) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if ffmpegProcess == nil {
		params := DefaultParams()
		ffmpegProcess = NewFFmpegProcess(params)
	}

	ffmpegProcess.ResetTimer()
	w.Write([]byte("FFmpeg process triggered or reset timer."))
}

func main() {
	var ffmpegProcess *FFmpegProcess
	http.HandleFunc("/trigger-ffmpeg", func(w http.ResponseWriter, r *http.Request) {
		handleFFmpegRequest(w, r, ffmpegProcess)
	})

	fmt.Println("Server started on :38080")
	err := http.ListenAndServe(":38080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
