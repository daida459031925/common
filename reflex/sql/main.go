package sql

import (
	"fmt"
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/reflex"
	"reflect"
	"strconv"
)

/**
解决问题：手动拼接sql字符串 中的 key value 字符串问题
*/

// RawField 返还struce 类型所有0:key 1:value
func RawField(in any, index int) ([]string, error) {

	//v := reflect.ValueOf(in)
	ref := reflex.GetRef(in)
	//目的是将指针类型的转换成实在的数据
	v := ref.GetPointerData(false)

	//如果make 切片 类型
	if reflex.IsSlice(v) {
		return getField(v, index), nil
		//panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	//如果时结构体
	if reflex.IsStruct(v) {
		//panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
		return getNumField(v, index), nil
	}

	//fmt.Sprintf("field %d, type %s, key %s ,value %v",
	//				i, t.Field(i).Type, t.Field(i).Name, v.FieldByName(t.Field(i).Name))
	return nil, err.New("无法获取")
}

// 返还slice 类型所有values
func getField(in reflect.Value, index int) []string {
	out := make([]string, 0)
	for i := 0; i < in.Len(); i++ {
		value := in.Index(i)
		out = append(getNumField(value, index))
	}
	return out
}

// 返还struce 类型所有0:key 1:value
func getNumField(value reflect.Value, index int) []string {
	out := make([]string, 0)
	//找到这个对象的类型
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		//获取对象中获取key value gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(fmt.Sprintf("%s", reflex.Db))
		switch tagv {
		case "-":
			continue
		case "":
			//强制执行default
			//fallthrough
		default:
			// get tag name with the tag opton, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			switch index {
			case 0:
				out = getStructKey(value, fi.Name, tagv, out)
			case 1:
				out = getStructValue(value, fi.Name, out)
			case 3:
				out = getStructAllKey(tagv, out)
			case 4:
				out = getStructAllValue(value, fi.Name, out)
			}
		}
	}
	return out
}

func getStructKey(value reflect.Value, fiName string, tagv string, out []string) []string {
	//如果里面返还的是结构体，那么执行里面内容
	a, b, _ := isStructValue(value, fiName)
	if a || b {
		out = getStructAllKey(tagv, out)
	}
	return out
}

// 返还所有key
func getStructAllKey(tagv string, out []string) []string {
	out = append(out, tagv)
	return out
}

func getStructValue(value reflect.Value, fiName string, out []string) []string {
	//这里面就有三种状态了 1.struct对象返还内容 2.基本类型返还值  3.struct对象内容返还
	tag, b, val := isStructValue(value, fiName)
	if tag {
		out = append(out, fmt.Sprintf("%v", value.FieldByName(fiName)))
	}

	if b {
		out = append(out, val)
	}

	return out
}

func getStructAllValue(value reflect.Value, fiName string, out []string) []string {
	//这里面就有三种状态了 1.struct对象返还内容 2.基本类型返还值  3.struct对象内容返还
	tag, b, val := isStructValue(value, fiName)
	//目前按照if else判断，后期改成单个判断
	if tag {
		out = append(out, fmt.Sprintf("%v", value.FieldByName(fiName)))
	} else {
		if b {
			out = append(out, val)
		} else {
			out = append(out, "null")
		}
	}

	return out
}

// 针对db内容
// 第一个返回值当前这个对象是基本类型，第二个返回值 这个对象参数是否可用， 第三个参数对象返还值内容
func isStructValue(value reflect.Value, fiName string) (bool, bool, string) {
	if reflex.IsStruct(value.FieldByName(fiName)) {
		//如果里面返还的是结构体，那么执行里面内容
		refStruct := value.FieldByName(fiName)

		refNumfield := refStruct.NumField()
		if 2 == refNumfield {
			keys := refStruct.Field(0)
			vals := refStruct.Field(1)

			if b, err := strconv.ParseBool(fmt.Sprintf("%v", vals)); err == nil {
				return false, b, fmt.Sprintf("%v", keys)
			}
		}
	}
	return true, false, ""
}
