package handler

import (
	"gTCP/api"
	"gTCP/utils"
	"net"
)

var numPool = utils.GlobalConfig.GPool

// 缓存 clientfd 的管道
var clientfdChan = make(chan *net.TCPConn, numPool)

var clientfdPool = make([]api.GConnfd, numPool)

func init() {
	go initPool()
}

// 惰性启动
func initPool() {
	var i uint8
	for ; i < numPool; i++ {
		clientfdPool[i] = Clientfd(<-clientfdChan)
		// 正式开启 可复用 协程
		go clientfdPool[i].Closefd()
	}
}

func HandleClientfd(fd *net.TCPConn) {
	clientfdChan <- fd
}
