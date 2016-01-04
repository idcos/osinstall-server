package mysqlrepo

import (
	"model"
)

//SystemConfig相关
func (repo *MySQLRepo) AddSystemConfig(name string, content string) (*model.SystemConfig, error) {
	systemConfig := model.SystemConfig{Name: name, Content: content}
	err := repo.db.Create(&systemConfig).Error
	return &systemConfig, err
}

func (repo *MySQLRepo) UpdateSystemConfigById(id uint, name string, content string) (*model.SystemConfig, error) {
	systemConfig := model.SystemConfig{Name: name, Content: content}
	err := repo.db.First(&systemConfig, id).Update("name", name).Update("content", content).Error
	return &systemConfig, err
}

func (repo *MySQLRepo) DeleteSystemConfigById(id uint) (*model.SystemConfig, error) {
	systemConfig := model.SystemConfig{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&systemConfig).Error
	return &systemConfig, err
}

func (repo *MySQLRepo) CountSystemConfigByName(name string) (uint, error) {
	mod := model.SystemConfig{Name: name}
	var count uint
	err := repo.db.Model(mod).Where("name = ?", name).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountSystemConfigByNameAndId(name string, id uint) (uint, error) {
	mod := model.SystemConfig{}
	var count uint
	err := repo.db.Model(mod).Where("name = ? and id != ?", name, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountSystemConfig() (uint, error) {
	mod := model.SystemConfig{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetSystemConfigListWithPage(limit uint, offset uint) ([]model.SystemConfig, error) {
	var mods []model.SystemConfig
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetSystemConfigById(id uint) (*model.SystemConfig, error) {
	var mod model.SystemConfig
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetSystemConfigIdByName(name string) (uint, error) {
	mod := model.SystemConfig{Name: name}
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return mod.ID, err
}
