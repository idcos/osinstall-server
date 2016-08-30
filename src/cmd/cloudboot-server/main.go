package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net"
	"net/http"
	"os"
	"server/osinstallserver"
	"server/osinstallserver/route"
	"server/osinstallserver/util"
	"time"
)

var date = time.Now().Format("2006-01-02")
var version = "v1.3.1 (" + date + ")"
var name = "cloudboot-server"
var description = "cloudboot server"
var usage = "CloudJ server install tool"
var configFile = "/etc/cloudboot-server/cloudboot-server.conf"

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: configFile,
			Usage: "config file",
		},
	}
	app.Action = func(c *cli.Context) {
		configFile = c.String("c")
		runServer(c)
	}

	app.Run(os.Args)
}

func runServer(c *cli.Context) {
	if !util.FileExist(configFile) {
		fmt.Println("The config file(" + configFile + ") is not exist!")
		return
	}

	srvr, err := osinstallserver.NewServer(configFile, osinstallserver.DevPipeline)
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

	//cron
	route.CloudBootCron(srvr.Conf, srvr.Log, srvr.Repo)

	if err := http.Serve(l4, srvr); err != nil {
		srvr.Log.Error(err)
		return
	}
}
