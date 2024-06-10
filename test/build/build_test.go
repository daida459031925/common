package build

import (
	"github.com/daida459031925/common/build"
	"testing"
)

// 需要指定到go的main方法中才能进行打包
func TestBuild(t *testing.T) {
	build.GetWindows("D:\\GolandProjects\\go_cloud\\common\\test\\main", "./asd", build.WindowsI386)
}
