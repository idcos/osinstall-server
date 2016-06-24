package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"utils/ping"

	"github.com/astaxie/beego/logs"
)

var Logger *logs.BeeLogger
var rootPath = "c:/firstboot"
var logPath = path.Join(rootPath, "log")
var logFile = path.Join(logPath, "setup.log")

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func InitFileLog() {
	os.MkdirAll(logPath, 0666)
	Logger = logs.NewLogger(1000)
	Logger.SetLogger("file", `{"filename":"`+logFile+`"}`)
}

func InitConsoleLog() {
	Logger = logs.NewLogger(1000)
	Logger.SetLogger("console", "")
}

// ExecCmd 执行 command
func ExecCmd(scriptFile, cmdString string) ([]byte, error) {

	// 生成临时文件
	if CheckFileIsExist(scriptFile) {
		os.Remove(scriptFile)
	}
	file, err := os.Create(scriptFile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if _, err = file.WriteString(cmdString); err != nil {
		return nil, err
	}
	file.Close()

	return ExecScript(scriptFile)
}

// ExecScript exec script
func ExecScript(scriptPath string) ([]byte, error) {

	var cmd = exec.Command("cmd", "/c", scriptPath)
	return cmd.Output()
}

// CallRestAPI 调用restful api
func CallRestAPI(url string, jsonReq interface{}) ([]byte, error) {
	var req *http.Request
	var resp *http.Response
	var err error
	var reqBody []byte

	if reqBody, err = json.Marshal(jsonReq); err != nil {
		return nil, err
	}

	fmt.Printf("Request BODY: %s \n", string(reqBody))
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

// ReportProgress 上报执行结果
func ReportProgress(installProgress float64, sn, title, installLog string, host string) bool {
	var url = "http://" + host + "/api/osinstall/v1/report/deviceInstallInfo"
	Logger.Debug("ReportProgress url:%s\n", url)
	var jsonReq struct {
		Sn              string
		InstallProgress float64
		InstallLog      string
		Title           string
	}
	jsonReq.Sn = sn
	jsonReq.InstallProgress = installProgress
	jsonReq.Title = title
	jsonReq.InstallLog = base64.StdEncoding.EncodeToString([]byte(installLog)) // base64编码
	Logger.Debug("SN: %s\n", jsonReq.Sn)
	Logger.Debug("InstallProgress: %f\n", jsonReq.InstallProgress)
	Logger.Debug("InstallLog: %s\n", jsonReq.InstallLog)
	Logger.Debug("Title: %s\n", jsonReq.Title)

	var jsonResp struct {
		Status  string
		Message string
		Content struct {
			Result string
		}
	}

	var ret, err = CallRestAPI(url, jsonReq)
	Logger.Debug("ReportProgress api result:%s\n", string(ret))
	if err != nil {
		Logger.Error(err.Error())
		return false
	}

	if err := json.Unmarshal(ret, &jsonResp); err != nil {
		Logger.Error(err.Error())
		return false
	}

	if jsonResp.Status != "success" {
		return false
	}
	return true
}

// PingLoop return when success
func PingLoop(host string, pkgCnt int, timeout int) bool {
	for i := 0; i < pkgCnt; i++ {
		if ping.Ping(host, timeout) {
			return true
		}
	}
	return false
}

func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		files = append(files, dirPth+PthSep+fi.Name())
	}
	return files, nil
}

func ListFiles(dirPth string, suffix string, onlyReturnFileName bool) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			if onlyReturnFileName == true {
				files = append(files, fi.Name())
			} else {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	return files, nil
}
