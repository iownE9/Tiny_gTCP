package bean

import (
	"gTCP/api"
)

type gMessage struct {
	tag    uint32 // 消息的 tag
	length uint32 // 消息的长度
	value  []byte // 消息的内容
}

func NewMessage(tag uint32, value []byte) api.GMessage {
	return &gMessage{
		tag:    tag,
		length: uint32(len(value)),
		value:  value,
	}
}

// 获取消息数据段长度 length
func (msg *gMessage) GetLength() uint32 {
	return msg.length
}

// 获取消息 Tag
func (msg *gMessage) GetTag() uint32 {
	return msg.tag
}

// 获取消息内容 value
func (msg *gMessage) GetValue() []byte {
	return msg.value
}

// 设置消息 tag
func (msg *gMessage) SetTag(tag uint32) {
	msg.tag = tag
}

// 设置消息数据段长度 length
func (msg *gMessage) SetLength(len uint32) {
	msg.length = len
}

// 设置消息内容 value
func (msg *gMessage) SetValue(vale []byte) {
	msg.value = vale
}
