golang proxy server
---
[![Go](https://github.com/Crtrpt/gps/actions/workflows/go.yml/badge.svg)](https://github.com/Crtrpt/gps/actions/workflows/go.yml)
---
golang 实现的外网代理回调服务器 主要让外网访问内网api用途

## 运行
```golang
git clone git@github.com:Crtrpt/gps.git
cd gps
go mod tidy
go run cmd/gps/main.go
```


## 流程图
![流程图](./flow.svg "工作流程图")

## 特性
- ssh隧道代理
- 本地代理
- 远程代理

## 问题
- static 安全访问
- ssh不稳定问题 断开的问题
- 文档
- 错误处理的问题
- 增加对udp tcp的支持

## 注意
如果需要sshd远程隧道 需要开启 sshd 隧道功能 否则无法监听外部端口
当前只支持key访问
```
GatewayPorts yes
```
