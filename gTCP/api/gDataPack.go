package api

import "io"

type GDataPack interface {
	HeadLen() uint32                         // 获取包头长度方法
	Pack(sendMsg GMessage) ([]byte, error)   // 封包方法
	Unpack(conn io.Reader) (GMessage, error) // 拆包方法
}
