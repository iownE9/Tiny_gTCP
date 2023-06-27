package api

// 将请求的一个消息封装到 GMessage中，定义抽象层接口
type GMessage interface {
	GetTag() uint32    // 获取消息 Tag
	GetLength() uint32 // 获取消息数据段长度
	GetValue() []byte  // 获取消息内容

	SetTag(uint32)    // 设计消息 Tag
	SetLength(uint32) // 设置消息数据段长度
	SetValue([]byte)  // 设计消息内容
}
