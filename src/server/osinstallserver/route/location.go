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

func DeleteLocationById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	locationIds, err := repo.FormatChildLocationIdById(info.ID, "", ",")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	count, err := repo.CountDeviceByWhere("location_id in (" + locationIds + ")")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	if count > 0 {
		devices, err := repo.GetDeviceByWhere("location_id in (" + locationIds + ")")
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}
		var str string
		for k, device := range devices {
			if k != len(devices)-1 {
				str += device.Sn + ","
			} else {
				str += device.Sn
			}

		}
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "下列设备正在使用此位置，不能删除! <br />SN:(" + str + ")"})
		return
	}

	mod, err := repo.DeleteLocationById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetLocationTreeNameById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	locationName, err := repo.FormatLocationNameById(info.ID, "", ">")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": locationName})
}

func UpdateLocationById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID          uint
		Pid         uint
		Name        string
		AccessToken string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}
	info.Name = strings.TrimSpace(info.Name)

	info.AccessToken = strings.TrimSpace(info.AccessToken)
	_, errVerify := VerifyAccessPurview(info.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	if info.Name == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	count, err := repo.CountLocationByNameAndPidAndId(info.Name, info.Pid, info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该位置名已存在!"})
		return
	}

	mod, err := repo.UpdateLocationById(info.ID, info.Pid, info.Name)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetLocationById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mod, err := repo.GetLocationById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetLocationList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Limit  uint
		Offset uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mods, err := repo.GetLocationListWithPage(info.Limit, info.Offset)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountLocation()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func GetLocationListByPid(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Limit  uint
		Offset uint
		Pid    uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mods, err := repo.GetLocationListByPidWithPage(info.Limit, info.Offset, info.Pid)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountLocationByPid(info.Pid)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func FormatLocationToTreeByPid(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Pid       uint
		SelectPid uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	//var initContent []map[string]interface{}
	mods, err := repo.FormatLocationToTreeByPid(info.Pid, nil, 0, info.SelectPid)
	//mods, err := repo.FormatLocationToTreeByPid(info.Pid)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mods})
}

//添加
func AddLocation(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Pid         uint
		Name        string
		AccessToken string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.Name = strings.TrimSpace(info.Name)

	info.AccessToken = strings.TrimSpace(info.AccessToken)
	_, errVerify := VerifyAccessPurview(info.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	if info.Name == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	count, err := repo.CountLocationByNameAndPid(info.Name, info.Pid)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该位置名已存在!"})
		return
	}

	_, errAdd := repo.AddLocation(info.Pid, info.Name)
	if errAdd != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAdd.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
