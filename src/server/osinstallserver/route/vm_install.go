package route

import (
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"math/rand"
	"middleware"
	"model"
	"regexp"
	"server/osinstallserver/util"
	"strconv"
	"strings"
	"time"
	"utils"
)

func CreateVmDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	type DeviceParam struct {
		ID uint
		Sn string
	}

	var info struct {
		Devices        []DeviceParam
		VmNumber       int
		OsID           uint
		CpuCoresNumber uint
		MemoryCurrent  uint
		DiskSize       uint
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if len(info.Devices) <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请选择要操作的物理机!"})
		return
	}

	for _, v := range info.Devices {
		if v.ID <= uint(0) && v.Sn != "" {
			count, err := repo.CountDeviceBySn(v.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
				return
			}
			if count > 0 {
				v.ID, err = repo.GetDeviceIdBySn(v.Sn)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
					return
				}
			}
		}

		if v.ID <= uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该物理机不存在(SN:" + v.Sn + ")!"})
			return
		}

		device, err := repo.GetDeviceById(v.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
			return
		}

		if device.Status != "success" {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "SN:" + device.Sn + "不能安装虚拟机!"})
			return
		}
	}

	if info.VmNumber <= 0 || info.OsID <= uint(0) {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "虚拟机个数和操作系统版本参数不能为空!", "Content": ""})
		return
	}

	if info.CpuCoresNumber <= uint(0) {
		info.CpuCoresNumber = uint(1)
	}

	if info.MemoryCurrent <= uint(0) {
		info.MemoryCurrent = uint(1024)
	}

	if info.DiskSize <= uint(0) {
		info.MemoryCurrent = uint(60)
	}

	//循环安装虚拟机
	var currentDeviceIndex int
	currentDeviceIndex = 0
	fix := time.Now().Format("20060102150405") + fmt.Sprintf("%d", rand.Intn(100))
	var result []model.VmDevice
	for i := 0; i < info.VmNumber; i++ {
		var deviceId uint
		deviceId = info.Devices[currentDeviceIndex].ID
		if deviceId <= uint(0) && info.Devices[currentDeviceIndex].Sn != "" {
			count, err := repo.CountDeviceBySn(info.Devices[currentDeviceIndex].Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
				return
			}
			if count > 0 {
				deviceId, err = repo.GetDeviceIdBySn(info.Devices[currentDeviceIndex].Sn)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
					return
				}
			}
		}

		if deviceId <= uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该物理机不存在(SN:" + info.Devices[currentDeviceIndex].Sn + ")!"})
			return
		}

		//deviceId := info.Devices[currentDeviceIndex].ID

		currentDeviceIndex++
		if currentDeviceIndex >= len(info.Devices) {
			currentDeviceIndex = 0
		}

		device, err := repo.GetDeviceById(deviceId)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
			return
		}
		ip, err := repo.AssignNewIpByNetworkId(device.NetworkID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
			return
		}

		var row model.VmDevice
		row.DeviceID = deviceId
		row.Ip = ip
		row.Hostname = fix + fmt.Sprintf("%d", i)
		row.Mac = util.CreateNewMacAddress()
		row.NetworkID = device.NetworkID
		row.OsID = info.OsID
		row.CpuCoresNumber = info.CpuCoresNumber
		row.CpuHotPlug = "No"
		row.CpuPassthrough = "No"
		row.CpuTopSockets = 0
		row.CpuTopCores = 0
		row.CpuTopThreads = 0
		row.CpuPinning = ""
		row.MemoryCurrent = info.MemoryCurrent
		row.MemoryMax = info.MemoryCurrent
		row.MemoryKsm = "No"
		row.DiskType = "raw"
		row.DiskSize = info.DiskSize
		row.DiskBusType = "virtio"
		row.DiskCacheMode = "writeback"
		row.DiskIoMode = "default"
		row.NetworkType = "bridge"
		row.NetworkDeviceType = "virtio"
		row.DisplayType = "serialPorts"
		row.DisplayPassword = ""
		row.DisplayUpdatePassword = "No"
		row.Status = "pre_install"

		resultAdd, errAdd := repo.AddVmDevice(row.DeviceID,
			row.Hostname,
			row.Mac,
			row.Ip,
			row.NetworkID,
			row.OsID,
			row.CpuCoresNumber,
			row.CpuHotPlug,
			row.CpuPassthrough,
			row.CpuTopSockets,
			row.CpuTopCores,
			row.CpuTopThreads,
			row.CpuPinning,
			row.MemoryCurrent,
			row.MemoryMax,
			row.MemoryKsm,
			row.DiskType,
			row.DiskSize,
			row.DiskBusType,
			row.DiskCacheMode,
			row.DiskIoMode,
			row.NetworkType,
			row.NetworkDeviceType,
			row.DisplayType,
			row.DisplayPassword,
			row.DisplayUpdatePassword,
			row.Status)
		if errAdd != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}

		result = append(result, resultAdd)
	}
	/*
		count, err := repo.CountVmDeviceByMac(info.Mac)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
			return
		}

		if count > 0 {
			w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该Mac已存在，继续填写会覆盖旧的数据!"})
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "Mac填写正确!"})
		}
	*/
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功!", "Content": result})

}

func CreateNewMacAddress(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	fmt.Println(repo)

	mac := util.CreateNewMacAddress()
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mac})
}

func DeleteVmDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	osConfig, err := repo.DeleteVmDeviceById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": osConfig})
}

func GetVmDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	osConfig, err := repo.GetVmDeviceById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": osConfig})
}

func GetFullVmDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mod, err := repo.GetFullVmDeviceById(info.ID)
	type VmDeviceFullWithTime struct {
		ID                    uint
		DeviceID              uint
		DeviceSn              string
		Hostname              string
		Mac                   string
		Ip                    string
		NetworkID             uint
		NetworkName           string
		OsID                  uint
		OsName                string
		CpuCoresNumber        uint
		CpuHotPlug            string
		CpuPassthrough        string
		CpuTopSockets         uint
		CpuTopCores           uint
		CpuTopThreads         uint
		CpuPinning            string
		MemoryCurrent         uint
		MemoryMax             uint
		MemoryKsm             string
		DiskType              string
		DiskSize              uint
		DiskBusType           string
		DiskCacheMode         string
		DiskIoMode            string
		NetworkType           string
		NetworkDeviceType     string
		DisplayType           string
		DisplayPassword       string
		DisplayUpdatePassword string
		Status                string
		CreatedAt             utils.ISOTime
		UpdatedAt             utils.ISOTime
	}
	var vm VmDeviceFullWithTime
	vm.ID = mod.ID
	vm.DeviceID = mod.DeviceID
	vm.DeviceSn = mod.DeviceSn
	vm.Hostname = mod.Hostname
	vm.Mac = mod.Mac
	vm.Ip = mod.Ip
	vm.NetworkID = mod.NetworkID
	vm.NetworkName = mod.NetworkName
	vm.OsID = mod.OsID
	vm.OsName = mod.OsName
	vm.CpuCoresNumber = mod.CpuCoresNumber
	vm.CpuHotPlug = mod.CpuHotPlug
	vm.CpuPassthrough = mod.CpuPassthrough
	vm.CpuTopSockets = mod.CpuTopSockets
	vm.CpuTopCores = mod.CpuTopCores
	vm.CpuTopThreads = mod.CpuTopThreads
	vm.CpuPinning = mod.CpuPinning
	vm.MemoryCurrent = mod.MemoryCurrent
	vm.MemoryMax = mod.MemoryMax
	vm.MemoryKsm = mod.MemoryKsm
	vm.DiskType = mod.DiskType
	vm.DiskSize = mod.DiskSize
	vm.DiskBusType = mod.DiskBusType
	vm.DiskCacheMode = mod.DiskCacheMode
	vm.DiskIoMode = mod.DiskIoMode
	vm.NetworkType = mod.NetworkType
	vm.NetworkDeviceType = mod.NetworkDeviceType
	vm.DisplayType = mod.DisplayType
	vm.DisplayPassword = mod.DisplayPassword
	vm.DisplayUpdatePassword = mod.DisplayUpdatePassword
	vm.Status = mod.Status

	vm.CreatedAt = utils.ISOTime(mod.CreatedAt)
	vm.UpdatedAt = utils.ISOTime(mod.UpdatedAt)

	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": vm})
}

func GetVmDeviceList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
		Status         string
		StartUpdatedAt string
		EndUpdatedAt   string
		DeviceID       int
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Keyword = strings.TrimSpace(info.Keyword)
	info.Status = strings.TrimSpace(info.Status)

	var where string
	where = " where t1.id > 0 "
	if info.OsID > 0 {
		where += " and t1.os_id = " + strconv.Itoa(info.OsID)
	}

	if info.Status != "" {
		where += " and t1.status = '" + info.Status + "'"
	}

	if info.StartUpdatedAt != "" {
		where += " and t1.updated_at >= '" + info.StartUpdatedAt + "'"
	}

	if info.EndUpdatedAt != "" {
		where += " and t1.updated_at <= '" + info.EndUpdatedAt + "'"
	}

	if info.DeviceID > 0 {
		where += " and t1.device_id = " + strconv.Itoa(info.DeviceID)
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
			where += str + " t1.mac = '" + v + "' or t1.hostname = '" + v + "' or t1.ip = '" + v + "'"
		}
		where += " ) "
	}

	osConfigs, err := repo.GetVmDeviceListWithPage(info.Limit, info.Offset, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = osConfigs

	//总条数
	count, err := repo.CountVmDevice(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func GetVmDeviceListByHostSn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var sn string
	sn = r.FormValue("sn")
	sn = strings.TrimSpace(sn)
	if sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "宿主机SN不能为空!"})
		return
	}

	countDevice, err := repo.CountDeviceBySn(sn)
	if countDevice <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "宿主机不存在!"})
		return
	}
	deviceId, err := repo.GetDeviceIdBySn(sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	var where string
	where = fmt.Sprintf("where device_id = %d", deviceId)

	mods, err := repo.GetVmDeviceListWithPage(1000000, 0, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountVmDevice(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func BatchAddVmDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	fmt.Println(repo)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var infos []model.VmDevice

	if err := r.DecodeJSONPayload(&infos); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	//先批量检测传入数据是否有问题
	for _, info := range infos {
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.Mac = strings.TrimSpace(info.Mac)
		info.CpuHotPlug = strings.TrimSpace(info.CpuHotPlug)
		info.CpuPassthrough = strings.TrimSpace(info.CpuPassthrough)
		info.CpuPinning = strings.TrimSpace(info.CpuPinning)
		info.MemoryKsm = strings.TrimSpace(info.MemoryKsm)
		info.DiskType = strings.TrimSpace(info.DiskType)
		info.DiskBusType = strings.TrimSpace(info.DiskBusType)
		info.DiskCacheMode = strings.TrimSpace(info.DiskCacheMode)
		info.DiskIoMode = strings.TrimSpace(info.DiskIoMode)
		info.NetworkType = strings.TrimSpace(info.NetworkType)
		info.NetworkDeviceType = strings.TrimSpace(info.NetworkDeviceType)
		info.DisplayType = strings.TrimSpace(info.DisplayType)
		info.DisplayPassword = strings.TrimSpace(info.DisplayPassword)
		info.DisplayUpdatePassword = strings.TrimSpace(info.DisplayUpdatePassword)

		if info.Hostname == "" || info.Ip == "" || info.Mac == "" || info.DiskType == "" ||
			info.NetworkType == "" || info.NetworkDeviceType == "" || info.DisplayType == "" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
			return
		}

		if info.DeviceID <= uint(0) || info.OsID <= uint(0) || info.CpuCoresNumber <= uint(0) || info.MemoryCurrent <= uint(0) ||
			info.MemoryMax <= uint(0) || info.DiskSize <= uint(0) {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
			return
		}

		count, err := repo.CountVmDeviceByMac(info.Mac)
		if count > 0 {
			vmDeviceId, err := repo.GetVmDeviceIdByMac(info.Mac)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			count, err := repo.CountVmDeviceByHostnameAndId(info.Hostname, vmDeviceId)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if count > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
				return
			}

			countIp, err := repo.CountVmDeviceByIpAndId(info.Ip, vmDeviceId)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
				return
			}
		} else {
			count, err := repo.CountVmDeviceByHostname(info.Hostname)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if count > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Hostname + " 该主机名已存在!"})
				return
			}

			countIp, err := repo.CountVmDeviceByIp(info.Ip)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + " 该IP已存在!"})
				return
			}
		}

		//物理机是否使用
		countHostname, err := repo.CountDeviceByHostname(info.Hostname)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countHostname > 0 {
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
		info.Hostname = strings.TrimSpace(info.Hostname)
		info.Ip = strings.TrimSpace(info.Ip)
		info.Mac = strings.TrimSpace(info.Mac)
		info.CpuHotPlug = strings.TrimSpace(info.CpuHotPlug)
		info.CpuPassthrough = strings.TrimSpace(info.CpuPassthrough)
		info.CpuPinning = strings.TrimSpace(info.CpuPinning)
		info.MemoryKsm = strings.TrimSpace(info.MemoryKsm)
		info.DiskType = strings.TrimSpace(info.DiskType)
		info.DiskBusType = strings.TrimSpace(info.DiskBusType)
		info.DiskCacheMode = strings.TrimSpace(info.DiskCacheMode)
		info.DiskIoMode = strings.TrimSpace(info.DiskIoMode)
		info.NetworkType = strings.TrimSpace(info.NetworkType)
		info.NetworkDeviceType = strings.TrimSpace(info.NetworkDeviceType)
		info.DisplayType = strings.TrimSpace(info.DisplayType)
		info.DisplayPassword = strings.TrimSpace(info.DisplayPassword)
		info.DisplayUpdatePassword = strings.TrimSpace(info.DisplayUpdatePassword)
		info.Status = "pre_install"
		//Mac已存在的情况下，要覆盖原数据
		count, err := repo.CountVmDeviceByMac(info.Mac)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		//覆盖
		if count > 0 {
			id, err := repo.GetVmDeviceIdByMac(info.Mac)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			_, errUpdate := repo.UpdateVmDeviceById(id,
				info.DeviceID,
				info.Hostname,
				info.Mac,
				info.Ip,
				info.NetworkID,
				info.OsID,
				info.CpuCoresNumber,
				info.CpuHotPlug,
				info.CpuPassthrough,
				info.CpuTopSockets,
				info.CpuTopCores,
				info.CpuTopThreads,
				info.CpuPinning,
				info.MemoryCurrent,
				info.MemoryMax,
				info.MemoryKsm,
				info.DiskType,
				info.DiskSize,
				info.DiskBusType,
				info.DiskCacheMode,
				info.DiskIoMode,
				info.NetworkType,
				info.NetworkDeviceType,
				info.DisplayType,
				info.DisplayPassword,
				info.DisplayUpdatePassword,
				info.Status)

			if errUpdate != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
				return
			}

		} else {
			_, err := repo.AddVmDevice(info.DeviceID,
				info.Hostname,
				info.Mac,
				info.Ip,
				info.NetworkID,
				info.OsID,
				info.CpuCoresNumber,
				info.CpuHotPlug,
				info.CpuPassthrough,
				info.CpuTopSockets,
				info.CpuTopCores,
				info.CpuTopThreads,
				info.CpuPinning,
				info.MemoryCurrent,
				info.MemoryMax,
				info.MemoryKsm,
				info.DiskType,
				info.DiskSize,
				info.DiskBusType,
				info.DiskCacheMode,
				info.DiskIoMode,
				info.NetworkType,
				info.NetworkDeviceType,
				info.DisplayType,
				info.DisplayPassword,
				info.DisplayUpdatePassword,
				info.Status)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
				return
			}
		}

	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

//重装
func BatchReInstallVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
		_, err := repo.ReInstallVmDeviceById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchDeleteVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

		_, errDevice := repo.DeleteVmDeviceById(info.ID)
		if errDevice != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDevice.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func ValidateMac(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Mac string
	}
	info.Mac = strings.TrimSpace(info.Mac)

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if info.Mac == "" {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "Mac参数不能为空!", "Content": ""})
		return
	}

	count, err := repo.CountVmDeviceByMac(info.Mac)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "参数错误!"})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "该Mac已存在，继续填写会覆盖旧的数据!"})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "Mac填写正确!"})
	}

}
