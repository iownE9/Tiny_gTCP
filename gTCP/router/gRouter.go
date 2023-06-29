package router

import (
	"fmt"
	"gTCP/api"
	"log"
	"sync"
)

// 将注册函数转换为该函数 实现 GHandler 接口
type HandlerFunc func(msg api.GMessage) api.GMessage

// 实现 GHandler 接口 换函数名
func (f HandlerFunc) HandlerTagMsg(msg api.GMessage) api.GMessage {
	return f(msg)
}

// =================

type gRouter struct {
	mu     sync.RWMutex
	router map[int]api.GHandler
}

// 全局变量
var GRouter *gRouter

func init() {
	GRouter = &gRouter{
		router: make(map[int]api.GHandler),
	}
	log.Println("GRouter init ok")
}

func AddHandleFunc(tag int, handler api.GHandler) {
	GRouter.AddHandleFunc(tag, handler)
}

func (r *gRouter) AddHandleFunc(tag int, handler api.GHandler) {
	// 写锁
	r.mu.Lock()
	defer r.mu.Unlock()

	// 0 是默认值 要排除
	if tag == 0 {
		panic("AddHandleFunc: tag is 0")
	}
	if handler == nil {
		panic("AddHandleFunc: nil handler")
	}

	r.router[tag] = handler
}

func (r *gRouter) getHandle(tag int) api.GHandler {
	// 读锁
	r.mu.RLock()
	defer r.mu.RUnlock()

	f, ok := r.router[tag]
	if !ok {
		panic(fmt.Sprintf("router[%d] is not install", tag))
	}
	return f
}

// 根据 tag 适配调用相应 Func
func (r *gRouter) HandlerTagMsg(msg api.GMessage) api.GMessage {
	tag := msg.GetTag()
	if tag == 0 {
		panic("handlerMsg: tag is 0")
	}
	f := r.getHandle(int(tag))

	return f.HandlerTagMsg(msg)
}
