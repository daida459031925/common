package main

import (
	"fmt"
	"github.com/lxn/win"
	"github.com/shirou/gopsutil/process"
	"os"
	"syscall"
	"testing"
)

func TestOpenExeLocal(t *testing.T) {
	//argv := []string{}
	//attr := os.ProcAttr{}
	//r, e := os.StartProcess("chromium", argv, &attr)
	//r, e := os.FindProcess(4995)
	//r/*, e */:= os.newProcess(5200, 66666)
	//if e != nil {
	//	return
	//}
	//fmt.Println(r.Pid)
	//r.Kill()

	p, err := os.StartProcess("/home/sga/goland-2021.3.1/GoLand-2021.3.1/bin/goland.sh", []string{"go"}, &os.ProcAttr{})
	if err != nil {
		t.Errorf("starting test process: %v", err)
	}
	fmt.Println(p.Pid)
	fmt.Println(p.Release())
	fmt.Println(p)
	//p.Kill()
	p.Wait()

	//if got := p.Signal(Kill); got != ErrProcessDone {
	//	t.Errorf("got %v want %v", got, ErrProcessDone)
	//}
}

func main() {
	println(FindWindow(`QQ`))
}

func FindWindow(str string) win.HWND {
	return win.FindWindow(nil, syscall.StringToUTF16Ptr(str))
}

/*
获取所有进程id,以数组返回
*/
func ProcessId() (pid []int32) {
	pids,_ := process.Pids()
	for _,p := range pids {
		pid = append(pid,p)
	}
	return pid
}
/*
获取所有进程名，以数组返回
*/
func ProcessName() (pname []string) {

	pids,pid := range pids {
		pn,_ := process.NewProcess(pid)
		pName,_ :=pn.Name()
		pname = append(pname,pName)
	}
	return pname
}

func main() {
	pName := ProcessName()
	for _,v := range pName {
		fmt.Println("进程名:",v)
	}
}
