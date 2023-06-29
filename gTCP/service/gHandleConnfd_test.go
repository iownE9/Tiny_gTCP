package service

import (
	"gTCP/api"
	"gTCP/handler"
	"gTCP/test"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// go test -v -run=TestHandlerClientfd gTCP/service -race -args numClient
func TestHandlerClientfd(t *testing.T) {
	// ERROR info in log.txt

	// 装配客户端处理函数
	test.HandleMsgRouter = test.CustomHandleMsg

	// 替换消息路由函数
	saved := handler.HandleMsgRouter
	defer func() { handler.HandleMsgRouter = saved }()
	handler.HandleMsgRouter = handleMsgRouterTest

	// 开启服务器
	go func() {
		err := ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // 确保先开启 TCPserver

	// 开启 client 数量
	num, _ := strconv.Atoi(os.Args[len(os.Args)-1])

	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		serverfd, err := test.ClientDialTest()

		if err != nil {
			t.Error("test.ClientDialTest() err", err)
			continue
		}
		wg.Add(1)
		go func(serverfd net.Conn) {
			defer wg.Done()
			test.HandleServerfd(serverfd)
		}(serverfd)
	}

	wg.Wait() // 确保客户端完全结束
}

// 回写到客户端
func handleMsgRouterTest(getMsg api.GMessage) api.GMessage {
	return getMsg
}
