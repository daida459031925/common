package main

import (
	"fmt"
	"github.com/robotn/gohook"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	s := hook.Start()
	defer hook.End()

	tout := time.After(time.Hour * 2)
	done := false
	for !done {
		select {
		case i := <-s:
			if i.Kind >= hook.KeyDown && i.Kind <= hook.KeyUp {
				if i.Keychar == 'q' {
					tout = time.After(0)
				}

				fmt.Printf("%v key: %c:%v\n", i.Kind, i.Keychar, i.Rawcode)
			} else if i.Kind >= hook.MouseDown && i.Kind < hook.MouseWheel {
				//fmt.Printf("x: %v, y: %v, button: %v\n", i.X, i.Y, i.Button)
			} else if i.Kind == hook.MouseWheel {
				//fmt.Printf("x: %v, y: %v, button: %v, wheel: %v, rotation: %v\n", i.X, i.Y, i.Button, i.Amount, i.Rotation)
			} else {
				fmt.Printf("%+v\n", i)
			}

		case <-tout:
			fmt.Print("Done.")
			done = true
			break
		}
	}
	//className := "ss" // 替换为记事本窗口的类名
	//className := "1.txt - 记事本" // 替换为记事本窗口的类名
	//className := "LaTale Client" // 替换为记事本窗口的类名
	//
	//r, e := syscall.UTF16PtrFromString(className)
	//if e != nil {
	//	fmt.Println("找不到记事本窗口")
	//	return
	//}
	//// 使用FindWindow函数获取记事本窗口的句柄
	//parentHandle, _, _ := FindWindowW.Call(
	//	uintptr(0),
	//	uintptr(unsafe.Pointer(r)),
	//)
	//if parentHandle == 0 {
	//	fmt.Println("找不到记事本窗口")
	//	return
	//}
	//i := 0

	// 定义回调函数，用于处理每个子窗口
	//callback := syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
	//	// 在这里对子窗口进行处理，例如获取句柄、标题等信息
	//	// 这里只打印子窗口句柄和标题
	//	fmt.Printf("Child Handle: %X\n", hwnd)
	//	windowTitle := getWindowText(hwnd)
	//	fmt.Println("Title:", windowTitle)
	//
	//	// 获取子窗口的属性
	//	//style := getWindowLongPtr(hwnd, GWL_STYLE)
	//	//if style != 0 {
	//	//fmt.Printf("Style: %X\n", style)
	//
	//	// 检查属性是否符合文本编辑的条件
	//	//if isTextEditor(style) {
	//	//	fmt.Println("找到文本编辑子窗口")
	//	// 这里可以对文本编辑子窗口进行操作
	//
	//	// 模拟按键操作
	//	//sendKey(hwnd, VK_SHIFT, 'A')
	//	postKey(hwnd, 'S')
	//	postKey(hwnd, 'S')
	//	//postKey(hwnd, 'A')
	//	//sendKey(hwnd, VK_CONTROL, 'C')
	//	//postKey(hwnd, 'C')
	//	//postKey(hwnd, 'C')
	//	//}
	//	//}
	//
	//	return 1 // 返回1继续枚举下一个子窗口，返回0终止枚举
	//})

	//postKey(1184026, 's')
	//sendKey(1184026, 's')
	// 调用EnumChildWindows函数枚举子窗口
	//EnumChildWindows.Call(
	//	parentHandle,
	//	callback,
	//	0,
	//)

	//fmt.Println(i)
}

// getWindowText 获取窗口标题
func getWindowText(hwnd uintptr) string {
	const bufferSize = 256
	var buffer [bufferSize]uint16

	GetWindowText.Call(
		hwnd,
		uintptr(unsafe.Pointer(&buffer[0])),
		bufferSize,
	)

	return syscall.UTF16ToString(buffer[:])
}

// getWindowLongPtr 获取子窗口的属性
func getWindowLongPtr(hwnd uintptr, index int) int64 {
	ret, _, _ := GetWindowLongPtr.Call(
		hwnd,
		uintptr(index),
	)
	return int64(ret)
}

// isTextEditor 检查子窗口属性是否符合文本编辑条件
func isTextEditor(style int64) bool {
	// 判断是否具有可编辑属性的文本框
	editable := style&(ES_MULTILINE|ES_AUTOVSCROLL|ES_AUTOHSCROLL|ES_READONLY|ES_PASSWORD|ES_WANTRETURN) == ES_MULTILINE
	return editable
}

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// sendKey 模拟按键操作
func sendKey(hwnd uintptr, keys ...uintptr) {
	for _, key := range keys {
		SendMessageW.Call(
			hwnd,
			WM_KEYDOWN,
			key,
			1<<30, // 设置[原状态]已按下为1
		)
		SendMessageW.Call(
			hwnd,
			WM_KEYUP,
			key,
			1<<31, // 设置[原状态]已按下为1，[状态切换]为1
		)
	}
}

// postKey 模拟按键操作
func postKey(hwnd uintptr, keys ...uintptr) {
	for _, key := range keys {
		PostMessageW.Call(
			hwnd,
			WM_KEYDOWN,
			key,
			1<<30, // 设置[原状态]已按下为1
		)

		PostMessageW.Call(
			hwnd,
			WM_KEYUP,
			key,
			1<<31, // 设置[原状态]已按下为1，[状态切换]为1
		)
	}
}

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	GetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	SendInput                = user32.NewProc("SendInput")
	GetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	SetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	SetWindowsHookExW        = user32.NewProc("SetWindowsHookExW")
	CallNextHookEx           = user32.NewProc("CallNextHookEx")
	UnhookWindowsHookEx      = user32.NewProc("UnhookWindowsHookEx")
	GetModuleHandleW         = user32.NewProc("GetModuleHandleW")
	GetMessageW              = user32.NewProc("GetMessageW")
	PostMessageW             = user32.NewProc("PostMessageW")
	SendMessageW             = user32.NewProc("SendMessageW")
	FindWindowW              = user32.NewProc("FindWindowW")
	EnumChildWindows         = user32.NewProc("EnumChildWindows")
	GetWindowText            = user32.NewProc("GetWindowTextW")
	GetWindowLongPtr         = user32.NewProc("GetWindowLongPtrW")
)

const (
	WM_CHAR  = 0x0102
	VK_ENTER = 0x0D

	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2

	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002

	GWL_STYLE      = -16
	ES_MULTILINE   = 4
	ES_AUTOVSCROLL = 0x0020
	ES_AUTOHSCROLL = 0x0080
	ES_READONLY    = 0x0800
	ES_PASSWORD    = 0x0020
	ES_WANTRETURN  = 0x1000

	WM_KEYDOWN = 0x0100
	WM_KEYUP   = 0x0101
)

const (
	VK_LBUTTON  = 0x01 // 鼠标左键
	VK_RBUTTON  = 0x02 // 鼠标右键
	VK_CANCEL   = 0x03 // Ctrl+Break
	VK_MBUTTON  = 0x04 // 鼠标中键
	VK_BACK     = 0x08 // Backspace 键
	VK_TAB      = 0x09 // Tab 键
	VK_RETURN   = 0x0D // Enter 键
	VK_SHIFT    = 0x10 // Shift 键
	VK_CONTROL  = 0x11 // Ctrl 键
	VK_MENU     = 0x12 // Alt 键
	VK_PAUSE    = 0x13 // Pause/Break 键
	VK_CAPITAL  = 0x14 // Caps Lock 键
	VK_ESCAPE   = 0x1B // Esc 键
	VK_SPACE    = 0x20 // Spacebar 键
	VK_PRIOR    = 0x21 // Page Up 键
	VK_NEXT     = 0x22 // Page Down 键
	VK_END      = 0x23 // End 键
	VK_HOME     = 0x24 // Home 键
	VK_LEFT     = 0x25 // 左箭头键
	VK_UP       = 0x26 // 上箭头键
	VK_RIGHT    = 0x27 // 右箭头键
	VK_DOWN     = 0x28 // 下箭头键
	VK_SELECT   = 0x29 // Select 键
	VK_PRINT    = 0x2A // Print 键
	VK_EXECUTE  = 0x2B // Execute 键
	VK_SNAPSHOT = 0x2C // Print Screen 键
	VK_INSERT   = 0x2D // Insert 键
	VK_DELETE   = 0x2E // Delete 键
	VK_HELP     = 0x2F // Help 键
	VK_0        = 0x30 // '0' 键
	VK_1        = 0x31 // '1' 键
	VK_2        = 0x32 // '2' 键
	VK_3        = 0x33 // '3' 键
	VK_4        = 0x34 // '4' 键
	VK_5        = 0x35 // '5' 键
	VK_6        = 0x36 // '6' 键
	VK_7        = 0x37 // '7' 键
	VK_8        = 0x38 // '8' 键
	VK_9        = 0x39 // '9' 键
	VK_A        = 0x41 // 'A' 键
	VK_B        = 0x42 // 'B' 键
	VK_C        = 0x43 // 'C' 键
	VK_D        = 0x44 // 'D' 键
	VK_E        = 0x45 // 'E' 键
	VK_F        = 0x46 // 'F' 键
	VK_G        = 0x47 // 'G' 键
	VK_H        = 0x48 // 'H' 键
	VK_I        = 0x49 // 'I' 键
	VK_J        = 0x4A // 'J' 键
	VK_K        = 0x4B // 'K' 键
	VK_L        = 0x4C // 'L' 键
	VK_M        = 0x4D // 'M' 键
	VK_N        = 0x4E // 'N' 键
	VK_O        = 0x4F // 'O' 键
	VK_P        = 0x50 // 'P' 键
	VK_Q        = 0x51 // 'Q' 键
	VK_R        = 0x52 // 'R' 键
	VK_S        = 0x53 // 'S' 键
	VK_T        = 0x54 // 'T' 键
	VK_U        = 0x55 // 'U' 键
	VK_V        = 0x56 // 'V' 键
	VK_W        = 0x57 // 'W' 键
	VK_X        = 0x58 // 'X' 键
	VK_Y        = 0x59 // 'Y' 键
	VK_Z        = 0x5A // 'Z' 键
)

const (
	VK_NUMPAD0   = 0x60 // 小键盘 0
	VK_NUMPAD1   = 0x61 // 小键盘 1
	VK_NUMPAD2   = 0x62 // 小键盘 2
	VK_NUMPAD3   = 0x63 // 小键盘 3
	VK_NUMPAD4   = 0x64 // 小键盘 4
	VK_NUMPAD5   = 0x65 // 小键盘 5
	VK_NUMPAD6   = 0x66 // 小键盘 6
	VK_NUMPAD7   = 0x67 // 小键盘 7
	VK_NUMPAD8   = 0x68 // 小键盘 8
	VK_NUMPAD9   = 0x69 // 小键盘 9
	VK_MULTIPLY  = 0x6A // 小键盘 *
	VK_ADD       = 0x6B // 小键盘 +
	VK_SEPARATOR = 0x6C // 小键盘 Separator
	VK_SUBTRACT  = 0x6D // 小键盘 -
	VK_DECIMAL   = 0x6E // 小键盘 .
	VK_DIVIDE    = 0x6F // 小键盘 /
)

const (
	VK_RSHIFT              = 0xA1 // 右 Shift
	VK_RCONTROL            = 0xA3 // 右 Ctrl
	VK_RMENU               = 0xA5 // 右 Alt
	VK_BROWSER_BACK        = 0xA6 // 浏览器后退
	VK_BROWSER_FORWARD     = 0xA7 // 浏览器前进
	VK_BROWSER_REFRESH     = 0xA8 // 浏览器刷新
	VK_BROWSER_STOP        = 0xA9 // 浏览器停止
	VK_BROWSER_SEARCH      = 0xAA // 浏览器搜索
	VK_BROWSER_FAVORITES   = 0xAB // 浏览器收藏夹
	VK_BROWSER_HOME        = 0xAC // 浏览器主页
	VK_VOLUME_MUTE         = 0xAD // 音量静音
	VK_VOLUME_DOWN         = 0xAE // 音量减小
	VK_VOLUME_UP           = 0xAF // 音量增大
	VK_MEDIA_NEXT_TRACK    = 0xB0 // 下一曲
	VK_MEDIA_PREV_TRACK    = 0xB1 // 上一曲
	VK_MEDIA_STOP          = 0xB2 // 媒体停止
	VK_MEDIA_PLAY_PAUSE    = 0xB3 // 播放/暂停
	VK_LAUNCH_MAIL         = 0xB4 // 启动邮件
	VK_LAUNCH_MEDIA_SELECT = 0xB5 // 启动媒体选择
	VK_LAUNCH_APP1         = 0xB6 // 启动应用程序 1
	VK_LAUNCH_APP2         = 0xB7 // 启动应用程序 2
	VK_OEM_1               = 0xBA // OEM 1
	VK_OEM_PLUS            = 0xBB // OEM 加号
	VK_OEM_COMMA           = 0xBC // OEM 逗号
	VK_OEM_MINUS           = 0xBD // OEM 减号
	VK_OEM_PERIOD          = 0xBE // OEM 句点
	VK_OEM_2               = 0xBF // OEM 2
	VK_OEM_3               = 0xC0 // OEM 3
	VK_OEM_4               = 0xDB // OEM 4
	VK_OEM_5               = 0xDC // OEM 5
	VK_OEM_6               = 0xDD // OEM 6
	VK_OEM_7               = 0xDE // OEM 7
	VK_OEM_8               = 0xDF // OEM 8
)
