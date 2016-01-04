package route

import (
	//"encoding/base64"
	//"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	"strings"
	//"net/http"
)

func DeleteHardwareById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
		ID        uint
		Company   string
		Product   string
		ModelName string
		Raid      string
		Oob       string
		Bios      string
		Tpl       string
		Data      string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Company = strings.TrimSpace(info.Company)
	info.Product = strings.TrimSpace(info.Product)
	info.ModelName = strings.TrimSpace(info.ModelName)
	info.Tpl = strings.TrimSpace(info.Tpl)
	info.Data = strings.TrimSpace(info.Data)

	if info.Company == "" || info.Product == "" || info.ModelName == "" {
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

	mod, err := repo.UpdateHardwareById(info.ID, info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.Tpl, info.Data)
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
		Company string
		Product string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	where := " company = '" + info.Company + "' and product = '" + info.Product + "'"
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
	result["list"] = mods

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

	if info.Company == "" || info.Product == "" || info.ModelName == "" {
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

	_, errAdd := repo.AddHardware(info.Company, info.Product, info.ModelName, info.Raid, info.Oob, info.Bios, info.IsSystemAdd, info.Tpl, info.Data)
	if errAdd != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAdd.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
