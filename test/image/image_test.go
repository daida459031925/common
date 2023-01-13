package main

import (
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/image"
	"github.com/daida459031925/common/result"
	ima "image"
	"testing"
)

// go test -bench=方法名 -benchmem
// ns/op平均每次多少时间 1s=1000ms 1ms=1000us 1us=1000ns
// allocs/op进行多少次内存分配
// B/op标识每次操作分配多少字节
func BenchmarkImage(b *testing.B) {
	r := image.GetImageFromNet("http://cloudcache.tencent-cloud.com/open_proj/proj_qcloud_v2/community/portal/css/img/wechat-qr.jpg")
	fmt.Printlnf("%s", r.Msg)
	fmt.Println(r.Date)
}

// 普通测试
func TestImage(t *testing.T) {
	r := image.GetImageFromNet("https://bkimg.cdn.bcebos.com/pic/79f0f736afc37931207276aee1c4b74543a9111a")
	fmt.Printlnf("logs: %s %s", r.Msg, "?")
	fmt.Println(r.Date)
	if r.Status == result.OK {
		rimage := r.Data.(ima.Image)
		fmt.Println(rimage.Bounds().Min)
		fmt.Println(rimage.Bounds().Max)
	}
}
