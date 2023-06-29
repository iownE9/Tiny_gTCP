package msg

import (
	"bytes"
	"encoding/binary"
	"gTCP/api"
	"gTCP/utils"
	"io"
	"log"
)

type gDataPack struct{}

func DataPack() api.GDataPack {
	return &gDataPack{}
}

// 包头长度 tag uint32(4字节) +  length uint32(4字节)
func (dp *gDataPack) HeadLen() uint32 {
	return utils.GlobalConfig.HeadLen
}

// 封包方法(封装数据) 大端法
func (dp *gDataPack) Pack(sendMsg api.GMessage) ([]byte, error) {
	// 创建一个存放 bytes 字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写 sendMsg tag
	if err := binary.Write(dataBuff, binary.BigEndian, sendMsg.GetTag()); err != nil {
		log.Println("binary.Write sendMsg tag err", err)
		return nil, err
	}

	// 写 sendMsg Length
	if err := binary.Write(dataBuff, binary.BigEndian, sendMsg.GetLength()); err != nil {
		log.Println("binary.Write data Length err", err)
		return nil, err
	}

	// 写 sendMsg value
	if err := binary.Write(dataBuff, binary.BigEndian, sendMsg.GetValue()); err != nil {
		log.Println("binary.Write data value err", err)
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法,只解压 head 的信息，得到 Length 和 tag 大端法
// func (dp *gDataPack) Unpack(conn *net.TCPConn) (api.GMessage, error) {
func (dp *gDataPack) Unpack(conn io.Reader) (api.GMessage, error) {
	// 读取 msg head
	headData := make([]byte, dp.HeadLen())
	if _, err := io.ReadFull(conn, headData); err != nil {
		log.Println("read msg head error", err)
		return nil, err
	}

	// 创建一个从输入二进制数据的 ioReader
	dataBuff := bytes.NewReader(headData)
	msg := &gMessage{}

	// 解析 msg tag
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.tag); err != nil {
		log.Println("binary.Read msg tag err", err)
		return nil, err
	}

	// 解析 msg length
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.length); err != nil {
		log.Println("binary.Read msg length err", err)
		return nil, err
	}

	// 读 msg value
	if msg.length > 0 {
		msg.value = make([]byte, msg.length)
		if n, err := io.ReadFull(conn, msg.value); err != nil {
			log.Println("io.ReadFull msg data value error", err, "n =", n)
			return nil, err
		}
	}

	return msg, nil
}
