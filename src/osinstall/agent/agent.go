package agent

import (
	"bytes"
	"config"
	"config/iniconf"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	GetSNScript = `dmidecode -s system-serial-number 2>/dev/null | awk '/^[^#]/ { print $1 }'`
	// GetMacScript       = `ip addr show $(ip route get 10.0.0.0 | awk '/src/ { print $(NF-2) }') | awk '/ether/ { print $2 }'`
	GetCmdlineArgs     = `cat /proc/cmdline`
	RegexpServerAddr   = `SERVER_ADDR=([^ ]+)`
	RegexpLoopInterval = `LOOP_INTERVAL=([^ ]+)`
	RegexpDeveloper    = `DEVELOPER=([^ ]+)`
	RebootScript       = `ipmitool chassis bootdev pxe; ipmitool power reset`
	RebootScript2      = `fdisk -lu | awk '/^Disk.*bytes/ { gsub(/:/, ""); system("dd if=/dev/zero of="$2" bs=512 count=1") }'; reboot -f`
	InstallHWTools     = `rpm --quiet -q %s-hw-tools || yum -y install %s-hw-tools`
	PingIp             = `ping -c 4 -w 3 %s`

	APIVersion = "v1"

	ProductInfoScript       = `/usr/local/bin/sysinfo.sh`
	On                      = "1"
	Off                     = "0"
	PreInstallScript        = "/tmp/preInstall.cmd"
	PostInstallScript       = "/tmp/postInstall.cmd"
	RegexpPreInstallScript  = `PRE=([^ ]+)`
	RegexpPostInstallScript = `POST=([^ ]+)`
)

var (
	confContent = `
[Logger]
color = false
level = debug
logFile = /var/log/cloudboot-agent.log
`

	defaultLoopInterval = 60
	hardwareURL         = fmt.Sprintf("/api/osinstall/%s/device/getHardwareBySn", APIVersion)
	isHaveHardwareConf  = fmt.Sprintf("/api/osinstall/%s/device/getPrepareInstallInfo", APIVersion)
	installListURL      = fmt.Sprintf("/api/osinstall/%s/device/isInInstallList", APIVersion)
	productInfoURL      = fmt.Sprintf("/api/osinstall/%s/report/deviceProductInfo", APIVersion)
	installInfoURL      = fmt.Sprintf("/api/osinstall/%s/report/deviceInstallInfo", APIVersion)
	macInfoURL          = fmt.Sprintf("/api/osinstall/%s/report/deviceMacInfo", APIVersion)
	netInfoURL          = fmt.Sprintf("/api/osinstall/%s/device/getNetworkBySn", APIVersion)
)

// HardWareConf 硬件配置结构
type HardWareConf struct {
	Name    string
	Scripts []struct {
		Name   string
		Script string
	}
}

// OSInstallAgent agent data struct
type OSInstallAgent struct {
	Logger        logger.Logger
	Config        *config.Config
	Sn            string
	MacAddr       string
	ServerAddr    string
	LoopInterval  int
	DevelopeMode  string
	Company       string
	Product       string
	ModelName     string
	hardwareConfs []HardWareConf // base64 编码的硬件配置脚本
	IsVm          string         //Whether it is a virtual machine
}

type nicInfo struct {
	Name string
	Mac  string
	Ip   string
}
type cpuInfo struct {
	Model string
	Core  string
}
type diskInfo struct {
	Name string
	Size string
}
type memoryInfo struct {
	Name string
	Size string
}
type motherboardInfo struct {
	Name  string
	Model string
}

// New create agent
func New() (*OSInstallAgent, error) {
	// get config
	var conf, err = iniconf.NewContent([]byte(confContent)).Load()
	if err != nil {
		return nil, err
	}
	var log = logger.NewBeeLogger(conf)
	var agent = &OSInstallAgent{
		Config: conf,
		Logger: log,
	}

	// var data []byte
	// // get sn
	// agent.Logger.Debug("START to get SN")
	// if data, err = execScript(GetSNScript); err != nil {
	// 	agent.Logger.Error(data)
	// 	agent.Logger.Error(err)
	// 	return nil, fmt.Errorf("get SN error: \n#%s\n%v\n%s", GetSNScript, err, string(data))
	// }
	// agent.Sn = string(data)
	// agent.Sn = strings.Trim(agent.Sn, "\n")
	// agent.Logger.Debugf("SN: %s", agent.Sn)

	// get mac addr by sysinfo.sh
	// get mac addr
	// agent.Logger.Debug("START to get mac addr")
	// if data, err = execScript(GetMacScript); err != nil {
	// 	agent.Logger.Error(data)
	// 	log.Error(err)
	// 	return nil, fmt.Errorf("get mac addr error: \n#%s\n%v\n%s", GetMacScript, err, string(data))
	// }
	// agent.MacAddr = string(data)
	// agent.MacAddr = strings.Trim(agent.MacAddr, "\n")
	// agent.Logger.Debugf("Mac ADDR: %s", agent.MacAddr)

	var serverAddr = ""
	// get server addr
	agent.Logger.Debug("START to get SERVER_ADDR")
	if serverAddr, err = getCmdlineArgs(RegexpServerAddr); err != nil {
		return nil, err
	}
	agent.ServerAddr = serverAddr
	// agent.ServerAddr = "http://10.0.0.135:8083"
	agent.ServerAddr = strings.Trim(agent.ServerAddr, "\n")
	agent.Logger.Debugf("SERVER_ADDR: %s", agent.ServerAddr)

	// loop interval
	var interval string
	agent.Logger.Debug("START to get LOOP_INTERVAL")
	if interval, err = getCmdlineArgs(RegexpLoopInterval); err != nil {
		agent.Logger.Debug(err)
		agent.LoopInterval = defaultLoopInterval
	} else {
		agent.LoopInterval = parseInterval(interval)
	}
	agent.Logger.Debugf("LOOP_INTERVAL: %s", agent.LoopInterval)

	var developMode = ""
	agent.Logger.Debug("START to get DEVELOPER")
	if developMode, err = getCmdlineArgs(RegexpDeveloper); err != nil {
		agent.Logger.Debug(err)
		agent.DevelopeMode = ""
	}
	agent.DevelopeMode = developMode
	agent.DevelopeMode = strings.Trim(agent.DevelopeMode, "\n")
	agent.Logger.Debugf("DEVELOPER: %s", agent.DevelopeMode)

	return agent, nil
}

//run pre install script
func (agent *OSInstallAgent) RunPreInstallScript() {
	agent.Logger.Debug("START to get PreInstallScript")
	url, err := getCmdlineArgs(RegexpPreInstallScript)
	url = strings.Trim(url, "\n")
	url = strings.TrimSpace(url)
	if err != nil {
		agent.Logger.Error(err.Error())
		return
	}
	agent.Logger.Infof("script:%s", url)
	if url == "" {
		return
	}

	agent.Logger.Debugf("START to wget %s", url)
	script, err := wget(url)
	if err != nil {
		agent.Logger.Error(err.Error())
		return
	}
	agent.Logger.Debugf("script:%s", script)

	agent.Logger.Debugf("write to file %s:%s", PreInstallScript, script)
	var bytes = []byte(script)
	errWrite := ioutil.WriteFile(PreInstallScript, bytes, 0666)
	if errWrite != nil {
		agent.Logger.Error(errWrite.Error())
		return
	}

	//chmod 755 PreInstallScript
	cmd := `chmod 755 ` + PreInstallScript
	agent.Logger.Debugf("exec:%s", cmd)
	data, errRun := execScript(cmd)
	if errRun != nil {
		agent.Logger.Error(errRun.Error())
		return
	}
	agent.Logger.Debugf("result:%s", string(data))

	//run PreInstallScript
	cmd = PreInstallScript
	agent.Logger.Debugf("exec:%s", cmd)
	data, errRun = execScript(cmd)
	if errRun != nil {
		agent.Logger.Error(errRun.Error())
		return
	}
	agent.Logger.Debugf("result:%s", string(data))
	return
}

//run post install script
func (agent *OSInstallAgent) RunPostInstallScript() {
	agent.Logger.Debug("START to get PostInstallScript")
	url, err := getCmdlineArgs(RegexpPostInstallScript)
	url = strings.Trim(url, "\n")
	url = strings.TrimSpace(url)
	if err != nil {
		agent.Logger.Error(err.Error())
		return
	}
	agent.Logger.Infof("script:%s", url)
	if url == "" {
		return
	}

	agent.Logger.Debugf("START to wget %s", url)
	script, err := wget(url)
	if err != nil {
		agent.Logger.Error(err.Error())
		return
	}
	agent.Logger.Debugf("script:%s", script)

	agent.Logger.Debugf("write to file %s:%s", PostInstallScript, script)
	var bytes = []byte(script)
	errWrite := ioutil.WriteFile(PostInstallScript, bytes, 0666)
	if errWrite != nil {
		agent.Logger.Error(errWrite.Error())
		return
	}

	//chmod 755 PostInstallScript
	cmd := `chmod 755 ` + PostInstallScript
	agent.Logger.Debugf("exec:%s", cmd)
	data, errRun := execScript(cmd)
	if errRun != nil {
		agent.Logger.Error(errRun.Error())
		return
	}
	agent.Logger.Debugf("result:%s", string(data))

	//run PostInstallScript
	cmd = PostInstallScript
	agent.Logger.Debugf("exec:%s", cmd)
	data, errRun = execScript(cmd)
	if errRun != nil {
		agent.Logger.Error(errRun.Error())
		return
	}
	agent.Logger.Debugf("result:%s", string(data))
	return
}

// IsInInstallQueue 检查是否在装机队列中 （定时执行）
func (agent *OSInstallAgent) IsInInstallQueue() {
	// 轮询是否在装机队列中
	var t = time.NewTicker(time.Duration(agent.LoopInterval) * time.Second)
	var url = agent.ServerAddr + installListURL
	agent.Logger.Debugf("IsInPreInstallQueue url:%s\n", url)
	var jsonReq struct {
		Sn string
	}
	jsonReq.Sn = agent.Sn

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Result string
		}
	}

LOOP:
	for {
		select {
		case <-t.C:
			agent.Logger.Debugf("IsInPreInstallQueue request body: %v", jsonReq)
			var ret, err = callRestAPI(url, jsonReq)
			agent.Logger.Debugf("IsInPreInstallQueue api result:%s\n", string(ret))
			if err != nil {
				agent.Logger.Error(err)
				continue // 继续等待下次轮询
			}

			if err := json.Unmarshal(ret, &jsonResp); err != nil {
				agent.Logger.Error(err)
				continue // 继续等待下次轮询
			}

			if jsonResp.Content.Result == "true" {
				t.Stop()
				break LOOP
			}
		}
	}
}

// IsIPInUse 判断IP是否在使用中
func (agent *OSInstallAgent) IsIpInUse() error {
	var url = fmt.Sprintf("%s?sn=%s&type=json", agent.ServerAddr+netInfoURL, agent.Sn)
	agent.Logger.Debugf("IsIPInUse url:%s\n", url)
	var body []byte
	var err error
	var resp *http.Response
	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Ip string
		}
	}

	resp, err = http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	agent.Logger.Debug(string(body))

	if err = json.Unmarshal(body, &jsonResp); err != nil {
		return err
	}

	if jsonResp.Status != "success" {
		return fmt.Errorf("Status: %s, Message: %s", jsonResp.Status, jsonResp.Message)
	}

	var pingScript = fmt.Sprintf(PingIp, jsonResp.Content.Ip)
	if output, err := execScript(pingScript); err == nil {
		return fmt.Errorf("IsIpInUse error: \n#%s\n%v\n%s", pingScript, err, string(output))
	}

	return nil
}

//
func (agent *OSInstallAgent) ReportProductInfo() error {
	var url = agent.ServerAddr + productInfoURL
	agent.Logger.Debugf("ReportProductInfo url:%s\n", url)
	var jsonReq struct {
		Sn          string
		Company     string
		Product     string
		ModelName   string
		Ip          string
		Mac         string
		Nic         []nicInfo
		Cpu         cpuInfo
		Memory      []memoryInfo
		MemorySum   int
		DiskSum     int
		Disk        []diskInfo
		Motherboard motherboardInfo
		Raid        string
		Oob         string
		IsVm        string
		NicDevice   string
	}

	// get infoFull from script
	if output, err := execScript(ProductInfoScript); err != nil {
		return fmt.Errorf("ReportProductInfo error: \n#%s\n%v\n%s", ProductInfoScript, err, string(output))
	} else if err = json.Unmarshal(output, &jsonReq); err != nil {
		return err
	}

	// set company to agent
	agent.Logger.Debug("Start to get SN")
	agent.Sn = strings.Trim(jsonReq.Sn, "\n")
	agent.Sn = strings.TrimSpace(agent.Sn)
	agent.Logger.Debugf("SN: %s", agent.Sn)

	agent.Company = strings.ToLower(jsonReq.Company)
	agent.Product = jsonReq.Product
	agent.ModelName = jsonReq.ModelName

	//set whether it is a virtual machine
	jsonReq.IsVm = strings.TrimSpace(jsonReq.IsVm)
	if jsonReq.IsVm != "Yes" {
		jsonReq.IsVm = "No"
	}
	agent.IsVm = jsonReq.IsVm

	// set mac info to agent
	for _, nic := range jsonReq.Nic {
		if nic.Ip != "" {
			agent.MacAddr = nic.Mac
			break
		}
	}

	var jsonResp struct {
		Status  string
		Message string
	}

	agent.Logger.Debugf("ReportProductInfo request body: %v", jsonReq)
	var ret, err = callRestAPI(url, jsonReq)
	if err != nil {
		return err
	}
	agent.Logger.Debugf("ReportProductInfo api result:%s\n", string(ret))

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		return err
	}

	if jsonResp.Status != "success" {
		return fmt.Errorf("Status: %s, Message: %s", jsonResp.Status, jsonResp.Message)
	}

	return nil
}

// IsHaveHardWareConf 检查服务端是否此机器的硬件配置
func (agent *OSInstallAgent) IsHaveHardWareConf() (bool, error) {
	if agent.IsVm == "Yes" {
		return true, nil
	}

	var url = agent.ServerAddr + isHaveHardwareConf
	var skipHWConf = false
	agent.Logger.Debugf("IsHaveHardWareConf url:%s\n", url)
	var jsonReq struct {
		Sn        string
		Company   string
		Product   string
		ModelName string
	}

	jsonReq.Sn = agent.Sn
	jsonReq.Company = agent.Company
	jsonReq.Product = agent.Product
	jsonReq.ModelName = agent.ModelName

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			IsVerify             string
			IsSkipHardwareConfig string
		}
	}

	agent.Logger.Debugf("IsHaveHardWareConf request body: %v", jsonReq)
	var ret, err = callRestAPI(url, jsonReq)
	agent.Logger.Debugf("IsHaveHardWareConf api result:%s\n", string(ret))
	if err != nil {
		return skipHWConf, err
	}

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		return skipHWConf, err
	}

	if jsonResp.Status != "success" {
		return skipHWConf, fmt.Errorf("Status: %s, Message: %s", jsonResp.Status, jsonResp.Message)
	}

	// is skip hardware configuration
	if jsonResp.Content.IsSkipHardwareConfig == "true" {
		return true, nil
	}

	if jsonResp.Content.IsVerify == "false" && agent.DevelopeMode != On {
		return skipHWConf, errors.New("Verify is false AND developMode is not On")
	}

	return false, nil
}

// GetHardConf 获取硬件配置
func (agent *OSInstallAgent) GetHardWareConf() error {
	var url = agent.ServerAddr + hardwareURL
	agent.Logger.Debugf("GetHardWareConf url:%s\n", url)
	var jsonReq struct {
		Sn string
	}
	jsonReq.Sn = agent.Sn

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Company   string
			ModelName string
			Product   string
			Hardware  []HardWareConf
		}
	}

	agent.Logger.Debugf("GetHardWareConf request body: %v", jsonReq)
	var ret, err = callRestAPI(url, jsonReq)
	agent.Logger.Debugf("GetHardWareConf api result:%s\n", string(ret))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		return err
	}

	if jsonResp.Status != "success" {
		return fmt.Errorf("Status: %s, Message: %s", jsonResp.Status, jsonResp.Message)
	}

	agent.hardwareConfs = jsonResp.Content.Hardware

	return nil
}

// ImplementHardConf 实施硬件配置
func (agent *OSInstallAgent) ImplementHardConf() error {

	// 安装硬件配置工具包
	installHWScript := fmt.Sprintf(InstallHWTools, agent.Company, agent.Company)
	agent.Logger.Debugf("installScript: %s\n", installHWScript)
	if output, err := execScript(installHWScript); err != nil {
		return fmt.Errorf("ImplementHardConf error: \n#%s\n%v\n%s", installHWScript, err, string(output))
	}

	// 开始硬件配置
	agent.ReportProgress(0.3, "开始硬件配置", "")

	var progressDelta int
	if len(agent.hardwareConfs) != 0 {
		progressDelta = 10 / len(agent.hardwareConfs)
	}

	var curProgress = 0.3
	for _, hardwareConf := range agent.hardwareConfs {
		curProgress = curProgress + float64(progressDelta)/100.0

		for _, scriptB64 := range hardwareConf.Scripts {
			script, err := base64.StdEncoding.DecodeString(scriptB64.Script)
			agent.Logger.Debugf("Script: %s\n", script)
			if err != nil {
				return err
			}

			if output, err := execScript(string(script)); err != nil {
				return fmt.Errorf("execscript hardware script error: \n#%s\n%v\n%s", string(script), err, string(output))
			}
			agent.ReportProgress(curProgress, hardwareConf.Name+" - "+scriptB64.Name, "")
		}
		agent.ReportProgress(curProgress, hardwareConf.Name+" 配置完成", "")
	}
	agent.ReportProgress(0.4, "硬件配置结束", "硬件配置正常结束")
	return nil
}

// ReportProgress 上报执行结果
func (agent *OSInstallAgent) ReportProgress(installProgress float64, title, installLog string) bool {
	var url = agent.ServerAddr + installInfoURL
	agent.Logger.Debugf("ReportProgress url:%s\n", url)
	var jsonReq struct {
		Sn              string
		InstallProgress float64
		InstallLog      string
		Title           string
	}
	jsonReq.Sn = agent.Sn
	jsonReq.InstallProgress = installProgress
	jsonReq.Title = title
	jsonReq.InstallLog = base64.StdEncoding.EncodeToString([]byte(installLog)) // base64编码
	agent.Logger.Debugf("SN: %s\n", jsonReq.Sn)
	agent.Logger.Debugf("InstallProgress: %f\n", jsonReq.InstallProgress)
	agent.Logger.Debugf("InstallLog: %s\n", jsonReq.InstallLog)
	agent.Logger.Debugf("Title: %s\n", jsonReq.Title)

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Result string
		}
	}

	agent.Logger.Debugf("ReportProgress request body: %v", jsonReq)
	var ret, err = callRestAPI(url, jsonReq)
	agent.Logger.Debugf("ReportProgress api result:%s\n", string(ret))
	if err != nil {
		agent.Logger.Error(err.Error())
		return false
	}

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		agent.Logger.Error(err.Error())
		return false
	}

	if jsonResp.Status != "success" {
		return false
	}
	return true
}

// ReportMacInfo 上报 mac 地址
func (agent *OSInstallAgent) ReportMacInfo() error {
	var url = agent.ServerAddr + macInfoURL
	agent.Logger.Debugf("ReportMacInfo url:%s\n", url)
	var jsonReq struct {
		Sn  string
		Mac string
	}
	jsonReq.Sn = agent.Sn
	jsonReq.Mac = agent.MacAddr

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Result string
		}
	}

	agent.Logger.Debugf("ReportMacInfo request body: %v", jsonReq)
	var ret, err = callRestAPI(url, jsonReq)
	agent.Logger.Debugf("ReportMacInfo api result:%s\n", string(ret))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		return err
	}

	if jsonResp.Status != "success" {
		return fmt.Errorf("Status: %s, Message: %s", jsonResp.Status, jsonResp.Message)
	}
	return nil
}

// Reboot 重启系统
func (agent *OSInstallAgent) Reboot() error {
	if output, err := execScript(RebootScript); err != nil {
		if output2, err2 := execScript(RebootScript2); err2 != nil {
			return fmt.Errorf("reboot error: \n#%s\n%v\n%s\n\nsecond reboot error: \n#%s\n%v\n%s", RebootScript, err, string(output), RebootScript2, err2, string(output2))
		}
		//return fmt.Errorf("reboot error: \n#%s\n%v\n%s", RebootScript, err, string(output))
	}
	return nil
}

// getCmdlineArgs get options from cmdline
func getCmdlineArgs(r string) (string, error) {
	var data, err = execScript(GetCmdlineArgs)
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile(r)
	var regResult = reg.FindStringSubmatch(string(data))
	if regResult == nil || len(regResult) != 2 {
		return "", errors.New("Can't find " + r)
	}

	return regResult[1], nil
}

// callRestAPI 调用restful api
func callRestAPI(url string, jsonReq interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var reqBody []byte

	if reqBody, err = json.Marshal(jsonReq); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody)); err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
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

	return body, nil
}

func wget(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// execScript 执行脚本
func execScript(script string) ([]byte, error) {

	// 生成临时文件
	file, err := ioutil.TempFile("", "tmp-script")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if _, err = file.WriteString(script); err != nil {
		return nil, err
	}
	file.Close()

	var cmd = exec.Command("/bin/bash", file.Name())
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	err = cmd.Wait()
	return output.Bytes(), err
}

func parseInterval(interval string) int {
	var err error
	interval = strings.Trim(interval, "\n")
	var i int
	if i, err = strconv.Atoi(interval); err != nil {
		return defaultLoopInterval
	}
	if i > 0 {
		return i
	} else {
		return defaultLoopInterval
	}
}
