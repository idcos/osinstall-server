package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/urfave/cli"

	"idcos.io/osinstall/build"
	"idcos.io/osinstall/config"
	"idcos.io/osinstall/config/jsonconf"
	"idcos.io/osinstall/logger"
	"idcos.io/osinstall/server/osinstallserver"
	"idcos.io/osinstall/server/osinstallserver/util"
)

var configFile = "/etc/cloudboot-server/cloudboot-server.conf"

func main() {
	app := cli.NewApp()
	app.Name = "cloudboot-server"
	app.Description = "cloudboot server"
	app.Usage = "CloudJ server install tool"
	app.Version = build.Version()

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
	// remove sql upgrade
	//route.CloudBootCron(srvr.Conf, log, srvr.Repo)

	if err := http.Serve(l4, srvr); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
