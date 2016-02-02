package osinstallserver

import (
	//"encoding/base64"
	//"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	//"golang.org/x/net/context"
	//"middleware"
	//"net/http"
	"server/osinstallserver/route"
)

var routes []*rest.Route

func init() {
	routes = append(routes, rest.Post("/api/osinstall/v1/osConfig/add", route.AddOsConfig))
	routes = append(routes, rest.Post("/api/osinstall/v1/osConfig/list", route.GetOsConfigList))
	routes = append(routes, rest.Post("/api/osinstall/v1/osConfig/view", route.GetOsConfigById))
	routes = append(routes, rest.Post("/api/osinstall/v1/osConfig/update", route.UpdateOsConfigById))
	routes = append(routes, rest.Post("/api/osinstall/v1/osConfig/delete", route.DeleteOsConfigById))
	//SystemConfig
	routes = append(routes, rest.Post("/api/osinstall/v1/systemConfig/add", route.AddSystemConfig))
	routes = append(routes, rest.Post("/api/osinstall/v1/systemConfig/list", route.GetSystemConfigList))
	routes = append(routes, rest.Post("/api/osinstall/v1/systemConfig/view", route.GetSystemConfigById))
	routes = append(routes, rest.Post("/api/osinstall/v1/systemConfig/update", route.UpdateSystemConfigById))
	routes = append(routes, rest.Post("/api/osinstall/v1/systemConfig/delete", route.DeleteSystemConfigById))
	//Location
	routes = append(routes, rest.Post("/api/osinstall/v1/location/add", route.AddLocation))
	routes = append(routes, rest.Post("/api/osinstall/v1/location/list", route.GetLocationListByPid))
	routes = append(routes, rest.Post("/api/osinstall/v1/location/view", route.GetLocationById))
	routes = append(routes, rest.Post("/api/osinstall/v1/location/update", route.UpdateLocationById))
	routes = append(routes, rest.Post("/api/osinstall/v1/location/delete", route.DeleteLocationById))
	routes = append(routes, rest.Post("/api/osinstall/v1/location/tree", route.FormatLocationToTreeByPid))
	routes = append(routes, rest.Post("/api/osinstall/v1/location/getLocationTreeNameById", route.GetLocationTreeNameById))
	//Network
	routes = append(routes, rest.Post("/api/osinstall/v1/network/add", route.AddNetwork))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/list", route.GetNetworkList))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/view", route.GetNetworkById))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/update", route.UpdateNetworkById))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/delete", route.DeleteNetworkById))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/cidr/get", route.GetCidrInfoByNetwork))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/validateIp", route.ValidateIp))
	routes = append(routes, rest.Post("/api/osinstall/v1/network/getNotUsedIPListByNetworkId", route.GetNotUsedIPListByNetworkId))

	//Device
	routes = append(routes, rest.Post("/api/osinstall/v1/device/add", route.AddDevice))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/batchAdd", route.BatchAddDevice))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/batchUpdate", route.BatchUpdateDevice))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/list", route.GetDeviceList))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/view", route.GetDeviceById))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/viewFull", route.GetFullDeviceById))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/delete", route.DeleteDeviceById))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/getNumByStatus", route.GetDeviceNumByStatus))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/batchReInstall", route.BatchReInstall))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/batchDelete", route.BatchDelete))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/validateSn", route.ValidateSn))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/batchCancelInstall", route.BatchCancelInstall))

	//Hardware
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/add", route.AddHardware))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/list", route.GetHardwareList))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/view", route.GetHardwareById))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/update", route.UpdateHardwareById))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/delete", route.DeleteHardwareById))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/getCompanyByGroup", route.GetCompanyByGroup))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/getProductByWhereAndGroup", route.GetProductByWhereAndGroup))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/getModelNameByWhereAndGroup", route.GetModelNameByWhereAndGroup))
	routes = append(routes, rest.Get("/api/osinstall/v1/hardware/export", route.ExportHardware))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/uploadCompanyHardware", route.UploadCompanyHardware))
	routes = append(routes, rest.Post("/api/osinstall/v1/hardware/uploadHardware", route.UploadHardware))
	//DeviceLog
	routes = append(routes, rest.Post("/api/osinstall/v1/deviceLog/list", route.GetDeviceLogByDeviceIdAndType))

	//Agent
	routes = append(routes, rest.Post("/api/osinstall/v1/device/getHardwareBySn", route.GetHardwareBySn))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/isInInstallList", route.IsInPreInstallList))
	routes = append(routes, rest.Post("/api/osinstall/v1/report/deviceInstallInfo", route.ReportInstallInfo))
	routes = append(routes, rest.Post("/api/osinstall/v1/report/deviceMacInfo", route.ReportMacInfo))
	routes = append(routes, rest.Post("/api/osinstall/v1/report/deviceProductInfo", route.ReportProductInfo))
	routes = append(routes, rest.Get("/api/osinstall/v1/device/getSystemBySn", route.GetSystemBySn))
	routes = append(routes, rest.Get("/api/osinstall/v1/device/getNetworkBySn", route.GetNetworkBySn))

	//Import device
	routes = append(routes, rest.Post("/api/osinstall/v1/device/upload", route.UploadDevice))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/importPriview", route.ImportPriview))
	routes = append(routes, rest.Post("/api/osinstall/v1/device/importDevice", route.ImportDevice))

	//VM install
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/createNewMacAddress", route.CreateNewMacAddress))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/batchAdd", route.BatchAddVmDevice))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/list", route.GetVmDeviceList))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/view", route.GetVmDeviceById))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/viewFull", route.GetFullVmDeviceById))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/delete", route.DeleteVmDeviceById))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/validateMac", route.ValidateMac))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/batchReInstallVm", route.BatchReInstallVm))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/batchDeleteVm", route.BatchDeleteVm))
	routes = append(routes, rest.Post("/api/osinstall/v1/vm/create", route.CreateVmDevice))
}
