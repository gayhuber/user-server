# user-server

基于自用的 `open-server` 做的一个脚手架, 方便 php 进行 tcp 链接, 数据通过 json 通信, 目的在于 api 端无感知的切换到这个 server 上

# 环境
```bash
➜  user-server git:(master) ✗ go version
go version go1.9.2 darwin/amd64
```

# 启动
```bash
go run main.go
```



# 使用到的工具包
[viper](https://github.com/spf13/viper)