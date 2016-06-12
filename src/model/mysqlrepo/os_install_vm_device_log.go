package mysqlrepo

import (
	"fmt"
	"model"
)

//DeviceLog相关
func (repo *MySQLRepo) AddVmDeviceLog(deviceID uint, title string, logType string, content string) (*model.VmDeviceLog, error) {
	mod := model.VmDeviceLog{DeviceID: deviceID, Title: title, Type: logType, Content: content}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteVmDeviceLogById(id uint) (*model.VmDeviceLog, error) {
	mod := model.VmDeviceLog{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteVmDeviceLogByDeviceIDAndType(deviceID uint, logType string) (*model.VmDeviceLog, error) {
	mod := model.VmDeviceLog{}
	err := repo.db.Unscoped().Where("device_id = ? and type = ?", deviceID).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteVmDeviceLogByDeviceID(deviceID uint) (*model.VmDeviceLog, error) {
	mod := model.VmDeviceLog{}
	err := repo.db.Unscoped().Where("device_id = ?", deviceID).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateVmDeviceLogTypeByDeviceIdAndType(deviceId uint, logType string, newLogType string) ([]model.VmDeviceLog, error) {
	var result []model.VmDeviceLog
	sql := "update vm_device_logs set type = '" + newLogType + "' where device_id = " + fmt.Sprintf("%d", deviceId) + " and type = '" + logType + "'"
	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) CountVmDeviceLogByDeviceID(deviceID uint) (uint, error) {
	mod := model.VmDeviceLog{DeviceID: deviceID}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("device_id = ?", deviceID).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceLogByDeviceIDAndType(deviceID uint, logType string) (uint, error) {
	mod := model.VmDeviceLog{}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("device_id = ? and type = ?", deviceID, logType).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceLog() (uint, error) {
	mod := model.VmDeviceLog{}
	var count uint
	err := repo.db.Unscoped().Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetVmDeviceLogListByDeviceID(deviceID uint, order string) ([]model.VmDeviceLog, error) {
	var mod []model.VmDeviceLog
	err := repo.db.Unscoped().Where("device_id = ?", deviceID).Order(order).Find(&mod).Error
	return mod, err
}

func (repo *MySQLRepo) GetVmDeviceLogListByDeviceIDAndType(deviceID uint, logType string, order string, maxId uint) ([]model.VmDeviceLog, error) {
	var mod []model.VmDeviceLog
	err := repo.db.Unscoped().Where("id > ? and device_id = ? and type = ?", maxId, deviceID, logType).Order(order).Find(&mod).Error
	return mod, err
}

func (repo *MySQLRepo) GetVmDeviceLogById(id uint) (*model.VmDeviceLog, error) {
	var mod model.VmDeviceLog
	err := repo.db.Unscoped().Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetLastVmDeviceLogByDeviceID(deviceID uint) (model.VmDeviceLog, error) {
	var mod model.VmDeviceLog
	err := repo.db.Unscoped().Limit(1).Where("device_id = ?", deviceID).Order("id DESC").Find(&mod).Error
	return mod, err
}
