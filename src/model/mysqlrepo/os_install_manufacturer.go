package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) GetManufacturerCompanyByGroup(where string) ([]model.Manufacturer, error) {
	var result []model.Manufacturer
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select company from manufacturers " + condition + " group by company order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetManufacturerProductByGroup(where string) ([]model.Manufacturer, error) {
	var result []model.Manufacturer
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select product from manufacturers " + condition + " group by product order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetManufacturerModelNameByGroup(where string) ([]model.Manufacturer, error) {
	var result []model.Manufacturer
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select model_name from manufacturers " + condition + " group by model_name order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetManufacturerListWithPage(limit uint, offset uint, where string) ([]model.Manufacturer, error) {
	var mods []model.Manufacturer
	err := repo.db.Unscoped().Where(where).Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) CountManufacturerByWhere(where string) (uint, error) {
	mod := model.Manufacturer{}
	var count uint
	err := repo.db.Model(mod).Where(where).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) AddManufacturer(deviceId uint, company string, product string, modelName string, sn string, ip string, mac string, nic string, cpu string, memory string, disk string, motherboard string, raid string, oob string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{DeviceID: deviceId, Company: company, Product: product, ModelName: modelName, Sn: sn, Ip: ip, Mac: mac, Nic: nic, Cpu: cpu, Memory: memory, Disk: disk, Motherboard: motherboard, Raid: raid, Oob: oob}
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

func (repo *MySQLRepo) UpdateManufacturerById(id uint, company string, product string, modelName string, sn string, ip string, mac string, nic string, cpu string, memory string, disk string, motherboard string, raid string, oob string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{Company: company, Product: product, ModelName: modelName}
	err := repo.db.First(&mod, id).Update("company", company).Update("product", product).Update("model_name", modelName).Update("sn", sn).Update("ip", ip).Update("mac", mac).Update("nic", nic).Update("cpu", cpu).Update("memory", memory).Update("disk", disk).Update("motherboard", motherboard).Update("raid", raid).Update("oob", oob).Error
	return &mod, err
}

func (repo *MySQLRepo) CountManufacturerBySn(sn string) (uint, error) {
	mod := model.Manufacturer{Sn: sn}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("sn = ?", sn).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetManufacturerIdBySn(sn string) (uint, error) {
	mod := model.Manufacturer{Sn: sn}
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	return mod.ID, err
}
