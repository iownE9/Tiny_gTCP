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
	router map[uint32]api.GHandler
}

// 全局变量
var GRouter *gRouter

func init() {
	GRouter = &gRouter{
		router: make(map[uint32]api.GHandler),
	}
	log.Println("GRouter init ok")
}

func AddHandleFunc(tag uint32, handler api.GHandler) {
	GRouter.AddHandleFunc(tag, handler)
}

func (r *gRouter) AddHandleFunc(tag uint32, handler api.GHandler) {
	// 写锁
	r.mu.Lock()
	defer r.mu.Unlock()
	if tag == 0 {
		panic("tag: 0 is invalid")
	}
	if handler == nil {
		panic("AddHandle: nil handler")
	}

	r.router[tag] = handler
}

func (r *gRouter) getHandle(tag uint32) api.GHandler {
	// 读锁
	r.mu.RLock()
	defer r.mu.RUnlock()

	f, ok := r.router[tag]
	if !ok {
		panic(fmt.Sprintf("router[%d] is not install", tag))
	}
	return f
}

func (r *gRouter) HandlerTagMsg(msg api.GMessage) api.GMessage {
	tag := msg.GetTag()
	f := r.getHandle(tag)

	return f.HandlerTagMsg(msg)
}
