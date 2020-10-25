package main

import (
	config "github.com/karai/go-karai/configuration"
	"github.com/karai/go-karai/network"
)

// Hello Karai
func main() {
	c := config.Config_Init()
	flags(&c)
	var s network.Server
	go network.ProtocolInit(&c, &s)
	inputHandler(&s)
}