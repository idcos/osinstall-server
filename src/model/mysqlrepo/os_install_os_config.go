package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) GetHardwareBySn(sn string) (*model.Hardware, error) {
	var hardware model.Hardware
	err := repo.db.Joins("inner join devices on devices.hardware_id = hardwares.id").Where("devices.sn = ?", sn).Find(&hardware).Error
	return &hardware, err
}

func (repo *MySQLRepo) GetDeviceBySnAndStatus(sn string, status string) (*model.Device, error) {
	var device model.Device
	err := repo.db.Where("sn = ? and status = ?", sn, status).Find(&device).Error
	return &device, err
}

func (repo *MySQLRepo) AddOsConfig(name string, pxe string) (*model.OsConfig, error) {
	osConfig := model.OsConfig{Name: name, Pxe: pxe}
	err := repo.db.Create(&osConfig).Error
	return &osConfig, err
}

func (repo *MySQLRepo) UpdateOsConfigById(id uint, name string, pxe string) (*model.OsConfig, error) {
	osConfig := model.OsConfig{Name: name, Pxe: pxe}
	err := repo.db.First(&osConfig, id).Update("name", name).Update("pxe", pxe).Error
	return &osConfig, err
}

func (repo *MySQLRepo) DeleteOsConfigById(id uint) (*model.OsConfig, error) {
	osConfig := model.OsConfig{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&osConfig).Error
	return &osConfig, err
}

func (repo *MySQLRepo) CountOsConfigByName(name string) (uint, error) {
	osConfig := model.OsConfig{Name: name}
	var count uint
	err := repo.db.Model(osConfig).Where("name = ?", name).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountOsConfigByNameAndId(name string, id uint) (uint, error) {
	osConfig := model.OsConfig{}
	var count uint
	err := repo.db.Model(osConfig).Where("name = ? and id != ?", name, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetOsConfigIdByName(name string) (uint, error) {
	osConfig := model.OsConfig{Name: name}
	err := repo.db.Where("name = ?", name).Find(&osConfig).Error
	return osConfig.ID, err
}

func (repo *MySQLRepo) CountOsConfig() (uint, error) {
	osConfig := model.OsConfig{}
	var count uint
	err := repo.db.Model(osConfig).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetOsConfigListWithPage(limit uint, offset uint) ([]model.OsConfig, error) {
	var osConfigs []model.OsConfig
	err := repo.db.Limit(limit).Offset(offset).Find(&osConfigs).Error
	return osConfigs, err
}

func (repo *MySQLRepo) GetOsConfigById(id uint) (*model.OsConfig, error) {
	var osConfig model.OsConfig
	err := repo.db.Where("id = ?", id).Find(&osConfig).Error
	return &osConfig, err
}

func (repo *MySQLRepo) GetOsConfigByName(name string) (*model.OsConfig, error) {
	var osConfig model.OsConfig
	err := repo.db.Where("name = ?", name).Find(&osConfig).Error
	return &osConfig, err
}
