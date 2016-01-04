package main

import (
	"fmt"
	"net"
	"net/http"
	"server/osinstallserver"
)

func main() {
	srvr, err := osinstallserver.NewServer("idcos-os-install.json", osinstallserver.DevPipeline)
	if err != nil {
		fmt.Println(err)
		return
	}

	addr := ":8083"

	l4, err := net.Listen("tcp4", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := http.Serve(l4, srvr); err != nil {
		fmt.Println(err)
		return
	}
}
