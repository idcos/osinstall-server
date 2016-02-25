package route

import (
	//"encoding/base64"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
	"middleware"
	//"net/http"
	"model"
	"server/osinstallserver/util"
	"strings"
	"utils"
)

var store = sessions.NewCookieStore([]byte("idcos-osinstall"))

func GetSession(w rest.ResponseWriter, r *rest.Request) (model.User, error) {
	session, err := store.Get(r.Request, "user-authentication")
	var user model.User
	if err != nil {
		return user, err
	}
	if session.Values["ID"] != nil {
		user.ID = session.Values["ID"].(uint)
		user.Username = session.Values["Username"].(string)
		user.Name = session.Values["Name"].(string)
		user.Role = session.Values["Role"].(string)
	}
	return user, nil
}

func Login(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Username string
		Password string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Username = strings.TrimSpace(info.Username)
	info.Password = strings.TrimSpace(info.Password)

	if info.Username == "" || info.Password == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "用户名和密码不能为空"})
		return
	}

	count, err := repo.CountUserByUsername(info.Username)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if count <= 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "用户名或密码错误!"})
		return
	}

	user, err := repo.GetUserByUsername(info.Username)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if user.Status != "Enable" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该账号已被锁定，请联系管理员解封!"})
		return
	}

	encodePassword, err := util.EncodePassword(info.Password)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	if encodePassword != user.Password {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "用户名或密码错误"})
		return
	}

	session, err := store.Get(r.Request, "user-authentication")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	session.Values["ID"] = user.ID
	session.Values["Username"] = user.Username
	session.Values["Name"] = user.Name
	session.Values["Role"] = user.Role
	session.Save(r.Request, w)

	type userInfo struct {
		ID       uint
		Username string
		Name     string
		Role     string
	}
	var userinfo userInfo
	userinfo.ID = user.ID
	userinfo.Username = user.Username
	userinfo.Name = user.Name
	userinfo.Role = user.Role

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "登录成功", "Content": userinfo})
	return
}

func LoginOut(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	_, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	session, err := store.Get(r.Request, "user-authentication")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	delete(session.Values, "ID")
	delete(session.Values, "Username")
	delete(session.Values, "Name")
	delete(session.Values, "Role")
	session.Save(r.Request, w)

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}

func DeleteUserById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	osConfig, err := repo.DeleteUserById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": osConfig})
}

func UpdateUserById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID          uint
		Password    string
		Name        string
		PhoneNumber string
		Permission  string
		Status      string
		Role        string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Password = strings.TrimSpace(info.Password)
	info.Name = strings.TrimSpace(info.Name)
	info.PhoneNumber = strings.TrimSpace(info.PhoneNumber)
	info.Permission = strings.TrimSpace(info.Permission)
	info.Status = strings.TrimSpace(info.Status)
	info.Role = strings.TrimSpace(info.Role)

	if info.Password == "" {
		user, err := repo.GetUserById(info.ID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		info.Password = user.Password
	} else {
		encodePassword, err := util.EncodePassword(info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		info.Password = encodePassword
	}

	mod, err := repo.UpdateUserById(info.ID, info.Password, info.Name, info.PhoneNumber, info.Permission, info.Status, info.Role)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func UpdateMyInfo(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		ID          uint
		Password    string
		Name        string
		PhoneNumber string
		Permission  string
		Status      string
		Role        string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Password = strings.TrimSpace(info.Password)
	info.Name = strings.TrimSpace(info.Name)
	info.PhoneNumber = strings.TrimSpace(info.PhoneNumber)

	user, err := repo.GetUserById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	if info.Password == "" {
		info.Password = user.Password
	} else {
		encodePassword, err := util.EncodePassword(info.Password)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}
		info.Password = encodePassword
	}
	info.Permission = user.Permission
	info.Status = user.Status
	info.Role = user.Role

	mod, err := repo.UpdateUserById(info.ID, info.Password, info.Name, info.PhoneNumber, info.Permission, info.Status, info.Role)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": mod})
}

func GetUserById(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mod, err := repo.GetUserById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type UserWithTime struct {
		ID          uint
		Username    string
		Password    string
		Name        string
		PhoneNumber string
		Permission  string
		Status      string
		Role        string
		CreatedAt   utils.ISOTime
		UpdatedAt   utils.ISOTime
	}

	var user UserWithTime
	user.ID = mod.ID
	user.Username = mod.Username
	user.Name = mod.Name
	user.PhoneNumber = mod.PhoneNumber
	user.Permission = mod.Permission
	user.Status = mod.Status
	user.Role = mod.Role
	user.CreatedAt = utils.ISOTime(mod.CreatedAt)
	user.UpdatedAt = utils.ISOTime(mod.UpdatedAt)

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": user})
}

func GetUserList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
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

	mods, err := repo.GetUserListWithPage(info.Limit, info.Offset, "")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountUser("")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result["recordCount"] = count

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func AddUser(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}
	var info struct {
		Username    string
		Password    string
		Name        string
		PhoneNumber string
		Permission  string
		Status      string
		Role        string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误"})
		return
	}

	info.Username = strings.TrimSpace(info.Username)
	info.Password = strings.TrimSpace(info.Password)
	info.Name = strings.TrimSpace(info.Name)
	info.PhoneNumber = strings.TrimSpace(info.PhoneNumber)
	info.Permission = strings.TrimSpace(info.Permission)
	info.Status = strings.TrimSpace(info.Status)
	info.Role = strings.TrimSpace(info.Role)

	if info.Username == "" || info.Password == "" || info.Status == "" || info.Role == "" {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "请将信息填写完整!"})
		return
	}

	count, err := repo.CountUserByUsername(info.Username)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	if count > 0 {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "该用户名已存在!"})
		return
	}

	encodePassword, err := util.EncodePassword(info.Password)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	info.Password = encodePassword

	_, errAdd := repo.AddUser(info.Username, info.Password, info.Name, info.PhoneNumber, info.Permission, info.Status, info.Role)
	if errAdd != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAdd.Error()})
		return
	}

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
}
