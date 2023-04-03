# 抖音直播间数据分析

 [![standard-readme compliant](https://camo.githubusercontent.com/f116695412df39ab3c98d8291befdb93af123f56aecc79fff4b20c410a5b54c7/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f726561646d652532307374796c652d7374616e646172642d627269676874677265656e2e7376673f7374796c653d666c61742d737175617265)](https://github.com/RichardLitt/standard-readme)

## 关于

使用这个程序可以实现抖音直播间数据分析和直播互动小游戏。一般实现此类功能有两种方式：
1. 使用 selenium 直接运行 js 来达到目的
2. 使用 https 代理等网络代理直接获取数据处理
本程序采用的是后者，使用很简单。

### 需要

* [go1.15+](https://go.dev/dl/)



## 开始
下载源码后进入源码目录安装依赖
```
go mod tidy
```

### 配置

可以参考我的文章 [golang 使用 elazarl / goproxy 代理https请求](https://zhuanlan.zhihu.com/p/514004767) 将证书配置一下，代码中有我生成的证书在2024年4月后失效，自己生成证书则按照文章中生成。

### 使用

在项目根目录下运行
```
go run .
```
将默认监听 8080 端口，安装好证书后配置本地代理到 ```localhost:8080``` 打开任意直播间即可捕获数据收到礼物时可以看到输出：
```shell
WebcastGiftMessage
收到人气票价值0.10元
```
[parseData 函数下区别不同响应类型做出自己的业务处理](https://github.com/Feiir/dy-live/blob/main/main.go#L109)

## 声明
本代码仅供学习交流使用，任何问题请私信联系






