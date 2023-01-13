package image

import (
	"bufio"
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/result"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
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
		path, e := fmt.ParseUnPointer[string](a)
		err.RuntimeExceptionTF(e != nil, e)
		f, e := os.Open(path)
		err.RuntimeExceptionTF(e != nil, e)
		m, _, e := image.Decode(bufio.NewReader(f))
		err.RuntimeExceptionTF(e != nil, e)
		return m
	}).Exec()
}

// Rotate90 旋转90度
func Rotate90(m image.Image) image.Image {
	rotate90 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	// 矩阵旋转
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			//  设置像素点
			rotate90.Set(m.Bounds().Max.Y-x, y, m.At(y, x))
		}
	}
	return rotate90
}

// Rotate180 旋转180度
func Rotate180(m image.Image) image.Image {
	rotate180 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	// 矩阵旋转
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			//  设置像素点
			rotate180.Set(m.Bounds().Max.X-x, m.Bounds().Max.Y-y, m.At(x, y))
		}
	}
	return rotate180
}

// Rotate270 旋转270度
func Rotate270(m image.Image) image.Image {
	rotate270 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	// 矩阵旋转
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			// 设置像素点
			rotate270.Set(x, m.Bounds().Max.X-y, m.At(y, x))
		}
	}
	return rotate270

}

// CenterImage 还有个需求就是将长方形图片修改为正方形图片，原图片居中显示（这里是长>=宽）
// 将图片居中处理
func CenterImage(m image.Image) image.Image {
	// 现在图片是长>宽，将图片居中设置
	max := m.Bounds().Dx()
	// 居中后距离最底部的高度为(x-y)/2
	temp := (max - m.Bounds().Dy()) / 2
	centerImage := image.NewRGBA(image.Rect(0, 0, max, max))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			centerImage.Set(x, temp+y, m.At(x, y))
		}
	}
	return centerImage

}
