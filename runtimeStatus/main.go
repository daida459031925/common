package runtimeStatus

import (
	"fmt"
	"runtime"
)

//获取运行时的状态

func GetCurrentStatus() string {
	return GetStatus(0)
}

func GetStatus(status int) string {
	start := ""
	if pc, file, line, ok := runtime.Caller(status); ok {
		f := runtime.FuncForPC(pc)
		//这里是打印错误，还可以进行报警处理，例如微信，邮箱通知
		//panic(err)//退出程序
		start = fmt.Sprintf("异常信息: ？？:%s、文件位置:%s、行数：%d、方法名称：%s、errpr信息：%s", pc, file, line, f.Name())
	}

	return start
}

func GetErrorStatus(status int, err any) string {
	return fmt.Sprintf("%s、errpr信息：%s", GetStatus(status), err)
}
