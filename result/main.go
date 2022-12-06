package result

import (
	"time"
)

type Result struct {
	Status    int16        `json:"status"` //状态类型
	Msg       string       `json:"msg"`    //错误时候的返回信息
	Data      any          `json:"data"`   //返还的数据
	Date      string       `json:"date"`   //记录数据返回时间
	funcSlice []func() any //记录当前需要执行的所有任务
}

const (
	OK  = 200
	ERR = 500

	// LAYOUT 格式化时间默认值 相当于java 1970年01月01日
	LAYOUT = "2006-01-02 15:04:05.999999999"
)

func Build(status int16, msg string, data any) Result {
	now := time.Now()
	return Result{
		Status: status,
		Msg:    msg,
		Data:   data,
		Date:   now.Format(LAYOUT),
	}
}

// Ok 默认添加正常
func Ok() Result {
	return Build(OK, "", nil)
}

// OkData 默认添加正常添加数据
func OkData(data any) Result {
	return Build(OK, "", data)
}

// OkDataMsg 通过并返还数据以及信息
func OkDataMsg(data any, msg string) Result {
	return Build(OK, msg, data)
}

// Error 默认错误
func Error() Result {
	return Build(ERR, "", nil)
}

// ErrorMsg 默认错误
func ErrorMsg(msg string) Result {
	return Build(ERR, msg, nil)
}

// ErrorData 默认错误并返还数据
func ErrorData(msg string, data any) Result {
	return Build(ERR, msg, data)
}

// SetFunc 添加方法
func (result *Result) SetFunc(f func(a any) any) *Result {
	return result.SetFuncErr(f, nil)
}

// SetFuncErr 添加方法并设置错误信息
func (result *Result) SetFuncErr(f func(a any) any, e error) *Result {
	fun := func() any {
		defer func() {
			r := recover()
			if r != nil {
				result.Status = ERR
				result.Data = nil
				result.Msg = r.(error).Error()
				if e != nil {
					result.Msg = e.Error()
				}
			}
		}()
		return f(result.Data)
	}
	result.funcSlice = append(result.funcSlice, fun)
	return result
}

// Exec 执行方法
func (result *Result) Exec() Result {
	f := result.funcSlice
	if f != nil {
		for i := range f {
			if result.Status != ERR {
				result.Data = f[i]()
			}
		}
	}
	return *result
}
