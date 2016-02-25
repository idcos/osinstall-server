package mysqlrepo

import (
	"model"
)

func (repo *MySQLRepo) AddUser(username string, password string, name string, phoneNumber string, permission string, status string, role string) (*model.User, error) {
	mod := model.User{Username: username, Password: password, Name: name, PhoneNumber: phoneNumber, Permission: permission, Status: status, Role: role}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateUserById(id uint, password string, name string, phoneNumber string, permission string, status string, role string) (*model.User, error) {
	mod := model.User{Password: password, Name: name, PhoneNumber: phoneNumber, Permission: permission, Status: status, Role: role}
	err := repo.db.First(&mod, id).Update("password", password).Update("name", name).Update("phoneNumber", phoneNumber).Update("permission", permission).Update("status", status).Update("role", role).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteUserById(id uint) (*model.User, error) {
	mod := model.User{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountUser(where string) (uint, error) {
	mod := model.User{}
	var count uint
	err := repo.db.Model(mod).Where("id > 0 " + where).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountUserById(id uint) (uint, error) {
	mod := model.User{}
	var count uint
	err := repo.db.Model(mod).Where("id = ?", id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountUserByUsername(username string) (uint, error) {
	mod := model.User{}
	var count uint
	err := repo.db.Model(mod).Where("username = ?", username).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetUserListWithPage(limit uint, offset uint, where string) ([]model.User, error) {
	var mods []model.User
	err := repo.db.Limit(limit).Offset(offset).Where("id > 0 " + where).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetUserById(id uint) (*model.User, error) {
	var mod model.User
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetUserByUsername(username string) (*model.User, error) {
	var mod model.User
	err := repo.db.Where("username = ?", username).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetUserByWhere(where string) (*model.User, error) {
	mod := model.User{}
	err := repo.db.Where(where).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountUserByWhere(where string) (uint, error) {
	mod := model.User{}
	var count uint
	err := repo.db.Model(mod).Where(where).Count(&count).Error
	return count, err
}
