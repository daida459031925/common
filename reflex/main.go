package reflex

import (
	err "github.com/daida459031925/common/error"
	"reflect"
)

//1. Reflection goes from interface value to reflection Object.
//2. Reflection goes from refelction object to interface value.
//3. To modify a reflection object, the value must be settable.

const (
	XssEscape   XssType = 0
	XssUnEscape XssType = 1

	Json TagType = "json"
	Form TagType = "form"
	Db   TagType = "db"
)

type (
	XssType int

	TagType string

	//添加反射内容结构体
	ref struct {
		item any
	}
)

// Equals 深度对比两个对象是否相等
func Equals(type1 any, type2 any) bool {
	return reflect.DeepEqual(type1, type2)
}

func Equal(type1 any, type2 any) bool {
	return type1 == type2
}

// GetRef 获取需要处理的对象，提供一系列方法
func GetRef(data any) ref {
	return ref{
		item: data,
	}
}

// IsPointer 判断当前是否是指针类型 如果true:是 false:不是
func IsPointer(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr
}

// IsStruct 判断当前是否是Struct 如果true:是 false:不是
func IsStruct(value reflect.Value) bool {
	return value.Kind() == reflect.Struct
}

// IsSlice 判断当前是否是Slice 如果true:是 false:不是
func IsSlice(value reflect.Value) bool {
	return value.Kind() == reflect.Slice
}

// IsString 判断当前是否是Slice 如果true:是 false:不是
func IsString(value reflect.Value) bool {
	return value.Kind() == reflect.String
}

// IsFunc 判断当前是否是Func 如果true:是 false:不是
func IsFunc(value reflect.Value) bool {
	return value.Kind() == reflect.Func
}

// GetPointerData 是否拿回指针类型数据 如果true:是 false:不是
func (r ref) GetPointerData(tf bool) reflect.Value {
	value := reflect.ValueOf(r.item)
	if tf {
		if IsPointer(value) {
			return value
		}
		return reflect.ValueOf(&r.item)
	}

	if IsPointer(value) {
		return value.Elem()
	}
	return value
}

// 执行结构体中的方法
func (r ref) executeMethod(funcName string, datas ...any) ([]reflect.Value, error) {
	value := r.GetPointerData(true)
	//通过反射获取它对应的函数，然后通过call来调用
	f := value.MethodByName(funcName)
	if !IsFunc(f) {
		return nil, err.New("不是方法不能执行")
	}

	in := make([]reflect.Value, 0)
	if datas != nil && len(datas) > 0 {
		for i := range datas {
			in = append(in, reflect.ValueOf(datas[i]))
		}
	}
	return f.Call(in), nil
}

func isBasicField(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}
