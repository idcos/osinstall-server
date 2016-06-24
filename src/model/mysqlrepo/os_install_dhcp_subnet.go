package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddDhcpSubnet(startIp string, endIp string, gateway string) (*model.DhcpSubnet, error) {
	mod := model.DhcpSubnet{StartIp: startIp, EndIp: endIp, Gateway: gateway}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateDhcpSubnetById(id uint, startIp string, endIp string, gateway string) (*model.DhcpSubnet, error) {
	mod := model.DhcpSubnet{StartIp: startIp, EndIp: endIp, Gateway: gateway}
	err := repo.db.First(&mod, id).Update("start_ip", startIp).Update("end_ip", endIp).Update("gateway", gateway).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteDhcpSubnetById(id uint) (*model.DhcpSubnet, error) {
	mod := model.DhcpSubnet{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountDhcpSubnet() (uint, error) {
	mod := model.DhcpSubnet{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetDhcpSubnetListWithPage(limit uint, offset uint) ([]model.DhcpSubnet, error) {
	var mods []model.DhcpSubnet
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetDhcpSubnetById(id uint) (*model.DhcpSubnet, error) {
	var mod model.DhcpSubnet
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}
