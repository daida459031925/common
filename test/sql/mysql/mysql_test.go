package office

import (
	"database/sql"
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/reflex"
	"github.com/daida459031925/common/sql/mysql"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func TestMySql(t *testing.T) {
	c, e := mysql.NewDbConfig("..\\mysql.yml")
	if e != nil {
		fmt.Println(e)
	}
	r, e := c.ConnectCheck("..\\redis.yml")
	if e != nil {
		fmt.Println(e)
	}
	rdb, _ := r.Gdb.DB()
	redis := r.Redis
	fmt.Println(redis.Get("cache_key"))
	defer r.Close()

	var i int64
	r.Gdb.Raw("select count(*) FROM user").Find(&i)

	fmt.Println(i)
	fmt.Println(rdb.Stats().OpenConnections)
}

func TestColumns(t *testing.T) {
	//getColumns(Name{1, "", ""})

	//s := Name{} // 创建一个Student类型的实例
	//
	//getFields(&s) // 传入结构体值
	//getFields(s)  // 传入接口值

	//p := Person{
	//	name:    "Alice",
	//	Age:     20,
	//	Birth:   sql.NullTime{Valid: true, Time: time.Now()},
	//	friends: []string{"Bob", "Charlie"},
	//}
	//getStructInfo(p)
	//getStructInfo(&p)

	user := &User{
		name:    "张三",
		Age:     18,
		Birth:   sql.NullTime{Valid: true, Time: time.Now()},
		friends: []string{"Bob", "Charlie"},
	}

	baseModel := &BaseModel{
		Model: user,
	}

	if _, err := baseModel.GetDbColumns(); err != nil {
		fmt.Println("保存数据出错：", err)
	}
}

type Modeler interface {
	GetDbColumns() (string, error)
}

type BaseModel struct {
	Model Modeler
}

func (b *BaseModel) GetDbColumns() (string, error) {
	if _, err := b.Model.GetDbColumns(); err != nil {
		return "", err
	}
	getStructInfo(b.Model)
	fmt.Printlnf("2")
	// 其他通用逻辑
	return "", nil
}

type User struct {
	name    string
	Age     int
	Birth   sql.NullTime
	friends []string
	names   Name
}

func (u *User) GetDbColumns() (string, error) {
	// 保存用户数据
	fmt.Printlnf("1")
	return "", nil
}

type Name struct {
	id      int64
	Name    string
	AddTime string
	time    sql.NullTime
	Friends []string
}

func getColumns(data any) {
	t := reflect.TypeOf(data)           // 获取Student类型的类型
	v := reflect.ValueOf(data)          // 获取Student类型的值的指针版本
	for i := 0; i < t.NumField(); i++ { // 遍历所有字段
		f := t.Field(i)           // 获取第i个字段
		key := f.Name             // 获取字段的键
		tag := f.Tag              // 获取字段的标签
		value := v.Field(i)       // 获取字段的值
		if value.CanInterface() { // 检查是否可以获取字段的值
			fmt.Printlnf("%s = %v (tag: %v)\n", key, value.Interface(), tag) // 打印键、值和标签
		} else {
			if value.CanAddr() { // 检查是否可以获取值的地址
				addr := value.UnsafeAddr()                                            // 获取未导出字段的地址
				ptr := reflect.NewAt(value.Type(), unsafe.Pointer(addr))              // 创建一个指向该地址的可寻址的反射值
				fmt.Printlnf("%s = %v (tag: %v)\n", key, ptr.Elem().Interface(), tag) // 打印键、值和标签
			} else {
				fmt.Printlnf("%s is unaddressable\n", key) // 打印键是不可寻址的
			}
		}
	}
}

// getFields 函数用于获取任意类型的结构体中的字段的键和值
func getFields(obj any) {
	fmt.Println("Getting fields of", obj)
	t := reflect.TypeOf(obj)  // 获取对象的类型
	v := reflect.ValueOf(obj) // 获取对象的值
	if t.Kind() == reflect.Ptr {
		// 如果对象是指针或者接口，就获取它包含或者指向的值
		t = t.Elem()
		v = v.Elem()
	} else {
		// 如果对象是结构体，就转换成指针类型
		t = reflect.TypeOf(obj)  // 重新获取对象的类型
		v = reflect.ValueOf(obj) // 重新获取对象的值
	}

	for i := 0; i < t.NumField(); i++ { // 遍历所有字段
		f := t.Field(i)                                      // 获取第i个字段
		key := f.Name                                        // 获取字段的键
		tag := f.Tag                                         // 获取字段的标签
		value := v.Field(i)                                  // 获取字段的值
		fmt.Printlnf("%s = %v (tag: %v)\n", key, value, tag) // 打印键、值和标签
	}
	fmt.Println()
}

type Person struct {
	name    string
	Age     int
	Birth   sql.NullTime
	friends []string
}

func getStructInfo(any any) {
	t := reflect.TypeOf(any)
	v := reflect.ValueOf(any)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println("Type:", t.Name())
	fmt.Println("Fields:")
	keys := make([]string, t.NumField())
	values := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i)
		if reflex.IsStruct(val) {
			fmt.Println("val.NumField() :", val.NumField())
			fmt.Println("val.Type() :", val.Type())
		}
		fmt.Println("val.Type() :", val.Type())
		f := t.Field(i)
		keys[i] = fmt.Sprintf(f.Name)
		values[i] = fmt.Sprintf("%v", v.Field(i))
		fmt.Println(" ", f.Name, f.Type, v.Field(i))
	}

	fmt.Println("Modified value:", v.Interface())

	fmt.Println("Modified keys:", keys)
	fmt.Println("Modified value:", values)

}
