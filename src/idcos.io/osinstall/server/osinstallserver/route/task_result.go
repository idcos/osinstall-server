package route

import (
	"context"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"idcos.com/cloudboot/src/idcos.io/cloudboot/utils"
	"idcos.io/osinstall/middleware"
	"idcos.io/osinstall/model"
	"time"
)

func ReceiveCallback(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var req model.JobCallbackParam
	if err := r.DecodeJSONPayload(&req); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	logger.Debugf("[ReceiveCallback] receive callback info, %s", utils.ToJsonString(req))

	taskInfo, err := repo.GetTaskInfoByNo(req.ExecuteID)
	if err != nil || taskInfo == nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "查询结果错误" + err.Error()})
		return
	}
	taskResults, err := repo.GetTaskResultByTaskNo(req.ExecuteID)
	if err != nil || len(taskResults) == 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "查询结果错误" + err.Error()})
		return
	}

	taskInfo.TaskStatus = req.ExecuteStatus

	for _, result := range taskResults {
		for _, host := range req.HostResults {
			if result.IP == host.HostIP {
				result.EndTIme = time.Now()
				result.Status = host.Status

				if host.Stdout != "" {
					result.Content = host.Stdout
				}

				if host.Stderr != "" {
					result.Content = host.Stderr
				}
			}
		}
	}

	if err := repo.AddTasks(taskInfo, taskResults); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "保存结果信息异常" + err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": ""})
}
