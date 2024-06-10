package windows

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"syscall"
)

var (
	user32   = syscall.MustLoadDLL("user32.dll")
	kernel32 = syscall.MustLoadDLL("kernel32.dll")

	findWindowW              = user32.MustFindProc("FindWindowW")
	showWindow               = user32.MustFindProc("ShowWindow")
	getWindowTextW           = user32.MustFindProc("GetWindowTextW")
	getWindowTextLengthW     = user32.MustFindProc("GetWindowTextLengthW")
	enumWindows              = user32.MustFindProc("EnumWindows")
	getWindowThreadProcessId = user32.MustFindProc("GetWindowThreadProcessId")
	sendMessage              = user32.MustFindProc("SendMessageW")
	sendInput                = user32.MustFindProc("SendInput")
	attachThreadInput        = user32.MustFindProc("AttachThreadInput")
	keybdEvent               = user32.MustFindProc("keybd_event")
	setForegroundWindow      = user32.MustFindProc("SetForegroundWindow")
	postMessage              = user32.MustFindProc("PostMessageW")
	setWindowsHookEx         = user32.MustFindProc("SetWindowsHookExW")
	unhookWindowsHookEx      = user32.MustFindProc("UnhookWindowsHookEx")
	callNextHookEx           = user32.MustFindProc("CallNextHookEx")
	getMessage               = user32.MustFindProc("GetMessageW")
	translateMessage         = user32.MustFindProc("TranslateMessage")
	dispatchMessage          = user32.MustFindProc("DispatchMessageW")

	getCurrentThreadId = kernel32.MustFindProc("GetCurrentThreadId")
	openProcess        = kernel32.MustFindProc("OpenProcess")
	closeHandle        = kernel32.MustFindProc("CloseHandle")
)

// GetProcess 根据进程id获取进程句柄
func GetProcess(i int) {
	processID := int32(i) // 替换成你要查询的进程ID
	//processID := int32(os.Getpid()) // 替换成你要查询的进程ID

	p, err := process.NewProcess(processID)
	if err != nil {
		fmt.Printf("无法获取进程：%s\n", err)
		return
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		fmt.Printf("无法获取内存信息：%s\n", err)
		return
	}

	fmt.Printf("进程ID: %d\n", processID)
	fmt.Printf("内存信息: %+v\n", memInfo)
}
