package router

// 避免 包 循环导入

import (
	"gTCP/api"
	"gTCP/msg"
	"sync"
	"testing"
)

func f01(msg api.GMessage) api.GMessage {
	return msgs[0]
}
func f02(msg api.GMessage) api.GMessage {
	return msgs[1]
}
func f03(msg api.GMessage) api.GMessage {
	return msgs[2]
}

var msgs = []api.GMessage{
	msg.NewMessage(1, []byte("hello gTCP v1")),
	msg.NewMessage(2, []byte("hello gTCP v2")),
	msg.NewMessage(3, []byte("hello gTCP v3")),
}

// go test -run=TestGRouter -v -race gTCP/router
func TestGRouter(t *testing.T) {

	// 装配 handlerFunc
	AddHandleFunc(1, HandlerFunc(f01))
	AddHandleFunc(2, HandlerFunc(f02))
	AddHandleFunc(3, HandlerFunc(f03))

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		currRead(t)
	}()

	go func() {
		defer wg.Done()
		currRead(t)
	}()

	wg.Wait()
}

func currRead(t *testing.T) {
	for _, msg := range msgs {
		getMsg := GRouter.HandlerTagMsg(msg)
		
		got := string(getMsg.GetValue())
		tag := getMsg.GetTag()
		want := string(msg.GetValue())

		if want != got {
			t.Error("ERROR: tag", tag, "want", want, ",but got", got)
		}
	}
}
