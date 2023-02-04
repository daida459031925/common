package image

import (
	"bufio"
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/fmt"
	"github.com/daida459031925/common/result"
	"github.com/nfnt/resize"
	"golang.org/x/sync/errgroup"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
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

// 图片缩放代码
func main() {
	img := strconv.Itoa(int(time.Now().UnixNano())) + ".jpeg"
	f, err := os.Create(img)
	if err != nil {
		fmt.Println("Create", err)
		return
	}
	// 把二维码图片合到海报上
	rgba, err := MergeImageNew("http://img/16028328876241490.jpeg", "http://8134.jpg", 76, 436, 87)
	if err != nil {
		fmt.Println("MergeImageNew", err)
		return
	}
	err = jpeg.Encode(f, rgba, nil)
	if err != nil {
		fmt.Println("Encode", err)
		return
	}
}

// MergeImageNew 图片合并 baseUrl:原图图片地址，maskUrl：小图图片地址
func MergeImageNew(baseUrl, maskUrl string, paddingX int, paddingY int, width uint) (*image.RGBA, error) {
	eg := errgroup.Group{}
	var base image.Image
	eg.Go(func() error {
		var err error
		base, err = GetImageFromNet(baseUrl)
		return err
	})
	mask, err := ImageZoom(maskUrl, width)
	if err != nil {
		return nil, err
	}
	if err = eg.Wait(); err != nil {
		return nil, err
	}
	baseSrcBounds := base.Bounds().Max
	maskSrcBounds := mask.Bounds().Max

	newWidth := baseSrcBounds.X
	newHeight := baseSrcBounds.Y

	maskWidth := maskSrcBounds.X
	maskHeight := maskSrcBounds.Y

	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
	//首先将一个图片信息存入jpg
	draw.Draw(des, des.Bounds(), base, base.Bounds().Min, draw.Over)
	//将另外一张图片信息存入jpg
	fmt.Println(newHeight, paddingY, maskHeight)
	draw.Draw(des, image.Rect(paddingX, paddingY, (paddingX+maskWidth), (maskHeight+paddingY)), mask, image.Point{}, draw.Over)
	return des, nil
}

// GetImageFromNet 从远程读取图片
func GetImageFromNet(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	return m, err
}

// ImageZoom 按宽度缩放图片
func ImageZoom(url string, width uint) (image.Image, error) {
	m, err := GetImageFromNet(url)
	if err != nil {
		return nil, err
	}
	if width == 0 {
		return m, nil
	}
	thImg := resize.Resize(width, 0, m, resize.Lanczos3)
	return thImg, nil
}

// 图片缩放
func main() {
	f1, err := os.Open("1.jpg")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	f2, err := os.Open("2.jpg")
	if err != nil {
		panic(err)
	}
	defer f2.Close()
	f3, err := os.Create("3.jpg")
	if err != nil {
		panic(err)
	}
	defer f3.Close()
	m1, err := jpeg.Decode(f1)
	if err != nil {
		panic(err)
	}
	bounds := m1.Bounds()
	m2, err := jpeg.Decode(f2)
	if err != nil {
		panic(err)
	}
	m := image.NewRGBA(bounds)
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(m, bounds, &image.Uniform{white}, image.ZP, draw.Src)
	draw.Draw(m, bounds, m1, image.ZP, draw.Src)
	draw.Draw(m, image.Rect(100, 200, 300, 600), m2, image.Pt(250, 60), draw.Src)
	err = jpeg.Encode(f3, m, &jpeg.Options{90})
	if err != nil {
		panic(err)
	}
	fmt.Printf("okn")
}

// golang 图片处理，剪切，base64数据转换，文件存储
// base64 -> file
//
// ddd, _ := base64.StdEncoding.DecodeString(datasource) //成图片文件并把文件写入到buffer
// err2 := ioutil.WriteFile("./output.jpg", ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
//
// datasource base64 string
//
// base64 -> buffer
//
// ddd, _ := base64.StdEncoding.DecodeString(datasource) //成图片文件并把文件写入到buffer
// bbb := bytes.NewBuffer(ddd)                           // 必须加一个buffer 不然没有read方法就会报错
//
// 转换成buffer之后里面就有Reader方法了。才能被图片API decode
//
// buffer-> ImageBuff（图片裁剪,代码接上面）
//
// m, _, _ := image.Decode(bbb)                                       // 图片文件解码
// rgbImg := m.(*image.YCbCr)
// subImg := rgbImg.SubImage(image.Rect(0, 0, 200, 200)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
//
// img -> file(代码接上面)
//
// f, _ := os.Create("test.jpg")     //创建文件
// defer f.Close()                   //关闭文件
// jpeg.Encode(f, subImg, nil)       //写入文件
//
// img -> base64(代码接上面)
//
// emptyBuff := bytes.NewBuffer(nil)                  //开辟一个新的空buff
// jpeg.Encode(emptyBuff, subImg, nil)                //img写入到buff
// dist := make([]byte, 50000)                        //开辟存储空间
// base64.StdEncoding.Encode(dist, emptyBuff.Bytes()) //buff转成base64
// fmt.Println(string(dist))                          //输出图片base64(type = []byte)
// _ = ioutil.WriteFile("./base64pic.txt", dist, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
//
// imgFile -> base64
//
// ff, _ := ioutil.ReadFile("output2.jpg")               //我还是喜欢用这个快速读文件
// bufstore := make([]byte, 5000000)                     //数据缓存
// base64.StdEncoding.Encode(bufstore, ff)               // 文件转base64
// _ = ioutil.WriteFile("./output2.jpg.txt", dist, 0666) //直接写入到文件就ok完活了。
func calcResizedRect(width int, src image.Rectangle, height int, centerAlign bool) image.Rectangle {
	var dst image.Rectangle
	if width*src.Dy() < height*src.Dx() { // width/src.width < height/src.height
		ratio := float64(width) / float64(src.Dx())

		tH := int(float64(src.Dy()) * ratio)
		pad := 0
		if centerAlign {
			pad = (height - tH) / 2
		}
		dst = image.Rect(0, pad, width, pad+tH)
	} else {
		ratio := float64(height) / float64(src.Dy())
		tW := int(float64(src.Dx()) * ratio)
		pad := 0
		if centerAlign {
			pad = (width - tW) / 2
		}
		dst = image.Rect(pad, 0, pad+tW, height)
	}

	return dst
}

func resizePic(img image.Image, width int, height int, keepRatio bool, fill int, centerAlign bool) image.Image {
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))
	if !keepRatio {
		draw.BiLinear.Scale(outImg, outImg.Bounds(), img, img.Bounds(), draw.Over, nil)
		return outImg
	}

	if fill != 0 {
		fillColor := color.RGBA{R: uint8(fill), G: uint8(fill), B: uint8(fill), A: 255}
		draw.Draw(outImg, outImg.Bounds(), &image.Uniform{C: fillColor}, image.Point{}, draw.Src)
	}
	dst := calcResizedRect(width, img.Bounds(), height, centerAlign)
	draw.ApproxBiLinear.Scale(outImg, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return outImg
}
