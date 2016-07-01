package route

import (
	//"encoding/base64"
	//"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	//"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"server/osinstallserver/util"
	"strings"
)

func GetDhcpSubnetList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Limit  uint
		Offset uint
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mods, err := repo.GetDhcpSubnetListWithPage(info.Limit, info.Offset)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountDhcpSubnet()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func SaveDhcpSubnet(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID          uint
		StartIp     string
		EndIp       string
		Gateway     string
		AccessToken string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.StartIp = strings.TrimSpace(info.StartIp)
	info.EndIp = strings.TrimSpace(info.EndIp)
	info.Gateway = strings.TrimSpace(info.Gateway)

	isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.StartIp)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	if !isValidate {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "起始IP格式不正确!", "Content": ""})
		return
	}

	isValidateEndIp, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.EndIp)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	if !isValidateEndIp {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "结束IP格式不正确!", "Content": ""})
		return
	}

	isValidateGateway, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.Gateway)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	if !isValidateGateway {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "网关格式不正确!", "Content": ""})
		return
	}

	info.AccessToken = strings.TrimSpace(info.AccessToken)
	_, errVerify := VerifyAccessPurview(info.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	if info.StartIp == "" || info.EndIp == "" || info.Gateway == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	if info.ID > uint(0) {
		_, err := repo.UpdateDhcpSubnetById(info.ID, info.StartIp, info.EndIp, info.Gateway)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	} else {
		_, err := repo.AddDhcpSubnet(info.StartIp, info.EndIp, info.Gateway)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	//save to dhcp config file
	file := "/etc/dhcp/dhcpd.conf"
	if !util.FileExist(file) {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "数据已保存，DHCP配置文件(" + file + ")不存在，请手工配置!"})
		return
	}

	text, err := util.ReadFile(file)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	reg, err := regexp.Compile("(?s)subnet(\\s+)(.*?)(\\s+)(netmask)(\\s+)(.*?)(\\s+){(.*?)}")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	matchs := reg.FindAllStringSubmatch(text, -1)
	if len(matchs) <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "DHCP配置文件(" + file + ")未配置subnet节点，请手工操作！"})
		return
	}

	if len(matchs) > 1 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "DHCP配置文件(" + file + ")存在多个subnet节点，请手工操作！"})
		return
	}

	strFormat := `subnet %s netmask %s {
	range %s %s;
	option routers %s;
}`
	str := fmt.Sprintf(strFormat,
		matchs[0][2],
		matchs[0][6],
		info.StartIp,
		info.EndIp,
		info.Gateway)
	text = strings.Replace(text, matchs[0][0], str, -1)

	logger.Debugf("update dhcp config %s:%s", file, text)
	var bytes = []byte(text)
	errWrite := ioutil.WriteFile(file, bytes, 0666)
	if errWrite != nil {
		logger.Errorf("error:%s", errWrite.Error())
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errWrite.Error()})
		return
	}

	//restart dhcp service
	var cmd = "service dhcpd restart"
	logger.Debugf("restart dhcpd:%s", cmd)
	restartBytes, err := util.ExecScript(cmd)
	logger.Debugf("result:%s", string(restartBytes))

	var runResult = "<br>执行脚本:" + cmd
	runResult += "<br>" + "执行结果:" + string(restartBytes)

	if err != nil {
		runResult += "<br>" + "错误信息:" + err.Error()
		logger.Errorf("error:%s", err.Error())
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "DHCP重启失败!" + runResult})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
