package router_test

// 避免 包 循环导入

import (
	"gTCP/api"
	"gTCP/bean"
	"gTCP/router"
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
	bean.NewMessage(1, []byte("hello gTCP v1")),
	bean.NewMessage(2, []byte("hello gTCP v2")),
	bean.NewMessage(3, []byte("hello gTCP v3")),
}

// go test -run=TestGRouter -v -race gTCP/router
func TestGRouter(t *testing.T) {

	// 装配 handlerFunc
	router.AddHandleFunc(1, router.HandlerFunc(f01))
	router.AddHandleFunc(2, router.HandlerFunc(f02))
	router.AddHandleFunc(3, router.HandlerFunc(f03))

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
		getMsg := router.GRouter.HandlerTagMsg(msg)
		got := string(getMsg.GetValue())
		tag := getMsg.GetTag()
		want := string(msg.GetValue())

		if want != got {
			t.Error("ERROR: tag", tag, "want", want, ",but got", got)
		}
	}
}
