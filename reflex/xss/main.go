package xss

import (
	"fmt"
	"github.com/daida459031925/common/reflex"
	"html"
	"reflect"
)

/**
主要解决项目中xss在对象中字符串问题 xss转义
*/

type (
	xssUtil struct {
		Data    any
		XssType reflex.XssType
		TagType reflex.TagType
		Filter  []string
	}
)

func XssUtilInit(xssType reflex.XssType, tagType reflex.TagType, data any, filter ...string) xssUtil {
	return xssUtil{Data: data, XssType: xssType, TagType: tagType, Filter: filter}
}

func (x xssUtil) xssCode() any {
	//获取对象所有内容
	//t := reflect.TypeOf(data)
	ref := reflex.GetRef(x.Data)
	//目的是将指针类型的转换成实在的数据
	v := ref.GetPointerData(false)

	fieldStruct(v, x.XssType, x.TagType, x.Filter)
	fieldSlice(v, x.XssType, x.TagType, x.Filter)
	switch x.XssType {
	case reflex.XssEscape:
		escapeString(v, "", nil)
	case reflex.XssUnEscape:
		unEscapeString(v, "", nil)
	}
	return x.Data
}

func fieldStruct(value reflect.Value, xssType reflex.XssType, tagType reflex.TagType, filter []string) {
	if reflex.IsStruct(value) {
		for j := 0; j < value.NumField(); j++ {
			v := value.Field(j)
			//获取当前对象名称key
			tagName := fmt.Sprintf("%s", tagType)
			name := value.Type().Field(j).Tag.Get(tagName)
			switch xssType {
			case reflex.XssEscape:
				escapeString(v, name, filter)
			case reflex.XssUnEscape:
				unEscapeString(v, name, filter)
			}
			fieldSlice(v, xssType, tagType, filter)
			fieldStruct(v, xssType, tagType, filter)
		}
	}
}

func fieldSlice(value reflect.Value, xssType reflex.XssType, tagType reflex.TagType, filter []string) {
	if reflex.IsSlice(value) {
		for i := 0; i < value.Len(); i++ {
			v := value.Index(i)
			fieldStruct(v, xssType, tagType, filter)
		}
	}
}

// 编码 " => &#34;
func escapeString(value reflect.Value, name string, filter []string) {
	if reflex.IsString(value) {
		if filter != nil && len(filter) > 0 {
			for i := range filter {
				if filter[i] == name {
					return
				}
			}
		}
		value.SetString(html.EscapeString(value.String()))
	}
}

// 解码 &#34; => "
func unEscapeString(value reflect.Value, name string, filter []string) {
	if reflex.IsString(value) {
		if filter != nil && len(filter) > 0 {
			for i := range filter {
				if filter[i] == name {
					return
				}
			}
		}
		value.SetString(html.UnescapeString(value.String()))
	}
}
