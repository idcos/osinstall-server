package route

import (
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	//"server/osinstallserver/util"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"github.com/qiniu/iconv"
	"io"
	"os"
	"regexp"
	"server/osinstallserver/util"
	"strings"
	"time"
)

func UploadDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "text/html; charset=utf-8")
	r.ParseForm()
	file, handle, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte("{\"Message\":\"" + err.Error() + "\",\"Status\":\"error\"}"))
		return
	}

	cd, err := iconv.Open("UTF-8", "GBK")
	if err != nil {
		w.Write([]byte("{\"Message\":\"" + err.Error() + "\",\"Status\":\"error\"}"))
		return
	}
	defer cd.Close()

	dir := "/tmp/cloudboot-server/"
	if !util.FileExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			w.Write([]byte("{\"Message\":\"" + err.Error() + "\",\"Status\":\"error\"}"))
			return
		}
	}

	list := strings.Split(handle.Filename, ".")
	fix := list[len(list)-1]

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s", time.Now().UnixNano()) + handle.Filename))
	cipherStr := h.Sum(nil)
	md5 := fmt.Sprintf("%s", hex.EncodeToString(cipherStr))
	filename := "osinstall-upload-" + md5 + "." + fix

	result := make(map[string]interface{})
	result["result"] = filename

	if util.FileExist(dir + filename) {
		os.Remove(dir + filename)
	}

	f, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		w.Write([]byte("{\"Message\":\"" + err.Error() + "\",\"Status\":\"error\"}"))
		return
	}
	defer f.Close()
	defer file.Close()

	data := map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result}
	json, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte("{\"Message\":\"" + err.Error() + "\",\"Status\":\"error\"}"))
		return
	}
	w.Write([]byte(json))
	return
}

func ImportPriview(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
		Filename string
		Limit    uint
		Offset   uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	file := "/tmp/cloudboot-server/" + info.Filename

	cd, err := iconv.Open("utf-8", "gbk") // convert gbk to utf8
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	defer cd.Close()

	input, err := os.Open(file)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	bufSize := 1024 * 1024
	read := iconv.NewReader(cd, input, bufSize)

	r2 := csv.NewReader(read)

	ra, err := r2.ReadAll()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	length := len(ra)

	type Device struct {
		ID                uint
		BatchNumber       string
		Sn                string
		Hostname          string
		Ip                string
		Netmask           string
		Gateway           string
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
		Content           string
		UserID            uint
		IsSupportVm       string
	}
	var success []Device
	var failure []Device
	//var result []string
	for i := 1; i < length; i++ {
		//result = append(result, ra[i][0])
		var row Device
		if len(ra[i]) != 10 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "导入文件格式错误!"
			failure = append(failure, row)
			continue
		}

		row.Sn = strings.TrimSpace(ra[i][0])
		row.Hostname = strings.TrimSpace(ra[i][1])
		row.Ip = strings.TrimSpace(ra[i][2])
		row.OsName = strings.TrimSpace(ra[i][3])
		row.HardwareName = strings.TrimSpace(ra[i][4])
		row.SystemName = strings.TrimSpace(ra[i][5])
		row.Location = strings.TrimSpace(ra[i][6])
		row.AssetNumber = strings.TrimSpace(ra[i][7])
		row.ManageIp = strings.TrimSpace(ra[i][8])
		row.IsSupportVm = strings.TrimSpace(ra[i][9])
		row.UserID = session.ID

		if len(row.Sn) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN长度超过255限制!"
		}

		if len(row.Hostname) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "主机名长度超过255限制!"
		}

		if len(row.Location) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "位置长度超过255限制!"
		}

		if len(row.AssetNumber) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "财编长度超过255限制!"
		}

		if row.Sn == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN不能为空!"
		}

		if row.Hostname == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "主机名不能为空!"
		}

		if row.Ip == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "IP不能为空!"
		}

		if row.OsName == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "操作系统不能为空!"
		}

		if row.SystemName == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "系统安装模板不能为空!"
		}

		if row.Location == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "位置不能为空!"
		}

		if row.IsSupportVm != "" && row.IsSupportVm != "Yes" && row.IsSupportVm != "No" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "是否支持虚拟机的参数格式不正确!"
		}

		//match manufacturer
		countManufacturer, errCountManufacturer := repo.CountManufacturerBySn(row.Sn)
		if errCountManufacturer != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCountManufacturer.Error()})
			return
		}
		if countManufacturer <= 0 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "未在【资源池管理】里匹配到该SN，请先将该设备加电并进入BootOS!"
		}
		if countManufacturer > 0 {
			//validate user from manufacturer
			manufacturer, err := repo.GetManufacturerBySn(row.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}
			if session.Role != "Administrator" && manufacturer.UserID != session.ID {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "您无权操作其他人的设备!"
			}
		}

		//validate ip from vm device
		countVmIp, errVmIp := repo.CountVmDeviceByIp(row.Ip)
		if errVmIp != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVmIp.Error()})
			return
		}
		if countVmIp > 0 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + row.Ip + " 该IP已被虚拟机使用!"
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

			device, err := repo.GetDeviceBySn(row.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}

			if session.Role != "Administrator" && device.UserID != session.ID {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该设备已被其他人录入，不能重复录入!"
			}
			/*else {
				if device.Status == "success" {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "该设备已安装成功，请使用【单台录入】的功能重新录入并安装"
				}
			}
			*/

			//hostname
			countHostname, err := repo.CountDeviceByHostnameAndId(row.Hostname, ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
				return
			}
			if countHostname > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该主机名已存在!"
			}

			//IP
			countIp, err := repo.CountDeviceByIpAndId(row.Ip, ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该IP已存在!"
			}

			if row.ManageIp != "" {
				//IP
				countManageIp, err := repo.CountDeviceByManageIpAndId(row.ManageIp, ID)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if countManageIp > 0 {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "该管理IP已存在!"
				}
			}

			//validate host server info
			countVm, errVm := repo.CountVmDeviceByDeviceId(row.ID)
			if errVm != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVm.Error()})
				return
			}
			if countVm > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该物理机下(SN:" + device.Sn + ")还存留有虚拟机，不允许修改信息。请先销毁虚拟机后再操作！"
			}

		} else {
			//hostname
			countHostname, err := repo.CountDeviceByHostname(row.Hostname)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
				return
			}
			if countHostname > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该主机名已存在!"
			}

			//IP
			countIp, err := repo.CountDeviceByIp(row.Ip)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该IP已存在!"
			}

			if row.ManageIp != "" {
				//IP
				countManageIp, err := repo.CountDeviceByManageIp(row.ManageIp)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if countManageIp > 0 {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "该管理IP已存在!"
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
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "IP格式不正确!"
		}

		modelIp, err := repo.GetIpByIp(row.Ip)
		if err != nil {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "未匹配到网段!"
		} else {
			network, errNetwork := repo.GetNetworkById(modelIp.NetworkID)
			if errNetwork != nil {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "未匹配到网段!"
			}
			row.NetworkName = network.Network
			row.Netmask = network.Netmask
			row.Gateway = network.Gateway
		}

		if row.ManageIp != "" {
			//匹配网络
			isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", row.ManageIp)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}

			if !isValidate {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "管理IP格式不正确!"
			}

			modelIp, err := repo.GetManageIpByIp(row.ManageIp)
			if err != nil {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "未匹配到管理网段!"
			} else {
				network, errNetwork := repo.GetManageNetworkById(modelIp.NetworkID)
				if errNetwork != nil {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "未匹配到管理网段!"
				}
				row.ManageNetworkID = network.ID
				row.ManageNetworkName = network.Network
			}
		}

		//OSName
		countOs, err := repo.CountOsConfigByName(row.OsName)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countOs <= 0 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "未匹配到操作系统!"
		}

		//SystemName
		countSystem, err := repo.CountSystemConfigByName(row.SystemName)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if countSystem <= 0 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "未匹配到系统安装模板!"
		}

		if row.HardwareName != "" {
			//HardwareName
			countHardware, err := repo.CountHardwareWithSeparator(row.HardwareName)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countHardware <= 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "未匹配到硬件配置模板!"
			} else {
				hardware, err := repo.GetHardwareBySeaprator(row.HardwareName)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
				row.HardwareID = hardware.ID
			}
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
						var br string
						if row.Content != "" {
							br = "<br />"
						}
						row.Content += br + "硬件配置模板的OOB网络类型为静态IP的方式，请填写管理IP!"
					}
				}
			}
		}

		/*
			if row.Location != "" {
				locationId, err := repo.GetLocationIdByName(row.Location)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
				if locationId <= 0 {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "未匹配到位置!"
				} else {
					row.LocationID = locationId
				}
			}
		*/

		if row.Content != "" {
			failure = append(failure, row)
		} else {
			success = append(success, row)
		}
	}

	var data []Device
	if len(failure) > 0 {
		data = failure
	} else {
		data = success
	}
	var result []Device
	for i := 0; i < len(data); i++ {
		if uint(i) >= info.Offset && uint(i) < (info.Offset+info.Limit) {
			result = append(result, data[i])
		}
	}

	if len(failure) > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "failure", "Message": "设备信息不正确", "recordCount": len(data), "Content": result})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "recordCount": len(data), "Content": result})
	}
}

func ImportDevice(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
		Filename string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	file := "/tmp/cloudboot-server/" + info.Filename

	cd, err := iconv.Open("utf-8", "gbk") // convert gbk to utf8
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	defer cd.Close()

	input, err := os.Open(file)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	bufSize := 1024 * 1024
	read := iconv.NewReader(cd, input, bufSize)

	r2 := csv.NewReader(read)
	ra, err := r2.ReadAll()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	length := len(ra)

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
	}

	batchNumber, err := repo.CreateBatchNumber()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	//var result []string
	for i := 1; i < length; i++ {
		//result = append(result, ra[i][0])
		var row Device

		if len(ra[i]) != 10 {
			continue
		}

		row.Sn = strings.TrimSpace(ra[i][0])
		row.Hostname = strings.TrimSpace(ra[i][1])
		row.Ip = strings.TrimSpace(ra[i][2])
		row.OsName = strings.TrimSpace(ra[i][3])
		row.HardwareName = strings.TrimSpace(ra[i][4])
		row.SystemName = strings.TrimSpace(ra[i][5])
		row.Location = strings.TrimSpace(ra[i][6])
		row.AssetNumber = strings.TrimSpace(ra[i][7])
		row.ManageIp = strings.TrimSpace(ra[i][8])
		row.IsSupportVm = strings.TrimSpace(ra[i][9])
		row.UserID = session.ID

		if len(row.Sn) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN长度超过255限制!"
		}

		if len(row.Hostname) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "主机名长度超过255限制!"
		}

		if len(row.Location) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "位置长度超过255限制!"
		}

		if len(row.AssetNumber) > 255 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "财编长度超过255限制!"
		}

		if row.Sn == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN不能为空!"
		}

		if row.Hostname == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "主机名不能为空!"
		}

		if row.Ip == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "IP不能为空!"
		}

		if row.OsName == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "操作系统不能为空!"
		}

		if row.SystemName == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "系统安装模板不能为空!"
		}

		if row.Location == "" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "位置不能为空!"
		}

		if row.IsSupportVm != "" && row.IsSupportVm != "Yes" && row.IsSupportVm != "No" {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "是否支持虚拟机的参数格式不正确!"
		}

		if row.IsSupportVm != "Yes" {
			row.IsSupportVm = "No"
		}

		//match manufacturer
		countManufacturer, errCountManufacturer := repo.CountManufacturerBySn(row.Sn)
		if errCountManufacturer != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errCountManufacturer.Error()})
			return
		}
		if countManufacturer <= 0 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "未在【资源池管理】里匹配到该SN，请先将该设备加电并进入BootOS!"
		}
		if countManufacturer > 0 {
			//validate user from manufacturer
			manufacturer, err := repo.GetManufacturerBySn(row.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}
			if session.Role != "Administrator" && manufacturer.UserID != session.ID {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "您无权操作其他人的设备!"
			}
		}

		//validate ip from vm device
		countVmIp, errVmIp := repo.CountVmDeviceByIp(row.Ip)
		if errVmIp != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVmIp.Error()})
			return
		}
		if countVmIp > 0 {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + row.Ip + " 该IP已被虚拟机使用!"
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

			device, err := repo.GetDeviceBySn(row.Sn)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}

			if session.Role != "Administrator" && device.UserID != session.ID {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + " 该设备已被录入，不能重复录入!"
			}

			//hostname
			countHostname, err := repo.CountDeviceByHostnameAndId(row.Hostname, ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
				return
			}
			if countHostname > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "该主机名已存在!"
			}

			//IP
			countIp, err := repo.CountDeviceByIpAndId(row.Ip, ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "该IP已存在!"
			}

			if row.ManageIp != "" {
				//IP
				countManageIp, err := repo.CountDeviceByManageIpAndId(row.ManageIp, ID)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if countManageIp > 0 {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "SN:" + row.Sn + "该管理IP已存在!"
				}
			}

			//validate host server info
			countVm, errVm := repo.CountVmDeviceByDeviceId(row.ID)
			if errVm != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVm.Error()})
				return
			}
			if countVm > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "该物理机下(SN:" + device.Sn + ")还存留有虚拟机，不允许修改信息。请先销毁虚拟机后再操作！"
			}
		} else {
			//hostname
			countHostname, err := repo.CountDeviceByHostname(row.Hostname)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误:" + err.Error()})
				return
			}
			if countHostname > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "该主机名已存在!"
			}

			//IP
			countIp, err := repo.CountDeviceByIp(row.Ip)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}

			if countIp > 0 {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "该IP已存在!"
			}

			if row.ManageIp != "" {
				//IP
				countManageIp, err := repo.CountDeviceByManageIp(row.ManageIp)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if countManageIp > 0 {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "SN:" + row.Sn + "该管理IP已存在!"
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
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN:" + row.Sn + "IP格式不正确!"
		}

		modelIp, err := repo.GetIpByIp(row.Ip)
		if err != nil {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN:" + row.Sn + "未匹配到网段!"
		}

		_, errNetwork := repo.GetNetworkById(modelIp.NetworkID)
		if errNetwork != nil {
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN:" + row.Sn + "未匹配到网段!"
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
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "管理IP格式不正确!"
			}

			modelIp, err := repo.GetManageIpByIp(row.ManageIp)
			if err != nil {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "未匹配到管理网段!"
			}

			_, errNetwork := repo.GetManageNetworkById(modelIp.NetworkID)
			if errNetwork != nil {
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "未匹配到管理网段!"
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
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN:" + row.Sn + "未匹配到操作系统!"
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
			var br string
			if row.Content != "" {
				br = "<br />"
			}
			row.Content += br + "SN:" + row.Sn + "未匹配到系统安装模板!"
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
				var br string
				if row.Content != "" {
					br = "<br />"
				}
				row.Content += br + "SN:" + row.Sn + "未匹配到硬件配置模板!"
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
						var br string
						if row.Content != "" {
							br = "<br />"
						}
						row.Content += br + "SN:" + row.Sn + "硬件配置模板的OOB网络类型为静态IP的方式，请填写管理IP!"
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
				/*
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "SN:" + row.Sn + " 未匹配到位置!"
				*/
				locationId, err := repo.ImportLocation(row.Location)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}

				if locationId <= uint(0) {
					var br string
					if row.Content != "" {
						br = "<br />"
					}
					row.Content += br + "SN:" + row.Sn + " 未匹配到位置!"
				}
				row.LocationID = locationId

			}
		}
		if row.Content != "" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": row.Content})
			return
		} else {
			status := "pre_install"
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
					_, errUpdate := repo.UpdateManufacturerDeviceIdById(manufacturerId, id)
					if errUpdate != nil {
						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error()})
						return
					}
				}

				//delete host server info
				_, errDeleteHost := repo.DeleteVmHostBySn(deviceOld.Sn)
				if errDeleteHost != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDeleteHost.Error()})
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
		}
	}

	//删除文件
	if util.FileExist(file) {
		err := os.Remove(file)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	}
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
