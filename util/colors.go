package util

var (
	nc            = "\033[0m"
	Brightblack   = "\033[1;30m"
	Brightred     = "\033[1;31m"
	Brightgreen   = "\033[1;32m"
	Brightyellow  = "\033[1;33m"
	Brightpurple  = "\033[1;34m"
	Brightmagenta = "\033[1;35m"
	Brightcyan    = "\033[1;36m"
	Brightwhite   = "\033[1;37m"
	black   = "\033[0;30m"
	red     = "\033[0;31m"
	green   = "\033[0;32m"
	yellow  = "\033[0;33m"
	purple  = "\033[0;34m"
	magenta = "\033[0;35m"
	cyan    = "\033[0;36m"
	white   = "\033[0;37m"
	Rcv     = Brightyellow + "[RECV]" + Brightred
	Send    = cyan + "[SEND]" + Brightred
)
