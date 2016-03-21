package route

import (
	//"encoding/base64"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
	"middleware"
	//"net/http"
	"errors"
	"fmt"
	"model"
	"server/osinstallserver/util"
	"strings"
	"time"
	"utils"
)

var store = sessions.NewCookieStore([]byte("idcos-osinstall"))

func GetSession(w rest.ResponseWriter, r *rest.Request) (model.UserWithToken, error) {
	session, err := store.Get(r.Request, "user-authentication")
	var user model.UserWithToken
	if err != nil {
		return user, err
	}
	if session.Values["ID"] != nil {
		user.ID = session.Values["ID"].(uint)
		user.Username = session.Values["Username"].(string)
		user.Name = session.Values["Name"].(string)
		user.Role = session.Values["Role"].(string)
		user.AccessToken = session.Values["AccessToken"].(string)
	}
	return user, nil
}

func VerifyAccessPurview(token string, ctx context.Context, isVerifyAdministratorRole bool, w rest.ResponseWriter, r *rest.Request) (model.UserWithToken, error) {
	var user model.UserWithToken
	session, errSession := GetSession(w, r)
	if errSession != nil {
		return user, errSession
	}

	if session.ID <= uint(0) {
		accessTokenUser, errAccessToken := VerifyAccessToken(token, ctx, isVerifyAdministratorRole)
		return accessTokenUser, errAccessToken
	}

	if session.Role == "" {
		return user, errors.New("请您先登录!")
	}

	if isVerifyAdministratorRole == true {
		if session.Role != "Administrator" {
			return user, errors.New("权限不足，请使用超级管理员账号登录!")
		} else {
			user.ID = session.ID
			user.Username = session.Username
			user.Name = session.Name
			user.Role = session.Role
			return user, nil
		}
	}
	return user, nil
}

func VerifyAccessToken(token string, ctx context.Context, isVerifyAdministratorRole bool) (model.UserWithToken, error) {
	var user model.UserWithToken
	token = strings.TrimSpace(token)
	if token == "" {
		return user, errors.New("AccessToken 不能为空!")
	}
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return user, errors.New("内部服务器错误")
	}

	count, err := repo.CountUserAccessTokenByToken(token)
	if err != nil {
		return user, err
	}

	if count != 1 {
		return user, errors.New("AccessToken 不正确!")
	}

	userInfo, err := repo.GetUserByAccessToken(token)
	if err != nil {
		return user, err
	}

	if isVerifyAdministratorRole == true {
		if userInfo.Role != "Administrator" {
			return user, errors.New("权限不足，请使用超级管理员账号登录!")
		}
	}

	user.ID = userInfo.ID
	user.Username = userInfo.Username
	user.Name = userInfo.Name
	user.Role = userInfo.Role
	user.AccessToken = userInfo.AccessToken
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

	token := fmt.Sprintf("%d", user.ID) + "_" + time.Now().String()
	accessToken, err := util.EncodePassword(token)
	accessToken = strings.ToUpper(accessToken)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	_, errToken := repo.AddUserAccessToken(user.ID, accessToken)
	if errToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errToken.Error()})
		return
	}

	session.Values["ID"] = user.ID
	session.Values["Username"] = user.Username
	session.Values["Name"] = user.Name
	session.Values["Role"] = user.Role
	session.Values["AccessToken"] = accessToken
	session.Save(r.Request, w)

	type userInfo struct {
		ID          uint
		Username    string
		Name        string
		Role        string
		AccessToken string
	}
	var userinfo userInfo
	userinfo.ID = user.ID
	userinfo.Username = user.Username
	userinfo.Name = user.Name
	userinfo.Role = user.Role
	userinfo.AccessToken = accessToken

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "登录成功", "Content": userinfo})
	return
}

func LoginOut(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		AccessToken string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.AccessToken = strings.TrimSpace(info.AccessToken)

	session, err := store.Get(r.Request, "user-authentication")
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	/*
		sessionUser, err := GetSession(w, r)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
			return
		}
	*/

	_, errToken := repo.DeleteUserAccessTokenByToken(info.AccessToken)
	if errToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errToken.Error()})
		return
	}

	delete(session.Values, "ID")
	delete(session.Values, "Username")
	delete(session.Values, "Name")
	delete(session.Values, "Role")
	delete(session.Values, "AccessToken")
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
		ID          uint
		AccessToken string
	}
	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	_, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, true)
	if errAccessToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
		return
	}

	osConfig, err := repo.DeleteUserById(info.ID)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	errAssign := repo.AssignManufacturerNewOnwer(uint(0), info.ID)
	if errAssign != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAssign.Error()})
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
		AccessToken string
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
	info.AccessToken = strings.TrimSpace(info.AccessToken)

	_, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, true)
	if errAccessToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
		return
	}

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
		AccessToken string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	info.Password = strings.TrimSpace(info.Password)
	info.Name = strings.TrimSpace(info.Name)
	info.PhoneNumber = strings.TrimSpace(info.PhoneNumber)
	info.AccessToken = strings.TrimSpace(info.AccessToken)

	_, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, false)
	if errAccessToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
		return
	}

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
		Limit       uint
		Offset      uint
		Status      string
		AccessToken string
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	session, errSession := GetSession(w, r)
	if errSession != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + errSession.Error()})
		return
	}
	if session.Username == "" {
		_, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, true)
		if errAccessToken != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
			return
		}
	}

	var where = ""
	if info.Status != "" {
		where += " and status = '" + info.Status + "'"
	}

	mods, err := repo.GetUserListWithPage(info.Limit, info.Offset, where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}
	result := make(map[string]interface{})
	result["list"] = mods

	//总条数
	count, err := repo.CountUser(where)
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
		AccessToken string
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
	info.AccessToken = strings.TrimSpace(info.AccessToken)

	_, errAccessToken := VerifyAccessToken(info.AccessToken, ctx, true)
	if errAccessToken != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errAccessToken.Error()})
		return
	}

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
