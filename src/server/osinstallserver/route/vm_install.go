package route

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"math"
	"middleware"
	"model"
	"os"
	"regexp"
	"server/osinstallserver/util"
	"strconv"
	"strings"
	"utils"
)

func AddVmDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		Hostname       string
		Mac            string
		Ip             string
		Sn             string
		NetworkID      uint
		OsID           uint
		SystemID       uint
		CpuCoresNumber uint
		MemoryCurrent  uint
		DiskSize       uint
		UserID         uint
		AccessToken    string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Hostname = strings.TrimSpace(info.Hostname)
	info.Mac = strings.TrimSpace(info.Mac)
	info.Mac = strings.ToLower(info.Mac)
	info.Ip = strings.TrimSpace(info.Ip)
	info.Sn = strings.TrimSpace(info.Sn)
	info.AccessToken = strings.TrimSpace(info.AccessToken)

	session, err := GetSession(w, r)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
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

	if info.Sn == "" || info.Hostname == "" || info.Ip == "" || info.OsID == uint(0) || info.SystemID == uint(0) {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	if info.NetworkID == uint(0) {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "未匹配到网段信息!"})
		return
	}

	if info.CpuCoresNumber <= uint(0) || info.MemoryCurrent <= uint(0) || info.DiskSize <= uint(0) {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "CPU、内存、磁盘参数格式不正确!"})
		return
	}

	//match network
	isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", info.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	if !isValidate {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": info.Ip + "IP格式不正确!"})
		return
	}

	//check host device
	count, err := repo.CountDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if count <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该宿主机不存在(SN:" + info.Sn + ")!"})
		return
	}
	countVmHost, err := repo.CountVmHostBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if countVmHost <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该宿主机没有可用资源(SN:" + info.Sn + ")!"})
		return
	}
	//get host device info
	device, err := repo.GetDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if device.Status != "success" || device.IsSupportVm != "Yes" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN:" + device.Sn + "不能安装虚拟机!"})
		return
	}

	vmHost, err := repo.GetVmHostBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if vmHost.IsAvailable != "Yes" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该宿主机无法安装虚拟机(SN:" + info.Sn + ")!"})
		return
	}

	//check mac
	countMac, err := repo.CountVmDeviceByMac(info.Mac)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if countMac > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该Mac地址已存在!"})
		return
	}

	//check hostname
	countHostname, err := repo.CountVmDeviceByHostname(info.Hostname)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if countHostname > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该主机名已存在!"})
		return
	}
	countDeviceHostname, err := repo.CountDeviceByHostname(info.Hostname)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if countDeviceHostname > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该主机名已存在(物理机已使用)!"})
		return
	}

	//check mac
	isValidateMac, err := regexp.MatchString("^([0-9a-fA-F]{2})(([/\\s:][0-9a-fA-F]{2}){5})$", info.Mac)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	if !isValidateMac {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "Mac地址格式不正确!"})
		return
	}
	if strings.Index(info.Mac, "52:54:00") != 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "Mac地址必须以\"52:54:00\"开头!"})
		return
	}

	//check ip
	countIp, err := repo.CountVmDeviceByIp(info.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if countIp > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该IP已存在!"})
		return
	}
	countDeviceIp, err := repo.CountDeviceByIp(info.Ip)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}
	if countDeviceIp > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该IP已存在(物理机已使用)!"})
		return
	}
	if info.NetworkID != device.NetworkID {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "虚拟机和宿主机不在同一网段!"})
		return
	}

	//check availability
	if info.MemoryCurrent >= vmHost.MemoryAvailable {
		memoryRound := float64(vmHost.MemoryAvailable / 1024)
		memory := int(math.Floor(memoryRound))
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("宿主机剩余内存不够分配(最大可用 %dG)!", memory)})
		return
	}
	if info.DiskSize >= vmHost.DiskAvailable {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("宿主机剩余磁盘空间不够分配(最大可用 %dG)!", vmHost.DiskAvailable)})
		return
	}
	where := fmt.Sprintf("device_id = %d", device.ID)

	//vnc port
	maxVncPort, err := repo.GetMaxVncPort(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
		return
	}
	var vncPort uint
	if maxVncPort <= uint(0) {
		vncPort = uint(5900)
	} else {
		vncPort = maxVncPort + 1
	}

	var row model.VmDevice
	row.DeviceID = device.ID
	row.Ip = info.Ip
	row.Hostname = info.Hostname
	row.Mac = info.Mac
	row.NetworkID = info.NetworkID
	row.OsID = info.OsID
	row.SystemID = info.SystemID
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
	row.Status = "pre_create"
	row.VncPort = fmt.Sprintf("%d", vncPort)
	row.RunStatus = ""
	row.UserID = info.UserID

	resultAdd, errAdd := repo.AddVmDevice(row.DeviceID,
		row.Hostname,
		row.Mac,
		row.Ip,
		row.NetworkID,
		row.OsID,
		row.SystemID,
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
		row.Status,
		row.UserID,
		row.VncPort,
		row.RunStatus)
	if errAdd != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errAdd.Error()})
		return
	}

	var infoHost *model.VmHost
	infoHost = vmHost
	//cpu update
	//cpu used sum
	infoHost.CpuUsed, err = repo.GetCpuUsedSum(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
		return
	}
	cpuAvailable := int(infoHost.CpuSum - infoHost.CpuUsed)
	if cpuAvailable <= 0 {
		cpuAvailable = 0
	}
	infoHost.CpuAvailable = uint(cpuAvailable)
	//memory update
	infoHost.MemoryUsed, err = repo.GetMemoryUsedSum(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
		return
	}
	memoryAvailable := int(infoHost.MemorySum - infoHost.MemoryUsed)
	if memoryAvailable <= 0 {
		memoryAvailable = 0
	}
	infoHost.MemoryAvailable = uint(memoryAvailable)
	//update disk
	infoHost.DiskUsed, err = repo.GetDiskUsedSum(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
		return
	}
	diskAvailable := int(infoHost.DiskSum - infoHost.DiskUsed)
	if diskAvailable < 0 {
		diskAvailable = 0
	}
	infoHost.DiskAvailable = uint(diskAvailable)

	infoHost.VmNum, err = repo.CountVmDeviceByDeviceId(device.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
		return
	}
	//update host
	_, errUpdate := repo.UpdateVmHostCpuMemoryDiskVmNumById(vmHost.ID, infoHost.CpuSum, infoHost.CpuUsed, infoHost.CpuAvailable, infoHost.MemorySum, infoHost.MemoryUsed, infoHost.MemoryAvailable, infoHost.DiskSum, infoHost.DiskUsed, infoHost.DiskAvailable, infoHost.VmNum, infoHost.IsAvailable)
	if errUpdate != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
		return
	}
	vmDeviceId := resultAdd.ID
	//update status
	_, errUpdateStatus := repo.UpdateVmInstallInfoById(vmDeviceId, "creating", 0)
	if errUpdateStatus != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateStatus.Error()})
		return
	}
	//create vm vol
	errCreateVol := RunCreateVol(ctx, vmDeviceId)
	var logTitle string
	var installLog string
	if errCreateVol != nil {
		logTitle = "逻辑卷创建失败"
		installLog = errCreateVol.Error()
	} else {
		logTitle = "逻辑卷创建成功"
		installLog = "逻辑卷创建成功"
	}
	_, errAddLog := repo.AddVmDeviceLog(vmDeviceId, logTitle, "install", installLog)
	if errAddLog != nil {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "数据添加成功，虚拟机创建失败:" + errAddLog.Error()})
		return
	}
	if errCreateVol != nil {
		//update status
		_, errUpdateStatus := repo.UpdateVmInstallInfoById(vmDeviceId, "create_failure", 0)
		if errUpdateStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateStatus.Error()})
			return
		}
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "数据添加成功，虚拟机创建失败:" + errCreateVol.Error()})
		return
	}

	//create vm
	errCreateVm := RunCreateVm(ctx, vmDeviceId)
	if errCreateVm != nil {
		logTitle = "虚拟机创建失败"
		installLog = errCreateVm.Error()
	} else {
		logTitle = "虚拟机创建成功"
		installLog = "虚拟机创建成功"
	}
	_, errAddLogVm := repo.AddVmDeviceLog(vmDeviceId, logTitle, "install", installLog)
	if errAddLogVm != nil {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "数据添加成功，虚拟机创建失败:" + errAddLogVm.Error()})
		return
	}
	if errCreateVm != nil {
		//update status
		_, errUpdateStatus := repo.UpdateVmInstallInfoById(vmDeviceId, "create_failure", 0)
		if errUpdateStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateStatus.Error()})
			return
		}
		//destroy vol
		errDestroyVol := RunDestroyVol(ctx, vmDeviceId)
		var logTitle string
		var installLog string
		if errDestroyVol != nil {
			logTitle = "逻辑卷销毁失败"
			installLog = errDestroyVol.Error()
		} else {
			logTitle = "逻辑卷销毁成功"
			installLog = "逻辑卷销毁成功"
		}
		_, errAddLogDestory := repo.AddVmDeviceLog(vmDeviceId, logTitle, "install", installLog)
		if errAddLogDestory != nil {
			w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "数据添加成功，虚拟机创建失败:" + errAddLogDestory.Error()})
			return
		}
		if errDestroyVol != nil {
			w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "数据添加成功，虚拟机创建失败:" + errDestroyVol.Error()})
			return
		}
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "数据添加成功，虚拟机创建失败:" + errCreateVm.Error()})
		return
	}

	//update status
	_, errUpdateStatus2 := repo.UpdateVmInstallInfoById(vmDeviceId, "pre_install", 0)
	if errUpdateStatus2 != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateStatus2.Error()})
		return
	}

	//update run status
	_, errUpdateRunStatus := repo.UpdateVmRunStatusById(vmDeviceId, "running")
	if errUpdateRunStatus != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateRunStatus.Error()})
		return
	}

	//create pxe file
	errPxe := CreatePxeFile(ctx, info.Mac)
	if errPxe != nil {
		logger.Error("Pxe文件生成失败:" + errPxe.Error())
	}

	//create novnc file
	errNovnc := RunCreateVmNoVncTokenFile(repo, logger, vmDeviceId)
	if errPxe != nil {
		logger.Error("noVnc文件生成失败:" + errNovnc.Error())
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功!", "Content": resultAdd})
}

func CreatePxeFile(ctx context.Context, mac string) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	if conf.OsInstall.PxeConfigDir == "" {
		return errors.New("Pxe配置文件目录没有指定")
	}

	var info struct {
		Mac string
	}

	info.Mac = strings.TrimSpace(mac)
	if info.Mac == "" {
		return errors.New("Mac地址参数不能为空")
	}

	//mac 大写转为 小写
	info.Mac = strings.ToLower(info.Mac)

	device, err := repo.GetVmDeviceByMac(info.Mac)
	if err != nil {
		return errors.New("该设备不存在")
	}

	osConfig, err := repo.GetOsConfigById(device.OsID)
	if err != nil {
		return errors.New("PXE信息没有配置" + err.Error())
	}

	//替换占位符
	osConfig.Pxe = strings.Replace(osConfig.Pxe, "{sn}", info.Mac, -1)
	osConfig.Pxe = strings.Replace(osConfig.Pxe, "\r\n", "\n", -1)

	pxeFileName := util.GetPxeFileNameByMac(info.Mac)
	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}
	logger.Debugf("Create pxe file: %s", conf.OsInstall.PxeConfigDir+"/"+pxeFileName)

	errCreatePxeFile := util.CreatePxeFile(conf.OsInstall.PxeConfigDir, pxeFileName, osConfig.Pxe)
	if errCreatePxeFile != nil {
		logger.Debugf("配置文件生成失败" + err.Error())
		return errors.New("配置文件生成失败" + err.Error())
	}

	return nil
}

func CreateNewMacAddress(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	_, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

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
		SystemName            string
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
		RunStatus             string
		VncPort               string
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
	vm.SystemName = mod.SystemName
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
	vm.RunStatus = mod.RunStatus
	vm.Status = mod.Status
	vm.VncPort = mod.VncPort

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
		RunStatus      string
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
	info.RunStatus = strings.TrimSpace(info.RunStatus)

	var where string
	where = " where t1.id > 0 "
	if info.OsID > 0 {
		where += " and t1.os_id = " + strconv.Itoa(info.OsID)
	}

	if info.Status != "" {
		where += " and t1.status = '" + info.Status + "'"
	}

	if info.RunStatus != "" {
		where += " and t1.run_status = '" + info.RunStatus + "'"
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

//重装
func BatchReInstallVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	logger, ok := middleware.LoggerFromContext(ctx)
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

		vmDevice, errInfo := repo.GetVmDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}
		if session.Role != "Administrator" && vmDevice.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "您无权操作其他人的设备!(Mac:" + vmDevice.Mac + ")"})
			return
		}
		//destroy vm
		errRun := RunDestroyVm(ctx, vmDevice.ID)
		var logTitle string
		var installLog string
		if errRun != nil {
			logTitle = "虚拟机销毁失败"
			installLog = errRun.Error()
		} else {
			logTitle = "虚拟机销毁成功"
			installLog = "虚拟机销毁成功"
		}
		_, errAddLog := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "install", installLog)
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
		if errRun != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errRun.Error()})
			return
		}
		//update status
		_, errUpdateStatusCreate := repo.UpdateVmInstallInfoById(vmDevice.ID, "pre_create", 0)
		if errUpdateStatusCreate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateStatusCreate.Error()})
			return
		}
		//update run status
		_, errUpdateRunStatus := repo.UpdateVmRunStatusById(vmDevice.ID, "")
		if errUpdateRunStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateRunStatus.Error()})
			return
		}
		//destroy vol
		errDestroyVol := RunDestroyVol(ctx, vmDevice.ID)
		if errDestroyVol != nil {
			logTitle = "逻辑卷销毁失败"
			installLog = errDestroyVol.Error()
		} else {
			logTitle = "逻辑卷销毁成功"
			installLog = "逻辑卷销毁成功"
		}
		_, errAddLogDestory := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "install", installLog)
		if errAddLogDestory != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLogDestory.Error()})
			return
		}
		if errDestroyVol != nil {
			//update status
			_, errUpdateStatusCreate := repo.UpdateVmInstallInfoById(vmDevice.ID, "create_failure", 0)
			if errUpdateStatusCreate != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateStatusCreate.Error()})
				return
			}
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDestroyVol.Error()})
			return
		}
		//update status
		_, errUpdateStatus := repo.UpdateVmInstallInfoById(vmDevice.ID, "creating", 0)
		if errUpdateStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateStatus.Error()})
			return
		}
		//create vol
		errCreateVol := RunCreateVol(ctx, vmDevice.ID)
		if errCreateVol != nil {
			logTitle = "逻辑卷创建失败"
			installLog = errCreateVol.Error()
		} else {
			logTitle = "逻辑卷创建成功"
			installLog = "逻辑卷创建成功"
		}
		_, errAddLogCreateVol := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "install", installLog)
		if errAddLogCreateVol != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLogCreateVol.Error()})
			return
		}
		if errCreateVol != nil {
			//update status
			_, errUpdateStatus := repo.UpdateVmInstallInfoById(vmDevice.ID, "create_failure", 0)
			if errUpdateStatus != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateStatus.Error()})
				return
			}
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCreateVol.Error()})
			return
		}
		//create vm
		errRunCreate := RunCreateVm(ctx, vmDevice.ID)
		if errRunCreate != nil {
			logTitle = "虚拟机创建失败"
			installLog = errRunCreate.Error()
		} else {
			logTitle = "虚拟机创建成功"
			installLog = "虚拟机创建成功"
		}
		_, errAddLogCreate := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "install", installLog)
		if errAddLogCreate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLogCreate.Error()})
			return
		}
		if errRunCreate != nil {
			//update status
			_, errUpdateStatusCreate := repo.UpdateVmInstallInfoById(vmDevice.ID, "create_failure", 0)
			if errUpdateStatusCreate != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateStatusCreate.Error()})
				return
			}

			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errRunCreate.Error()})
			return
		}
		//update status
		_, errUpdateStatus2 := repo.UpdateVmInstallInfoById(vmDevice.ID, "pre_install", 0)
		if errUpdateStatus2 != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdateStatus2.Error()})
			return
		}

		//update run status
		_, errUpdateRunStatus2 := repo.UpdateVmRunStatusById(vmDevice.ID, "running")
		if errUpdateRunStatus2 != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateRunStatus2.Error()})
			return
		}

		//create pxe file
		errPxe := CreatePxeFile(ctx, vmDevice.Mac)
		if errPxe != nil {
			logger.Error("Pxe文件生成失败:" + errPxe.Error())
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

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	conf, ok := middleware.ConfigFromContext(ctx)
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

		vmDevice, errInfo := repo.GetVmDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}
		if session.Role != "Administrator" && vmDevice.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "您无权操作其他人的设备!(Mac:" + vmDevice.Mac + ")"})
			return
		}

		//get host device info
		device, err := repo.GetDeviceById(vmDevice.DeviceID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "宿主机不存在!(ID:" + fmt.Sprintf("%d", vmDevice.DeviceID) + ")"})
			return
		}

		//delete vol
		if vmDevice.Status != "pre_create" && vmDevice.Status != "create_failure" {
			errDestoryVm := RunDestroyVm(ctx, vmDevice.ID)
			if errDestoryVm != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "虚拟机销毁失败:" + errDestoryVm.Error()})
				return
			}

			errDestoryVol := RunDestroyVol(ctx, vmDevice.ID)
			if errDestoryVol != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "逻辑卷销毁失败:" + errDestoryVol.Error()})
				return
			}
		}

		//remove pxe config file
		pxeFileName := util.GetPxeFileNameByMac(vmDevice.Mac)
		confDir := conf.OsInstall.PxeConfigDir
		if util.FileExist(confDir + "/" + pxeFileName) {
			err := os.Remove(confDir + "/" + pxeFileName)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
		}

		//delete vm novnc token
		errDeleteVncToken := RunDeleteVmNoVncTokenFile(repo, logger, vmDevice.ID)
		if errDeleteVncToken != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDeleteVncToken.Error()})
			return
		}

		//delete vm device
		_, errDelete := repo.DeleteVmDeviceById(info.ID)
		if errDelete != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDelete.Error()})
			return
		}

		//delete vm device log
		_, errDeleteLog := repo.DeleteVmDeviceLogByDeviceID(vmDevice.ID)
		if errDeleteLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDeleteLog.Error()})
			return
		}

		//update host resource info
		//get host info
		vmHost, err := repo.GetVmHostBySn(device.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
			return
		}
		//condition
		where := fmt.Sprintf("device_id = %d", device.ID)
		var infoHost *model.VmHost
		infoHost = vmHost
		//cpu update
		//cpu used sum
		infoHost.CpuUsed, err = repo.GetCpuUsedSum(where)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}
		cpuAvailable := int(infoHost.CpuSum - infoHost.CpuUsed)
		if cpuAvailable <= 0 {
			cpuAvailable = 0
		}
		infoHost.CpuAvailable = uint(cpuAvailable)
		//memory update
		infoHost.MemoryUsed, err = repo.GetMemoryUsedSum(where)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}
		memoryAvailable := int(infoHost.MemorySum - infoHost.MemoryUsed)
		if memoryAvailable <= 0 {
			memoryAvailable = 0
		}
		infoHost.MemoryAvailable = uint(memoryAvailable)
		//update disk
		infoHost.DiskUsed, err = repo.GetDiskUsedSum(where)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}
		diskAvailable := int(infoHost.DiskSum - infoHost.DiskUsed)
		if diskAvailable < 0 {
			diskAvailable = 0
		}
		infoHost.DiskAvailable = uint(diskAvailable)

		infoHost.VmNum, err = repo.CountVmDeviceByDeviceId(device.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + err.Error()})
			return
		}
		//update host resource info
		_, errUpdate := repo.UpdateVmHostCpuMemoryDiskVmNumById(vmHost.ID, infoHost.CpuSum, infoHost.CpuUsed, infoHost.CpuAvailable, infoHost.MemorySum, infoHost.MemoryUsed, infoHost.MemoryAvailable, infoHost.DiskSum, infoHost.DiskUsed, infoHost.DiskAvailable, infoHost.VmNum, infoHost.IsAvailable)
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdate.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchStartVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

		vmDevice, errInfo := repo.GetVmDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}
		if session.Role != "Administrator" && vmDevice.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "您无权操作其他人的设备!(Mac:" + vmDevice.Mac + ")"})
			return
		}

		/*
			if vmDevice.Status != "success" {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备未完成安装，无法启动!(主机名:" + vmDevice.Hostname + ")"})
				return
			}
		*/

		errRun := RunStartVm(ctx, vmDevice.ID)
		//log
		var logTitle string
		var installLog string
		if errRun != nil {
			logTitle = "虚拟机启动失败"
			installLog = errRun.Error()
		} else {
			logTitle = "虚拟机启动成功"
			installLog = "虚拟机启动成功"
		}
		_, errAddLog := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "operate", installLog)
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
		if errRun != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errRun.Error()})
			return
		}
		//update run status
		_, errUpdateRunStatus := repo.UpdateVmRunStatusById(vmDevice.ID, "running")
		if errUpdateRunStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateRunStatus.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchStopVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

		vmDevice, errInfo := repo.GetVmDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}
		if session.Role != "Administrator" && vmDevice.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "您无权操作其他人的设备!(Mac:" + vmDevice.Mac + ")"})
			return
		}

		/*
			if vmDevice.Status != "success" {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备未完成安装，无法停止!(主机名:" + vmDevice.Hostname + ")"})
				return
			}
		*/

		errRun := RunStopVm(ctx, vmDevice.ID)
		//log
		var logTitle string
		var installLog string
		if errRun != nil {
			logTitle = "虚拟机停止失败"
			installLog = errRun.Error()
		} else {
			logTitle = "虚拟机停止成功"
			installLog = "虚拟机停止成功"
		}
		_, errAddLog := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "operate", installLog)
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
		if errRun != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errRun.Error()})
			return
		}

		//update run status
		_, errUpdateRunStatus := repo.UpdateVmRunStatusById(vmDevice.ID, "stop")
		if errUpdateRunStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateRunStatus.Error()})
			return
		}
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func BatchReStartVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

		vmDevice, errInfo := repo.GetVmDeviceById(info.ID)
		if errInfo != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errInfo.Error()})
			return
		}
		if session.Role != "Administrator" && vmDevice.UserID != info.UserID {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "您无权操作其他人的设备!(Mac:" + vmDevice.Mac + ")"})
			return
		}

		/*
			if vmDevice.Status != "success" {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备未完成安装，无法重启!(主机名:" + vmDevice.Hostname + ")"})
				return
			}
		*/

		errRun := RunReStartVm(ctx, vmDevice.ID)
		//log
		var logTitle string
		var installLog string
		if errRun != nil {
			logTitle = "虚拟机重启失败"
			installLog = errRun.Error()
		} else {
			logTitle = "虚拟机重启成功"
			installLog = "虚拟机重启成功"
		}
		_, errAddLog := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "operate", installLog)
		if errAddLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
			return
		}
		if errRun != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errRun.Error()})
			return
		}

		//update run status
		_, errUpdateRunStatus := repo.UpdateVmRunStatusById(vmDevice.ID, "running")
		if errUpdateRunStatus != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "操作失败:" + errUpdateRunStatus.Error()})
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
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "Mac参数不能为空!", "Content": ""})
		return
	}

	count, err := repo.CountVmDeviceByMac(info.Mac)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误!"})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该Mac已存在，继续填写会覆盖旧的数据!"})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "Mac填写正确!"})
	}

}

func GetSystemBySnForVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mod, err := repo.GetSystemByVmMac(info.Sn)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}

		return
	}
	mod.Content = strings.Replace(mod.Content, "\r\n", "\n", -1)

	if info.Type == "raw" {
		w.Header().Add("Content-type", "text/html; charset=utf-8")
		w.Write([]byte(mod.Content))
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "成功获取system信息", "Content": mod})
	}
}

func GetNetworkBySnForVm(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	vmDeviceId, err := repo.GetVmDeviceIdByMac(info.Sn)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}
		return
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}
		return
	}

	mod, err := repo.GetNetworkByVmMac(info.Sn)
	if err != nil {
		if info.Type == "raw" {
			w.Write([]byte(""))
		} else {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": ""})
		}
		return
	}

	mac := vmDevice.Mac

	mod.Vlan = strings.Replace(mod.Vlan, "\r\n", "\n", -1)
	mod.Trunk = strings.Replace(mod.Trunk, "\r\n", "\n", -1)
	mod.Bonding = strings.Replace(mod.Bonding, "\r\n", "\n", -1)

	result := make(map[string]interface{})
	result["Hostname"] = vmDevice.Hostname
	result["Ip"] = vmDevice.Ip
	result["Netmask"] = mod.Netmask
	result["Gateway"] = mod.Gateway
	result["Vlan"] = mod.Vlan
	result["Trunk"] = mod.Trunk
	result["Bonding"] = mod.Bonding
	result["HWADDR"] = mac
	if info.Type == "raw" {
		w.Header().Add("Content-type", "text/html; charset=utf-8")
		var str string
		if vmDevice.Hostname != "" {
			str += "HOSTNAME=\"" + vmDevice.Hostname + "\""
		}
		if vmDevice.Ip != "" {
			str += "\nIPADDR=\"" + vmDevice.Ip + "\""
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
		str += "\nHWADDR=\"" + mac + "\""
		w.Write([]byte(str))
	} else {
		w.Header().Add("Content-type", "application/json; charset=utf-8")
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "成功获取network信息", "Content": result})
	}
}

func IsInPreInstallListForVm(ctx context.Context, w rest.ResponseWriter, mac string) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误", "Content": ""})
		return
	}

	mac = strings.TrimSpace(mac)
	if mac == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空!"})
		return
	}

	result := make(map[string]string)
	count, err := repo.CountVmDeviceByMac(mac)
	if err != nil || count <= 0 {
		result["Result"] = "false"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备不在安装列表里", "Content": result})
		return
	}

	vmDevice, err := repo.GetVmDeviceByMac(mac)
	if err != nil {
		result["Result"] = "false"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备不在安装列表里", "Content": result})
		return
	}

	if vmDevice.Status == "pre_install" || vmDevice.Status == "installing" {
		result["Result"] = "true"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备在安装列表里", "Content": result})
	} else {
		result["Result"] = "false"
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "该设备不在安装列表里", "Content": result})
	}
}

func ReportInstallInfoForVm(ctx context.Context, w rest.ResponseWriter, mac string, title string, installProgress float64, log string) {
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

	info.Sn = strings.TrimSpace(mac)
	info.Title = strings.TrimSpace(title)
	info.InstallProgress = installProgress
	info.InstallLog = strings.TrimSpace(log)
	if info.Sn == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "SN参数不能为空!"})
		return
	}

	vmDeviceId, err := repo.GetVmDeviceIdByMac(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该设备不存在!"})
		return
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if vmDevice.Status != "pre_install" && vmDevice.Status != "installing" {
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
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "安装进度参数不正确!"})
		return
	}

	_, errUpdate := repo.UpdateVmInstallInfoById(vmDevice.ID, status, info.InstallProgress)
	if errUpdate != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
		return
	}

	//删除PXE配置文件
	if info.InstallProgress == 1 {
		pxeFileName := util.GetPxeFileNameByMac(vmDevice.Mac)
		confDir := conf.OsInstall.PxeConfigDir
		if util.FileExist(confDir + "/" + pxeFileName) {
			err := os.Remove(confDir + "/" + pxeFileName)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
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

	_, errAddLog := repo.AddVmDeviceLog(vmDevice.ID, logTitle, "install", installLog)
	if errAddLog != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAddLog.Error()})
		return
	}

	//add report
	if info.InstallProgress == 1 {
		errReportLog := repo.CopyVmDeviceToInstallReport(vmDevice.ID)
		if errReportLog != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errReportLog.Error()})
			return
		}
	}

	result := make(map[string]string)
	result["Result"] = "true"
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}
