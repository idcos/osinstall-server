package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddMac(deviceId uint, mac string) (*model.Mac, error) {
	mod := model.Mac{Mac: mac, DeviceID: deviceId}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteMacById(id uint) (*model.Mac, error) {
	mod := model.Mac{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteMacByDeviceId(deviceId uint) (*model.Mac, error) {
	mod := model.Mac{}
	err := repo.db.Unscoped().Where("device_id = ?", deviceId).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountMacByMac(mac string) (uint, error) {
	mod := model.Mac{Mac: mac}
	var count uint
	err := repo.db.Model(mod).Where("mac = ?", mac).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountMacByMacAndDeviceID(mac string, deviceId uint) (uint, error) {
	mod := model.Mac{Mac: mac, DeviceID: deviceId}
	var count uint
	err := repo.db.Model(mod).Where("mac = ? and device_id = ?", mac, deviceId).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetMacListByDeviceID(deviceId uint) ([]model.Mac, error) {
	var mods []model.Mac
	err := repo.db.Where("device_id = ?", deviceId).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetMacById(id uint) (*model.Mac, error) {
	var mod model.Mac
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}
