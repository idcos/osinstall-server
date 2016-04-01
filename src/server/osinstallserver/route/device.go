package route

import (
	"encoding/base64"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	"regexp"
	"server/osinstallserver/util"
	"strconv"
	"strings"
	//"net/http"
	"encoding/json"
	"model"
	"os"
	"utils"
)

//Device
func DeleteDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, errSession := GetSession(w, r)
	if errSession != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + errSession.Error()})
		return
	}

	var info struct {
		ID          uint
		AccessToken string
		UserID      uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	device, errDevice := repo.GetDeviceById(info.ID)
	if errDevice != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDevice.Error()})
		return
	}

	if session.ID <= uint(0) {
		accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
		if errAccessToken != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
			return
		}
		info.UserID = accessTokenUser.ID
		session.ID = accessTokenUser.ID
		session.Role = accessTokenUser.Role
	} else {
		info.UserID = session.ID
	}

	if session.Role != "Administrator" && device.UserID != info.UserID {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "您无权操作其他人的设备!"})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	//删除PXE配置文件
	macs, errMac := repo.GetMacListByDeviceID(device.ID)
	if errMac != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errMac.Error()})
		return
	}
	for _, mac := range macs {
		pxeFileName := util.GetPxeFileNameByMac(mac.Mac)
		confDir := conf.OsInstall.PxeConfigDir
		if util.FileExist(confDir + "/" + pxeFileName) {
			err := os.Remove(confDir + "/" + pxeFileName)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
		}
	}

	//删除mac
	_, err := repo.DeleteMacByDeviceId(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	//删除LOG
	_, errLog := repo.DeleteDeviceLogByDeviceID(info.ID)
	if errLog != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
		return
	}

	errCopy := repo.CopyDeviceToHistory(info.ID)
	if errCopy != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCopy.Error()})
		return
	}

	//删除设备关联的硬件信息
	_, errManufacturer := repo.DeleteManufacturerBySn(device.Sn)
	if errManufacturer != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errManufacturer.Error()})
		return
	}

	mod, err := repo.DeleteDeviceById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

//重装
func BatchReInstall(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var infos []struct {
		ID          uint
		AccessToken string
		UserID      uint
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	for _, info := range infos {
		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		} else {
			info.UserID = session.ID
		}

		//log
		device, errDevice := repo.GetDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDevice.Error()})
			return
		}

		if session.Role != "Administrator" && device.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "您无权操作其他人的设备!"})
			return
		}

		_, err := repo.ReInstallDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		//删除PXE配置文件
		macs, errMac := repo.GetMacListByDeviceID(device.ID)
		if errMac != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errMac.Error()})
			return
		}
		for _, mac := range macs {
			pxeFileName := util.GetPxeFileNameByMac(mac.Mac)
			confDir := conf.OsInstall.PxeConfigDir
			if util.FileExist(confDir + "/" + pxeFileName) {
				err := os.Remove(confDir + "/" + pxeFileName)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
			}
		}

		logContent := make(map[string]interface{})
		logContent["data"] = device
		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(info.ID, "重装设备", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}

		_, errLog := repo.UpdateDeviceLogTypeByDeviceIdAndType(info.ID, "install", "install_history")
		if errLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
			return
		}

		/*
			//删除LOG
			_, errLog := repo.DeleteDeviceLogByDeviceID(info.ID)
			if errLog != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
				return
			}
		*/
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

//取消安装
func BatchCancelInstall(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var infos []struct {
		ID          uint
		AccessToken string
		UserID      uint
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	for _, info := range infos {
		device, err := repo.GetDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		} else {
			info.UserID = session.ID
		}

		if session.Role != "Administrator" && device.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "您无权操作其他人的设备!"})
			return
		}

		//安装成功的设备不允许取消安装
		if device.Status == "success" {
			continue
		}

		_, errCancel := repo.CancelInstallDeviceById(info.ID)
		if errCancel != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCancel.Error()})
			return
		}

		//删除PXE配置文件
		macs, errMac := repo.GetMacListByDeviceID(device.ID)
		if errMac != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errMac.Error()})
			return
		}
		for _, mac := range macs {
			pxeFileName := util.GetPxeFileNameByMac(mac.Mac)
			confDir := conf.OsInstall.PxeConfigDir
			if util.FileExist(confDir + "/" + pxeFileName) {
				err := os.Remove(confDir + "/" + pxeFileName)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
			}
		}

		logContent := make(map[string]interface{})
		logContent["data"] = device
		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(info.ID, "取消安装设备", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}

		_, errLog := repo.UpdateDeviceLogTypeByDeviceIdAndType(info.ID, "install", "install_history")
		if errLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
			return
		}

		/*
			//删除LOG
			_, errLog := repo.DeleteDeviceLogByDeviceID(info.ID)
			if errLog != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
				return
			}
		*/
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchDelete(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var infos []struct {
		ID          uint
		AccessToken string
		UserID      uint
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	for _, info := range infos {
		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		} else {
			info.UserID = session.ID
		}

		device, errInfo := repo.GetDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}

		if session.Role != "Administrator" && device.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "您无权操作其他人的设备!"})
			return
		}

		//删除PXE配置文件
		macs, errMac := repo.GetMacListByDeviceID(device.ID)
		if errMac != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errMac.Error()})
			return
		}
		for _, mac := range macs {
			pxeFileName := util.GetPxeFileNameByMac(mac.Mac)
			confDir := conf.OsInstall.PxeConfigDir
			if util.FileExist(confDir + "/" + pxeFileName) {
				err := os.Remove(confDir + "/" + pxeFileName)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
			}
		}

		//删除mac
		_, err := repo.DeleteMacByDeviceId(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		/*
			//删除LOG
			_, errLog := repo.DeleteDeviceLogByDeviceID(info.ID)
			if errLog != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
				return
			}
		*/
		//删除设备关联的硬件信息
		_, errManufacturer := repo.DeleteManufacturerBySn(device.Sn)
		if errManufacturer != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errManufacturer.Error()})
			return
		}

		errCopy := repo.CopyDeviceToHistory(info.ID)
		if errCopy != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCopy.Error()})
			return
		}
		_, errUpdate := repo.UpdateHistoryDeviceStatusById(info.ID, "delete")
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
			return
		}

		_, errDevice := repo.DeleteDeviceById(info.ID)
		if errDevice != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDevice.Error()})
			return
		}

		//log
		logContent := make(map[string]interface{})
		logContent["data"] = device
		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(device.ID, "删除设备信息", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func GetDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mod, err := repo.GetDeviceById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetDeviceBySn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		Sn string
	}
	info.Sn = r.FormValue("sn")
	info.Sn = strings.TrimSpace(info.Sn)

	count, err := repo.CountDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	if count <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "设备不存在!"})
		return
	}

	mod, err := repo.GetDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetFullDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mod, err := repo.GetFullDeviceById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type DeviceWithTime struct {
		ID                uint
		BatchNumber       string
		Sn                string
		Hostname          string
		Ip                string
		ManageIp          string
		NetworkID         uint
		ManageNetworkID   uint
		OsID              uint
		HardwareID        uint
		SystemID          uint
		Location          string
		LocationID        uint
		AssetNumber       string
		Status            string
		InstallProgress   float64
		InstallLog        string
		NetworkName       string
		ManageNetworkName string
		OsName            string
		HardwareName      string
		SystemName        string
		LocationName      string
		IsSupportVm       string
		UserID            uint
		CreatedAt         utils.ISOTime
		UpdatedAt         utils.ISOTime
	}

	var device DeviceWithTime
	device.ID = mod.ID
	device.BatchNumber = mod.BatchNumber
	device.Sn = mod.Sn
	device.Hostname = mod.Hostname
	device.Ip = mod.Ip
	device.ManageIp = mod.ManageIp
	device.NetworkID = mod.NetworkID
	device.ManageNetworkID = mod.ManageNetworkID
	device.OsID = mod.OsID
	device.HardwareID = mod.HardwareID
	device.SystemID = mod.SystemID
	device.Location = mod.Location
	device.LocationID = mod.LocationID
	device.AssetNumber = mod.AssetNumber
	device.Status = mod.Status
	device.InstallProgress = mod.InstallProgress
	device.InstallLog = mod.InstallLog
	device.NetworkName = mod.NetworkName
	device.ManageNetworkName = mod.ManageNetworkName
	device.OsName = mod.OsName
	device.HardwareName = mod.HardwareName
	device.SystemName = mod.SystemName
	device.IsSupportVm = mod.IsSupportVm
	device.UserID = mod.UserID
	device.LocationName, err = repo.FormatLocationNameById(mod.LocationID, "", "-")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	device.CreatedAt = utils.ISOTime(mod.CreatedAt)
	device.UpdatedAt = utils.ISOTime(mod.UpdatedAt)

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": device})
}

func GetDeviceList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Limit          uint
		Offset         uint
		Keyword        string
		OsID           int
		HardwareID     int
		SystemID       int
		Status         string
		BatchNumber    string
		StartUpdatedAt string
		EndUpdatedAt   string
		UserID         int
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Keyword = strings.TrimSpace(info.Keyword)
	info.Status = strings.TrimSpace(info.Status)
	info.BatchNumber = strings.TrimSpace(info.BatchNumber)

	var where string
	where = " where t1.id > 0 "
	if info.OsID > 0 {
		where += " and t1.os_id = " + strconv.Itoa(info.OsID)
	}
	if info.HardwareID > 0 {
		where += " and t1.hardware_id = " + strconv.Itoa(info.HardwareID)
	}
	if info.SystemID > 0 {
		where += " and t1.system_id = " + strconv.Itoa(info.SystemID)
	}
	if info.Status != "" {
		where += " and t1.status = '" + info.Status + "'"
	}
	if info.BatchNumber != "" {
		where += " and t1.batch_number = '" + info.BatchNumber + "'"
	}

	if info.StartUpdatedAt != "" {
		where += " and t1.updated_at >= '" + info.StartUpdatedAt + "'"
	}

	if info.EndUpdatedAt != "" {
		where += " and t1.updated_at <= '" + info.EndUpdatedAt + "'"
	}

	if info.UserID > 0 {
		where += " and t1.user_id = " + strconv.Itoa(info.UserID)
	}

	if info.Keyword != "" {
		where += " and ( "
		info.Keyword = strings.Replace(info.Keyword, "\n", ",", -1)
		info.Keyword = strings.Replace(info.Keyword, ";", ",", -1)
		list := strings.Split(info.Keyword, ",")
		for k, v := range list {
			var str string
			v = strings.TrimSpace(v)
			if k == 0 {
				str = ""
			} else {
				str = " or "
			}
			where += str + " t1.sn = '" + v + "' or t1.batch_number = '" + v + "' or t1.hostname = '" + v + "' or t1.ip = '" + v + "'"
		}
		where += " ) "
	}

	mods, err := repo.GetDeviceListWithPage(info.Limit, info.Offset, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type DeviceWithTime struct {
		ID              uint
		BatchNumber     string
		Sn              string
		Hostname        string
		Ip              string
		NetworkID       uint
		OsID            uint
		HardwareID      uint
		SystemID        uint
		Location        string
		LocationID      uint
		AssetNumber     string
		Status          string
		InstallProgress float64
		InstallLog      string
		NetworkName     string
		OsName          string
		HardwareName    string
		SystemName      string
		LocationName    string
		IsSupportVm     string
		UserID          uint
		OwnerName       string
		CreatedAt       utils.ISOTime
		UpdatedAt       utils.ISOTime
	}
	var rows []DeviceWithTime
	for _, v := range mods {
		var device DeviceWithTime
		device.ID = v.ID
		device.BatchNumber = v.BatchNumber
		device.Sn = v.Sn
		device.Hostname = v.Hostname
		device.Ip = v.Ip
		device.NetworkID = v.NetworkID
		device.OsID = v.OsID
		device.HardwareID = v.HardwareID
		device.SystemID = v.SystemID
		device.Location = v.Location
		device.LocationID = v.LocationID
		device.AssetNumber = v.AssetNumber
		device.Status = v.Status
		device.InstallProgress = v.InstallProgress
		device.InstallLog = v.InstallLog
		device.NetworkName = v.NetworkName
		device.OsName = v.OsName
		device.HardwareName = v.HardwareName
		device.SystemName = v.SystemName
		device.IsSupportVm = v.IsSupportVm
		device.UserID = v.UserID
		device.OwnerName = v.OwnerName
		/*
			device.LocationName, err = repo.FormatLocationNameById(v.LocationID, "", "-")
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
		*/

		device.CreatedAt = utils.ISOTime(v.CreatedAt)
		device.UpdatedAt = utils.ISOTime(v.UpdatedAt)

		deviceLog, _ := repo.GetLastDeviceLogByDeviceID(v.ID)
		device.InstallLog = deviceLog.Title
		rows = append(rows, device)
	}

	result := make(map[string]interface{})
	result["list"] = rows

	//总条数
	count, err := repo.CountDevice(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func GetDeviceNumByStatus(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Status string
		UserID int
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Status = strings.TrimSpace(info.Status)

	var where string
	where = " where t1.id > 0 "
	where += " and t1.status = '" + info.Status + "'"
	if info.UserID > 0 {
		where += " and t1.user_id = " + strconv.Itoa(info.UserID)
	}

	//总条数
	count, err := repo.CountDevice(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["count"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

//添加
func AddDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	var info struct {
		BatchNumber     string
		Sn              string
		Hostname        string
		Ip              string
		ManageIp        string
		NetworkID       uint
		ManageNetworkID uint
		OsID            uint
		HardwareID      uint
		SystemID        uint
		LocationID      uint
		AssetNumber     string
		IsSupportVm     string
		Status          string
		UserID          uint
		AccessToken     string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.BatchNumber = strings.TrimSpace(info.BatchNumber)
	info.Sn = strings.TrimSpace(info.Sn)
	info.Hostname = strings.TrimSpace(info.Hostname)
	info.Ip = strings.TrimSpace(info.Ip)
	info.ManageIp = strings.TrimSpace(info.ManageIp)
	info.AssetNumber = strings.TrimSpace(info.AssetNumber)
	info.IsSupportVm = strings.TrimSpace(info.IsSupportVm)
	info.Status = strings.TrimSpace(info.Status)
	info.AccessToken = strings.TrimSpace(info.AccessToken)
	info.UserID = session.ID

	if session.ID <= uint(0) {
		accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
		if errAccessToken != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
			return
		}
		info.UserID = accessTokenUser.ID
		session.ID = accessTokenUser.ID
		session.Role = accessTokenUser.Role
	}

	if info.Sn == "" || info.Hostname == "" || info.Ip == "" || info.NetworkID == uint(0) || info.SystemID == uint(0) || info.OsID == uint(0) {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	if info.IsSupportVm == "" {
		info.IsSupportVm = "Yes"
	}

	countDevice, err := repo.CountDeviceBySn(info.Sn)
	if countDevice > 0 {
		device, err := repo.GetDeviceBySn(info.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if session.Role != "Administrator" && device.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该设备已被录入，不能重复录入!"})
			return
		}

		deviceId, err := repo.GetDeviceIdBySn(info.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		count, err := repo.CountDeviceByHostnameAndId(info.Hostname, deviceId)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if count > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
			return
		}

		countIp, err := repo.CountDeviceByIpAndId(info.Ip, deviceId)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countIp > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
			return
		}

		if info.ManageIp != "" {
			countManageIp, err := repo.CountDeviceByManageIpAndId(info.ManageIp, deviceId)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countManageIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + " 该IP已存在!"})
				return
			}
		}
	} else {
		count, err := repo.CountDeviceByHostname(info.Hostname)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if count > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
			return
		}

		countIp, err := repo.CountDeviceByIp(info.Ip)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countIp > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
			return
		}

		if info.ManageIp != "" {
			countManageIp, err := repo.CountDeviceByManageIp(info.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countManageIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
				return
			}
		}
	}

	//匹配网络
	isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if !isValidate {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "IP格式不正确!"})
		return
	}

	modelIp, err := repo.GetIpByIp(info.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
		return
	}

	_, errNetwork := repo.GetNetworkById(modelIp.NetworkID)
	if errNetwork != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
		return
	}

	if info.ManageIp != "" {
		//匹配网络
		isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.ManageIp)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if !isValidate {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "IP格式不正确!"})
			return
		}

		modelIp, err := repo.GetManageIpByIp(info.ManageIp)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "未匹配到网段!"})
			return
		}

		_, errNetwork := repo.GetManageNetworkById(modelIp.NetworkID)
		if errNetwork != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "未匹配到网段!"})
			return
		}
	}

	//校验是否使用OOB静态IP及管理IP是否填写
	if info.HardwareID > uint(0) {
		hardware, err := repo.GetHardwareById(info.HardwareID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if hardware.Data != "" {
			if strings.Contains(hardware.Data, "<{manage_ip}>") || strings.Contains(hardware.Data, "<{manage_netmask}>") || strings.Contains(hardware.Data, "<{manage_gateway}>") {
				if info.ManageIp == "" || info.ManageNetworkID <= uint(0) {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备(SN:" + info.Sn + ")使用的硬件配置模板的OOB网络类型为静态IP的方式，请填写管理IP!"})
					return
				}
			}
		}
	}

	location := ""
	//SN已存在的情况下，要覆盖原数据
	count, err := repo.CountDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	//生成安装批次号
	batchNumber, err := repo.CreateBatchNumber()

	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
		return
	}
	status := "pre_install"
	//覆盖
	if count > 0 {
		id, err := repo.GetDeviceIdBySn(info.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		deviceOld, err := repo.GetDeviceById(id)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		_, errLog := repo.UpdateDeviceLogTypeByDeviceIdAndType(id, "install", "install_history")
		if errLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
			return
		}

		//log
		logContent := make(map[string]interface{})
		logContent["data_old"] = deviceOld

		device, errUpdate := repo.UpdateDeviceById(id, batchNumber, info.Sn, info.Hostname, info.Ip, info.ManageIp, info.NetworkID, info.ManageNetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm, info.UserID)
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
			return
		}

		logContent["data"] = device

		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(device.ID, "修改设备信息", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
	} else {
		device, err := repo.AddDevice(batchNumber, info.Sn, info.Hostname, info.Ip, info.ManageIp, info.NetworkID, info.ManageNetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm, info.UserID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		//log
		logContent := make(map[string]interface{})
		logContent["data"] = device
		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(device.ID, "录入新设备", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}

		//init manufactures device_id
		countManufacturer, errCountManufacturer := repo.CountManufacturerBySn(info.Sn)
		if errCountManufacturer != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCountManufacturer.Error()})
			return
		}
		if countManufacturer > 0 {
			manufacturerId, errGetManufacturerBySn := repo.GetManufacturerIdBySn(info.Sn)
			if errGetManufacturerBySn != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errGetManufacturerBySn.Error()})
				return
			}
			_, errUpdate := repo.UpdateManufacturerDeviceIdById(manufacturerId, device.ID)
			if errUpdate != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
				return
			}
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

//添加
func BatchAddDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var infos []struct {
		BatchNumber     string
		Sn              string
		Hostname        string
		Ip              string
		ManageIp        string
		NetworkID       uint
		ManageNetworkID uint
		OsID            uint
		HardwareID      uint
		SystemID        uint
		LocationID      uint
		AssetNumber     string
		IsSupportVm     string
		Status          string
		UserID          uint
		AccessToken     string
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	//先批量检测传入数据是否有问题
	for _, info := range infos {
		info.BatchNumber = strings.TrimSpace(info.BatchNumber)
		info.Sn = strings.TrimSpace(info.Sn)
		info.Sn = strings.Replace(info.Sn, "	", "", -1)
		info.Sn = strings.Replace(info.Sn, " ", "", -1)
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.ManageIp = strings.TrimSpace(info.ManageIp)
		info.AssetNumber = strings.TrimSpace(info.AssetNumber)
		info.Status = strings.TrimSpace(info.Status)
		info.AccessToken = strings.TrimSpace(info.AccessToken)
		info.UserID = session.ID

		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		}

		if info.Sn == "" || info.Hostname == "" || info.Ip == "" || info.NetworkID == uint(0) || info.OsID == uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
			return
		}

		count, err := repo.CountDeviceBySn(info.Sn)
		if count > 0 {
			device, err := repo.GetDeviceBySn(info.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}

			if session.Role != "Administrator" && device.UserID != session.ID {
				w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该设备已被录入，不能重复录入!"})
				return
			}

			deviceId, err := repo.GetDeviceIdBySn(info.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
			count, err := repo.CountDeviceByHostnameAndId(info.Hostname, deviceId)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if count > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
				return
			}

			countIp, err := repo.CountDeviceByIpAndId(info.Ip, deviceId)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
				return
			}

			if info.ManageIp != "" {
				countManageIp, err := repo.CountDeviceByManageIpAndId(info.ManageIp, deviceId)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if countManageIp > 0 {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + " 该IP已存在!"})
					return
				}
			}
		} else {
			count, err := repo.CountDeviceByHostname(info.Hostname)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if count > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
				return
			}

			countIp, err := repo.CountDeviceByIp(info.Ip)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
				return
			}

			if info.ManageIp != "" {
				countManageIp, err := repo.CountDeviceByManageIp(info.ManageIp)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if countManageIp > 0 {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
					return
				}
			}
		}

		//匹配网络
		isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.Ip)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if !isValidate {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "IP格式不正确!"})
			return
		}

		modelIp, err := repo.GetIpByIp(info.Ip)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
			return
		}

		_, errNetwork := repo.GetNetworkById(modelIp.NetworkID)
		if errNetwork != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
			return
		}

		if info.ManageIp != "" {
			//匹配网络
			isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}

			if !isValidate {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "IP格式不正确!"})
				return
			}

			modelIp, err := repo.GetManageIpByIp(info.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "未匹配到网段!"})
				return
			}

			_, errNetwork := repo.GetManageNetworkById(modelIp.NetworkID)
			if errNetwork != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
				return
			}
		}

		//校验是否使用OOB静态IP及管理IP是否填写
		if info.HardwareID > uint(0) {
			hardware, err := repo.GetHardwareById(info.HardwareID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
			if hardware.Data != "" {
				if strings.Contains(hardware.Data, "<{manage_ip}>") || strings.Contains(hardware.Data, "<{manage_netmask}>") || strings.Contains(hardware.Data, "<{manage_gateway}>") {
					if info.ManageIp == "" || info.ManageNetworkID <= uint(0) {
						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备(SN:" + info.Sn + ")使用的硬件配置模板的OOB网络类型为静态IP的方式，请填写管理IP!"})
						return
					}
				}
			}
		}
	}

	//生成安装批次号
	batchNumber, err := repo.CreateBatchNumber()

	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
		return
	}
	status := "pre_install"
	for _, info := range infos {
		info.BatchNumber = strings.TrimSpace(info.BatchNumber)
		info.Sn = strings.TrimSpace(info.Sn)
		info.Sn = strings.Replace(info.Sn, "	", "", -1)
		info.Sn = strings.Replace(info.Sn, " ", "", -1)
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.ManageIp = strings.TrimSpace(info.ManageIp)
		info.AssetNumber = strings.TrimSpace(info.AssetNumber)
		info.Status = strings.TrimSpace(info.Status)
		info.UserID = session.ID
		location := ""

		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
			session.ID = accessTokenUser.ID
		}

		info.IsSupportVm = strings.TrimSpace(info.IsSupportVm)
		if info.IsSupportVm == "" {
			info.IsSupportVm = "Yes"
		}

		//SN已存在的情况下，要覆盖原数据
		count, err := repo.CountDeviceBySn(info.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		//覆盖
		if count > 0 {
			id, err := repo.GetDeviceIdBySn(info.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			deviceOld, err := repo.GetDeviceById(id)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			_, errLog := repo.UpdateDeviceLogTypeByDeviceIdAndType(id, "install", "install_history")
			if errLog != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
				return
			}

			//log
			logContent := make(map[string]interface{})
			logContent["data_old"] = deviceOld

			device, errUpdate := repo.UpdateDeviceById(id, batchNumber, info.Sn, info.Hostname, info.Ip, info.ManageIp, info.NetworkID, info.ManageNetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm, info.UserID)
			if errUpdate != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
				return
			}
			logContent["data"] = device

			json, err := json.Marshal(logContent)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
				return
			}

			_, errAddLog := repo.AddDeviceLog(device.ID, "修改设备信息", "operate", string(json))
			if errAddLog != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
				return
			}
		} else {
			device, err := repo.AddDevice(batchNumber, info.Sn, info.Hostname, info.Ip, info.ManageIp, info.NetworkID, info.ManageNetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm, info.UserID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
				return
			}

			//log
			logContent := make(map[string]interface{})
			logContent["data"] = device
			json, err := json.Marshal(logContent)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
				return
			}

			_, errAddLog := repo.AddDeviceLog(device.ID, "录入新设备", "operate", string(json))
			if errAddLog != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
				return
			}

			//init manufactures device_id
			countManufacturer, errCountManufacturer := repo.CountManufacturerBySn(info.Sn)
			if errCountManufacturer != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCountManufacturer.Error()})
				return
			}
			if countManufacturer > 0 {
				manufacturerId, errGetManufacturerBySn := repo.GetManufacturerIdBySn(info.Sn)
				if errGetManufacturerBySn != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errGetManufacturerBySn.Error()})
					return
				}
				_, errUpdate := repo.UpdateManufacturerDeviceIdById(manufacturerId, device.ID)
				if errUpdate != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
					return
				}
			}
		}

	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchUpdateDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	var infos []struct {
		ID              uint
		Hostname        string
		Ip              string
		ManageIp        string
		NetworkID       uint
		ManageNetworkID uint
		OsID            uint
		HardwareID      uint
		SystemID        uint
		LocationID      uint
		IsSupportVm     string
		AssetNumber     string
		UserID          uint
		AccessToken     string
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	//先批量检测传入数据是否有问题
	for _, info := range infos {
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.ManageIp = strings.TrimSpace(info.ManageIp)
		info.AssetNumber = strings.TrimSpace(info.AssetNumber)
		info.AccessToken = strings.TrimSpace(info.AccessToken)
		info.UserID = session.ID

		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		}

		if info.Hostname == "" || info.Ip == "" || info.NetworkID == uint(0) || info.OsID == uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
			return
		}

		device, err := repo.GetDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if session.Role != "Administrator" && device.UserID != session.ID {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该设备已被录入，不能重复录入!"})
			return
		}

		count, err := repo.CountDeviceByHostnameAndId(info.Hostname, info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if count > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
			return
		}

		countIp, err := repo.CountDeviceByIpAndId(info.Ip, info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countIp > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
			return
		}

		if info.ManageIp != "" {
			countManageIp, err := repo.CountDeviceByManageIpAndId(info.ManageIp, info.ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countManageIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + " 该管理IP已存在!"})
				return
			}
		}

		//匹配网络
		isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.Ip)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if !isValidate {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "IP格式不正确!"})
			return
		}

		modelIp, err := repo.GetIpByIp(info.Ip)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
			return
		}

		_, errNetwork := repo.GetNetworkById(modelIp.NetworkID)
		if errNetwork != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "未匹配到网段!"})
			return
		}

		if info.ManageIp != "" {
			//匹配网络
			isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}

			if !isValidate {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "IP格式不正确!"})
				return
			}

			modelIp, err := repo.GetManageIpByIp(info.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "未匹配到网段!"})
				return
			}

			_, errNetwork := repo.GetManageNetworkById(modelIp.NetworkID)
			if errNetwork != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.ManageIp + "未匹配到网段!"})
				return
			}
		}

		//校验是否使用OOB静态IP及管理IP是否填写
		if info.HardwareID > uint(0) {
			hardware, err := repo.GetHardwareById(info.HardwareID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
			if hardware.Data != "" {
				if strings.Contains(hardware.Data, "<{manage_ip}>") || strings.Contains(hardware.Data, "<{manage_netmask}>") || strings.Contains(hardware.Data, "<{manage_gateway}>") {
					if info.ManageIp == "" || info.ManageNetworkID <= uint(0) {
						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备(SN:" + device.Sn + ")使用的硬件配置模板的OOB网络类型为静态IP的方式，请填写管理IP!"})
						return
					}
				}
			}
		}

	}

	for _, info := range infos {
		location := ""
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.ManageIp = strings.TrimSpace(info.ManageIp)
		info.AssetNumber = strings.TrimSpace(info.AssetNumber)
		info.AccessToken = strings.TrimSpace(info.AccessToken)
		info.UserID = session.ID

		if session.ID <= uint(0) {
			accessTokenUser, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
			if errAccessToken != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
				return
			}
			info.UserID = accessTokenUser.ID
			session.ID = accessTokenUser.ID
			session.Role = accessTokenUser.Role
		}

		info.IsSupportVm = strings.TrimSpace(info.IsSupportVm)
		if info.IsSupportVm == "" {
			info.IsSupportVm = "Yes"
		}

		device, err := repo.GetDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		//log
		logContent := make(map[string]interface{})
		logContent["data_old"] = device

		deviceNew, errUpdate := repo.UpdateDeviceById(info.ID, device.BatchNumber, device.Sn, info.Hostname, info.Ip, info.ManageIp, info.NetworkID, info.ManageNetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, device.Status, info.IsSupportVm, info.UserID)
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
			return
		}

		_, errLog := repo.UpdateDeviceLogTypeByDeviceIdAndType(info.ID, "install", "install_history")
		if errLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
			return
		}

		logContent["data"] = deviceNew

		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(device.ID, "修改设备信息", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

//上报安装进度
func ReportInstallInfo(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	if conf.OsInstall.PxeConfigDir == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "Pxe配置文件目录没有指定"})
		return
	}

	var info struct {
		Sn              string
		Title           string
		InstallProgress float64
		InstallLog      string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.Sn = strings.TrimSpace(info.Sn)
	info.Title = strings.TrimSpace(info.Title)
	if info.Sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空!"})
		return
	}

	deviceId, err := repo.GetDeviceIdBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备不存在!"})
		return
	}

	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if device.Status != "pre_install" && device.Status != "installing" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备不在安装列表里!"})
		return
	}

	var status string
	var logTitle string

	if info.InstallProgress == -1 {
		status = "failure"
		info.InstallProgress = 0
		logTitle = info.Title
	} else if info.InstallProgress >= 0 && info.InstallProgress < 1 {
		status = "installing"
		logTitle = info.Title + "(" + fmt.Sprintf("安装进度 %.1f", info.InstallProgress*100) + "%)"
	} else if info.InstallProgress == 1 {
		status = "success"
		logTitle = info.Title + "(" + fmt.Sprintf("安装进度 %.1f", info.InstallProgress*100) + "%)"
		//logTitle = "安装成功"
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "安装进度参数不正确!"})
		return
	}

	/*
		if device.InstallLog != "" {
			info.InstallLog = device.InstallLog + "\n" + info.InstallLog
		}
	*/
	_, errUpdate := repo.UpdateInstallInfoById(device.ID, status, info.InstallProgress)
	if errUpdate != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
		return
	}

	//删除PXE配置文件
	if info.InstallProgress == 1 {
		macs, err := repo.GetMacListByDeviceID(device.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		for _, mac := range macs {
			pxeFileName := util.GetPxeFileNameByMac(mac.Mac)
			confDir := conf.OsInstall.PxeConfigDir
			if util.FileExist(confDir + "/" + pxeFileName) {
				err := os.Remove(confDir + "/" + pxeFileName)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
			}
		}
	}

	var installLog string
	byteDecode, err := base64.StdEncoding.DecodeString(info.InstallLog)
	if err != nil {
		installLog = ""
	} else {
		installLog = string(byteDecode)
	}

	_, errAddLog := repo.AddDeviceLog(device.ID, logTitle, "install", installLog)
	if errAddLog != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
		return
	}

	//add report
	if info.InstallProgress == 1 {
		errReportLog := repo.CopyDeviceToInstallReport(device.ID)
		if errReportLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errReportLog.Error()})
			return
		}
	}

	result := make(map[string]string)
	result["Result"] = "true"
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

//上报Mac信息，生成Pxe文件
func ReportMacInfo(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	if conf.OsInstall.PxeConfigDir == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "Pxe配置文件目录没有指定"})
		return
	}

	var info struct {
		Sn  string
		Mac string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.Sn = strings.TrimSpace(info.Sn)
	info.Mac = strings.TrimSpace(info.Mac)

	if info.Sn == "" || info.Mac == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN和Mac参数不能为空!"})
		return
	}

	//mac 大写转为 小写
	info.Mac = strings.ToLower(info.Mac)

	deviceId, err := repo.GetDeviceIdBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备不存在!"})
		return
	}

	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	osConfig, err := repo.GetOsConfigById(device.OsID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "PXE信息没有配置" + err.Error()})
		return
	}

	//录入Mac信息
	count, err := repo.CountMacByMacAndDeviceID(info.Mac, device.ID)
	if count <= 0 {
		count, err := repo.CountMacByMac(info.Mac)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if count > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该MAC地址已被其他设备录入"})
			return
		}

		_, errAddMac := repo.AddMac(device.ID, info.Mac)
		if errAddMac != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddMac.Error()})
			return
		}
	}

	//替换占位符
	osConfig.Pxe = strings.Replace(osConfig.Pxe, "{sn}", info.Sn, -1)

	pxeFileName := util.GetPxeFileNameByMac(info.Mac)
	errCreatePxeFile := util.CreatePxeFile(conf.OsInstall.PxeConfigDir, pxeFileName, osConfig.Pxe)
	if errCreatePxeFile != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "配置文件生成失败" + err.Error()})
		return
	}

	result := make(map[string]string)
	result["Result"] = "true"
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func IsInPreInstallList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误", "Content": ""})
		return
	}
	var info struct {
		Sn string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误", "Content": ""})
		return
	}

	info.Sn = strings.TrimSpace(info.Sn)

	if info.Sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空!"})
		return
	}

	deviceId, err := repo.GetDeviceIdBySn(info.Sn)
	result := make(map[string]string)
	if err != nil {
		result["Result"] = "false"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备不在安装列表里", "Content": result})
		return
	}

	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		result["Result"] = "false"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备不在安装列表里", "Content": result})
		return
	}

	if device.Status == "pre_install" || device.Status == "installing" {
		result["Result"] = "true"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备在安装列表里", "Content": result})
	} else {
		result["Result"] = "false"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备不在安装列表里", "Content": result})
	}
}

func GetHardwareBySn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	//repo := middleware.RepoFromContext(ctx)
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		//rest.Error(w, " ,", http.StatusInternalServerError)
		//w.WriteHeader(http.StatusFound)
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误", "Content": ""})
		return
	}
	var info struct {
		Sn string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		//rest.Error(w, " ", http.status)
		//w.WriteHeader(http.StatusFound)
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误", "Content": ""})
		return
	}

	info.Sn = strings.TrimSpace(info.Sn)

	if info.Sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空!"})
		return
	}

	device, err := repo.GetDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		return
	}

	var manageNetwork model.ManageNetwork
	if device.ManageNetworkID > 0 {
		manageNetworkDetail, err := repo.GetManageNetworkById(device.ManageNetworkID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
			return
		}
		manageNetwork.Netmask = manageNetworkDetail.Netmask
		manageNetwork.Gateway = manageNetworkDetail.Gateway
	}

	hardware, err := repo.GetHardwareBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		return
	}

	type ChildData struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	}

	type ScriptData struct {
		Name string       `json:"Name"`
		Data []*ChildData `json:"Data"`
	}

	var data []*ScriptData
	var result2 []map[string]interface{}
	if hardware.Data != "" {
		bytes := []byte(hardware.Data)
		errData := json.Unmarshal(bytes, &data)
		if errData != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
			return
		}

		for _, v := range data {
			result3 := make(map[string]interface{})
			result3["Name"] = v.Name
			var result5 []map[string]interface{}
			for _, v2 := range v.Data {
				result4 := make(map[string]interface{})
				if strings.Contains(v2.Value, "<{manage_ip}>") {
					v2.Value = strings.Replace(v2.Value, "<{manage_ip}>", device.ManageIp, -1)
				}
				if strings.Contains(v2.Value, "<{manage_netmask}>") {
					v2.Value = strings.Replace(v2.Value, "<{manage_netmask}>", manageNetwork.Netmask, -1)
				}
				if strings.Contains(v2.Value, "<{manage_gateway}>") {
					v2.Value = strings.Replace(v2.Value, "<{manage_gateway}>", manageNetwork.Gateway, -1)
				}

				result4["Name"] = v2.Name
				result4["Script"] = base64.StdEncoding.EncodeToString([]byte(v2.Value))
				result5 = append(result5, result4)
			}
			result3["Scripts"] = result5
			result2 = append(result2, result3)
		}
	}

	result := make(map[string]interface{})
	result["Company"] = hardware.Company
	result["Product"] = hardware.Product
	result["ModelName"] = hardware.ModelName

	result["Hardware"] = result2

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "成功获取hardware信息", "Content": result})
}

func GetSystemBySn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	var info struct {
		Sn   string
		Type string
	}

	info.Sn = r.FormValue("sn")
	info.Type = r.FormValue("type")
	info.Sn = strings.TrimSpace(info.Sn)
	info.Type = strings.TrimSpace(info.Type)

	if info.Type == "" {
		info.Type = "raw"
	}

	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误", "Content": ""})
		}
		return
	}

	if info.Sn == "" {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空"})
		}
		return
	}

	mod, err := repo.GetSystemBySn(info.Sn)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}

		return
	}

	if info.Type == "raw" {
		w.Header().Add("Content-type", "text/html; charset=utf-8")
		w.Write([]byte(mod.Content))
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "成功获取system信息", "Content": mod})
	}
}

func GetNetworkBySn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	var info struct {
		Sn   string
		Type string
	}

	info.Sn = r.FormValue("sn")
	info.Type = r.FormValue("type")
	info.Sn = strings.TrimSpace(info.Sn)
	info.Type = strings.TrimSpace(info.Type)

	if info.Type == "" {
		info.Type = "raw"
	}

	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误", "Content": ""})
		}
		return
	}

	if info.Sn == "" {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空"})
		}
		return
	}

	deviceId, err := repo.GetDeviceIdBySn(info.Sn)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}
		return
	}

	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}
		return
	}

	mod, err := repo.GetNetworkBySn(info.Sn)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}
		return
	}

	result := make(map[string]interface{})
	result["Hostname"] = device.Hostname
	result["Ip"] = device.Ip
	result["Netmask"] = mod.Netmask
	result["Gateway"] = mod.Gateway
	result["Vlan"] = mod.Vlan
	result["Trunk"] = mod.Trunk
	result["Bonding"] = mod.Bonding
	if info.Type == "raw" {
		w.Header().Add("Content-type", "text/html; charset=utf-8")
		var str string
		if device.Hostname != "" {
			str += "HOSTNAME=\"" + device.Hostname + "\""
		}
		if device.Ip != "" {
			str += "\nIPADDR=\"" + device.Ip + "\""
		}
		if mod.Netmask != "" {
			str += "\nNETMASK=\"" + mod.Netmask + "\""
		}
		if mod.Gateway != "" {
			str += "\nGATEWAY=\"" + mod.Gateway + "\""
		}
		if mod.Vlan != "" {
			str += "\nVLAN=\"" + mod.Vlan + "\""
		}
		if mod.Trunk != "" {
			str += "\nTrunk=\"" + mod.Trunk + "\""
		}
		if mod.Bonding != "" {
			str += "\nBonding=\"" + mod.Bonding + "\""
		}
		w.Write([]byte(str))
	} else {
		w.Header().Add("Content-type", "application/json; charset=utf-8")
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "成功获取network信息", "Content": result})
	}
}

func ValidateSn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Sn string
	}
	info.Sn = strings.TrimSpace(info.Sn)

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if info.Sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "SN参数不能为空!", "Content": ""})
		return
	}

	count, err := repo.CountDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
		return
	}

	if count > 0 {
		session, err := GetSession(w, r)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		device, err := repo.GetDeviceBySn(info.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if session.Role != "Administrator" {
			if device.UserID != session.ID {
				w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该设备已被录入，不能重复录入!"})
				return
			}
		}

		if device.Status == "success" {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该设备已安装成功，确定要重装？"})
			return
		}

		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该SN已存在，继续填写会覆盖旧的数据!"})
		return

	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "SN填写正确!"})
		return
	}

}

func ImportDeviceForOpenApi(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	type Device struct {
		ID              uint
		BatchNumber     string
		Sn              string
		Hostname        string
		Ip              string
		ManageIp        string
		NetworkID       uint
		ManageNetworkID uint
		OsID            uint
		HardwareID      uint
		SystemID        uint
		Location        string
		LocationID      uint
		AssetNumber     string
		Status          string
		InstallProgress float64
		InstallLog      string
		NetworkName     string
		OsName          string
		HardwareName    string
		SystemName      string
		Content         string
		IsSupportVm     string
		UserID          uint
		AccessToken     string
	}

	var row Device
	if err := r.DecodeJSONPayload(&row); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error(), "Content": ""})
		return
	}

	accessTokenUser, errAccessToken := VerifyAccessToken(row.AccessToken, ctx, false)
	if errAccessToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
		return
	}
	row.UserID = accessTokenUser.ID

	row.Sn = strings.TrimSpace(row.Sn)
	row.Hostname = strings.TrimSpace(row.Hostname)
	row.Ip = strings.TrimSpace(row.Ip)
	row.ManageIp = strings.TrimSpace(row.ManageIp)
	row.HardwareName = strings.TrimSpace(row.HardwareName)
	row.SystemName = strings.TrimSpace(row.SystemName)
	row.OsName = strings.TrimSpace(row.OsName)
	row.AssetNumber = strings.TrimSpace(row.AssetNumber)

	batchNumber, err := repo.CreateBatchNumber()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if len(row.Sn) > 255 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN长度超过255限制"})
		return
	}

	if len(row.Hostname) > 255 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "主机名长度超过255限制"})
		return
	}

	if len(row.Location) > 255 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "位置长度超过255限制"})
		return
	}

	if len(row.AssetNumber) > 255 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "财编长度超过255限制"})
		return
	}

	if row.Sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN不能为空"})
		return
	}

	if row.Hostname == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "主机名不能为空"})
		return
	}

	if row.Ip == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "IP不能为空"})
		return
	}

	if row.OsName == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作系统模板名称不能为空"})
		return
	}

	if row.SystemName == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "系统安装模板不能为空"})
		return
	}

	if row.Location == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "位置不能为空"})
		return
	}

	countDevice, err := repo.CountDeviceBySn(row.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if countDevice > 0 {
		ID, err := repo.GetDeviceIdBySn(row.Sn)
		row.ID = ID
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		device, errDevice := repo.GetDeviceBySn(row.Sn)
		if errDevice != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if accessTokenUser.Role != "Administrator" && device.UserID != accessTokenUser.ID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备已被其他人录入，不能重复录入"})
			return
		} else {
			if device.Status == "success" {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备已安装成功，请使用【单台录入】的功能重新录入并安装"})
				return
			}
		}

		//hostname
		countHostname, err := repo.CountDeviceByHostnameAndId(row.Hostname, ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
			return
		}
		if countHostname > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该主机名已存在"})
			return
		}

		//IP
		countIp, err := repo.CountDeviceByIpAndId(row.Ip, ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countIp > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该IP已存在"})
			return
		}

		if row.ManageIp != "" {
			//IP
			countManageIp, err := repo.CountDeviceByManageIpAndId(row.ManageIp, ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countManageIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该管理IP已存在"})
				return
			}
		}
	} else {
		//hostname
		countHostname, err := repo.CountDeviceByHostname(row.Hostname)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
			return
		}
		if countHostname > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该主机名已存在"})
			return
		}

		//IP
		countIp, err := repo.CountDeviceByIp(row.Ip)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countIp > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该IP已存在"})
			return
		}

		if row.ManageIp != "" {
			//IP
			countManageIp, err := repo.CountDeviceByManageIp(row.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countManageIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该管理IP已存在"})
				return
			}
		}
	}

	//匹配网络
	isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", row.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if !isValidate {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "IP格式不正确"})
		return
	}

	modelIp, err := repo.GetIpByIp(row.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到网段"})
		return
	}

	_, errNetwork := repo.GetNetworkById(modelIp.NetworkID)
	if errNetwork != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到网段"})
		return
	}

	row.NetworkID = modelIp.NetworkID

	if row.ManageIp != "" {
		//匹配网络
		isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", row.ManageIp)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}

		if !isValidate {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "管理IP格式不正确"})
			return
		}

		modelIp, err := repo.GetManageIpByIp(row.ManageIp)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到管理网段"})
			return
		}

		_, errNetwork := repo.GetManageNetworkById(modelIp.NetworkID)
		if errNetwork != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到管理网段"})
			return
		}

		row.ManageNetworkID = modelIp.NetworkID
	}

	//OSName
	countOs, err := repo.CountOsConfigByName(row.OsName)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if countOs <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到操作系统"})
		return
	}
	mod, err := repo.GetOsConfigByName(row.OsName)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	row.OsID = mod.ID

	//SystemName
	countSystem, err := repo.CountSystemConfigByName(row.SystemName)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if countSystem <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到系统安装模板"})
		return
	}

	systemId, err := repo.GetSystemConfigIdByName(row.SystemName)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	row.SystemID = systemId

	if row.HardwareName != "" {
		//HardwareName
		countHardware, err := repo.CountHardwareWithSeparator(row.HardwareName)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countHardware <= 0 {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到硬件配置模板"})
			return
		}

		hardware, err := repo.GetHardwareBySeaprator(row.HardwareName)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		row.HardwareID = hardware.ID
	}

	if row.HardwareID > uint(0) {
		hardware, err := repo.GetHardwareById(row.HardwareID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if hardware.Data != "" {
			if strings.Contains(hardware.Data, "<{manage_ip}>") || strings.Contains(hardware.Data, "<{manage_netmask}>") || strings.Contains(hardware.Data, "<{manage_gateway}>") {
				if row.ManageIp == "" || row.ManageNetworkID <= uint(0) {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN:" + row.Sn + "硬件配置模板的OOB网络类型为静态IP的方式，请填写管理IP!"})
					return
				}
			}
		}
	}

	if row.Location != "" {
		countLocation, err := repo.CountLocationByName(row.Location)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		if countLocation > 0 {
			locationId, err := repo.GetLocationIdByName(row.Location)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
			row.LocationID = locationId
		}
		if row.LocationID <= uint(0) {
			locationId, err := repo.ImportLocation(row.Location)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if locationId <= uint(0) {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到位置"})
				return
			}
			row.LocationID = locationId
		}
	}
	status := "pre_install"
	row.IsSupportVm = "Yes"
	if countDevice > 0 {
		id, err := repo.GetDeviceIdBySn(row.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		deviceOld, err := repo.GetDeviceById(id)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		_, errLog := repo.UpdateDeviceLogTypeByDeviceIdAndType(id, "install", "install_history")
		if errLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errLog.Error()})
			return
		}

		device, errUpdate := repo.UpdateDeviceById(id, batchNumber, row.Sn, row.Hostname, row.Ip, row.ManageIp, row.NetworkID, row.ManageNetworkID, row.OsID, row.HardwareID, row.SystemID, "", row.LocationID, row.AssetNumber, status, row.IsSupportVm, row.UserID)
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
			return
		}

		//log
		logContent := make(map[string]interface{})
		logContent["data_old"] = deviceOld
		logContent["data"] = device

		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(device.ID, "修改设备信息", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
	} else {
		device, err := repo.AddDevice(batchNumber, row.Sn, row.Hostname, row.Ip, row.ManageIp, row.NetworkID, row.ManageNetworkID, row.OsID, row.HardwareID, row.SystemID, "", row.LocationID, row.AssetNumber, status, row.IsSupportVm, row.UserID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		//log
		logContent := make(map[string]interface{})
		logContent["data"] = device
		json, err := json.Marshal(logContent)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		_, errAddLog := repo.AddDeviceLog(device.ID, "录入新设备", "operate", string(json))
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}

		//init manufactures device_id
		countManufacturer, errCountManufacturer := repo.CountManufacturerBySn(row.Sn)
		if errCountManufacturer != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCountManufacturer.Error()})
			return
		}
		if countManufacturer > 0 {
			manufacturerId, errGetManufacturerBySn := repo.GetManufacturerIdBySn(row.Sn)
			if errGetManufacturerBySn != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errGetManufacturerBySn.Error()})
				return
			}
			_, errUpdate := repo.UpdateManufacturerDeviceIdById(manufacturerId, device.ID)
			if errUpdate != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
				return
			}
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
