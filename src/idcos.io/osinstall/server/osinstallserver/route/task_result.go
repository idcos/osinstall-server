package route

import (
	"context"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"idcos.io/osinstall/middleware"
	"idcos.io/osinstall/model"
	"idcos.io/osinstall/utils"
	"strconv"
	"strings"
	"time"
)

type TaskResultPageReq struct {
	TaskID      uint   `json:"TaskID"`
	TaskNo      string `json:"TaskNo"`
	AccessToken string `json:"AccessToken"`
	Limit       uint
	Offset      uint
}

func GetTaskResultPage(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var req TaskResultPageReq

	if err := r.DecodeJSONPayload(&req); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mods, err := repo.GetTaskResultPage(req.Limit, req.Offset, getResultsConditions(req))
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountTaskResult(getResultsConditions(req))
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func getResultsConditions(req TaskResultPageReq) string {
	var where []string
	if req.TaskID > 0 {
		where = append(where, fmt.Sprintf("task_result.task_id = %d", req.TaskID))
	}
	if req.TaskNo != "" {
		where = append(where, fmt.Sprintf("task_result.task_no like %s", "'%"+req.TaskNo+"%'"))
	}

	whereStr := strings.Join(where, " and ")

	if len(where) > 0 {
		whereStr = " where " + whereStr
	}

	return whereStr
}

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

	taskID, _ := strconv.Atoi(req.ExecuteID)

	taskInfo, err := repo.GetTaskInfoByID(uint(taskID))
	if err != nil || len(taskInfo) == 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "查询结果错误" + err.Error()})
		return
	}
	taskResults, err := repo.GetTaskResultByTaskID(uint(taskID))
	if err != nil || len(taskResults) == 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "查询结果错误" + err.Error()})
		return
	}

	task := &taskInfo[0]

	task.TaskStatus = req.ExecuteStatus

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

				if host.Status == "fail" {
					result.Content = host.Message
				}

			}
		}
	}

	if err := repo.AddTasks(task, taskResults); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "保存结果信息异常" + err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": ""})
}
