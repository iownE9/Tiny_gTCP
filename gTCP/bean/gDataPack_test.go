package bean

import (
	"gTCP/api"
	"io"
	"strings"
	"testing"
)

// go test -v -run=TestDataPack
func TestDataPack(t *testing.T) {
	var wants = []string{
		"sendData1111",
		"sendData2222",
	}
	var tests = []api.GMessage{
		NewMessage(0, []byte(wants[0])),
		NewMessage(1, []byte(wants[1])),
	}
	var dp api.GDataPack = DataPack()
	var r io.Reader

	for i, test := range tests {
		bytes, err := dp.Pack(test)
		if err != nil {
			t.Error("dp.Pack(test)", err)
		}

		r = strings.NewReader(string(bytes))

		msg, err := dp.Unpack(r)
		if err != nil {
			t.Error("dp.Unpack(r)", err)
		}
		got := string(msg.GetValue())
		if got != wants[i] {
			t.Error("want", wants[i], "but got", got)
		}
	}
}
