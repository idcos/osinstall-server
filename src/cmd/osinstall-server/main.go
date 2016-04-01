package main

import (
	"github.com/codegangsta/cli"
	"net"
	"net/http"
	"os"
	"server/osinstallserver"
)

var version = "v1.2.1 (2016-03-31)"
var name = "osinstall-server"
var description = "osinstall server"
var usage = "CloudJ X86 server install tool"

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	app.Action = func(c *cli.Context) {
		runServer(c)
	}

	app.Run(os.Args)
}

func runServer(c *cli.Context) {
	srvr, err := osinstallserver.NewServer("/etc/osinstall-server/osinstall-server.conf", osinstallserver.DevPipeline)
	if err != nil {
		srvr.Log.Error(err)
		return
	}

	addr := ":8083"

	l4, err := net.Listen("tcp4", addr)
	if err != nil {
		srvr.Log.Error(err)
		return
	}

	srvr.Log.Info("The server is running.")

	if err := http.Serve(l4, srvr); err != nil {
		srvr.Log.Error(err)
		return
	}
}
