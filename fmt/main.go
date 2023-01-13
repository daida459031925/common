package fmt

import (
	"fmt"
	err "github.com/daida459031925/common/error"
	"strconv"
)

// 控制台输出
func Println(a ...any) {
	fmt.Println(a...)
}

// 控制台输出
func Printlnf(s string, a ...any) {
	fmt.Printf(s, a...)
	fmt.Println()
}

//func FmtValue(value reflect.Value) string {
//	return fmt.Sprintf("%v", value)
//}
//
//func FmtValue(i int) string  {
//	return fmt.Sprintf("%d", i)
//}

// StringToInt 字符串转换成int
func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringToInt64 字符串转换成int64
func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// IntToString int转换成字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Int64ToString int转换成字符串
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// ParseBool 1、0、t、f、T、F、true、false、True、False、TRUE、FALSE；否则返回错误。
func ParseBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

// ParseInt base指定进制（2到36），如果base为0，则会从字符串前置判断，”0x”是16进制，”0”是8进制，否则是10进制；
// bitSize指定结果必须能无溢出赋值的整数类型，0、8、16、32、64 分别代表 int、int8、int16、int32、int64；
func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// ParseUnit 与ParseInt类似 只不过都是无符号正整数
func ParseUnit(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// ParseFloat bitSize指定了期望的接收类型，32是float32（返回值可以不改变精确值的赋值给float32），64是float64
func ParseFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// ParseUnPointer 使用条件，传入一个任意类型的值，返还指定泛型内容，并拆开指针返还真实值 返还的对象是一个新的对象
func ParseUnPointer[T any](data any) (T, error) {

	var e error

	r, ok := data.(T)
	if ok {
		return r, e
	}

	r1, ok1 := data.(*T)
	if ok1 {
		var aaa any = *r1
		r2, ok2 := aaa.(T)
		if ok2 {
			return r2, e
		}
	}
	var t T
	return t, err.New("解析失败")
}
