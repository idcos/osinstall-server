package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddPlatformConfig(name string, content string) (*model.PlatformConfig, error) {
	mod := model.PlatformConfig{Name: name, Content: content}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdatePlatformConfigById(id uint, name string, content string) (*model.PlatformConfig, error) {
	mod := model.PlatformConfig{Name: name, Content: content}
	err := repo.db.First(&mod, id).Update("name", name).Update("content", content).Error
	return &mod, err
}

func (repo *MySQLRepo) DeletePlatformConfigById(id uint) (*model.PlatformConfig, error) {
	mod := model.PlatformConfig{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountPlatformConfigByName(name string) (uint, error) {
	mod := model.PlatformConfig{Name: name}
	var count uint
	err := repo.db.Model(mod).Where("name = ?", name).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountPlatformConfigByNameAndId(name string, id uint) (uint, error) {
	mod := model.PlatformConfig{}
	var count uint
	err := repo.db.Model(mod).Where("name = ? and id != ?", name, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetPlatformConfigIdByName(name string) (uint, error) {
	mod := model.PlatformConfig{Name: name}
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return mod.ID, err
}

func (repo *MySQLRepo) CountPlatformConfig() (uint, error) {
	mod := model.PlatformConfig{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetPlatformConfigListWithPage(limit uint, offset uint) ([]model.PlatformConfig, error) {
	var mods []model.PlatformConfig
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetPlatformConfigById(id uint) (*model.PlatformConfig, error) {
	var mod model.PlatformConfig
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetPlatformConfigByName(name string) (*model.PlatformConfig, error) {
	var mod model.PlatformConfig
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return &mod, err
}
