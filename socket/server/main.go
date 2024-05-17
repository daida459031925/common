package server

import (
	"flag"
	"github.com/daida459031925/common/fmt"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/net/gtcp"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	startFlag bool
	stopFlag  bool
	pidFile   = "pid.txt" // 存放 PID 的文件路径
)

func init() {
	flag.BoolVar(&startFlag, "start", false, "启动程序")
	flag.BoolVar(&stopFlag, "stop", false, "关闭程序")
	flag.Parse()
}

func Server() {
	if startFlag {
		start()
	} else if stopFlag {
		stop()
	} else {
		fmt.Println("Usage: program -start to start, program -stop to stop")
	}
}

func start() {
	fmt.Println("程序已启动，按 Ctrl+C 或发送 kill 命令退出")

	// 写入当前进程的 PID 到文件   0644是权限代码
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		fmt.Println("无法写入 PID 文件:", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("收到退出信号，正在关闭程序...")
		stop()
	}()

	// 模拟程序运行
	select {}
}

func stop() {
	// 读取存放 PID 的文件
	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("无法读取 PID 文件:", err)
		os.Exit(1)
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		fmt.Println("PID 文件内容不合法:", err)
		os.Exit(1)
	}

	// 根据 PID 杀死进程
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("无法找到进程:", err)
		os.Exit(1)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		fmt.Println("无法关闭进程:", err)
		os.Exit(1)
	}

	// 删除 PID 文件
	if err := os.Remove(pidFile); err != nil {
		fmt.Println("无法删除 PID 文件:", err)
		os.Exit(1)
	}

	fmt.Println("程序已关闭")
	os.Exit(0)
}

// StartServer windows 后台启动 start /b server.exe
// linux 后台启动 nohup ./server.exe &
func StartServer(ip string, port int, set *gset.Set, handlers map[string]*Handler) *gtcp.Server {
	server := gtcp.NewServer(fmt.Sprintf("%s:%d", ip, port), func(conn *gtcp.Conn) {
		defer func(conn *gtcp.Conn) {
			err := conn.Close()
			if err != nil {
				// 处理错误
				fmt.Println(err)
			}
		}(conn)
		for {
			//这个地方估计不能那么写，或者说不能单独这么写，需要根据实际情况调整
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				fmt.Println(string(data))
				set.Add(conn)

				// 假设业务名称在数据中的某个位置,目前设定前30个字符串默认是这个
				businessName := string(data[0:30])
				handler, ok := handlers[businessName]
				if !ok {
					fmt.Sprintf("Unsupported business:%s", businessName)
					continue
				}
				// 调用相应的处理程序
				g := *handler
				g.Handle(conn)
			}
			if err != nil {
				set.Remove(conn)
				break
			}
		}
	})

	go func() {
		err := server.Run()
		if err != nil {
			fmt.Sprintf("Server stopped with error: %s\n", err)
		}
	}()

	return server
}

// AddHandler 添加需要处理的handler任务
func AddHandler(handlers *map[string]*Handler, name string, handler *Handler) *map[string]*Handler {
	h := *handlers
	h[name] = handler
	return handlers
}

// RemoveHandler 删除某个任务
func RemoveHandler(handlers *map[string]*Handler, name string) *map[string]*Handler {
	delete(*handlers, name)
	return handlers
}

// Handler 定义一个结构用于其他调用我们这边的进行实现
type Handler interface {
	Handle(conn *gtcp.Conn)
}
