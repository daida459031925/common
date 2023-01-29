package main

import (
	"fmt"
	"os"
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

	p, err := os.StartProcess("/home/sga/IdeaProjects/给戴达.docx", []string{"go"}, &os.ProcAttr{})
	if err != nil {
		t.Errorf("starting test process: %v", err)
	}
	fmt.Println(p.Pid)
	fmt.Println(p.Release())
	//p.Kill()
	p.Wait()

	//if got := p.Signal(Kill); got != ErrProcessDone {
	//	t.Errorf("got %v want %v", got, ErrProcessDone)
	//}
}
