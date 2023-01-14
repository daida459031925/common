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
var imageDecodes = make(map[string]func(r io.Reader) (image.Image, error))
var imageEncodes = make(map[string]func(w io.Writer, m image.Image) error)

func init() {
	imageDecode()
	imageEncode()
}

func imageDecode() {
	jpg := func(r io.Reader) (image.Image, error) {
		return jpeg.Decode(r)
	}

	png := func(r io.Reader) (image.Image, error) {
		return png.Decode(r)
	}

	gif := func(r io.Reader) (image.Image, error) {
		return gif.Decode(r)
	}

	imageDecodes["JPG"] = jpg
	imageDecodes["JPEG"] = jpg
	imageDecodes["PNG"] = png
	imageDecodes["GIF"] = gif

	imageDecodes["jpg"] = jpg
	imageDecodes["jpeg"] = jpg
	imageDecodes["png"] = png
	imageDecodes["gif"] = gif
}

func imageEncode() {
	jpg := func(w io.Writer, m image.Image) error {
		return jpeg.Encode(w, m, nil)
	}

	png := func(w io.Writer, m image.Image) error {
		return png.Encode(w, m)
	}

	gif := func(w io.Writer, m image.Image) error {
		return gif.Encode(w, m, nil)
	}

	imageEncodes["JPG"] = jpg
	imageEncodes["JPEG"] = jpg
	imageEncodes["PNG"] = png
	imageEncodes["GIF"] = gif

	imageEncodes["jpg"] = jpg
	imageEncodes["jpeg"] = jpg
	imageEncodes["png"] = png
	imageEncodes["gif"] = gif
}

// GetImageFromNet 从远程读取图片
func GetImageFromNet(url string) result.Result {
	return result.OkData(url).SetFuncErr(func(a any) any {
		s, e := fmt.ParseUnPointer[string](a)
		err.RuntimeExceptionTF(e != nil, e)
		res, e := http.Get(s)
		err.RuntimeExceptionTF(e != nil || res.StatusCode != result.OK, e)
		//从网络获取照片不需要知道图片类型直接可以使用image解析，目前只能解析jpg、png、gif
		defer res.Body.Close()
		m, _, e := image.Decode(res.Body)
		err.RuntimeExceptionTF(e != nil, e)
		return m
	}, err.New("远程读取文件失败")).SetFunc(func(a any) any {
		m, e := fmt.ParseUnPointer[image.Image](a)
		err.RuntimeExceptionTF(e != nil, e)
		//可以做其他处理
		return m
	}).Exec()

}

func GetImageLoad(url string) (result.Result, string) {
	imageType := ""
	return result.OkData(url).SetFunc(func(a any) any {
		path, e := fmt.ParseUnPointer[string](a)
		err.RuntimeExceptionTF(e != nil, e)
		f, e := os.Open(path)
		err.RuntimeExceptionTF(e != nil, e)
		m, imageTypeName, e := image.Decode(bufio.NewReader(f))
		imageType = imageTypeName
		err.RuntimeExceptionTF(e != nil, e)
		return m
	}).Exec(), imageType
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

func NewImage(typeName, path string, srcImage image.Image) {
	r := result.OkData(typeName).SetFuncErr(func(a any) any {
		r, e := fmt.ParseUnPointer[string](a)
		err.RuntimeExceptionTF(e != nil || len(r) <= 0, e)
		return r
	}, err.New("图片文件格式不正确")).Exec()
	if r.Status == result.OK {
		result.OkData(path).SetFuncErr(func(a any) any {
			r, e := fmt.ParseUnPointer[string](a)
			err.RuntimeExceptionTF(e != nil || len(r) <= 0, e)
			fun := imageEncodes[typeName]
			err.RuntimeExceptionTF(fun == nil, e)
			f, e := os.Create(r)
			err.RuntimeExceptionTF(e != nil, e)
			e = fun(bufio.NewWriter(f), srcImage)
			err.RuntimeExceptionTF(e != nil, e)
			return r
		}, err.New("图片生成失败")).Exec()
	}

}
