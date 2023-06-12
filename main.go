package main

import (
	"os"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/server"
)

func main() {
	server.DefaultBootstrapServer(os.Args, common.EmptyExecContext())
}
