package route

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"idcos.io/osinstall/middleware"
	"idcos.io/osinstall/model"
	"idcos.io/osinstall/server/osinstallserver/exec"
	"idcos.io/osinstall/server/osinstallserver/util"
	"idcos.io/osinstall/utils"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Extend struct {
	SrcFile    string `json:"SrcFile"`
	DestFile   string `json:"DestFile"`
	ScriptType string `json:"ScriptType"`
	Script     string `json:"Script"`
}
type TaskInfoReq struct {
	TaskName    string   `json:"TaskName"`
	TaskType    string   `json:"TaskType"`    //file or  script
	TaskChannel string   `json:"TaskChannel"` //ssh or salt
	RunAs       string   `json:"RunAs"`
	Timeout     uint     `json:"Timeout"`
	Extend      Extend   `json:"Extend"`
	SNs         []string `json:"SNs"`
	Password    string   `json:"Password"`
	AccessToken string   `json:"AccessToken"`
}

func AddTaskInfo(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	config, _ := middleware.ConfigFromContext(ctx)
	logger, _ := middleware.LoggerFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var req TaskInfoReq
	if err := r.DecodeJSONPayload(&req); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	req.AccessToken = strings.TrimSpace(req.AccessToken)
	user, errVerify := VerifyAccessPurview(req.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	taskInfo := model.TaskInfo{
		TaskNo:      time.Now().Format("20060102150304000"),
		TaskName:    req.TaskName,
		TaskType:    req.TaskType,
		TaskChannel: req.TaskChannel,
		RunAs:       req.RunAs,
		Timeout:     req.Timeout,
		Extend:      utils.ToJsonString(req.Extend),
		Creator:     user.Username,
		IsActive:    true,
	}
	taskConf, err := repo.AddTaskInfoAndResult(&taskInfo, req.SNs, req.Password)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	taskConf.Callback = config.OsInstall.LocalServer + "/callback"
	taskConf.ExecParam.ScriptType = stConvert(req)
	taskConf.ExecParam.Params = paramConvert(req)
	taskConf.ExecParam.Script = scConvert(req, config.OsInstall.LocalServer)
	taskConf.Callback = config.OsInstall.LocalServer + model.CallbackURL

	go exec.TaskExec(logger, taskConf, config.Act2.URL)

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": ""})
}

func stConvert(req TaskInfoReq) string {
	if req.TaskType == model.File {
		return "url"
	}
	return req.Extend.ScriptType
}

func paramConvert(req TaskInfoReq) map[string]interface{} {
	var param = make(map[string]interface{})

	param["args"] = ""

	if req.TaskType == model.File {
		param["target"] = filepath.Dir(req.Extend.DestFile)
		param["fileName"] = filepath.Base(req.Extend.DestFile)
	}

	return param
}

func scConvert(req TaskInfoReq, url string) string {
	if req.TaskType == model.File {
		v, _ := json.Marshal([]string{fmt.Sprintf("%s/%s", url, strings.TrimLeft(req.Extend.SrcFile, model.Root))})
		return string(v)
	}

	return req.Extend.Script
}

func DeleteTaskInfoByID(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var req struct {
		ID          uint   `json:"id"`
		AccessToken string `json:"AccessToken"`
	}

	if err := r.DecodeJSONPayload(&req); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	req.AccessToken = strings.TrimSpace(req.AccessToken)
	_, errVerify := VerifyAccessPurview(req.AccessToken, ctx, true, w, r)
	if errVerify != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errVerify.Error()})
		return
	}

	if err := repo.DeleteTaskInfo(req.ID); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": ""})
}

type TaskInfoPageReq struct {
	ID          uint   `json:"id"`
	TaskNo      string `json:"TaskNo"`
	AccessToken string `json:"AccessToken"`
	Keyword     string `json:"keyword"`
	Limit       uint
	Offset      uint
}

func GetTaskInfoPage(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var req TaskInfoPageReq

	if err := r.DecodeJSONPayload(&req); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	mods, err := repo.GetTaskInfoPage(req.Limit, req.Offset, getInfoConditions(req))
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountTaskInfo(getInfoConditions(req))
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func UploadScript(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "text/html; charset=utf-8")
	r.ParseForm()
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": fmt.Sprintf("读取文件内容异常，%s" + err.Error())})
		return
	}

	defer file.Close()

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": string(b)})
	return
}

func UploadFile(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	w.Header().Add("Content-type", "text/html; charset=utf-8")
	r.ParseForm()
	file, handle, err := r.FormFile("file")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	rootDir := model.Root
	if !util.FileExist(rootDir) {
		err := os.MkdirAll(rootDir, 0777)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
	}

	filename := path.Join(rootDir, handle.Filename)

	result := make(map[string]interface{})
	result["result"] = filename

	if util.FileExist(filename) {
		os.Remove(filename)
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)

	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	defer f.Close()
	defer file.Close()

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
	return
}

func getInfoConditions(req TaskInfoPageReq) string {
	var where []string
	if req.ID > 0 {
		where = append(where, fmt.Sprintf("task_info.id = %d", req.ID))
	}
	if req.TaskNo != "" {
		where = append(where, fmt.Sprintf("task_info.task_no like %s", "'%"+req.TaskNo+"%'"))
	}

	if req.Keyword != "" {
		where = append(where, fmt.Sprintf("task_info.task_no like %s or task_info.task_name like %s", "'%"+req.Keyword+"%'", "'%"+req.Keyword+"%'"))
	}

	where = append(where, "task_info.is_active= 1")

	return " where " + strings.Join(where, " and ")
}
