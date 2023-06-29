package test

import (
	"gTCP/api"
	"gTCP/msg"
	"gTCP/utils"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// 对应服务器对象
type serverfd struct {
	conn       net.Conn // serverfd 文件描述符
	read       chan api.GMessage
	pack       chan api.GMessage
	send       chan []byte
	exit       chan struct{}
	exitSendOk chan struct{}
	exitReadOk chan struct{}
}

// TLV 拆包装包 句柄
var dp api.GDataPack = msg.DataPack()

// channel 容量
var chanCap uint8 = utils.GlobalConfig.ChanCap

// 待发送消息
var msgsValue []string = []string{
	"hello gTCP v1",
	"hello gTCP v2",
	"hello gTCP v3",
	"hello gTCP v4",
	"hello gTCP v5",
}
var wants []string = []string{
	"",
	"hello gTCP v5",
	"hello gTCP v4",
	"hello gTCP v3",
	"hello gTCP v2",
	"hello gTCP v1",
}

// 待接收消息总数 -> 5 个 模拟数据 粘包 发两遍 -> 20 需要接收 20 个
var totalRead int = len(msgsValue) * 4

// +++++++++++++++
var HandleMsgRouter func(api.GMessage) api.GMessage

// 处理 msg router 方便测试替换
// gHandleConnfd_test.go
// var HandleMsgRouter func(api.GMessage) api.GMessage = CustomHandleMsgRouter

// gRouter_test.go
// var HandleMsgRouter func(api.GMessage) api.GMessage = CustomHandleMsg

// +++++++++++++++

// 新增 匹配回复的消息
func CustomHandleMsgRouter(msg api.GMessage) api.GMessage {
	tag := msg.GetTag()
	got := string(msg.GetValue())

	if got != wants[tag] {
		log.Println("ERROR: want", wants[tag], ", but got", got)
	}
	return nil
}

// 新增 匹配回复的消息
func CustomHandleMsg(msg api.GMessage) api.GMessage {
	tag := msg.GetTag() - 1
	got := string(msg.GetValue())

	if got != msgsValue[tag] {
		log.Println("ERROR: want", msgsValue[tag], ", but got", got)
	}
	return nil
}

func Serverfd(conn net.Conn) api.GConnfd {
	// 描述符实例对象
	fd := &serverfd{
		conn: conn,
		read: make(chan api.GMessage, chanCap),
		pack: make(chan api.GMessage, chanCap),
		send: make(chan []byte, chanCap),
		// 容量为 >=2， 读写都有权限通知关闭 connfd
		// +1 消息发送完 接收完 后 主动关闭
		exit:       make(chan struct{}, chanCap),
		exitSendOk: make(chan struct{}),
		exitReadOk: make(chan struct{}),
	}
	return fd
}

// 处理 Serverfd
func HandleServerfd(conn net.Conn) {
	fd := Serverfd(conn)
	fd.Closefd()
}

// 接收消息  关闭 serverfd 后 退出
func (fd *serverfd) ReadMsg() {
	defer close(fd.read) // 关闭 fd.handleMsg 协程

	for total := totalRead; ; {
		readMsg, err := dp.Unpack(fd.conn)
		if err != nil {
			if err == io.EOF {
				log.Println("INFO: conn is EOF")
			} else {
				log.Println("ERROR: Unpack err", err)
			}
			fd.exit <- struct{}{} //异常 退出
			break
		}
		fd.read <- readMsg
		total--
		if total == 0 {
			close(fd.exitReadOk) // 读取完全
			break
		}
	}
}

// 处理消息 修改
func (fd *serverfd) HandleMsg() {
	if HandleMsgRouter == nil {
		panic("please install HandleMsgRouter")
	}
	for msg := range fd.read {
		go HandleMsgRouter(msg)
	}
}

// 生成消息 新增
func (fd *serverfd) TLVMsgs() {
	defer close(fd.pack) // 关闭 fd.packMsg 协程

	var tag uint32 = 1
	n := len(msgsValue)
	var testMsg = make([]api.GMessage, n)
	for i := 0; i < n; i++ {
		testMsg[i] = msg.NewMessage(tag, []byte(msgsValue[i]))
		tag++
	}
	tag = 1

	for _, msg := range testMsg {
		fd.pack <- msg
	}

	time.Sleep(100 * time.Millisecond)

	for _, msg := range testMsg {
		fd.pack <- msg
	}
}

// 打包回复
func (fd *serverfd) PackMsg() {
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

// 发送消息 新增粘包数据
func (fd *serverfd) SendMsg() {
	// 新增粘包数据 一定要这样写 适配 Route_test 和 Connfd_test
	catData, _ := dp.Pack(msg.NewMessage(3, []byte(msgsValue[2])))

	for sendData := range fd.send {
		sendData = append(sendData, catData...) // 新增粘包数据
		_, err := fd.conn.Write(sendData)
		if err != nil {
			log.Println("ERROR: conn.Write err", err)
			fd.exit <- struct{}{} //异常 退出
			break
		}
	}

	// 1. 因异常跳出上面 forrange 消耗掉未发送的 msg
	// 2. fd.pack 因关闭而执行到此 模拟的消息已发送完
	for range fd.send {
	}

	close(fd.exitSendOk) // 发送完全
}

// 连接关闭
func (fd *serverfd) Closefd() {

	// 读 msg
	go fd.ReadMsg()

	// 处理 msg
	go fd.HandleMsg()

	// 生成消息 add
	go fd.TLVMsgs()

	// 打包 msg
	go fd.PackMsg()

	// 发送 msg
	go fd.SendMsg()

	// 读写完全
	go fd.exitOk()

	// 只执行一次
	select {
	case <-fd.exit:
		fd.conn.Close() // 关闭 serverfd
		// close(fd.exit) // ERROR
		// 让 容量  去存放 ReadMsg ReadMsg 中慢的一方
	}
}

func (fd *serverfd) exitOk() {
	_ = <-fd.exitSendOk
	_ = <-fd.exitReadOk
	fd.exit <- struct{}{} // 读写完全 主动 退出
}
