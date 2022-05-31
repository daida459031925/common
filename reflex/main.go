package reflex

import "reflect"

//深度对比两个对象是否相等
func Equals(type1 any, type2 any) bool {
	return reflect.DeepEqual(type1, type2)
}

func Equal(type1 any, type2 any) bool {
	return type1 == type2
}
