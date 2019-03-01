package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"

	"idcos.io/osinstall/build"
	"idcos.io/osinstall/utils"
)

var xmlPath = "X:\\Windows\\System32\\unattended.xml"
var rootPath = "X:\\Windows\\System32"
var scriptFile = path.Join(rootPath, "temp-script.cmd")
var serverHost = "osinstall" //cloudboot server host

func main() {
	app := cli.NewApp()
	app.Version = build.Version()
	app.Action = func(c *cli.Context) {
		if err := run(c); err != nil {
			utils.Logger.Error(err.Error())
		}
	}

	app.Run(os.Args)
}

func run(c *cli.Context) (err error) {
	utils.InitConsoleLog()

	var sn = getSN()
	isVM := isVirtualMachine()
	if isVM {
		sn = getMacAddress()
	}

	if sn == "" {
		return errors.New("get sn failed")
	}

	//init cloudboot server host
	serverIP := getDomainLookupIP(serverHost)
	if serverIP != "" {
		serverHost = serverIP
	}

	if !utils.PingLoop(serverHost, 30, 2) {
		return errors.New("ping timeout")
	}

	// get xml for install windows
	if err := getXmlFile(sn, serverHost); err != nil {
		return err
	}

	// mount samba
	retries := 3 // Samba挂载失败时重试次数
	for i := 0; i < retries; i++ {
		if err = mountSamba(serverHost); err != nil {
			time.Sleep(time.Duration((i+1)*5) * time.Second)
			continue
		}
		break
	}
	if err != nil {
		utils.Logger.Error(err.Error())
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

// mountSamba 挂载Samba目录
func mountSamba(host string) (err error) {
	if sambaMounted() {
		return nil
	}
	cmdAndArgs := fmt.Sprintf(`net use Z: \\%s\image`, host)
	_, err = execOutput(cmdAndArgs)
	return err
}

// sambaMounted 判断Samba是否已经挂载
func sambaMounted() (mounted bool) {
	_, err := execOutput(`net use Z:`)
	return err == nil
}

// execOutput windows系统下，执行命令字符串cmdAndArgs，并将命令执行的标准输出和标准错误输出都通过字节切片output返回。
func execOutput(cmdAndArgs string) (output []byte, err error) {
	scriptFile, err := genTempScript([]byte(cmdAndArgs))
	if err != nil {
		return nil, err
	}
	defer os.Remove(scriptFile)

	output, err = exec.Command("cmd", "/c", scriptFile).Output()
	utils.Logger.Debug("%s ==>\n%s\n", cmdAndArgs, output)

	return output, err
}

// genTempScript 在系统临时目录生成bat脚本文件
func genTempScript(content []byte) (scriptFile string, err error) {
	scriptFile = filepath.Join(os.TempDir(), fmt.Sprintf("%d.bat", time.Now().UnixNano()))
	if err = ioutil.WriteFile(scriptFile, content, 0744); err != nil {
		return "", err
	}
	return scriptFile, nil
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
	var dir = "C:/firstboot"
	if !utils.CheckFileIsExist("c:/windows") {
		dir = "D:/firstboot"
	}
	var cmd = fmt.Sprintf(`xcopy /s /e /y /i Z:\windows\firstboot %s`, dir)
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
