package main

import (
	config "github.com/karai/go-karai/configuration"
	"github.com/karai/go-karai/network"
	"github.com/karai/go-karai/util/flags"
	"github.com/karai/go-karai/util/menu"
)

func main() {
	c := config.InitConfig()
	flags.Flags(&c)
	var s network.Server
	go network.ProtocolInit(&c, &s)
	menu.InputHandler(&s)
}