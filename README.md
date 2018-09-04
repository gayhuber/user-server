# user-server

基于自用的 `open-server` 做的一个脚手架, 方便 php 进行 tcp 链接, 数据通过 json 通信, 目的在于 api 端无感知的切换到这个 server 上

# 环境
```bash
➜  user-server git:(master) ✗ go version
go version go1.9.2 darwin/amd64
```

# 启动
```bash
cp config-example.json config.json && go run main.go
```

# 使用到的工具包
* [viper](https://github.com/spf13/viper)
* [beego-log](https://github.com/astaxie/beego/logs)

这里的 vendor 包使用的是 [govendor](https://github.com/kardianos/govendor) 工具来管理,
```bash
govendor add -tree xxx   // 加载本地库到 vendor 中
govendor fetch xxx       // 从远程拉取库到 vendor
```
 当想添加需要的库时,可以运行`govendor get xxx`来把本地的库加载到 vendor 目录中,