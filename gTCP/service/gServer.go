package service

import (
	"gTCP/handler"
	"gTCP/utils"
	"log"
	"net"
)

type GServer struct {
	ipVersion string
	ip        string
	port      string
	name      string
}

// 启动端口监听，处理链接
func (s *GServer) ListenAndServe() error {
	log.Printf("[%s] listenner at IP: %s, Port: %s, is starting\n", s.name, s.ip, s.port)

	// 获取套接字地址
	addr, err := net.ResolveTCPAddr(s.ipVersion, (s.ip + ":" + s.port))
	if err != nil {
		log.Println("ResolveTCPAddr err", err)
		return err
	}

	// 获取监听描述符
	listener, err := net.ListenTCP(s.ipVersion, addr)
	if err != nil {
		log.Println("ListenTCP  err", err)
		return err
	}

	log.Println("start", s.name, "server succ, listenning...")

	// 处理链接请求
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("listenfd.AcceptTCP() err", err)
			continue
		}

		// 处理连接
		go handler.HandlerEchoTLVMsg(conn)
	}
}

func ListenAndServe(name string) error {
	if name != "" {
		utils.GlobalConfig.Name = name
	}
	server := &GServer{
		ipVersion: utils.GlobalConfig.IpVersion,
		ip:        utils.GlobalConfig.Host,
		port:      utils.GlobalConfig.TcpPort,
		name:      utils.GlobalConfig.Name,
	}
	return server.ListenAndServe()
}
