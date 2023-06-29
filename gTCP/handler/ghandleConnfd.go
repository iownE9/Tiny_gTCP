package handler

import (
	"gTCP/api"
	"gTCP/msg"
	"gTCP/router"
	"gTCP/utils"
	"io"
	"log"
	"net"
)

// 对应 client 连接
type clientfd struct {
	conn *net.TCPConn
	read chan api.GMessage
	pack chan api.GMessage
	send chan []byte

	exit    chan struct{}
	restart chan struct{}
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

		exit:    make(chan struct{}),
		restart: make(chan struct{}),
	}
	return fd
}

// 接收消息
func (fd *clientfd) ReadMsg() {
	for {
		readMsg, err := dp.Unpack(fd.conn)
		if err != nil {
			if err == io.EOF {
				log.Println("INFO: conn is EOF")
			} else {
				log.Println("ERROR: Unpack err", err)
			}

			fd.exit <- struct{}{} // ReadMsg 协程已暂停
			fd.read <- nil        // 发送终止信号
			<-fd.restart          // 等待 获取 新conn success
			continue
		}
		fd.read <- readMsg
	}
}

// 处理消息
func (fd *clientfd) HandleMsg() {
	for msg := range fd.read {
		// 终止信号
		if msg == nil {
			fd.pack <- msg
			continue
		}

		go func(msg api.GMessage) {
			sendMsg := HandleMsgRouter(msg)

			if sendMsg == nil {
				log.Println("ERROR: HandlerTagMsg(msg) is nil")
			} else {
				fd.pack <- sendMsg
			}
		}(msg)
	}
}

// 打包回复
func (fd *clientfd) PackMsg() {
	for resp := range fd.pack {
		// 终止信号
		if resp == nil {
			fd.send <- nil
			continue
		}

		go func(resp api.GMessage) {
			sendData, err := dp.Pack(resp)
			if err != nil {
				log.Println("ERROR: PackMsg() err")
			} else {
				fd.send <- sendData
			}
		}(resp)
	}
}

// 发送消息
func (fd *clientfd) SendMsg() {
	for sendData := range fd.send {
		if sendData == nil {
			// 写异常 晚于终止信号
			fd.exit <- struct{}{} // SendMsg 协程已处理完 旧 fd 的数据
			<-fd.restart          // 等待 获取 新conn success
		} else {

			_, err := fd.conn.Write(sendData) // 正常执行区域

			// 写异常 早于终止信号
			if err != nil {
				log.Println("ERROR: conn.Write err", err)
				//  因异常 消耗掉未发送的 msg
				for sendData := range fd.send {
					if sendData == nil {
						fd.exit <- struct{}{} // SendMsg 协程已处理完 旧 fd 的数据
						break
					}
				}
				<-fd.restart // 等待 获取 新conn success
			} // 写异常 早于终止信号
		}
	}
}

// 连接关闭
func (fd *clientfd) Closefd() {
	log.Println("正式开启 可复用 协程")

	// 读 msg
	go fd.ReadMsg()

	// 处理 msg
	go fd.HandleMsg()

	// 打包 msg
	go fd.PackMsg()

	// 发送 msg
	go fd.SendMsg()

	// 获取新 clientfd
	for {
		// 读写 都已 ok
		_ = <-fd.exit
		_ = <-fd.exit

		fd.conn.Close() //  读写完全 关闭 clientfd

		// 获取新 clientfd
		fd.conn = <-clientfdChan
		log.Println("协程 复用")

		// 发送两个 通知 读写
		fd.restart <- struct{}{}
		fd.restart <- struct{}{}
	}
}
