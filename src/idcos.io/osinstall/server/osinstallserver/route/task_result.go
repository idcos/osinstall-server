package route

import (
	"context"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"idcos.io/osinstall/model"
)

func ReceiveCallback(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	//repo, ok := middleware.RepoFromContext(ctx)
	//config, _ := middleware.ConfigFromContext(ctx)
	//logger, _ := middleware.LoggerFromContext(ctx)
	//if !ok {
	//	w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
	//	return
	//}

	fmt.Println("callback success")

	var req model.JobCallbackParam
	if err := r.DecodeJSONPayload(&req); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": ""})
}
