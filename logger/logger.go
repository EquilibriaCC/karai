package logger

import (
	config "github.com/karai/go-karai/configuration"
	"log"
)

var Brightblack = "\033[1;30m"
var Brightred = "\033[1;31m"
var Brightgreen = "\033[1;32m"
var Brightyellow = "\033[1;33m"
var Brightpurple = "\033[1;34m"
var Brightmagenta = "\033[1;35m"
var Brightcyan = "\033[1;36m"
var Brightwhite = "\033[1;37m"
var black = "\033[0;30m"
var red = "\033[0;31m"
var green = "\033[0;32m"
var yellow = "\033[0;33m"
var purple = "\033[0;34m"
var magenta = "\033[0;35m"
var cyan = "\033[0;36m"
var white = "\033[0;37m"

func Receive(msg string) {
	if config.InfoLogging {
		log.Println(Brightgreen + "[INFO]" + Brightyellow + " [RECEIVE]" + white + msg)
	}
}

func Send(msg string) {
	if config.InfoLogging {
		log.Println(Brightgreen + "[INFO]" + cyan + " [SEND]" + white + msg)
	}
}

func Info(msg string) {
	if config.InfoLogging {
		log.Println(Brightgreen + "[INFO]" + white + msg + white)
	}
}

func Error(msg string) {
	if config.ErrorLogging {
		log.Println(Brightred + "[ERROR]" + white + msg + white)
	}
}

func Warning(msg string) {
	if config.WarningLogging {
		log.Println(Brightyellow + "[WARNING]" + white + msg + white)
	}
}

//func Success_log_array(msg string) {
//	log.Print(Brightgreen + msg + white)
//}
