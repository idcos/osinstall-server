package route

import (
	//"encoding/base64"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	//"json"
	"encoding/json"
	"middleware"
	//"strconv"
	"strings"
	//"net/http"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/qiniu/iconv"
	"io"
	"io/ioutil"
	"model"
	"net/http"
	"os"
	"server/osinstallserver/util"
	"time"
	"utils"
)

func DeleteHardwareById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID          uint
		AccessToken string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.AccessToken = strings.TrimSpace(info.AccessToken)
	_, errVerify := VerifyAccessPurview(info.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	hardware, err := repo.GetHardwareById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if hardware.IsSystemAdd == "Yes" {
		//w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "系统添加的配置不允许删除!"})
		//return
	}

	mod, err := repo.DeleteHardwareById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func UpdateHardwareById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID          uint
		Company     string
		Product     string
		ModelName   string
		Raid        string
		Oob         string
		Bios        string
		Tpl         string
		Data        string
		Source      string
		Version     string
		Status      string
		AccessToken string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.AccessToken = strings.TrimSpace(info.AccessToken)
	_, errVerify := VerifyAccessPurview(info.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)
	info.ModelName = strings.TrimSpace(info.ModelName)
	info.Tpl = strings.TrimSpace(info.Tpl)
	info.Data = strings.TrimSpace(info.Data)
	info.Source = strings.TrimSpace(info.Source)
	info.Version = strings.TrimSpace(info.Version)
	info.Status = strings.TrimSpace(info.Status)

	if info.Company == "" || info.ModelName == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	count, err := repo.CountHardwareByCompanyAndProductAndNameAndId(info.Company, info.Product, info.ModelName, info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该硬件型号已存在!"})
		return
	}

	hardware, err := repo.GetHardwareById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if hardware.IsSystemAdd == "Yes" {
		//w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "系统添加的配置不允许修改!"})
		//return
	}

	mod, err := repo.UpdateHardwareById(info.ID, info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.Tpl, info.Data, info.Source, info.Version, info.Status)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetHardwareById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mod, err := repo.GetHardwareById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetCompanyByGroup(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	mod, err := repo.GetCompanyByGroup()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func ExportHardware(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	/*
		var info struct {
			Ids []int
		}

			if err := r.DecodeJSONPayload(&info); err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
				return
			}
	*/

	var where string
	where = " and is_system_add = 'Yes' "
	idsParam := r.FormValue("ids")
	if idsParam != "" {
		ids := strings.Split(idsParam, ",")
		if len(ids) > 0 {
			/*
				for _, id := range info.Ids {
					ids = append(ids, strconv.Itoa(id))
				}
			*/
			where += " and id in (" + strings.Join(ids, ",") + ")"
		}
	}

	company := r.FormValue("company")
	if company != "" {
		where += " and company = '" + company + "' "
	}

	product := r.FormValue("product")
	if product != "" {
		where += " and product = '" + product + "' "
	}

	modelName := r.FormValue("modelName")
	if modelName != "" {
		where += " and model_name = '" + modelName + "' "
	}

	mods, err := repo.GetHardwareListWithPage(10000, 0, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	var result []map[string]interface{}
	for _, v := range mods {
		result2 := make(map[string]interface{})
		result2["Company"] = v.Company
		result2["Product"] = v.Product
		result2["ModelName"] = v.ModelName
		result2["IsSystemAdd"] = v.IsSystemAdd
		result2["Tpl"] = v.Tpl
		result2["Data"] = v.Data
		result = append(result, result2)
	}

	filename := "idcos-osinstall-hardware.json"
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename='%s';filename*=utf-8''%s", filename, filename))
	w.Header().Add("Content-Type", "application/octet-stream")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProductByWhereAndGroup(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		Company string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	where := " company = '" + info.Company + "'"
	mod, err := repo.GetProductByWhereAndGroup(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetModelNameByWhereAndGroup(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		Company     string
		Product     string
		IsSystemAdd string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	where := " company = '" + info.Company + "'"
	if info.Product != "" {
		where += " and product = '" + info.Product + "'"
	}

	if info.IsSystemAdd != "" {
		where += " and is_system_add = '" + info.IsSystemAdd + "'"
	}

	mod, err := repo.GetModelNameByWhereAndGroup(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetHardwareList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Limit       uint
		Offset      uint
		Company     string
		Product     string
		ModelName   string
		IsSystemAdd string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	var where string
	if info.Company != "" {
		where += " and company = '" + info.Company + "' "
	}
	if info.Product != "" {
		where += " and product = '" + info.Product + "' "
	}
	if info.ModelName != "" {
		where += " and model_name = '" + info.ModelName + "' "
	}
	if info.IsSystemAdd != "" {
		where += " and is_system_add = '" + info.IsSystemAdd + "' "
	}

	mods, err := repo.GetHardwareListWithPage(info.Limit, info.Offset, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})

	type Hardware struct {
		ID          uint
		ShowName    string
		Company     string
		Product     string
		ModelName   string
		Raid        string
		Oob         string
		Bios        string
		IsSystemAdd string
		Tpl         string
		Data        string
		Source      string
		Version     string
		Status      string
	}
	var hardwares []Hardware
	for _, v := range mods {
		var hardware Hardware
		hardware.ID = v.ID
		hardware.ShowName = v.Company + "-" + v.ModelName
		hardware.Company = v.Company
		hardware.Product = v.Product
		hardware.ModelName = v.ModelName
		hardware.Raid = v.Raid
		hardware.Oob = v.Oob
		hardware.Bios = v.Bios
		hardware.IsSystemAdd = v.IsSystemAdd
		hardware.Tpl = v.Tpl
		hardware.Data = v.Data
		hardware.Source = v.Source
		hardware.Version = v.Version
		hardware.Status = v.Status
		hardwares = append(hardwares, hardware)
	}
	result["list"] = hardwares

	//总条数
	count, err := repo.CountHardware(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

//添加
func AddHardware(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Company     string
		Product     string
		ModelName   string
		Raid        string
		Oob         string
		Bios        string
		IsSystemAdd string
		Tpl         string
		Data        string
		Source      string
		Version     string
		Status      string
		AccessToken string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)
	info.ModelName = strings.TrimSpace(info.ModelName)
	info.IsSystemAdd = strings.TrimSpace(info.IsSystemAdd)
	info.Tpl = strings.TrimSpace(info.Tpl)
	info.Data = strings.TrimSpace(info.Data)
	info.Source = strings.TrimSpace(info.Source)
	info.Version = strings.TrimSpace(info.Version)
	info.Status = strings.TrimSpace(info.Status)

	info.AccessToken = strings.TrimSpace(info.AccessToken)
	_, errVerify := VerifyAccessPurview(info.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	if info.Company == "" || info.ModelName == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	count, err := repo.CountHardwareByCompanyAndProductAndName(info.Company, info.Product, info.ModelName)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该硬件型号已存在!"})
		return
	}

	_, errAdd := repo.AddHardware(info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.IsSystemAdd, info.Tpl, info.Data, info.Source, info.Version, info.Status)
	if errAdd != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAdd.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func UploadCompanyHardware(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()
	file, handle, err := r.FormFile("file")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	cd, err := iconv.Open("UTF-8", "GBK")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	defer cd.Close()

	dir := "./upload/"
	if !util.FileExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	}

	list := strings.Split(handle.Filename, ".")
	fix := list[len(list)-1]

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s", time.Now().UnixNano()) + handle.Filename))
	cipherStr := h.Sum(nil)
	md5 := fmt.Sprintf("%s", hex.EncodeToString(cipherStr))
	filename := md5 + "." + fix

	result := make(map[string]interface{})
	result["result"] = filename

	if util.FileExist(dir + filename) {
		os.Remove(dir + filename)
	}

	f, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	defer f.Close()
	defer file.Close()

	fileHandle, err := os.Open(dir + filename)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	var data []map[string]interface{}
	bytes, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	errDecode := json.Unmarshal(bytes, &data)
	if errDecode != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "文件格式错误:" + errDecode.Error()})
		return
	}
	for _, v := range data {
		if v["Company"] == "" || v["ModelName"] == "" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "文件格式错误!"})
			return
		}
	}

	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	for _, v := range data {
		var info struct {
			Company     string
			Product     string
			ModelName   string
			Raid        string
			Oob         string
			Bios        string
			IsSystemAdd string
			Tpl         string
			Data        string
			Source      string
			Version     string
			Status      string
		}

		info.Company = strings.TrimSpace(v["Company"].(string))
		info.Product = strings.TrimSpace(v["Product"].(string))
		info.ModelName = strings.TrimSpace(v["ModelName"].(string))
		info.IsSystemAdd = "Yes"
		info.Tpl = strings.TrimSpace(v["Tpl"].(string))
		info.Data = strings.TrimSpace(v["Data"].(string))
		info.Source = strings.TrimSpace(info.Source)
		info.Version = strings.TrimSpace(info.Version)
		info.Status = "Pending"

		where := "company = '" + info.Company + "' and product = '" + info.Product + "' and model_name = '" + info.ModelName + "' and is_system_add = 'Yes'"

		count, err := repo.CountHardwareByWhere(where)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if count > 0 {
			hardware, err := repo.GetHardwareByWhere(where)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
			var isUpdate bool
			isUpdate = false
			if hardware.Data != info.Data {
				isUpdate = true
			}
			if hardware.Tpl != info.Tpl {
				isUpdate = true
			}
			if isUpdate == true {
				_, err := repo.UpdateHardwareById(hardware.ID, info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.Tpl, info.Data, info.Source, info.Version, info.Status)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
			}
		} else {
			_, errAdd := repo.AddHardware(info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.IsSystemAdd, info.Tpl, info.Data, info.Source, info.Version, info.Status)
			if errAdd != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAdd.Error()})
				return
			}
		}
	}
	//delete tmp file
	errDeleteFile := os.Remove(dir + filename)
	if errDeleteFile != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDeleteFile.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
	return
}

func UploadHardware(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()
	file, handle, err := r.FormFile("file")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	cd, err := iconv.Open("UTF-8", "GBK")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	defer cd.Close()

	dir := "./upload/"
	if !util.FileExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	}

	list := strings.Split(handle.Filename, ".")
	fix := list[len(list)-1]

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s", time.Now().UnixNano()) + handle.Filename))
	cipherStr := h.Sum(nil)
	md5 := fmt.Sprintf("%s", hex.EncodeToString(cipherStr))
	filename := md5 + "." + fix

	result := make(map[string]interface{})
	result["result"] = filename

	if util.FileExist(dir + filename) {
		os.Remove(dir + filename)
	}

	f, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	defer f.Close()
	defer file.Close()

	fileHandle, err := os.Open(dir + filename)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	var data []map[string]interface{}
	bytes, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	errDecode := json.Unmarshal(bytes, &data)
	if errDecode != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "文件格式错误:" + errDecode.Error()})
		return
	}
	for _, v := range data {
		if v["Company"] == "" || v["ModelName"] == "" {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "文件格式错误!"})
			return
		}
	}

	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	for _, v := range data {
		var info struct {
			Company     string
			Product     string
			ModelName   string
			Raid        string
			Oob         string
			Bios        string
			IsSystemAdd string
			Tpl         string
			Data        string
			Source      string
			Version     string
			Status      string
		}

		info.Company = strings.TrimSpace(v["Company"].(string))
		info.Product = strings.TrimSpace(v["Product"].(string))
		info.ModelName = strings.TrimSpace(v["ModelName"].(string))
		info.IsSystemAdd = "Yes"
		info.Tpl = strings.TrimSpace(v["Tpl"].(string))
		info.Data = strings.TrimSpace(v["Data"].(string))
		info.Source = strings.TrimSpace(info.Source)
		info.Version = strings.TrimSpace(info.Version)
		info.Status = "Success"

		where := "company = '" + info.Company + "' and product = '" + info.Product + "' and model_name = '" + info.ModelName + "' and is_system_add = 'Yes'"

		count, err := repo.CountHardwareByWhere(where)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		if count > 0 {
			hardware, err := repo.GetHardwareByWhere(where)
			if err != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
				return
			}
			var isUpdate bool
			isUpdate = false
			if hardware.Data != info.Data {
				isUpdate = true
			}
			if hardware.Tpl != info.Tpl {
				isUpdate = true
			}
			if isUpdate == true {
				_, err := repo.UpdateHardwareById(hardware.ID, info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.Tpl, info.Data, info.Source, info.Version, info.Status)
				if err != nil {
					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
					return
				}
			}
		} else {
			_, errAdd := repo.AddHardware(info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.IsSystemAdd, info.Tpl, info.Data, info.Source, info.Version, info.Status)
			if errAdd != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAdd.Error()})
				return
			}
		}
	}
	//delete tmp file
	errDeleteFile := os.Remove(dir + filename)
	if errDeleteFile != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDeleteFile.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
	return
}

func CheckOnlineUpdate(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	lastestHardware, err := repo.GetLastestVersionHardware()
	type HardwareVersionWithDate struct {
		Version string
		Date    utils.ISOTime
	}
	var hardwareVersionWithDate HardwareVersionWithDate
	hardwareVersionWithDate.Version = lastestHardware.Version
	hardwareVersionWithDate.Date = utils.ISOTime(lastestHardware.CreatedAt)

	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	result := make(map[string]interface{})
	result["Version"] = lastestHardware.Version
	b, err := json.Marshal(result)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	body := bytes.NewBuffer([]byte(b))
	resp, err := http.Post("http://open.idcos.com/api/x86/check-online-update", "application/json;charset=utf-8", body)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "网络连接故障，在线更新配置库失败！", "CurrentVersion": hardwareVersionWithDate})
		return
	}
	reportResult, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type content struct {
		Version string `json:"Version"`
		Date    string `json:"Date"`
	}
	type Response struct {
		Status  string    `json:"Status"`
		Message string    `json:"Message"`
		Content []content `json:"Content"`
	}
	var response Response
	errJson := json.Unmarshal(reportResult, &response)
	if errJson != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errJson.Error()})
		return
	}
	if response.Status == "success" {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "", "Content": response.Content, "CurrentVersion": hardwareVersionWithDate})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作失败！"})
	}
}

func RunOnlineUpdate(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	lastestHardware, err := repo.GetLastestVersionHardware()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	result := make(map[string]interface{})
	result["Version"] = lastestHardware.Version

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	body := bytes.NewBuffer([]byte(b))
	resp, err := http.Post("http://open.idcos.com/api/x86/run-online-update", "application/json;charset=utf-8", body)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "网络连接故障，操作失败！您可以通过线下方式分享给我们，谢谢！"})
		return
	}
	reportResult, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type Response struct {
		Status  string           `json:"Status"`
		Message string           `json:"Message"`
		Content []model.Hardware `json:"Content"`
	}
	var response Response
	errJson := json.Unmarshal(reportResult, &response)
	if errJson != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errJson.Error()})
		return
	}
	if response.Status == "success" {
		var fix = "_" + time.Now().Format("20060102")
		errBack := repo.CreateHardwareBackupTable(fix)
		if errBack != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errBack.Error()})
			return
		}

		for _, info := range response.Content {
			info.Company = strings.TrimSpace(info.Company)
			info.Product = strings.TrimSpace(info.Product)
			info.ModelName = strings.TrimSpace(info.ModelName)
			info.IsSystemAdd = strings.TrimSpace(info.IsSystemAdd)
			info.Tpl = strings.TrimSpace(info.Tpl)
			info.Data = strings.TrimSpace(info.Data)
			info.Source = strings.TrimSpace(info.Source)
			info.Version = strings.TrimSpace(info.Version)
			info.Status = strings.TrimSpace(info.Status)

			if info.Status == "Success" {
				if info.Company == "" || info.ModelName == "" {
					continue
				}

				count, err := repo.CountHardwareByCompanyAndProductAndName(info.Company, info.Product, info.ModelName)
				if err != nil {
					errRollback := repo.RollbackHardwareFromBackupTable(fix)
					message := err.Error()
					if errRollback != nil {
						message += errRollback.Error()
					}

					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
					return
				}

				if count > 0 {
					hardwareId, err := repo.GetHardwareIdByCompanyAndProductAndName(info.Company, info.Product, info.ModelName)
					if err != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := err.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
					hardware, err := repo.GetHardwareById(hardwareId)
					if err != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := err.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
					if hardware.IsSystemAdd != "Yes" {
						continue
					}

					_, errUpdate := repo.UpdateHardwareById(hardware.ID, info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.Tpl, info.Data, info.Source, info.Version, info.Status)
					if errUpdate != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := errUpdate.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
				} else {
					_, errAdd := repo.AddHardware(info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.IsSystemAdd, info.Tpl, info.Data, info.Source, info.Version, info.Status)
					if errAdd != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := errAdd.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
				}
			} else if info.Status == "Deleted" {
				count, err := repo.CountHardwareByCompanyAndProductAndName(info.Company, info.Product, info.ModelName)
				if err != nil {
					errRollback := repo.RollbackHardwareFromBackupTable(fix)
					message := err.Error()
					if errRollback != nil {
						message += errRollback.Error()
					}

					w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
					return
				}

				if count > 0 {
					hardwareId, err := repo.GetHardwareIdByCompanyAndProductAndName(info.Company, info.Product, info.ModelName)
					if err != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := err.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
					hardware, err := repo.GetHardwareById(hardwareId)
					if err != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := err.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
					if hardware.IsSystemAdd != "Yes" {
						continue
					}

					_, errDelete := repo.DeleteHardwareById(hardware.ID)
					if errDelete != nil {
						errRollback := repo.RollbackHardwareFromBackupTable(fix)
						message := errDelete.Error()
						if errRollback != nil {
							message += errRollback.Error()
						}

						w.WriteJSON(map[string]interface{}{"Status": "error", "Message": message})
						return
					}
				}
			}
		}
		/*
			errDrop := repo.DropHardwareBackupTable()
			if errDrop != nil {
				w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errDrop.Error()})
				return
			}
		*/

		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": response.Content})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作失败！您可以通过线下方式分享给我们，谢谢！"})
	}
}
