#!/bin/bash

# v4
# 单元测试 路由
go test -run=TestGRouter -v -race gTCP/router_test

# 测试 clientfd 读写分离 v4
go test -v -run=TestGServer gTCP/service

# 查找日志 v4
grep "got" gTCP/static/log/log_v4.txt -i -n

grep "ERROR:" gTCP/static/log/log_v4.txt -i -n

# # v3
# # 测试 clientfd 读写分离 v3
# go test -v -run=TestGServer gTCP/service

# # 查找日志 成功关闭的 v3
# grep "error eof" gTCP/static/log/log_v3.txt -i -n

# # v2 消息封装
# 测试 TLV 拆包打包 v2
# go test -v -run=TestDataPack gTCP/bean
