package route

import (
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"github.com/qiniu/iconv"
	"golang.org/x/net/context"
	"middleware"
	"strings"
	"utils"
)

func GetScanDeviceList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Limit      uint
		Offset     uint
		Keyword    string
		Company    string
		Product    string
		ModelName  string
		CpuRule    string
		Cpu        string
		MemoryRule string
		Memory     string
		DiskRule   string
		Disk       string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	info.Keyword = strings.TrimSpace(info.Keyword)
	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)
	info.ModelName = strings.TrimSpace(info.ModelName)
	info.CpuRule = strings.TrimSpace(info.CpuRule)
	info.Cpu = strings.TrimSpace(info.Cpu)
	info.MemoryRule = strings.TrimSpace(info.MemoryRule)
	info.Memory = strings.TrimSpace(info.Memory)
	info.DiskRule = strings.TrimSpace(info.DiskRule)
	info.Disk = strings.TrimSpace(info.Disk)

	var where string
	where = "device_id = 0 "

	if info.Company != "" {
		where += " and company = '" + info.Company + "'"
	}
	if info.Product != "" {
		where += " and product = '" + info.Product + "'"
	}
	if info.ModelName != "" {
		where += " and model_name = '" + info.ModelName + "'"
	}
	if info.CpuRule != "" && info.Cpu != "" {
		where += " and cpu " + info.CpuRule + info.Cpu
	}
	if info.MemoryRule != "" && info.Memory != "" {
		where += " and memory " + info.MemoryRule + info.Memory
	}
	if info.DiskRule != "" && info.Disk != "" {
		where += " and disk " + info.DiskRule + info.Disk
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
			where += str + " sn = '" + v + "' or ip = '" + v + "' or company = '" + v + "' or product = '" + v + "' or model_name = '" + v + "'"
		}
		where += " ) "
	}

	mods, err := repo.GetManufacturerListWithPage(info.Limit, info.Offset, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountManufacturerByWhere(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func ExportScanDeviceList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Keyword    string
		Company    string
		Product    string
		ModelName  string
		CpuRule    string
		Cpu        string
		MemoryRule string
		Memory     string
		DiskRule   string
		Disk       string
	}

	info.Keyword = r.FormValue("Keyword")
	info.Company = r.FormValue("Company")
	info.Product = r.FormValue("Product")
	info.ModelName = r.FormValue("ModelName")
	info.CpuRule = r.FormValue("CpuRule")
	info.Cpu = r.FormValue("Cpu")
	info.MemoryRule = r.FormValue("MemoryRule")
	info.Memory = r.FormValue("Memory")
	info.DiskRule = r.FormValue("DiskRule")
	info.Disk = r.FormValue("Disk")

	info.Keyword = strings.TrimSpace(info.Keyword)
	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)
	info.ModelName = strings.TrimSpace(info.ModelName)
	info.CpuRule = strings.TrimSpace(info.CpuRule)
	info.Cpu = strings.TrimSpace(info.Cpu)
	info.MemoryRule = strings.TrimSpace(info.MemoryRule)
	info.Memory = strings.TrimSpace(info.Memory)
	info.DiskRule = strings.TrimSpace(info.DiskRule)
	info.Disk = strings.TrimSpace(info.Disk)

	var where string
	where = "device_id = 0 "

	idsParam := r.FormValue("ids")
	if idsParam != "" {
		ids := strings.Split(idsParam, ",")
		if len(ids) > 0 {
			where += " and id in (" + strings.Join(ids, ",") + ")"
		}
	}

	if info.Company != "" {
		where += " and company = '" + info.Company + "'"
	}
	if info.Product != "" {
		where += " and product = '" + info.Product + "'"
	}
	if info.ModelName != "" {
		where += " and model_name = '" + info.ModelName + "'"
	}
	if info.CpuRule != "" && info.Cpu != "" {
		where += " and cpu " + info.CpuRule + info.Cpu
	}
	if info.MemoryRule != "" && info.Memory != "" {
		where += " and memory " + info.MemoryRule + info.Memory
	}
	if info.DiskRule != "" && info.Disk != "" {
		where += " and disk " + info.DiskRule + info.Disk
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
			where += str + " sn = '" + v + "' or ip = '" + v + "' or company = '" + v + "' or product = '" + v + "' or model_name = '" + v + "'"
		}
		where += " ) "
	}

	mods, err := repo.GetManufacturerListWithPage(1000000, 0, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	var str string
	str = "SN(必填),主机名(必填),IP(必填),操作系统(必填),硬件配置模板,系统安装模板(必填),位置(必填),财编\n"
	for _, device := range mods {
		str += device.Sn + ","
		str += ","
		str += ","
		str += ","
		str += ","
		str += ","
		str += ","
		str += "\n"
	}
	fmt.Println(str)

	cd, err := iconv.Open("gbk", "utf-8") // convert utf-8 to gbk
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	defer cd.Close()
	gbkStr := cd.ConvString(str)

	bytes := []byte(gbkStr)

	filename := "idcos-osinstall-scan-device.csv"
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename='%s';filename*=utf-8''%s", filename, filename))
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Write(bytes)
}

func GetScanDeviceById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mod, err := repo.GetManufacturerById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type DeviceWithTime struct {
		ID          uint
		DeviceID    uint
		Company     string
		Product     string
		ModelName   string
		Sn          string
		Ip          string
		Mac         string
		Nic         string
		Cpu         string
		Memory      string
		Disk        string
		Motherboard string
		Raid        string
		Oob         string
		CreatedAt   utils.ISOTime
		UpdatedAt   utils.ISOTime
	}

	var device DeviceWithTime
	device.ID = mod.ID
	device.DeviceID = mod.DeviceID
	device.Company = mod.Company
	device.Product = mod.Product
	device.ModelName = mod.ModelName
	device.Sn = mod.Sn
	device.Ip = mod.Ip
	device.Mac = mod.Mac
	device.Nic = mod.Nic
	device.Cpu = mod.Cpu
	device.Memory = mod.Memory
	device.Disk = mod.Disk
	device.Motherboard = mod.Motherboard
	device.Raid = mod.Raid
	device.Oob = mod.Oob

	device.CreatedAt = utils.ISOTime(mod.CreatedAt)
	device.UpdatedAt = utils.ISOTime(mod.UpdatedAt)

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": device})
}

func GetScanDeviceCompany(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var where string
	where = "device_id = 0"
	mod, err := repo.GetManufacturerCompanyByGroup(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetScanDeviceProduct(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		Company string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}
	info.Company = strings.TrimSpace(info.Company)

	var where string
	where = "device_id = 0 and company = '" + info.Company + "'"
	mod, err := repo.GetManufacturerProductByGroup(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetScanDeviceModelName(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		Company string
		Product string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}
	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)

	var where string
	where = "device_id = 0 and company = '" + info.Company + "' and product = '" + info.Product + "'"
	mod, err := repo.GetManufacturerModelNameByGroup(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
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
		Sn          string
		Company     string
		Product     string
		ModelName   string
		Ip          string
		Mac         string
		Nic         string
		Cpu         string
		Memory      string
		Disk        string
		Motherboard string
		Raid        string
		Oob         string
		DeviceID    uint
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

	countDevice, err := repo.CountDeviceBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	if countDevice > 0 {
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

		info.DeviceID = device.ID
		/*
			count, err := repo.CountManufacturerByDeviceID(device.ID)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
		*/

	} else {
		info.DeviceID = uint(0)
	}

	result := make(map[string]string)

	count, err := repo.CountManufacturerBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	if count > 0 {
		id, err := repo.GetManufacturerIdBySn(info.Sn)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error(), "Content": result})
			return
		}

		_, errUpdate := repo.UpdateManufacturerById(id, info.Company, info.Product, info.ModelName, info.Sn, info.Ip, info.Mac, info.Nic, info.Cpu, info.Memory, info.Disk, info.Motherboard, info.Raid, info.Oob)
		if errUpdate != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errUpdate.Error(), "Content": result})
			return
		}

	} else {
		_, err := repo.AddManufacturer(info.DeviceID, info.Company, info.Product, info.ModelName, info.Sn, info.Ip, info.Mac, info.Nic, info.Cpu, info.Memory, info.Disk, info.Motherboard, info.Raid, info.Oob)
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
