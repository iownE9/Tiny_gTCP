package api

// 注册函数接口
type GHandler interface {
	HandlerTagMsg(msg GMessage) GMessage
}
