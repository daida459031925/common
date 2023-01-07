package main

import (
	"fmt"
	"github.com/daida459031925/common/image"
	"testing"
)

// go test -bench=方法名 -benchmem
// ns/op平均每次多少时间 1s=1000ms 1ms=1000us 1us=1000ns
// allocs/op进行多少次内存分配
// B/op标识每次操作分配多少字节
func BenchmarkImage(b *testing.B) {
	r := image.GetImageFromNet("http://cloudcache.tencent-cloud.com/open_proj/proj_qcloud_v2/community/portal/css/img/wechat-qr.jpg")
	fmt.Println(r.Msg)
	fmt.Println(r.Date)
}

// 普通测试
func TestImage(t *testing.T) {
	r := image.GetImageFromNet("http://cloudcache.tencent-cloud.com/open_proj/proj_qcloud_v2/community/portal/css/img/wechat-qr.jpg")
	fmt.Println(r.Msg)
	fmt.Println(r.Date)
}
