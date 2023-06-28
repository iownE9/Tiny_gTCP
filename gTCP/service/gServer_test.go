package service

import (
	"gTCP/api"
	"gTCP/bean"
	"gTCP/client"
	"gTCP/router"
	"log"
	"sync"
	"testing"
	"time"
)

// 开启 client 数量
var num int = 8

// go test -v -run=TestGServer
func TestGServer(t *testing.T) {
	// 装配 handlerFunc
	router.AddHandleFunc(1, router.HandlerFunc(f01))
	router.AddHandleFunc(2, router.HandlerFunc(f02))
	router.AddHandleFunc(3, router.HandlerFunc(f03))
	router.AddHandleFunc(4, router.HandlerFunc(f04))
	router.AddHandleFunc(5, router.HandlerFunc(f05))

	// 开启服务器
	go func() { log.Fatal(ListenAndServe("gTCP v4 并发")) }()

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

var msgs = []api.GMessage{
	bean.NewMessage(0, []byte("")), // 占位
	bean.NewMessage(1, []byte("hello gTCP v1")),
	bean.NewMessage(2, []byte("hello gTCP v2")),
	bean.NewMessage(3, []byte("hello gTCP v3")),
	bean.NewMessage(4, []byte("hello gTCP v4")),
	bean.NewMessage(5, []byte("hello gTCP v5")),
}

func f01(msg api.GMessage) api.GMessage {
	return msgs[1]
}
func f02(msg api.GMessage) api.GMessage {
	return msgs[2]
}
func f03(msg api.GMessage) api.GMessage {
	return msgs[3]
}
func f04(msg api.GMessage) api.GMessage {
	return msgs[4]
}
func f05(msg api.GMessage) api.GMessage {
	return msgs[5]
}
