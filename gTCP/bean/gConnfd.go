package bean

import (
	"gTCP/api"
	"gTCP/router"
	"gTCP/utils"
	"io"
	"log"
	"net"
	"sync"
)

// v4 服务器 对 clientfd 的多个 msg 读写分离 并发处理
type clientfd struct {
	conn net.Conn
	read chan api.GMessage
	pack chan api.GMessage
	send chan []byte
	exit chan struct{}
}

func Clientfd(conn *net.TCPConn) api.GConnfd {
	chanCap := utils.GlobalConfig.ChanCap

	// 描述符实例对象
	fd := &clientfd{
		conn: conn,
		read: make(chan api.GMessage, chanCap),
		pack: make(chan api.GMessage, chanCap),
		send: make(chan []byte, chanCap),
		// 容量为 >=2， 读写都有权限通知关闭 connfd
		exit: make(chan struct{}, chanCap),
	}
	return fd
}

// 接收消息
func (fd *clientfd) ReadMsg() {
	defer close(fd.read) // 关闭 fd.handleMsg 协程
	dp := DataPack()
	for {
		readMsg, err := dp.Unpack(fd.conn)
		if err != nil {
			if err == io.EOF {
				log.Println("INFO: conn is EOF")
			} else {
				log.Println("ERROR: Unpack err", err)
			}
			fd.exit <- struct{}{}
			return
		}
		fd.read <- readMsg
	}
}

// 处理消息
func (fd *clientfd) HandleMsg() {
	defer close(fd.pack) // 关闭 fd.packMsg 协程

	var wg sync.WaitGroup
	for msg := range fd.read {
		wg.Add(1)
		go func(msg api.GMessage) {
			defer wg.Done()

			sendMsg := router.GRouter.HandlerTagMsg(msg)
			if sendMsg == nil {
				log.Println("ERROR: HandlerTagMsg(msg) is nil")
			} else {
				fd.pack <- sendMsg
			}
		}(msg)
	}
	wg.Wait() // 防止 fd.pack 关闭后 还向它发送
}

// 打包回复
func (fd *clientfd) PackMsg() {
	defer close(fd.send) // 关闭 fd.sendMsg 协程

	var wg sync.WaitGroup
	dp := DataPack()
	for resp := range fd.pack {
		wg.Add(1)
		go func(resp api.GMessage) {
			defer wg.Done()
			sendData, err := dp.Pack(resp)
			if err != nil {
				log.Println("ERROR: PackMsg() err")
			} else {
				fd.send <- sendData
			}
		}(resp)
	}
	wg.Wait() // 防止 fd.send 关闭后 还向它发送
}

// 发送 msg
func (fd *clientfd) SendMsg() {
	for sendData := range fd.send {
		_, err := fd.conn.Write(sendData)
		if err != nil {
			log.Println("ERROR: conn.Write err", err)
			fd.exit <- struct{}{}
			break
		}
	}

	// 消耗掉未发送的 msg
	for range fd.send {
	}
}

// 连接关闭
func (fd *clientfd) Closefd() {
	// 只执行一次
	select {
	case <-fd.exit:
		fd.conn.Close() // 关闭 clientfd
		// close(fd.exit) // ERROR
		// 让 容量  去存放 ReadMsg ReadMsg 中慢的一方
	}
}
