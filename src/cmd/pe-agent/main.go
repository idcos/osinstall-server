package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"utils"

	"github.com/codegangsta/cli"
)

var xmlPath = "X:\\Windows\\System32\\unattended.xml"
var rootPath = "X:\\Windows\\System32"
var scriptFile = path.Join(rootPath, "temp-script.cmd")

func main() {

	app := cli.NewApp()
	app.Version = "v2016.02.23"
	app.Action = func(c *cli.Context) {
		if err := run(c); err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	utils.InitConsoleLog()

	var sn = getSN()
	// sn = "214245856"
	if sn == "" {
		return errors.New("get sn failed")
	}

	if !utils.PingLoop("osinstall.", 30, 2) {
		return errors.New("ping timeout")
	}

	if err := getXmlFile(sn); err != nil {
		return err
	}

	utils.ReportProgress(0.6, sn, "开始安装", "start install")

	return nil
}

// 查看本机 SN
func getSN() string {
	var cmd = `wmic bios get SerialNumber /VALUE`
	var r = `SerialNumber=(.+)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecScript(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 2 {
		return ""
	}

	// fmt.Println(strings.Trim(regResult[1], "\r\n"))
	return strings.Trim(regResult[1], "\r\n")
}

func getXmlFile(sn string) error {
	var url = fmt.Sprintf("http://osinstall./api/osinstall/v1/device/getSystemBySn?sn=%s",
		sn)

	resp, err := http.Get(url)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("call url: %s failed", url)
	}

	return ioutil.WriteFile(xmlPath, body, 0666)
}
