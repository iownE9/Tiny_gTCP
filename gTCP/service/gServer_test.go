package service

import (
	"gTCP/test"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

// go test -v -run=TestGServer gTCP/service -args numClient
func TestGServer(t *testing.T) {
	// 替换 clientfd 处理函数
	saved := handleClientFunc
	defer func() { handleClientFunc = saved }()

	var gotClientAddr = make(chan string)
	handleClientFunc = func(clientfd *net.TCPConn) {
		gotClientAddr <- clientfd.RemoteAddr().String()
		defer clientfd.Close()
	}

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
	// 	log.Println(os.Args)
	//  [/tmp/go-build2202014767/b001/service.test -test.paniconexit0 -test.timeout=10m0s -test.v=true -test.run=TestGServer 3]

	// 开启客户端测试
	for i := 0; i < num; i++ {
		serverfd, err := test.ClientDialTest()
		if err != nil {
			t.Error("test.ClientDialTest() err", err)
			continue
		}
		want := serverfd.LocalAddr().String()
		serverfd.Close()

		got := <-gotClientAddr

		if want != got {
			t.Errorf("want %s,but got %s\n", want, got)
		}
	}
}
