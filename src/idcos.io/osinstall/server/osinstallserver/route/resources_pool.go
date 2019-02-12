package route

import (
	"fmt"
	"idcos.io/osinstall/middleware"
	"idcos.io/osinstall/server/osinstallserver/util"
	"strings"
	"time"

	"github.com/AlexanderChen1989/go-json-rest/rest"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/context"
)

var c = cache.New(5*time.Minute, 30*time.Second)

type BatchOperateInfo struct {
	Sn          string
	AccessToken string
	UserID      uint
	OobIp       string
	Username    string
	Password    string
}

func BatchPowerOn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, _ := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)

	var infos []BatchOperateInfo
	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("参数错误: %s", err.Error())})
		return
	}

	//check permission
	isValidated, infos := CheckPermissionForBatchOperate(ctx, w, r, infos)
	if isValidated != true {
		return
	}

	for _, info := range infos {
		info.Sn = strings.TrimSpace(info.Sn)
		_, errInfo := repo.GetManufacturerBySn(info.Sn)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}

		status, err := util.GetDevicePowerStatusFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if status != "off" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("该设备已开机，无法执行开机指令(SN: %s)", info.Sn)})
			return
		}

		runErr := util.PowerOnDeviceFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if runErr != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("通过IPMI开机失败: %s", runErr.Error()) + "(SN:" + info.Sn + ")"})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchReStart(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, _ := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)

	var infos []BatchOperateInfo
	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
		return
	}

	//check permission
	isValidated, infos := CheckPermissionForBatchOperate(ctx, w, r, infos)
	if isValidated != true {
		return
	}

	for _, info := range infos {
		info.Sn = strings.TrimSpace(info.Sn)
		_, errInfo := repo.GetManufacturerBySn(info.Sn)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}

		status, err := util.GetDevicePowerStatusFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if status != "on" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备已关机，无法执行重启指令(SN: " + info.Sn + ")"})
			return
		}

		runErr := util.RestartDeviceFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if runErr != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "通过IPMI重启失败: " + runErr.Error() + "(SN:" + info.Sn + ")"})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchStartFromPxe(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, _ := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)

	var infos []BatchOperateInfo
	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("参数错误: %s", err.Error())})
		return
	}

	//check permission
	isValidated, infos := CheckPermissionForBatchOperate(ctx, w, r, infos)
	if isValidated != true {
		return
	}

	for _, info := range infos {
		info.Sn = strings.TrimSpace(info.Sn)
		_, errInfo := repo.GetManufacturerBySn(info.Sn)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}

		status, err := util.GetDevicePowerStatusFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if status == "" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("从PXE启动失败，无法连接到设备(SN: %s)", info.Sn)})
			return
		}
		runErr := util.BootDeviceToPXEFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if runErr != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("通过IPMI从PXE启动失败: %s", runErr.Error()) + "(SN:" + info.Sn + ")"})
			return
		}
		if status == "on" {
			runErr := util.RestartDeviceFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
			if runErr != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("通过IPMI重启失败: %s", runErr.Error()) + "(SN:" + info.Sn + ")"})
				return
			}
		} else if status == "off" {
			runErr := util.PowerOnDeviceFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
			if runErr != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("通过IPMI开机失败: %s", runErr.Error()) + "(SN:" + info.Sn + ")"})
				return
			}
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchPowerOff(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, _ := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)

	var infos []BatchOperateInfo
	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("参数错误: %s", err.Error())})
		return
	}

	//check permission
	isValidated, infos := CheckPermissionForBatchOperate(ctx, w, r, infos)
	if isValidated != true {
		return
	}

	for _, info := range infos {
		info.Sn = strings.TrimSpace(info.Sn)
		_, errInfo := repo.GetManufacturerBySn(info.Sn)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}

		status, err := util.GetDevicePowerStatusFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if status != "on" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("该设备已关机，无法执行关机指令(SN: %s)", info.Sn)})
			return
		}

		runErr := util.PowerOffDeviceFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if runErr != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("通过IPMI关机失败: %s", runErr.Error()) + "(SN:" + info.Sn + ")"})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func CheckPermissionForBatchOperate(ctx context.Context, w rest.ResponseWriter, r *rest.Request, infos []BatchOperateInfo) (bool, []BatchOperateInfo) {
	repo, _ := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)
	// locale, _ := middleware.LocaleFromContext(ctx)

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误: " + err.Error()})
		return false, nil
	}

	var result []BatchOperateInfo
	for _, info := range infos {
		if info.OobIp == "" || info.Username == "" || info.Password == "" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "OOB IP、用户名、密码不能为空" + "(SN:" + info.Sn + ")"})
			return false, nil
		}
		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return false, nil
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		} else {
			info.UserID = session.ID
		}

		info.Sn = strings.TrimSpace(info.Sn)
		device, errInfo := repo.GetManufacturerBySn(info.Sn)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return false, nil
		}

		if session.Role != "Administrator" && device.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "您没有权限操作其他人的设备(SN:" + info.Sn + ")"})
			return false, nil
		}

		sn, err := util.GetDeviceSnFromIpmitool(repo, logger, info.OobIp, info.Username, info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error() + "(SN:" + info.Sn + ")!"})
			return false, nil
		}
		if sn != info.Sn {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("设备登记SN和实际SN不一致(登记SN:%s, 实际SN:%s)", info.Sn, sn)})
			return false, nil
		}

		result = append(result, info)
	}
	return true, result
}
