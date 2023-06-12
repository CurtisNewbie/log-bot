package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/log-bot/logbot"
)

func main() {
	server.OnServerBootstrapped(logbot.AfterServerBootstrapped)
	server.DefaultBootstrapServer(os.Args, common.EmptyExecContext())
}
