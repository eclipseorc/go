# go
go语言杂项工具


### go mod使用方法
  1、创建环境变量（系统变量）GO111MODULE=on/auto，如果已存在则忽略<br/>
  2、创建环境变量（系统变量）GOPROXY=goproxy.cn，因默认的代理无法访问导致下载包失败。如果已存在则忽略<br/>
  3、在需要使用go mod的文件下打开console，然后使用go mod init 模块名，来初始化go mod
  4、使用go mod tidy 来清理无用的包引用
  
  