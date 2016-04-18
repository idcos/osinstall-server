package mysqlrepo

import (
	"model"
)

//DeviceLog相关
func (repo *MySQLRepo) AddDeviceInstallCallback(DeviceID uint, CallbackType string, Content string, RunTime string, RunResult string, RunStatus string) (*model.DeviceInstallCallback, error) {
	mod := model.DeviceInstallCallback{DeviceID: DeviceID, CallbackType: CallbackType, Content: Content, RunTime: RunTime, RunResult: RunResult, RunStatus: RunStatus}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteDeviceInstallCallbackByID(id uint) (*model.DeviceInstallCallback, error) {
	mod := model.DeviceInstallCallback{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteDeviceInstallCallbackByDeviceID(deviceID uint) (*model.DeviceInstallCallback, error) {
	mod := model.DeviceInstallCallback{}
	err := repo.db.Unscoped().Where("device_id = ?", deviceID).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateDeviceInstallCallbackByID(Id uint, DeviceID uint, CallbackType string, Content string, RunTime string, RunResult string, RunStatus string) (*model.DeviceInstallCallback, error) {
	mod := model.DeviceInstallCallback{DeviceID: DeviceID, CallbackType: CallbackType, Content: Content, RunTime: RunTime, RunResult: RunResult, RunStatus: RunStatus}
	err := repo.db.Unscoped().First(&mod, Id).Update("device_id", DeviceID).Update("callback_type", CallbackType).Update("content", Content).Update("run_result", RunResult).Update("run_time", RunTime).Update("run_status", RunStatus).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateDeviceInstallCallbackRunInfoByID(Id uint, RunTime string, RunResult string, RunStatus string) (*model.DeviceInstallCallback, error) {
	mod := model.DeviceInstallCallback{RunTime: RunTime, RunResult: RunResult, RunStatus: RunStatus}
	err := repo.db.Unscoped().First(&mod, Id).Update("run_result", RunResult).Update("run_time", RunTime).Update("run_status", RunStatus).Error
	return &mod, err
}

func (repo *MySQLRepo) CountDeviceInstallCallbackByDeviceIDAndType(deviceID uint, callbackType string) (uint, error) {
	mod := model.DeviceInstallCallback{DeviceID: deviceID, CallbackType: callbackType}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("device_id = ? and callback_type = ?", deviceID, callbackType).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetDeviceInstallCallbackByWhere(where string, order string) ([]model.DeviceInstallCallback, error) {
	var mod []model.DeviceInstallCallback
	err := repo.db.Unscoped().Where(where).Order(order).Find(&mod).Error
	return mod, err
}

func (repo *MySQLRepo) GetDeviceInstallCallbackByDeviceIDAndType(deviceID uint, callbackType string) (*model.DeviceInstallCallback, error) {
	var mod model.DeviceInstallCallback
	err := repo.db.Unscoped().Where("device_id = ? and callback_type = ?", deviceID, callbackType).Find(&mod).Error
	return &mod, err
}
