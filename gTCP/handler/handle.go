package handler

import (
	"gTCP/api"
	"gTCP/bean"
	"io"
	"log"
	"net"
)

// 处理 connfd V3
func HandlerCurr(conn *net.TCPConn) {
	fd := bean.Clientfd(conn)

	// 读 msg
	go fd.ReadMsg()

	// 处理 msg
	go fd.HandleMsg()

	// 打包 msg
	go fd.PackMsg()

	// 发送 msg
	go fd.SendMsg()

	// 连接关闭
	go fd.Closefd()
}

// 对TLV消息回写 v2
func HandlerEchoTLVMsg(conn *net.TCPConn) {
	defer conn.Close()
	// TLV 拆包装包 句柄
	var dp api.GDataPack = bean.DataPack()

	for {
		// 读 msg
		getMsg, err := dp.Unpack(conn)
		if err == io.EOF {
			log.Println("conn is closed: EOF")
			return
		}
		if err != nil {
			log.Println("handlerEchoMsg continue")
			continue
		}

		// 处理 msg
		log.Println("server receive:", getMsg.GetTag(), string(getMsg.GetValue()))

		// 打包 msg
		sedMsg, err := dp.Pack(getMsg)
		if err != nil {
			log.Println("handlerEchoMsg continue")
			continue
		}

		// 发送 msg
		if _, err := conn.Write(sedMsg); err != nil {
			log.Println("conn.Write(sedMsg) err", err)
			return
		}

		log.Println("server echo msg ok")
	}
}

// Echo 回写 v1
func Echo(conn *net.TCPConn) {
	defer conn.Close()

	for {
		buf := make([]byte, 512)

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("clientfd is closed: EOF")
			} else {
				log.Println("conn.Read err", err)
			}
			return
		}

		log.Println("server receive:", string(buf[:n]))

		if _, err := conn.Write(buf[:n]); err != nil {
			log.Println("conn.Write(buf[:n]) err", err)
			return
		}

		log.Println("server echo ok")
	}
}
