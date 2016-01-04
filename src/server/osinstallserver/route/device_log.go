package route

import (
	//"encoding/base64"
	//"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	"utils"
	//"net/http"
)

func GetDeviceLogByDeviceIdAndType(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		DeviceID uint
		Type     string
		Order    string
		MaxID    uint
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	if info.Order == "" {
		info.Order = "id DESC"
	}

	mods, err := repo.GetDeviceLogListByDeviceIDAndType(info.DeviceID, info.Type, info.Order, info.MaxID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type log struct {
		ID        uint
		DeviceID  uint
		Title     string
		Type      string
		Content   string
		CreatedAt utils.ISOTime
		UpdatedAt utils.ISOTime
	}

	var result []log
	for _, v := range mods {
		var row log
		row.ID = v.ID
		row.DeviceID = v.DeviceID
		row.Title = v.Title
		row.Type = v.Type
		row.Content = v.Content
		row.CreatedAt = utils.ISOTime(v.CreatedAt)
		row.UpdatedAt = utils.ISOTime(v.UpdatedAt)
		result = append(result, row)
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}
