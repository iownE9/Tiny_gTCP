package service

import (
	"gTCP/client"
	"log"
	"sync"
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

	var wg sync.WaitGroup
	// 开启客户端测试
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.ClientTLVTest()
		}()
	}

	// 等待退出
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

}
