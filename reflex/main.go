package reflex

import "reflect"

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
)

//深度对比两个对象是否相等
func Equals(type1 any, type2 any) bool {
	return reflect.DeepEqual(type1, type2)
}

func Equal(type1 any, type2 any) bool {
	return type1 == type2
}
