package build

import (
	err "github.com/daida459031925/common/error"
	"github.com/daida459031925/common/file"
	"github.com/daida459031925/common/fmt"
	"os"
	"os/exec"
	"runtime"
)

type Goarch string

const (
	amd64    = "amd64"
	arm64    = "arm64"
	arm      = "arm"
	i386     = "386"
	mips64le = "mips64le"
	mipsle   = "mipsle"

	LinuxAMD64    Goarch = amd64
	LinuxARM64    Goarch = arm64
	LinuxARM      Goarch = arm
	LinuxI386     Goarch = i386
	LinuxMIPS64LE Goarch = mips64le
	LinuxMIPSLE   Goarch = mipsle
	WindowsAMD64  Goarch = amd64
	WindowsARM64  Goarch = arm64
	WindowsI386   Goarch = i386
	MacOsAMD64    Goarch = amd64
	MacOsARM64    Goarch = arm64
	FreeBSDAMD64  Goarch = amd64
	FreeBSDARM64  Goarch = arm64
	FreeBSDARM    Goarch = arm
	FreeBSDI386   Goarch = i386
)

// PackageOptions 打包选项
type PackageOptions struct {
	SourceDir string // 源码目录
	OutputDir string // 输出目录
	Goos      string //	系统 windows linux
	Goarch    string //	架构 arm 386等
	FileName  string //文件名字
}

// myPackage 执行打包操作
func myPackage(opts PackageOptions) error {
	// 如果未指定目标平台，则使用当前系统作为默认平台
	if opts.Goos == "" {
		opts.Goos = runtime.GOOS
		e := os.Setenv("GOOS", opts.Goos)
		if e != nil {
			return e
		}
	}

	if opts.Goarch == "" {
		opts.Goarch = runtime.GOARCH
		e := os.Setenv("GOARCH", opts.Goarch)
		if e != nil {
			return e
		}
	}

	if opts.FileName == "" {
		opts.FileName = fmt.Sprintf("myProject-%s-%s", opts.Goos, opts.Goarch)
	}

	// 创建并执行命令行 SourceDir 最后一个参数指定是哪个位置的程序需要打包
	cmd := exec.Command("go", "build", "-o", opts.FileName, opts.SourceDir)
	//输出位置 OutputDir
	cmd.Dir = opts.OutputDir
	e := cmd.Run()
	if e != nil {
		return e
	}

	fmt.Printlnf("Package created: %s\n", opts.OutputDir)
	return nil
}

// 公共方法打包项目
func base(srcDir, outDir, arch, os string, win ...string) error {
	if srcDir == "" {
		return err.New("项目源文件路径不能为空")
	}

	if outDir == "" {
		outDir = "./myProject"
	}
	e := file.MkdirAll(outDir)
	if e != nil {
		fmt.Println(e)
		return err.New("创建目录异常")
	}

	// 打包 Windows 64 位版本
	opts := PackageOptions{
		SourceDir: srcDir,
		OutputDir: outDir,
		Goos:      os,
		Goarch:    arch,
	}

	if win != nil {
		opts.FileName = fmt.Sprintf("myProject-%s-%s%s", opts.Goos, opts.Goarch, win[0])
	}

	e = myPackage(opts)
	if e != nil {
		fmt.Println(e)
		return err.New("项目打包失败")
	}

	return nil
}

// GetWindows 打包 Windows 版本
func GetWindows(srcDir, outDir string, arch Goarch) {
	e := base(srcDir, outDir, string(arch), "windows", ".exe")
	if e != nil {
		fmt.Println(e)
	}
}

// GetLinux 打包 Linux 版本
func GetLinux(srcDir, outDir string, arch Goarch) {
	e := base(srcDir, outDir, string(arch), "linux")
	if e != nil {
		fmt.Println(e)
	}
}

// GetMacOs 打包 macOS 版本
func GetMacOs(srcDir, outDir string, arch Goarch) {
	e := base(srcDir, outDir, string(arch), "darwin")
	if e != nil {
		fmt.Println(e)
	}
}

// GetFreeBSD 打包 FreeBSD 版本
func GetFreeBSD(srcDir, outDir string, arch Goarch) {
	e := base(srcDir, outDir, string(arch), "freebsd")
	if e != nil {
		fmt.Println(e)
	}
}
