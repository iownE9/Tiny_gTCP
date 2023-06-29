package api

type GConnfd interface {
	ReadMsg()   // 接收消息
	HandleMsg() // 处理消息
	PackMsg()   // 打包回复
	SendMsg()   // 发送消息
	Closefd()   // 关闭连接
}
