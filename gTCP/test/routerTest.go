package test

import (
	"gTCP/api"
	"gTCP/msg"
)

var Msgs = []api.GMessage{
	msg.NewMessage(0, []byte("")), // 占位
	msg.NewMessage(1, []byte("hello gTCP v1")),
	msg.NewMessage(2, []byte("hello gTCP v2")),
	msg.NewMessage(3, []byte("hello gTCP v3")),
	msg.NewMessage(4, []byte("hello gTCP v4")),
	msg.NewMessage(5, []byte("hello gTCP v5")),
}

func F01(msg api.GMessage) api.GMessage {
	return Msgs[1]
}
func F02(msg api.GMessage) api.GMessage {
	return Msgs[2]
}
func F03(msg api.GMessage) api.GMessage {
	return Msgs[3]
}
func F04(msg api.GMessage) api.GMessage {
	return Msgs[4]
}
func F05(msg api.GMessage) api.GMessage {
	return Msgs[5]
}

// // 装配 handlerFunc
// 	router.AddHandleFunc(1, router.HandlerFunc(f01))
// 	router.AddHandleFunc(2, router.HandlerFunc(f02))
// 	router.AddHandleFunc(3, router.HandlerFunc(f03))
// 	router.AddHandleFunc(4, router.HandlerFunc(f04))
// 	router.AddHandleFunc(5, router.HandlerFunc(f05))
