package test

import (
	"gTCP/utils"
	"log"
	"net"
)

type gClient struct {
	network string
	address string
}

var gc = &gClient{
	network: utils.GlobalConfig.IpVersion,
	address: utils.GlobalConfig.Host + ":" +
		utils.GlobalConfig.TcpPort,
}

// 与客户端建立连接
func ClientDialTest() (net.Conn, error) {
	// 与服务器建立连接
	conn, err := net.Dial(gc.network, gc.address)
	if err != nil {
		log.Println("ERROR: net.Dial(", gc.network, gc.address, ") err")
		return nil, err
	}
	return conn, nil
}
