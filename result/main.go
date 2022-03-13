package result

import (
	"time"
)

type Result struct {
	Status int16  `json:"status"` //状态类型
	Msg    string `json:"msg"`    //错误时候的返回信息
	Data   any    `json:"data"`   //返还的数据
	Date   string `json:"date"`   //记录数据返回时间
}

func Build(status int16, msg string, data any) Result {
	now := time.Now()
	return Result{
		Status: status,
		Msg:    msg,
		Data:   data,
		Date:   now.Format("2006-01-02 15:04:05.999999999"),
	}
}

func Ok(data any) Result {
	return Build(200, "", data)
}

func Error(msg string) Result {
	return Build(500, msg, nil)
}
