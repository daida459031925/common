# common 常用方法汇总避免重复造轮子

# 1.try CatchAll

# 2.雪花id生成 idgenerator

# 3.go反射 reflex 

## 获取db模型中 key value 为自定义开发sql做准备

# 4.result 自定义返还对象

# 5.time时间工具类

# 6.translate 翻译工具

# 7.zip 文件压缩


安装断点
为了解决这个问题，您需要更新您的Delve调试器到一个支持Go 1.23.4的版本。您可以通过以下步骤来更新Delve：
卸载当前版本的Delve：
go uninstall github.com/go-delve/delve/cmd/dlv
安装最新版本的Delve：
go install github.com/go-delve/delve/cmd/dlv@latest
确保您的PATH环境变量包含了Go的bin目录，以便能够找到新安装的Delve。