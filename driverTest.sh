#!/bin/bash

# 测试 TLV 拆包打包
go test -v -run=TestDataPack gTCP/bean

# 启动测试 v2
go test -v -run=TestGServer gTCP/service

# 查找报错日志
grep "error" gTCP/static/log/log_v2.txt -i -n