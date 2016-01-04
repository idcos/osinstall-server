package mysqlrepo

import (
	"fmt"
	"model"
)

func (repo *MySQLRepo) GetCompanyByGroup() ([]model.Hardware, error) {
	var result []model.Hardware
	err := repo.db.Raw("select company from hardwares group by company order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetProductByWhereAndGroup(where string) ([]model.Hardware, error) {
	var result []model.Hardware
	err := repo.db.Raw("select product from hardwares where " + where + " group by product order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetModelNameByWhereAndGroup(where string) ([]model.Hardware, error) {
	var result []model.Hardware
	err := repo.db.Raw("select model_name from hardwares where " + where + " group by model_name order by count(*) DESC").Scan(&result).Error
	return result, err
}

//Hardware
func (repo *MySQLRepo) AddHardware(company string, product string, modelName string, raid string, oob string, bios string, isSystemAdd string, tpl string, data string) (*model.Hardware, error) {
	fmt.Println(isSystemAdd)
	mod := model.Hardware{Company: company, Product: product, ModelName: modelName, Raid: raid, Oob: oob, Bios: bios, IsSystemAdd: isSystemAdd, Tpl: tpl, Data: data}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateHardwareById(id uint, company string, product string, modelName string, raid string, oob string, bios string, tpl string, data string) (*model.Hardware, error) {
	mod := model.Hardware{Company: company, Product: product, ModelName: modelName, Raid: raid, Oob: oob, Bios: bios}
	err := repo.db.First(&mod, id).Update("company", company).Update("product", product).Update("model_name", modelName).Update("raid", raid).Update("oob", oob).Update("bios", bios).Update("tpl", tpl).Update("data", data).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteHardwareById(id uint) (*model.Hardware, error) {
	mod := model.Hardware{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountHardwareByCompanyAndProductAndName(company string, product string, modelName string) (uint, error) {
	mod := model.Hardware{}
	var count uint
	err := repo.db.Model(mod).Where("company = ? and product = ? and model_name = ?", company, product, modelName).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountHardwarrWithSeparator(name string) (uint, error) {
	mod := model.Hardware{}
	var count uint
	err := repo.db.Model(mod).Where("CONCAT(company,'-',product,'-',model_name) = ?", name).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetHardwareBySeaprator(name string) (*model.Hardware, error) {
	var mod model.Hardware
	err := repo.db.Where("CONCAT(company,'-',product,'-',model_name) = ?", name).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountHardwareByCompanyAndProductAndNameAndId(company string, product string, modelName string, id uint) (uint, error) {
	mod := model.Hardware{}
	var count uint
	err := repo.db.Model(mod).Where("company = ? and product = ? and model_name = ? and id != ?", company, product, modelName, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountHardware(where string) (uint, error) {
	mod := model.Hardware{}
	var count uint
	err := repo.db.Model(mod).Where("id > 0 " + where).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetHardwareListWithPage(limit uint, offset uint, where string) ([]model.Hardware, error) {
	var mods []model.Hardware
	err := repo.db.Limit(limit).Offset(offset).Where("id > 0 " + where).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetHardwareById(id uint) (*model.Hardware, error) {
	var mod model.Hardware
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetHardwareIdByCompanyAndProductAndName(company string, product string, modelName string) (uint, error) {
	mod := model.Hardware{}
	err := repo.db.Where("company = ? and product = ? and model_name = ?", company, product, modelName).Find(&mod).Error
	return mod.ID, err
}
