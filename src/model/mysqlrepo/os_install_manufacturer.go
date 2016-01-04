package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddManufacturer(deviceId uint, company string, product string, modelName string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{DeviceID: deviceId, Company: company, Product: product, ModelName: modelName}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteManufacturerById(id uint) (*model.Manufacturer, error) {
	mod := model.Manufacturer{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountManufacturerByDeviceID(deviceId uint) (uint, error) {
	mod := model.Manufacturer{DeviceID: deviceId}
	var count uint
	err := repo.db.Model(mod).Where("device_id = ?", deviceId).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetManufacturerById(id uint) (*model.Manufacturer, error) {
	var mod model.Manufacturer
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetManufacturerByDeviceID(deviceId uint) (*model.Manufacturer, error) {
	var mod model.Manufacturer
	err := repo.db.Where("device_id = ?", deviceId).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateManufacturerById(id uint, company string, product string, modelName string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{Company: company, Product: product, ModelName: modelName}
	err := repo.db.First(&mod, id).Update("company", company).Update("product", product).Update("model_name", modelName).Error
	return &mod, err
}
