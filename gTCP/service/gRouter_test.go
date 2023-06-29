package service

import (
	"gTCP/api"
	"gTCP/msg"
	"gTCP/router"
	"gTCP/test"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// go test -v -run=TestRouterHandle gTCP/service -race -args 3
func TestRouterHandle(t *testing.T) {
	// ERROR info in log.txt

	// 装配客户端处理函数
	test.HandleMsgRouter = test.CustomHandleMsgRouter

	// 装配 handlerFunc
	router.AddHandleFunc(1, router.HandlerFunc(f01))
	router.AddHandleFunc(2, router.HandlerFunc(f02))
	router.AddHandleFunc(3, router.HandlerFunc(f03))
	router.AddHandleFunc(4, router.HandlerFunc(f04))
	router.AddHandleFunc(5, router.HandlerFunc(f05))

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

// 1-5 2-4 3-3
func f01(msg api.GMessage) api.GMessage {
	return msgs[0]
}
func f02(msg api.GMessage) api.GMessage {
	return msgs[1]
}
func f03(msg api.GMessage) api.GMessage {
	return msgs[2]
}
func f04(msg api.GMessage) api.GMessage {
	return msgs[3]
}
func f05(msg api.GMessage) api.GMessage {
	return msgs[4]
}

var msgs = []api.GMessage{
	msg.NewMessage(5, []byte("hello gTCP v1")),
	msg.NewMessage(4, []byte("hello gTCP v2")),
	msg.NewMessage(3, []byte("hello gTCP v3")),
	msg.NewMessage(2, []byte("hello gTCP v4")),
	msg.NewMessage(1, []byte("hello gTCP v5")),
}
