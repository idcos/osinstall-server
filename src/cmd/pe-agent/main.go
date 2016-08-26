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
	"time"
	"utils"

	"github.com/codegangsta/cli"
)

var xmlPath = "X:\\Windows\\System32\\unattended.xml"
var rootPath = "X:\\Windows\\System32"
var scriptFile = path.Join(rootPath, "temp-script.cmd")
var serverHost = "osinstall" //cloudboot server host

var date = time.Now().Format("2006-01-02")
var version = "v1.3.1 (" + date + ")"

func main() {

	app := cli.NewApp()
	app.Version = version
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
	isVm := isVirtualMachine()
	if isVm {
		sn = getMacAddress()
	}

	if sn == "" {
		return errors.New("get sn failed")
	}

	//init cloudboot server host
	serverIp := getDomainLookupIP(serverHost)
	if serverIp != "" {
		serverHost = serverIp
	}

	if !utils.PingLoop(serverHost, 30, 2) {
		return errors.New("ping timeout")
	}

	// get xml for install windows
	if err := getXmlFile(sn, serverHost); err != nil {
		return err
	}

	// mount samba
	if err := mountSamba(); err != nil {
		return err
	}

	//load drive
	loadDrive()

	utils.ReportProgress(0.6, sn, "开始安装", "start install", serverHost)
	// install windows
	if err := installWindows(); err != nil {
		return err
	}

	// xcopy firstboot
	if err := copyFirstBoot(); err != nil {
		return err
	}

	// reboot
	return reboot()
}

// 查看本机 SN
func getSN() string {
	var cmd = `wmic bios get SerialNumber /VALUE`
	var r = `SerialNumber=(\S+)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
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
	var result string
	result = strings.Trim(regResult[1], "\r\n")
	result = strings.TrimSpace(result)
	return result
}

//是否是虚拟机
func isVirtualMachine() bool {
	var cmd = `systeminfo`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	isValidate, err := regexp.MatchString(`(?i)VMware|VirtualBox|KVM|Xen|Parallels`, output)
	if err != nil {
		utils.Logger.Error(err.Error())
		return false
	}

	if isValidate {
		return true
	} else {
		return false
	}
}

// 获取Mac地址
func getMacAddress() string {
	var cmd = `wmic nic where "NetConnectionStatus=2" get MACAddress /VALUE`
	var r = `(?i)MACAddress=(\S+)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
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

	var result string
	result = strings.Trim(regResult[1], "\r\n")
	result = strings.TrimSpace(result)
	return result
}

//get domain's lookup ip
func getDomainLookupIP(domain string) string {
	var cmd = `ping ` + domain
	var r = `(.+)(\s)(\d+)\.(\d+)\.(\d+)\.(\d+)([:|\s])(.+)TTL`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 9 {
		return ""
	}

	var result = fmt.Sprintf("%s.%s.%s.%s", strings.TrimSpace(regResult[3]),
		strings.TrimSpace(regResult[4]),
		strings.TrimSpace(regResult[5]),
		strings.TrimSpace(regResult[6]))
	return result
}

func loadDrive() {
	var cmd = `Z:`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
	}
	utils.Logger.Info(output)
	//load drive
	dirs, err := utils.ListDir("Z:\\windows\\drivers\\winpe")
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	for _, dir := range dirs {
		files, err := utils.ListFiles(dir, ".inf", true)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		for _, file := range files {
			//cd dir
			cmd = `cmd /c "cd /d ` + dir + ` && drvload ` + file + `"`
			utils.Logger.Debug(cmd)
			if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
				utils.Logger.Error(err.Error())
			} else {
				output = string(outputBytes)
				utils.Logger.Info(output)
			}
			/*
				//load
				cmd = `drvload ` + file
				utils.Logger.Debug(cmd)
				if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
					utils.Logger.Error(err.Error())
				} else {
					output = string(outputBytes)
					utils.Logger.Info(output)
				}
			*/
		}
	}
	//go back to X:
	cmd = `X:`
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		output = string(outputBytes)
	}
	utils.Logger.Info(output)
	return
}

func getXmlFile(sn string, host string) error {
	var url = fmt.Sprintf("http://%s/api/osinstall/v1/device/getSystemBySn?sn=%s",
		host,
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

func mountSamba() error {
	var cmd = `net use Z:`
	utils.Logger.Debug(cmd)
	if _, err := utils.ExecCmd(scriptFile, cmd); err == nil {
		return nil
	} else {
		utils.Logger.Debug(err.Error())
	}

	cmd = `net use Z: \\osinstall\image`
	utils.Logger.Debug(cmd)
	if _, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func installWindows() error {
	// get setup.exe path from xmlPath
	var output []byte
	var err error
	if output, err = ioutil.ReadFile(xmlPath); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	var r = `<Path>(.*)\\install.wim</Path>`
	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(string(output))
	if regResult == nil || len(regResult) != 2 {
		return fmt.Errorf("Can't found %s", "install.wim")
	}
	utils.Logger.Debug("setup path: %s", regResult[1])

	var cmd = fmt.Sprintf(`%s\\setup.exe /unattend:unattended.xml /noreboot`, regResult[1])
	utils.Logger.Debug(cmd)
	if _, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func copyFirstBoot() error {
	var cmd = `xcopy /s /e /y /i Z:\windows\firstboot C:\firstboot`
	utils.Logger.Debug(cmd)
	if _, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func reboot() error {
	var cmd = `wpeutil reboot`
	utils.Logger.Debug(cmd)
	if _, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}
