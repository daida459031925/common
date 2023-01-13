package image

import (
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/result"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
)

// 初始化jpeg、png、gif默认使用image解析默认就要加载所有格式图片。莫名其妙解决image bug问题
// 由于已经知道了image.Decode(res.Body)可以解析网络图片但是目前只能解析jpg、png、gif
// 为了实现从字符串中拿到这个是什么类型的图片调用对应的方法
var imageParses = make(map[string]func(r io.Reader) (image.Image, error))

func init() {
	jpg := func(r io.Reader) (image.Image, error) {
		return jpeg.Decode(r)
	}

	png := func(r io.Reader) (image.Image, error) {
		return png.Decode(r)
	}

	gif := func(r io.Reader) (image.Image, error) {
		return gif.Decode(r)
	}

	imageParses["JPG"] = jpg
	imageParses["JPEG"] = jpg
	imageParses["PNG"] = png
	imageParses["GIF"] = gif
}

// GetImageFromNet 从远程读取图片
func GetImageFromNet(url string) result.Result {
	return result.OkData(url).SetFunc(func(a any) any {
		s, e := fmt.ParseUnPointer[string](a)
		err.RuntimeExceptionTF(e != nil, e)
		res, e := http.Get(s)
		err.RuntimeExceptionTF(e != nil || res.StatusCode != result.OK, e)
		//从网络获取照片不需要知道图片类型直接可以使用image解析，目前只能解析jpg、png、gif
		defer res.Body.Close()
		m, _, e := image.Decode(res.Body)
		err.RuntimeExceptionTF(e != nil, e)
		return m
	}).SetFunc(func(a any) any {
		m, e := fmt.ParseUnPointer[image.Image](a)
		err.RuntimeExceptionTF(e != nil, e)
		//可以做其他处理
		return m
	}).Exec()

}

func GetImageLoad(url string) result.Result {
	return result.OkData(url).SetFunc(func(a any) any {
		//r := a.(string)

		return a
	}).Exec()
}
