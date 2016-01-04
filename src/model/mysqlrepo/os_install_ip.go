package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddIp(networkId uint, ip string) (*model.Ip, error) {
	mod := model.Ip{NetworkID: networkId, Ip: ip}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteIpByNetworkId(networkId uint) (*model.Ip, error) {
	mod := model.Ip{}
	err := repo.db.Unscoped().Where("network_id = ?", networkId).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountIpByIp(ip string) (uint, error) {
	mod := model.Ip{Ip: ip}
	var count uint
	err := repo.db.Model(mod).Where("ip = ?", ip).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetIpByIp(ip string) (*model.Ip, error) {
	var mod model.Ip
	err := repo.db.Unscoped().Where("ip = ?", ip).Find(&mod).Error
	return &mod, err
}
