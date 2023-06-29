#!/bin/bash

# v6 协程池

# gTCP/service
# 协程池整体性功能测试
go test -v -run=TestRouterHandle gTCP/service -race -args 30

# 查看日志 ERROR
# -n：显示匹配行的行号 number
# -c：只打印匹配的行数 count
grep "ERROR" -n  gTCP/static/log/log_v6.txt
grep "正式开启" -c gTCP/static/log/log_v6.txt
grep "协程 复用" -c gTCP/static/log/log_v6.txt

# ======================= #
# v5 正式版

# gTCP/service
# 整体性功能测试
go test -v -run=TestRouterHandle gTCP/service -race -args 3

# 查看日志 ERROR
grep "ERROR" -n  gTCP/static/log/log_v5.txt

# ======================= #

# gTCP/router 

# 路由 函数注册  单元测试 
go test -v -run=TestGRouter gTCP/router -race

# 覆盖率
go test -v -run=TestGRouter gTCP/router -race -coverprofile=gTCP/static/cover/router.out

# 覆盖率 html 
go tool cover -html=gTCP/static/cover/router.out -o gTCP/static/cover/router.html

# ======================= #

# gTCP/handler ghandleConnfd.go

# 每个 clientfd 对消息的 拆包读 处理 打包 发送 单元测试 
go test -v -run=TestHandlerClientfd gTCP/service -race -args 3

# 查看日志 ERROR
grep "ERROR" -n  gTCP/static/log/log_v5.txt

# ======================= #

# gTCP/msg

# 对消息的 TLV 拆包打包 单元测试 
go test -v -run=TestDataPack gTCP/msg -race

# 覆盖率
go test -v -run=TestDataPack gTCP/msg -coverprofile=gTCP/static/cover/dataPack.out

# 覆盖率 html 
go tool cover -html=gTCP/static/cover/dataPack.out -o gTCP/static/cover/dataPack.html

# 基准测试 -> 打包 拆包 耗时
go test -v -bench=Datapack gTCP/msg

# ======================= #

# gTCP/service

# TCP 服务器启动并接收请求 单元测试 数据竞争检测
go test -v -run=TestGServer gTCP/service -race -args 3

# 覆盖率
go test -v -run=TestGServer gTCP/service -coverprofile=gTCP/static/cover/service.out

# 覆盖率 html 
go tool cover -html=gTCP/static/cover/service.out -o gTCP/static/cover/service.html
