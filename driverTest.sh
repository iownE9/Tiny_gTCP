#!/bin/bash

# 测试 clientfd 读写分离 v3
go test -v -run=TestGServer gTCP/service

# 查找日志 成功关闭的 v3
grep "error eof" gTCP/static/log/log_v3.txt -i -n



# 
# 

# 测试 TLV 拆包打包 
# go test -v -run=TestDataPack gTCP/bean
