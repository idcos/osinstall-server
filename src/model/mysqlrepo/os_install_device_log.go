package mysqlrepo

import (
	"fmt"
	"model"
)

//DeviceLog相关
func (repo *MySQLRepo) AddDeviceLog(deviceID uint, title string, logType string, content string) (*model.DeviceLog, error) {
	mod := model.DeviceLog{DeviceID: deviceID, Title: title, Type: logType, Content: content}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteDeviceLogById(id uint) (*model.DeviceLog, error) {
	mod := model.DeviceLog{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteDeviceLogByDeviceIDAndType(deviceID uint, logType string) (*model.DeviceLog, error) {
	mod := model.DeviceLog{}
	err := repo.db.Unscoped().Where("device_id = ? and type = ?", deviceID).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteDeviceLogByDeviceID(deviceID uint) (*model.DeviceLog, error) {
	mod := model.DeviceLog{}
	err := repo.db.Unscoped().Where("device_id = ?", deviceID).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateDeviceLogTypeByDeviceIdAndType(deviceId uint, logType string, newLogType string) ([]model.DeviceLog, error) {
	var result []model.DeviceLog
	sql := "update device_logs set type = '" + newLogType + "' where device_id = " + fmt.Sprintf("%d", deviceId) + " and type = '" + logType + "'"
	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) CountDeviceLogByDeviceID(deviceID uint) (uint, error) {
	mod := model.DeviceLog{DeviceID: deviceID}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("device_id = ?", deviceID).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceLogByDeviceIDAndType(deviceID uint, logType string) (uint, error) {
	mod := model.DeviceLog{}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("device_id = ? and type = ?", deviceID, logType).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceLog() (uint, error) {
	mod := model.DeviceLog{}
	var count uint
	err := repo.db.Unscoped().Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetDeviceLogListByDeviceID(deviceID uint, order string) ([]model.DeviceLog, error) {
	var mod []model.DeviceLog
	err := repo.db.Unscoped().Limit(1000).Where("device_id = ?", deviceID).Order(order).Find(&mod).Error
	return mod, err
}

func (repo *MySQLRepo) GetDeviceLogListByDeviceIDAndType(deviceID uint, logType string, order string, maxId uint) ([]model.DeviceLog, error) {
	var mod []model.DeviceLog
	err := repo.db.Unscoped().Limit(1000).Where("id > ? and device_id = ? and type = ?", maxId, deviceID, logType).Order(order).Find(&mod).Error
	return mod, err
}

func (repo *MySQLRepo) GetDeviceLogById(id uint) (*model.DeviceLog, error) {
	var mod model.DeviceLog
	err := repo.db.Unscoped().Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetLastDeviceLogByDeviceID(deviceID uint) (model.DeviceLog, error) {
	var mod model.DeviceLog
	err := repo.db.Unscoped().Limit(1).Where("device_id = ?", deviceID).Order("id DESC").Find(&mod).Error
	return mod, err
}
