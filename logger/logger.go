package logger

import "log"

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
	log.Println(Brightyellow + "[RECEIVE]" + white)
}

func Send(msg string) {
	log.Println(cyan + "[SEND]" + white)
}

func Success_log(msg string) {
	log.Println(Brightgreen + "[SUCCESS]" + msg + white)
}

func Info(msg string) {
	log.Println(Brightgreen + "[INFO]" + msg + white)
}

func Error_log(msg string) {
	log.Println(Brightred + "[ERROR]" + white + msg + white)
}

func Warning_log(msg string) {
	log.Println(Brightyellow + "[WARNING]" + white + msg + white)
}

func Success_log_array(msg string) {
	log.Print(Brightgreen + msg + white)
}