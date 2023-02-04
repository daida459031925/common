package jsteg

//隐写术
//jsteg是一个用于在jpeg文件中隐藏数据的包，
//这是一种被称为隐写术的技术。
//这是通过将数据的每一位复制到图像的最低有效位来实现的。
//可以隐藏的数据量取决于jpeg的文件大小；
//存储隐藏数据的每个字节需要大约10-14字节的jpeg。
import (
	"github.com/lukechampine/jsteg"
	"image/jpeg"
	"os"
)

func main() {
	// open an existing jpeg
	f, _ := os.Open(filename)
	img, _ := jpeg.Decode(f)

	// add hidden data to it
	out, _ := os.Create(outfilename)
	data := []byte("my secret data")
	jsteg.Hide(out, img, data, nil)

	// read hidden data:
	hidden, _ := jsteg.Reveal(out)
}
