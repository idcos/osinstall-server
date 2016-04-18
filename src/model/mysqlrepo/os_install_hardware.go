package mysqlrepo

import (
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
func (repo *MySQLRepo) AddHardware(company string, product string, modelName string, raid string, oob string, bios string, isSystemAdd string, tpl string, data string, source string, version string, status string) (*model.Hardware, error) {
	mod := model.Hardware{Company: company, Product: product, ModelName: modelName, Raid: raid, Oob: oob, Bios: bios, IsSystemAdd: isSystemAdd, Tpl: tpl, Data: data, Source: source, Version: version, Status: status}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateHardwareById(id uint, company string, product string, modelName string, raid string, oob string, bios string, tpl string, data string, source string, version string, status string) (*model.Hardware, error) {
	mod := model.Hardware{Company: company, Product: product, ModelName: modelName, Raid: raid, Oob: oob, Bios: bios}
	err := repo.db.First(&mod, id).Update("company", company).Update("product", product).Update("model_name", modelName).Update("raid", raid).Update("oob", oob).Update("bios", bios).Update("tpl", tpl).Update("data", data).Update("source", source).Update("version", version).Update("status", status).Error
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

func (repo *MySQLRepo) CountHardwareWithSeparator(name string) (uint, error) {
	mod := model.Hardware{}
	var count uint
	err := repo.db.Model(mod).Where("CONCAT(company,'-',model_name) = ?", name).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetHardwareBySeaprator(name string) (*model.Hardware, error) {
	var mod model.Hardware
	err := repo.db.Where("CONCAT(company,'-',model_name) = ?", name).Find(&mod).Error
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

func (repo *MySQLRepo) GetHardwareByWhere(where string) (*model.Hardware, error) {
	mod := model.Hardware{}
	err := repo.db.Where(where).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountHardwareByWhere(where string) (uint, error) {
	mod := model.Hardware{}
	var count uint
	err := repo.db.Model(mod).Where(where).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) ValidateHardwareProductModel(company string, product string, modelName string) (bool, error) {
	sql := "select count(*) from hardwares where (company = '" + company + "' or '" + company + "' like concat(\"%\",company,\"%\")) and model_name = '" + modelName + "'"
	row := repo.db.DB().QueryRow(sql)
	var count = -1
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (repo *MySQLRepo) GetLastestVersionHardware() (model.Hardware, error) {
	var result model.Hardware
	err := repo.db.Raw("select * from `hardwares` where `is_system_add` = 'Yes' order by `version` DESC limit 1").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) CreateHardwareBackupTable(fix string) error {
	sql := "DROP TABLE IF EXISTS `hardwares_back" + fix + "`"
	err := repo.db.Exec(sql).Error
	if err != nil {
		return err
	}

	sqlBack := "create table `hardwares_back" + fix + "` select * from `hardwares`"
	errBack := repo.db.Exec(sqlBack).Error
	return errBack
}

func (repo *MySQLRepo) RollbackHardwareFromBackupTable(fix string) error {
	sqlTruncate := "truncate table `hardwares`"
	errTruncate := repo.db.Exec(sqlTruncate).Error
	if errTruncate != nil {
		return errTruncate
	}

	sql := "insert into `hardwares` select * from `hardwares_back" + fix + "`"
	err := repo.db.Exec(sql).Error
	return err
}

func (repo *MySQLRepo) DropHardwareBackupTable(fix string) error {
	sql := "DROP TABLE IF EXISTS `hardwares_back" + fix + "`"
	err := repo.db.Exec(sql).Error
	return err
}
