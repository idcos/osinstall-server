package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
	"utils"

	"github.com/axgle/mahonia"
	"github.com/codegangsta/cli"
)

type RestInfo struct {
	Bonding  string
	Gateway  string
	Hostname string
	Ip       string
	Netmask  string
	Trunk    string
	Vlan     string
	HWADDR   string
}

var date = time.Now().Format("2006-01-02")
var version = "v1.3.1 (" + date + ")"

func main() {

	app := cli.NewApp()
	app.Version = version
	app.Action = func(c *cli.Context) {
		run(c)
	}

	app.Run(os.Args)
}

var rootPath = "c:/firstboot"
var scriptFile = path.Join(rootPath, "temp-script.cmd")
var preInstallScript = path.Join(rootPath, "preInstall.cmd")
var postInstallScript = path.Join(rootPath, "postInstall.cmd")
var serverHost = "osinstall" //cloudboot server host

func run(c *cli.Context) error {
	// init log
	utils.InitFileLog()

	if utils.CheckFileIsExist(preInstallScript) {
		if _, err := utils.ExecScript(preInstallScript); err != nil {
			utils.Logger.Error("preinstall error: %s", err.Error())
		}
	}
	//init cloudboot server host
	serverIp := getDomainLookupIP(serverHost)
	if serverIp != "" {
		serverHost = serverIp
	}

	if !utils.PingLoop(serverHost, 300, 2) {
		return errors.New("ping timeout")
	}

	var sn = getSN()
	isVm := isVirtualMachine()
	if isVm {
		sn = getMacAddress()
	}

	if sn == "" {
		utils.Logger.Error("get sn failed!")
	}

	restInfo, err := getRestInfo(sn, serverHost)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	// fmt.Println(restInfo)

	var nicInterfaceIndex = getNicInterfaceIndex(restInfo.HWADDR)
	if nicInterfaceIndex == "" {
		utils.Logger.Error("get nic interface index failed")
	}

	var dns = getDNS()
	if dns == "" {
		utils.Logger.Error("get dns failed")
	}

	if err = diskpart(); err != nil {
		utils.Logger.Error(err.Error())
	}
	utils.ReportProgress(0.7, sn, "分区", "diskpart", serverHost)

	if err = changeHostname(restInfo.Hostname); err != nil {
		utils.Logger.Error(err.Error())
	}
	utils.ReportProgress(0.75, sn, "修改主机名", "change hostname", serverHost)

	if err = changeIP(nicInterfaceIndex, restInfo.Ip, restInfo.Netmask, restInfo.Gateway); err != nil {
		utils.Logger.Error(err.Error())
	}

	if err = changeDNS(nicInterfaceIndex, dns); err != nil {
		utils.Logger.Error(err.Error())
	}

	time.Sleep(30 * time.Second)
	if !utils.PingLoop(serverHost, 300, 2) {
		return errors.New("ping timeout")
	}
	utils.ReportProgress(0.8, sn, "修改网络配置", "change network", serverHost)

	if err = changeReg(); err != nil {
		utils.Logger.Error(err.Error())
	}
	utils.ReportProgress(0.9, sn, "修改注册表", "change reg", serverHost)

	utils.ReportProgress(1.0, sn, "安装完成", "finish", serverHost)

	if utils.CheckFileIsExist(postInstallScript) {
		if _, err := utils.ExecScript(postInstallScript); err != nil {
			utils.Logger.Error("postInstall error: %s", err.Error())
		}
	}

	reboot()

	return nil
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

// 网卡名称
func getNicInterfaceIndex(mac string) string {
	var cmd = fmt.Sprintf(`wmic nic where (MACAddress="%s" AND netConnectionStatus=2) get InterfaceIndex /value`, mac)
	var r = `InterfaceIndex=(.*)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		enc := mahonia.NewDecoder("gbk")
		output = enc.ConvertString(string(outputBytes))
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 2 {
		return ""
	}
	utils.Logger.Info("Nic Interface Index:" + regResult[1])
	// fmt.Println(strings.Trim(regResult[1], "\r\n"))
	return regResult[1]
}

// http get 主机名，网络
func getRestInfo(sn string, host string) (*RestInfo, error) {
	var url = fmt.Sprintf("http://%s/api/osinstall/v1/device/getNetworkBySn?sn=%s&type=json",
		host,
		sn)

	utils.Logger.Debug(url)
	resp, err := http.Get(url)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("call url: %s failed", url)
	}

	var jsonResp struct {
		Status  string
		Message string
		Content RestInfo
	}

	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, err
	}

	return &jsonResp.Content, nil
}

// 分区
func diskpart() error {
	var cmd = `select disk 0
create partition extended
create partition logical
assign
format fs=ntfs quick`

	utils.Logger.Debug(cmd)
	diskpartFilePath := path.Join(rootPath, "disk.txt")

	if utils.CheckFileIsExist(diskpartFilePath) {
		os.Remove(diskpartFilePath)
	}
	file, err := os.Create(diskpartFilePath)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	defer file.Close()

	file.WriteString(cmd)
	file.Close()

	var diskCmd = exec.Command("diskpart", "/s", diskpartFilePath)
	if output, err := diskCmd.Output(); err != nil {
		utils.Logger.Error(err.Error())
		return err
	} else {
		utils.Logger.Debug(string(output))
		return nil
	}
}

// DNS
func getDNS() string {
	var cmd = `echo | nslookup`
	var r = `Address:[:blank:]*(.+)`
	var output string
	utils.Logger.Debug(cmd)
	if outputBytes, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
	} else {
		enc := mahonia.NewDecoder("gbk")
		output = enc.ConvertString(string(outputBytes))
		utils.Logger.Debug(output)
	}

	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(output)
	if regResult == nil || len(regResult) != 2 {
		return ""
	}

	var dns = strings.TrimSpace(regResult[1])
	// fmt.Println(dns)
	return dns
}

// 修改主机名
func changeHostname(hostname string) error {
	var cmd = `wmic ntdomain get Caption  /value`
	var r = `Caption=(.*)`
	var oldname string
	if output, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		oldname = string(output)
	}
	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(oldname)
	if regResult == nil || len(regResult) != 2 {
		return errors.New("not fount caption")
	}

	oldname = strings.TrimSpace(regResult[1])
	utils.Logger.Debug(oldname)

	cmd = fmt.Sprintf(`netdom renamecomputer %s /newname:%s /force`, strings.Trim(oldname, "\r\n"), hostname)
	utils.Logger.Debug(cmd)

	if output, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		// fmt.Println(string(output))
		utils.Logger.Debug(string(output))
	}
	return nil
}

// 修改 IP
func changeIP(nic, ip, netmask, gateway string) error {
	var cmd = fmt.Sprintf(`netsh interface ipv4 set address name="%s" source=static addr=%s mask=%s gateway=%s`, nic, ip, netmask, gateway)
	enc := mahonia.NewEncoder("gbk")
	cmd = enc.ConvertString(cmd)
	utils.Logger.Debug(cmd)
	if output, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		utils.Logger.Debug(string(output))
	}

	return nil
}

// 修改DNS
func changeDNS(nic, dns string) error {
	var cmd = fmt.Sprintf(`netsh interface ipv4 set dnsservers name="%s" static %s primary`, nic, dns)
	enc := mahonia.NewEncoder("gbk")
	cmd = enc.ConvertString(cmd)
	utils.Logger.Debug(cmd)
	if output, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		utils.Logger.Debug(string(output))
	}

	return nil
}

// 修改注册表
func changeReg() error {
	var cmd1 = `reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v AutoAdminLogon /t reg_sz /d 0 /f`
	utils.Logger.Debug(cmd1)
	if output, err := utils.ExecCmd(scriptFile, cmd1); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		utils.Logger.Debug(string(output))
	}

	var cmd2 = `reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v Defaultpassword /t reg_sz /d "" /f`
	utils.Logger.Debug(cmd2)
	if output, err := utils.ExecCmd(scriptFile, cmd2); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		utils.Logger.Debug(string(output))
	}
	return nil
}

// 重启
func reboot() error {
	var cmd = fmt.Sprintf(`shutdown -f -r -t 10`)
	utils.Logger.Debug(cmd)
	if output, err := utils.ExecCmd(scriptFile, cmd); err != nil {
		utils.Logger.Error(err.Error())
		return nil
	} else {
		utils.Logger.Debug(string(output))
	}

	return nil
}
