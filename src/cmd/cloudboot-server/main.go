package main

import (
	"build"
	"config"
	"config/jsonconf"
	"fmt"
	"logger"
	"net"
	"net/http"
	"os"
	"server/osinstallserver"
	"server/osinstallserver/route"
	"server/osinstallserver/util"

	"github.com/urfave/cli"
)

var date = "2017.01.05"
var version = "v1.4.1 (" + date + ")"
var name = "cloudboot-server"
var description = "cloudboot server"
var usage = "CloudJ server install tool"
var configFile = "/etc/cloudboot-server/cloudboot-server.conf"

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = build.Version("v1.3.1")

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: configFile,
			Usage: "config file",
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		configFile = c.String("c")
		if !util.FileExist(configFile) {
			return cli.NewExitError(fmt.Sprintf("The configuration file does not exist: %s", configFile), -1)
		}
		conf, err := jsonconf.New(configFile).Load()
		if err != nil {
			return cli.NewExitError(err.Error(), -1)
		}
		if err = runServer(conf); err != nil {
			return cli.NewExitError(err.Error(), -1)
		}
		return nil
	}

	app.Run(os.Args)
}

func runServer(conf *config.Config) (err error) {
	log := logger.NewBeeLogger(conf)

	srvr, err := osinstallserver.NewServer(log, conf, osinstallserver.DevPipeline)
	if err != nil {
		return err
	}

	port := 8083

	l4, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("The HTTP server is running at http://localhost:%d", port)

	//cron
	route.CloudBootCron(srvr.Conf, log, srvr.Repo)

	if err := http.Serve(l4, srvr); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
