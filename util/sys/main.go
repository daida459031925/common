package sys

import (
	err "github.com/daida459031925/common/error"
	"runtime"
)

// GetArch 获取当前系统是多少位
func GetArch() (int, error) {
	var arch int
	if runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64" {
		arch = 64
	} else if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		arch = 32
	} else {
		return arch, err.New("无法获取到对应系统位数")
	}
	return arch, nil
}
