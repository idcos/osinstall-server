package route

import (
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	"strconv"
	"strings"
)

func GetVmHostBySn(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Sn string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mod, err := repo.GetVmHostBySn(info.Sn)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetVmHostList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
		IsAvailable    string
		BatchNumber    string
		StartUpdatedAt string
		EndUpdatedAt   string
		UserID         int
		ID             int
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
	if info.ID > 0 {
		where += " and t1.id = " + strconv.Itoa(info.ID)
	}
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
	if info.IsAvailable != "" {
		where += " and t7.is_available = '" + info.IsAvailable + "'"
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

	mods, err := repo.GetVmHostListWithPage(info.Limit, info.Offset, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountVmHost(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func CollectAndUpdateVmHostResource(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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
	UpdateVmHostResource(logger, repo, conf, 0)
	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
