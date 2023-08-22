package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/log-bot/logbot"
)

func main() {
	server.PreServerBootstrap(logbot.BeforeServerBootstrapp)
	server.PostServerBootstrapped(logbot.AfterServerBootstrapped)
	server.BootstrapServer(os.Args)
}
