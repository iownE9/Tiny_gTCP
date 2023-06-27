package service

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

// 开启 client 数量
var num int = 9

// echo 内容
var echoContent string = "hello gTCP v1"

// go test -v -run=TestGServer
func TestGServer(t *testing.T) {

	// 开启服务器
	go func() { log.Fatal(ListenAndServe("gTCP v1 echoTest")) }()

	time.Sleep(100 * time.Millisecond) // 有必要

	// 开启客户端测试
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			if got := clientWRTest(i); echoContent != got {
				t.Errorf("clientWRTest %d: want %s,but got %s", i, echoContent, got)
			}
			defer wg.Done()
		}(i)
	}
	// 全部客户端结束
	wg.Wait()

}

// 模拟 client WR
func clientWRTest(i int) (got string) {
	log.Println("Client Test", i, "... start")

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	// defer conn.Close() // Error
	if err != nil {
		log.Println(`net.Dial("tcp", "127.0.0.1:8080") err, exit!`)
		return
	}

	_, err = conn.Write([]byte(echoContent))
	if err != nil {
		log.Println("conn.Write err", err)
		conn.Close()
		return
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read error")
		conn.Close()
		return
	}

	conn.Close()
	return string(buf[:n])
}

// net.Dial() 失败 -> conn 为 nil -> conn.Close() -> panic
// panic: runtime error: invalid memory address or nil pointer dereference
// [signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x52e8d0]
