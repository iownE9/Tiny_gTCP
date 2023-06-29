package msg

import (
	"gTCP/api"
	"io"
	"strings"
	"testing"
)

// go test -v -run=TestDataPack gTCP/bean
func TestDataPack(t *testing.T) {
	var wants = []string{
		"sendData1111",
		"sendData2222",
		"发送消息2222",
		"发送消息2112",
		"发送消息211&……￥#@@@！T（（2",
	}

	var tests = make([]api.GMessage, len(wants))
	for i := len(wants) - 1; i >= 0; i-- {
		tests[i] = NewMessage(uint32(i), []byte(wants[i]))
	}

	var dp api.GDataPack = DataPack()
	var r io.Reader

	for i, test := range tests {
		if test == nil {
			t.Error(i, "test == nil")
			continue
		}
		bytes, err := dp.Pack(test)
		if err != nil {
			t.Error("dp.Pack(test)", err)
			continue
		}

		r = strings.NewReader(string(bytes))

		msg, err := dp.Unpack(r)
		if err != nil {
			t.Error("dp.Unpack(r)", err)
			continue
		}
		got := string(msg.GetValue())
		if got != wants[i] {
			t.Error("want", wants[i], "but got", got)
		}
	}
}

// go test -v -bench=Datapack gTCP/msg
func BenchmarkDatapack(b *testing.B) {
	var dp api.GDataPack = DataPack()
	testdata := NewMessage(0, []byte("发送消息211&……￥#@@@！T（（2"))
	
	bytes, _ := dp.Pack(testdata)
	r := strings.NewReader(string(bytes))

	// 目的是测性能 忽略错误处理
	for i := 0; i < b.N; i++ {
		dp.Pack(testdata)

		dp.Unpack(r)
	}
}
