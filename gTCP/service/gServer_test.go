package service

import (
	"gTCP/client"
	"log"
	"testing"
	"time"
)

// 开启 client 数量
var num int = 6

// go test -v -run=TestGServer
func TestGServer(t *testing.T) {

	// 开启服务器
	go func() { log.Fatal(ListenAndServe("gTCP v2 TLV消息封装")) }()

	time.Sleep(100 * time.Millisecond) // 有必要

	// 开启客户端测试
	for i := 0; i < num; i++ {
		go client.ClientTLVTest()
	}

	// 主动退出
	time.Sleep(3 * time.Second)
}
