package main

import (
	"net"
	"net/http"
	"server/osinstallserver"
)

func main() {
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
