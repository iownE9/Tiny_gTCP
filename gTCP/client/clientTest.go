package client

import (
	"gTCP/api"
	"gTCP/bean"
	"gTCP/utils"
	"io"
	"log"
	"net"
	"time"
)

const sendData1 string = "hello gTCP v1"
const sendData2 string = "hello gTCP v2"

// 对应连接对象
type serverfd struct {
	conn net.Conn          // serverfd 文件描述符
	send chan []byte       // 发送消息管道
	read chan api.GMessage // 接收消息管道
	exit chan struct{}     // 关闭 fd，退出协程
}

// TLV 拆包装包 句柄
var dp api.GDataPack
var chanCap uint8 // channel 容量

// 初始化 client 包 全局只读变量
func init() {
	dp = bean.DataPack()
	chanCap = utils.GlobalConfig.ChanCap
	log.Println("ClientTest Init ok")
}

// 模拟 client TLV 格式消息发送
func ClientTLVTest() {
	// 模拟粘包数据
	sendData := msgTLV()
	if sendData == nil {
		log.Println("ERROR: msgTLV() err")
		return
	}

	// 与服务器建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	// defer conn.Close() // Error
	if err != nil {
		log.Println(`ERROR: net.Dial("tcp", "127.0.0.1:8080") err, exit!`)
		return
	}

	// 描述符实例对象
	fd := &serverfd{
		conn: conn,
		send: make(chan []byte, chanCap),
		read: make(chan api.GMessage, chanCap),
		// 容量 >=2， 读写都有权限通知关闭 connfd
		exit: make(chan struct{}, chanCap),
	}

	go func() { fd.send <- sendData }() // 缓存有无均可

	go func() {
		// 模拟测试 v2 的 一直阻塞在 read 的 bug
		time.Sleep(1 * time.Second)
		select {
		case <-fd.exit:
			// 若已经关闭，直接退出
		default:
			fd.send <- sendData
		}

		time.Sleep(1 * time.Second)
		fd.exit <- struct{}{}
	}()

	// 接收消息
	go fd.readMsg()
	// 处理消息
	go fd.handleMsg()
	// 发送消息
	go fd.sendMsg()

	// 一直阻塞直到连接关闭
	select {
	case <-fd.exit:
		fd.conn.Close() // 关闭 serverfd

		// 该 goroutine 为 fd.send 发送方
		// 消耗掉未发送的 msg
		close(fd.send)
		for range fd.send {
		}
	}
}

// 发送消息
func (fd *serverfd) sendMsg() {
	for sendData := range fd.send {
		_, err := fd.conn.Write(sendData)
		if err != nil {
			log.Println("INFO: conn.Write err", err)
			fd.exit <- struct{}{}
			return
		}
	}
}

// 接收消息
func (fd *serverfd) readMsg() {
	defer close(fd.read) // 关闭 fd.handleMsg() 协程

	for {
		revMsg, err := dp.Unpack(fd.conn)

		if err != nil {
			if err == io.EOF {
				log.Println("INFO: sever conn is EOF")
			} else {
				log.Println("ERROR: Unpack err", err)
			}
			fd.exit <- struct{}{}
			return
		}
		fd.read <- revMsg
	}
}

// 处理消息
func (fd *serverfd) handleMsg() {
	for msg := range fd.read {
		got := string(msg.GetValue())
		tag := msg.GetTag()
		switch tag {
		case 1:
			if sendData1 != got {
				log.Println("ERROR: tag 1 want", sendData1, "but got", got)
			}
		case 2:
			if sendData2 != got {
				log.Println("ERROR: tag 2 want", sendData2, "but got", got)
			}
		default:
			log.Println("ERROR: got message tag err", tag)
		}
	}
}

// 模拟粘包数据 TLV 打包
func msgTLV() []byte {
	// 配合 (fd *clientfd) PackMsg() 打包回复 tag 设 2
	sendMsg, err := dp.Pack(bean.NewMessage(2, []byte(sendData2)))
	if err != nil {
		log.Println("dp.Pack sendMsg err")
		return nil
	}
	sendMsg2, err := dp.Pack(bean.NewMessage(2, []byte(sendData2)))
	if err != nil {
		log.Println("dp.Pack sendMsg2 err")
		return nil
	}
	// 俩个消息粘在一起
	return append(sendMsg, sendMsg2...)
}
