package utils

import (
	"encoding/json"
	"log"
	"os"
)

var GlobalConfig *globalConfig

type globalConfig struct {
	HeadLen   uint32 // TLV 消息封装 头长度
	ChanCap   uint8  // channel 容量
	Name      string
	IpVersion string
	Host      string
	TcpPort   string
	LogPath   string
}

// 读取用户的配置文件
func (g *globalConfig) readConfig() error {
	data, err := os.ReadFile("../conf/gTCP-config.json")
	if err != nil {
		// 退出
		// log.Fatal(err)
		log.Println("gTCP/conf/gTCP-config.json", err)
		return err
	}

	// 将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Println("Unmarshal gTCP/conf/gTCP-config.json err", err)
		return err
	}

	return nil
}

func configInit() {
	GlobalConfig = &globalConfig{
		Name:      "gTCP",
		IpVersion: "tcp4",
		Host:      "localhost",
		TcpPort:   "8080",
		LogPath:   "../static/log/log_default.txt",
	}

	//从配置文件读取用户配置
	if err := GlobalConfig.readConfig(); err != nil {
		log.Println("gTCP/conf/gTCP-config.json err: use default args")
	}
}
