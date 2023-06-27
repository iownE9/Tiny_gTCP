package utils

import (
	"log"
	"os"
)

// 初始化 log 配置
func logInit() {
	filename := GlobalConfig.LogPath
	os.Remove(filename)

	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("open log file fail", err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
