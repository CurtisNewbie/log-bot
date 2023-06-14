package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/log-bot/logbot"
)

func main() {
	server.BeforeServerBootstrap(logbot.BeforeServerBootstrapp)
	server.OnServerBootstrapped(logbot.AfterServerBootstrapped)
	server.DefaultBootstrapServer(os.Args, common.EmptyExecContext())
}
