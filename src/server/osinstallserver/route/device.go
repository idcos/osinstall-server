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
	//"model"
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
	var info struct {
		ID uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
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
	var infos []struct {
		ID uint
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	for _, info := range infos {
		_, err := repo.ReInstallDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		//log
		device, err := repo.GetDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
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

func BatchDelete(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var infos []struct {
		ID uint
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	for _, info := range infos {
		device, errInfo := repo.GetDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
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
		CreatedAt       utils.ISOTime
		UpdatedAt       utils.ISOTime
	}

	var device DeviceWithTime
	device.ID = mod.ID
	device.BatchNumber = mod.BatchNumber
	device.Sn = mod.Sn
	device.Hostname = mod.Hostname
	device.Ip = mod.Ip
	device.NetworkID = mod.NetworkID
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
	device.OsName = mod.OsName
	device.HardwareName = mod.HardwareName
	device.SystemName = mod.SystemName
	device.IsSupportVm = mod.IsSupportVm
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
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Status = strings.TrimSpace(info.Status)

	var where string
	where = " where t1.id > 0 "
	where += " and t1.status = '" + info.Status + "'"

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
	var info struct {
		BatchNumber string
		Sn          string
		Hostname    string
		Ip          string
		NetworkID   uint
		OsID        uint
		HardwareID  uint
		SystemID    uint
		LocationID  uint
		AssetNumber string
		IsSupportVm string
		Status      string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.BatchNumber = strings.TrimSpace(info.BatchNumber)
	info.Sn = strings.TrimSpace(info.Sn)
	info.Hostname = strings.TrimSpace(info.Hostname)
	info.Ip = strings.TrimSpace(info.Ip)
	info.AssetNumber = strings.TrimSpace(info.AssetNumber)
	info.IsSupportVm = strings.TrimSpace(info.IsSupportVm)
	info.Status = strings.TrimSpace(info.Status)

	if info.Sn == "" || info.Hostname == "" || info.Ip == "" || info.NetworkID == uint(0) || info.SystemID == uint(0) || info.OsID == uint(0) {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	if info.IsSupportVm == "" {
		info.IsSupportVm = "Yes"
	}

	countDevice, err := repo.CountDeviceBySn(info.Sn)
	if countDevice > 0 {
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

		device, errUpdate := repo.UpdateDeviceById(id, batchNumber, info.Sn, info.Hostname, info.Ip, info.NetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm)
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
		device, err := repo.AddDevice(batchNumber, info.Sn, info.Hostname, info.Ip, info.NetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm)
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
		BatchNumber string
		Sn          string
		Hostname    string
		Ip          string
		NetworkID   uint
		OsID        uint
		HardwareID  uint
		SystemID    uint
		LocationID  uint
		AssetNumber string
		IsSupportVm string
		Status      string
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
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
		info.AssetNumber = strings.TrimSpace(info.AssetNumber)
		info.Status = strings.TrimSpace(info.Status)

		if info.Sn == "" || info.Hostname == "" || info.Ip == "" || info.NetworkID == uint(0) || info.OsID == uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
			return
		}

		count, err := repo.CountDeviceBySn(info.Sn)
		if count > 0 {
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

	}

	//生成安装批次号
	batchNumber, err := repo.CreateBatchNumber()

	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
		return
	}
	status := "pre_install"
	for _, info := range infos {
		location := ""

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

			device, errUpdate := repo.UpdateDeviceById(id, batchNumber, info.Sn, info.Hostname, info.Ip, info.NetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm)
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
			device, err := repo.AddDevice(batchNumber, info.Sn, info.Hostname, info.Ip, info.NetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, status, info.IsSupportVm)
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
	var infos []struct {
		ID          uint
		Hostname    string
		Ip          string
		NetworkID   uint
		OsID        uint
		HardwareID  uint
		SystemID    uint
		LocationID  uint
		IsSupportVm string
		AssetNumber string
	}

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	//先批量检测传入数据是否有问题
	for _, info := range infos {
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.AssetNumber = strings.TrimSpace(info.AssetNumber)

		if info.Hostname == "" || info.Ip == "" || info.NetworkID == uint(0) || info.OsID == uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
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

	}

	for _, info := range infos {
		location := ""

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

		deviceNew, errUpdate := repo.UpdateDeviceById(info.ID, device.BatchNumber, device.Sn, info.Hostname, info.Ip, info.NetworkID, info.OsID, info.HardwareID, info.SystemID, location, info.LocationID, info.AssetNumber, device.Status, info.IsSupportVm)
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
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该MAC地址已被其他机器录入"})
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

//上报厂商信息
func ReportProductInfo(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Sn        string
		Company   string
		Product   string
		ModelName string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.Sn = strings.TrimSpace(info.Sn)
	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)
	info.ModelName = strings.TrimSpace(info.ModelName)

	if info.Sn == "" || info.Company == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
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

	count, err := repo.CountManufacturerByDeviceID(device.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	result := make(map[string]string)
	if count > 0 {
		manufacturer, err := repo.GetManufacturerByDeviceID(device.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": result})
			return
		}
		_, errUpdate := repo.UpdateManufacturerById(manufacturer.ID, info.Company, info.Product, info.ModelName)
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error(), "Content": result})
			return
		}

	} else {
		_, err := repo.AddManufacturer(device.ID, info.Company, info.Product, info.ModelName)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": result})
			return
		}
	}

	//校验是否在配置库
	isValidate, err := repo.ValidateHardwareProductModel(info.Company, info.Product, info.ModelName)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": result})
		return
	}
	if isValidate == true {
		result["IsVerify"] = "true"
	} else {
		result["IsVerify"] = "false"
	}

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

	/*
		resultHardware := make(map[string]string)
		resultHardware["Raid"] = base64.StdEncoding.EncodeToString([]byte(hardware.Raid))
		resultHardware["Oob"] = base64.StdEncoding.EncodeToString([]byte(hardware.Oob))
		resultHardware["Bios"] = base64.StdEncoding.EncodeToString([]byte(hardware.Bios))
	*/
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
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该SN已存在，继续填写会覆盖旧的数据!"})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "SN填写正确!"})
	}

}
