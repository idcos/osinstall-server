package route

import (
	//"encoding/base64"
	"errors"
	"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	//"net/http"
	"model"
	"strings"
)

func GetDeviceInstallCallbackList(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	var info struct {
		DeviceId uint
	}

	if err := r.DecodeJSONPayload(&info); err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "参数错误" + err.Error()})
		return
	}

	where := fmt.Sprintf("device_id = %d", info.DeviceId)
	order := "id asc"
	list, err := repo.GetDeviceInstallCallbackByWhere(where, order)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	var rows []model.DeviceInstallCallback
	for _, callback := range list {
		device, err := repo.GetDeviceById(callback.DeviceID)
		if err != nil {
			w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
			return
		}

		callback.Content = strings.Replace(callback.Content, "<{sn}>", device.Sn, -1)
		callback.Content = strings.Replace(callback.Content, "<{hostname}>", device.Hostname, -1)
		callback.Content = strings.Replace(callback.Content, "<{ip}>", device.Ip, -1)
		callback.Content = strings.Replace(callback.Content, "<{manage_ip}>", device.ManageIp, -1)

		if device.NetworkID > uint(0) {
			network, _ := repo.GetNetworkById(device.NetworkID)
			callback.Content = strings.Replace(callback.Content, "<{gateway}>", network.Gateway, -1)
			callback.Content = strings.Replace(callback.Content, "<{netmask}>", network.Netmask, -1)
		}
		if device.ManageNetworkID > uint(0) {
			manageNetwork, _ := repo.GetManageNetworkById(device.ManageNetworkID)
			callback.Content = strings.Replace(callback.Content, "<{manage_gateway}>", manageNetwork.Gateway, -1)
			callback.Content = strings.Replace(callback.Content, "<{manage_netmask}>", manageNetwork.Netmask, -1)
		}
		callback.Content = strings.Replace(callback.Content, "\n", "<br />", -1)
		callback.RunResult = strings.Replace(callback.RunResult, "\n", "<br />", -1)
		rows = append(rows, callback)
	}

	result := make(map[string]interface{})
	result["list"] = rows

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func SaveDeviceInstallCallback(ctx context.Context, deviceId uint, callbackType string, content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}

	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	countCallback, errCount := repo.CountDeviceInstallCallbackByDeviceIDAndType(deviceId, callbackType)
	if errCount != nil {
		return errCount
	}
	if countCallback > 0 {
		callback, err := repo.GetDeviceInstallCallbackByDeviceIDAndType(deviceId, callbackType)
		if err != nil {
			return err
		}
		_, errUpdateCallback := repo.UpdateDeviceInstallCallbackByID(callback.ID, deviceId, callbackType, content, "", "", "")
		if errUpdateCallback != nil {
			return errUpdateCallback
		}
	} else {
		_, errAddCallback := repo.AddDeviceInstallCallback(deviceId, callbackType, content, "", "", "")
		if errAddCallback != nil {
			return errAddCallback
		}
	}
	return nil
}
