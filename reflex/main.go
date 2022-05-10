package reflex

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//可以拿到结构体中对应的数据，不管时什么结构体都能拿到
func PrintReflexBody(in any) {
	t := reflect.TypeOf(in)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		fmt.Printf("Check type eror not Struct")
	}

	v := reflect.ValueOf(in)

	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		fmt.Sprintf("field %d, type %s, key %s ,value %v",
			i, t.Field(i).Type, t.Field(i).Name, v.FieldByName(t.Field(i).Name))
	}
}

//根据传入的切片，通过反射拿到对应数据values,若为nil则不返还
//针对格式Account         sql.NullString `db:"account"`             // 账号
//返还struce 类型所有0:key 1:value
func RawField(in any, index int) ([]string, error) {

	v := reflect.ValueOf(in)
	//这不知道是干什么的 但是go-zero里面写了
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	//如果make 切片 类型
	if v.Kind() == reflect.Slice {
		return getField(v, index), nil
	}

	//如果时结构体
	if v.Kind() == reflect.Struct {
		return getNumField(v, index), nil
	}

	return nil, errors.New("无法获取")
}

//返还slice 类型所有values
func getField(in reflect.Value, index int) []string {
	out := make([]string, 0)
	for i := 0; i < in.Len(); i++ {
		value := in.Index(i)
		out = append(getNumField(value, index))
	}
	return out
}

//返还struce 类型所有0:key 1:value
func getNumField(value reflect.Value, index int) []string {
	out := make([]string, 0)
	//找到这个对象的类型
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		//获取对象中获取key value gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get("db")
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
			}
		}
	}
	return out
}

func getStructKey(value reflect.Value, fiName string, tagv string, out []string) []string {
	//如果里面返还的是结构体，那么执行里面内容
	a, b, _ := isStructValue(value, fiName)
	if a || b {
		out = append(out, tagv)
	}
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

//针对db内容
//第一个返回值当前这个对象是基本类型，第二个返回值 这个对象参数是否可用， 第三个参数对象返还值内容
func isStructValue(value reflect.Value, fiName string) (bool, bool, string) {
	if value.FieldByName(fiName).Kind() == reflect.Struct {
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
