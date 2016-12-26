package route

import (
	"fmt"
	"middleware"
	"server/osinstallserver/util"
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
	isValidated, infos := CheckPermissionForBatchOperate(ctx, w, r, infos, true, true)
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

		// TODO add device log
		// _, errAddLog := repo.AddDeviceLog(0, "设备开机", "manage", "设备开机", info.Sn) // use english
		// if errAddLog != nil {
		// 	w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
		// 	return
		// }
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
	isValidated, infos := CheckPermissionForBatchOperate(ctx, w, r, infos, true, true)
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

		// TODO add device log
		// _, errAddLog := repo.AddDeviceLog(0, "设备重启", "manage", "设备重启", info.Sn)
		// if errAddLog != nil {
		// 	w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
		// 	return
		// }
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func CheckPermissionForBatchOperate(ctx context.Context, w rest.ResponseWriter, r *rest.Request, infos []BatchOperateInfo, isCheckOnline bool, isCheckBatchOperateNum bool) (bool, []BatchOperateInfo) {
	repo, _ := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)
	conf, _ := middleware.ConfigFromContext(ctx)
	// locale, _ := middleware.LocaleFromContext(ctx)

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误: " + err.Error()})
		return false, nil
	}

	cacheKey := fmt.Sprintf("RecentOperateDeviceNum_%d", session.ID)
	//5分钟内超过操作次数限制
	var recentOperateDeviceNum int
	if isCheckBatchOperateNum == true {
		var maxNum int
		maxNum = conf.Device.MaxBatchOperateNum
		if maxNum <= 0 {
			maxNum = 5
		}

		if len(infos) > maxNum {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("批量操作设备台数超过限制，最大为%d台。", maxNum)})
			return false, nil
		}

		ch, found := c.Get(cacheKey)
		if ch != nil {
			logger.Debugf("get cache %s:%d,%t", cacheKey, ch, found)
		} else {
			logger.Debugf("get cache %s:no cache found!", cacheKey)
		}
		if found {
			recentOperateDeviceNum = ch.(int)
		} else {
			recentOperateDeviceNum = 0
		}

		if conf.Device.MaxOperateNumIn5Minutes > 0 &&
			(recentOperateDeviceNum+len(infos)) > conf.Device.MaxOperateNumIn5Minutes {
			if session.ID > 0 && session.Role != "Administrator" {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("批量操作设备台数超过限制，5分钟内最大为%d台，请稍候再试。", conf.Device.MaxOperateNumIn5Minutes)})
				return false, nil
			}
		}
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
		if isCheckOnline == true {
			if device.Status == "online" {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("请先将该设备下线后再进行操作(SN:%s)", info.Sn)})
				return false, nil
			}
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
	if isCheckBatchOperateNum == true {
		logger.Debugf("set cache %s:%d", cacheKey, len(infos)+recentOperateDeviceNum)
		c.Set(cacheKey, len(infos)+recentOperateDeviceNum, cache.DefaultExpiration)
	}
	return true, result
}
