package handler

import (
	"gTCP/api"
	"gTCP/msg"
	"gTCP/router"
	"gTCP/utils"
	"io"
	"log"
	"net"
	"sync"
)

// 对应 client 连接
type clientfd struct {
	conn *net.TCPConn
	read chan api.GMessage
	pack chan api.GMessage
	send chan []byte
	exit chan struct{}
}

// 消息打包拆包句柄
var dp api.GDataPack = msg.DataPack()

// channel 容量
var chanCap uint8 = utils.GlobalConfig.ChanCap

// 处理 msg router 方便测试替换
var HandleMsgRouter func(api.GMessage) api.GMessage = router.GRouter.HandlerTagMsg

func Clientfd(conn *net.TCPConn) api.GConnfd {
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

// 处理 Clientfd
func HandleClientfd(conn *net.TCPConn) {
	fd := Clientfd(conn)
	fd.Closefd()
}

// 接收消息
func (fd *clientfd) ReadMsg() {
	defer close(fd.read) // 关闭 fd.handleMsg 协程
	for {
		readMsg, err := dp.Unpack(fd.conn)
		if err != nil {
			if err == io.EOF {
				log.Println("INFO: conn is EOF")
			} else {
				log.Println("ERROR: Unpack err", err)
			}
			fd.exit <- struct{}{}
			break
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

			sendMsg := HandleMsgRouter(msg)

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

// 发送消息
func (fd *clientfd) SendMsg() {
	for sendData := range fd.send {
		_, err := fd.conn.Write(sendData)
		if err != nil {
			log.Println("ERROR: conn.Write err", err)
			fd.exit <- struct{}{}
			break
		}
	}

	// 因异常跳出上面 forrange 消耗掉未发送的 msg
	for range fd.send {
	}
}

// 连接关闭
func (fd *clientfd) Closefd() {

	// 读 msg
	go fd.ReadMsg()

	// 处理 msg
	go fd.HandleMsg()

	// 打包 msg
	go fd.PackMsg()

	// 发送 msg
	go fd.SendMsg()

	// 只执行一次
	select {
	case <-fd.exit:
		fd.conn.Close() // 关闭 clientfd
		// close(fd.exit) // ERROR
		// 让 容量  去存放 ReadMsg ReadMsg 中慢的一方
	}
}
