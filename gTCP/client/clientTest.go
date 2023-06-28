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

// 对应连接对象
type serverfd struct {
	conn net.Conn          // serverfd 文件描述符
	send chan []byte       // 发送消息管道
	read chan api.GMessage // 接收消息管道
	exit chan struct{}     // 关闭 fd，退出协程
}

// TLV 拆包装包 句柄
var dp api.GDataPack
var chanCap uint8       // channel 容量
var msgs []api.GMessage // 发送消息

// 初始化 client 包 全局只读变量
func init() {
	dp = bean.DataPack()
	chanCap = utils.GlobalConfig.ChanCap
	msgs = []api.GMessage{
		bean.NewMessage(1, []byte("hello gTCP v1")),
		bean.NewMessage(2, []byte("hello gTCP v2")),
		bean.NewMessage(3, []byte("hello gTCP v3")),
		bean.NewMessage(4, []byte("hello gTCP v4")),
		bean.NewMessage(5, []byte("hello gTCP v5")),
	}
	log.Println("ClientTest Init ok")
}

// 模拟 client TLV 格式消息发送
func ClientTLVTest() {

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

	// 打包消息
	go fd.msgTLVpack()

	go func() {
		time.Sleep(1 * time.Second)
		select {
		case <-fd.exit:
			// 若已经关闭，直接退出
		default:
			fd.msgTLVpack()
		}

		time.Sleep(1 * time.Second)
		// 主动结束
		fd.exit <- struct{}{}
	}()

	// 发送消息
	go fd.sendMsg()

	// 接收消息
	go fd.readMsg()

	// 处理消息
	go fd.handleMsg()

	// 一直阻塞直到连接关闭
	select {
	case <-fd.exit:
		fd.conn.Close() // 关闭 serverfd
		return
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
		want := string(msgs[tag-1].GetValue())
		if want != got {
			log.Println("ERROR: tag", tag, "want", want, ",but got", got)
		}
	}
}

// 模拟粘包数据 TLV 打包
func (fd *serverfd) msgTLVpack() {
	for _, msg := range msgs {
		go func(msg api.GMessage) {
			sendData, err := dp.Pack(msg)
			if err != nil {
				log.Println("ERROR: msgTLVpack() err")
			} else {
				fd.send <- sendData
			}
		}(msg)
	}
}

// 发送消息
func (fd *serverfd) sendMsg() {
	for sendData := range fd.send {
		_, err := fd.conn.Write(sendData)
		if err != nil {
			log.Println("INFO: conn.Write err", err)
			fd.exit <- struct{}{}
			break
		}
	}

	// 消耗掉未发送的数据，防止泄露
	for range fd.send {
		// 连接已关闭
	}
}
