package utils

import (
	"log"
)

func init() {
	configInit()
	log.Println("configInit init ok")

	logInit()
	log.Println("log file init ok")

	log.Println("global init ok")
}
