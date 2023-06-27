package handler

import (
	"io"
	"log"
	"net"
)

// Echo 回写
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
