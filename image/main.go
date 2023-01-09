package image

import (
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/result"
	"image"
	"image/jpeg"
	"net/http"
)

// GetImageFromNet 从远程读取图片
func GetImageFromNet(url string) result.Result {
	return result.OkData(url).SetFunc(func(a any) any {
		s, e := fmt.ParseUnPointer[string](a)
		err.RuntimeExceptionTF(e != nil, e)
		res, e := http.Get(s)
		err.RuntimeExceptionTF(e != nil || res.StatusCode != result.OK, e)
		//若这个地方需要判断jpg和png则可以开启第二个方法进行逻辑计算
		defer res.Body.Close()
		m, e := jpeg.Decode(res.Body)
		err.RuntimeExceptionTF(e != nil, e)
		return m
	}).SetFunc(func(a any) any {
		m, e := fmt.ParseUnPointer[image.Image](a)
		err.RuntimeExceptionTF(e != nil, e)
		//可以做其他处理
		return m
	}).Exec()

}
