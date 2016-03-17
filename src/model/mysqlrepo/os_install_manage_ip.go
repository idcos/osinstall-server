package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddManageIp(networkId uint, ip string) (*model.ManageIp, error) {
	mod := model.ManageIp{NetworkID: networkId, Ip: ip}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteManageIpByNetworkId(networkId uint) (*model.ManageIp, error) {
	mod := model.ManageIp{}
	err := repo.db.Unscoped().Where("network_id = ?", networkId).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountManageIpByIp(ip string) (uint, error) {
	mod := model.ManageIp{Ip: ip}
	var count uint
	err := repo.db.Model(mod).Where("ip = ?", ip).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetManageIpByIp(ip string) (*model.ManageIp, error) {
	var mod model.ManageIp
	err := repo.db.Unscoped().Where("ip = ?", ip).Find(&mod).Error
	return &mod, err
}
