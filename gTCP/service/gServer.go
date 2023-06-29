package service

import (
	"gTCP/handler"
	"gTCP/utils"
	"log"
	"net"
)

type gServer struct {
	ipVersion string
	ip        string
	port      string
	name      string
}

var gTCPServer *gServer

func init() {
	gTCPServer = &gServer{
		ipVersion: utils.GlobalConfig.IpVersion,
		ip:        utils.GlobalConfig.Host,
		port:      utils.GlobalConfig.TcpPort,
		name:      utils.GlobalConfig.Name,
	}
}

// 处理客户端连接 方便测试替换
var handleClientFunc = handler.HandleClientfd

// 监听并处理请求
func ListenAndServe() error {
	listener, err := gTCPServer.listen()
	if err != nil {
		return err
	}

	go gTCPServer.serve(listener)
	return nil
}

// 启动端口监听
func (s *gServer) listen() (*net.TCPListener, error) {
	// 获取套接字地址
	addr, err := net.ResolveTCPAddr(s.ipVersion, (s.ip + ":" + s.port))
	if err != nil {
		log.Println("ERROR: ResolveTCPAddr() err")
		return nil, err
	}

	// 获取监听描述符
	listener, err := net.ListenTCP(s.ipVersion, addr)
	if err != nil {
		log.Println("ERROR: ListenTCP() err")
		return nil, err
	}

	return listener, nil
}

// 处理处理请求
func (s *gServer) serve(listener *net.TCPListener) {
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("listenfd.AcceptTCP() err", err)
			continue
		}

		// 处理连接
		go func(conn *net.TCPConn) { handleClientFunc(conn) }(conn)
	}
}
