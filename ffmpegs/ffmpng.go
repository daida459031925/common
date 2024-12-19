package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

//生成的sh或者其他linux下的内容记得需要使用chmod +x 文件   给上执行权限
//https://github.com/datarhei/restreamer.git    灵感来自于

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
			ffmpegProcess.Stop(params.InputURL)
		}
	})
	return ffmpegProcess
}

// Start 启动FFmpeg命令
func (fp *FFmpegProcess) Start(params FFmpegParams) {
	//cmd := exec.Command(params.Name,
	//	"-i", params.InputURL,
	//	"-c:v", "libx264",
	//	"-c:a", "aac",
	//	"-f", "hls",
	//	"-hls_time", fmt.Sprintf("%d", params.HlsTime),
	//	"-hls_list_size", fmt.Sprintf("%d", params.HlsListSize),
	//	"-hls_flags", fmt.Sprintf("%d", params.HlsFlags),
	//	"-y", params.OutputFile,
	//)

	cmd := exec.Command(params.Name,
		"-loglevel", "level+info",
		"-err_detect", "ignore_err",
		"-y", "-fflags", "+genpts",
		"-thread_queue_size", "512",
		"-probesize", "5000000",
		"-analyzeduration", "5000000",
		"-timeout", "5000000",
		"-rtsp_transport", "tcp",
		"-i", params.InputURL,
		"-dn", "-sn",
		"-map", "0:0",
		"-codec:v", "copy",
		"-an", "-metadata",
		"title=http://192.168.20.172:48080/e0c9e586-4b61-4f4d-9e4e-3bd015615475/oembed.json",
		"-metadata", "service_provider=datarhei-Restreamer",
		"-f", "hls",
		"-start_number", "0",
		"-hls_time", "2",
		"-hls_list_size", "6",
		"-hls_flags", "append_list+delete_segments+program_date_time+temp_file",
		"-hls_delete_threshold", "4",
		"-master_pl_publish_rate", "2",
		params.OutputFile,
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
func (fp *FFmpegProcess) Stop(url string) {
	if fp.cmd != nil && fp.cmd.Process != nil {
		err := fp.cmd.Process.Kill()
		if err != nil {
			fmt.Printf("Failed to kill FFmpeg process: %v\n", err)
			return
		}
	}
	ffmpegProcessMap.Delete(url)
	fp.active = false
}

// ResetTimer 重置定时器，每次有新请求时调用该方法来重新计时60秒
func (fp *FFmpegProcess) ResetTimer() {
	if fp.timer != nil {
		fp.timer.Reset(60 * time.Second)
	}
}

func handleFFmpegRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var params FFmpegParams
	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if params.Name == "" || params.InputURL == "" || params.OutputFile == "" {
		params = DefaultParams()
	}

	fmt.Println(params.Name)
	fmt.Println(params.InputURL)
	fmt.Println(params.OutputFile)

	// 获取文件所在目录路径
	dir := filepath.Dir(params.OutputFile)
	// 创建目录（如果不存在）
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 创建文件（如果不存在）
	_, err = os.OpenFile(params.OutputFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.Chmod(params.OutputFile, 0666) // 设置为所有者和其他用户都有读写权限
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var ffmpegProcess *FFmpegProcess

	// 检查键 "arg1" 是否存在于 ffmpegProcessMap 中
	value, ok := ffmpegProcessMap.Load(params.InputURL)
	if ok {
		fmt.Printf("key is the ffmpegProcessMap ，to: %v\n", value)
		ffmpegProcess = value.(*FFmpegProcess)
	} else {
		ffmpegProcess = NewFFmpegProcess(params)
		ffmpegProcessMap.Store(params.InputURL, ffmpegProcess)
	}

	ffmpegProcess.ResetTimer()
	_, err = w.Write([]byte("FFmpeg process triggered or reset timer."))
	if err != nil {
		return
	}
}

// ffmpegProcessMap用于保存正在执行的FFmpeg请求相关信息，键是请求标识，值是对应的FFmpegProcess实例
var ffmpegProcessMap sync.Map

func main() {
	http.HandleFunc("/trigger-ffmpeg", func(w http.ResponseWriter, r *http.Request) {
		handleFFmpegRequest(w, r)
	})

	fmt.Println("Server started on :38080")
	err := http.ListenAndServe(":38080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
