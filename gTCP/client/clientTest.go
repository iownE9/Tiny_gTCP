package client

import (
	"gTCP/api"
	"gTCP/bean"
	"gTCP/utils"
	"io"
	"log"
	"net"
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
		exit: make(chan struct{}),
	}

	fd.send <- sendData // 必须有缓存，不然就一直阻塞

	// 处理消息
	go fd.handleMsg()

	for {
		// fd 限制在一个 goroutine 内，虽然它本身就并发安全
		select {
		// 连接关闭
		case <-fd.exit:
			fd.conn.Close() // 关闭 serverfd
			close(fd.read)  // 关闭 fd.handleMsg() 协程

			// 消耗掉未发送的
			close(fd.send)
			for range fd.send {
			}

			return // 退出 clientTest

		// 发送消息
		case sendData := <-fd.send:
			_, err := fd.conn.Write(sendData)
			if err != nil {
				log.Println("conn.Write err", err)
				close(fd.exit)
			}

		// 接收消息
		default:
			// bug: 若无回复内容就会一直阻塞
			// 需确保 fd.send <- sendData 先运行 故：
			// go func() { fd.send <- sendData }() // Error
			revMsg, err := dp.Unpack(fd.conn)
			if err != nil {
				if err == io.EOF {
					log.Println("sever conn is EOF")
				} else {
					log.Println("Unpack err", err)
				}
				close(fd.exit)
			} else {
				fd.read <- revMsg
			}
		}
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
				log.Println("ERROR: want", sendData1, "but got", got)
			}
		case 2:
			if sendData2 != got {
				log.Println("ERROR: want", sendData1, "but got", got)
			}
		default:
			log.Println("ERROR: got message tag err", tag)
		}
	}
}

// 模拟粘包数据 TLV 打包
func msgTLV() []byte {
	sendMsg, err := dp.Pack(bean.NewMessage(1, []byte(sendData1)))
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
