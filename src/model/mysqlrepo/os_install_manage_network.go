package mysqlrepo

import (
	"model"
)

//SystemConfig相关
func (repo *MySQLRepo) AddManageNetwork(network string, netmask string, gateway string, vlan string, trunk string, bonding string) (*model.ManageNetwork, error) {
	mod := model.ManageNetwork{Network: network, Netmask: netmask, Gateway: gateway, Vlan: vlan, Trunk: trunk, Bonding: bonding}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateManageNetworkById(id uint, network string, netmask string, gateway string, vlan string, trunk string, bonding string) (*model.ManageNetwork, error) {
	mod := model.ManageNetwork{Network: network, Netmask: netmask, Gateway: gateway, Vlan: vlan, Trunk: trunk, Bonding: bonding}
	err := repo.db.First(&mod, id).Update("network", network).Update("netmask", netmask).Update("gateway", gateway).Update("vlan", vlan).Update("trunk", trunk).Update("bonding", bonding).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteManageNetworkById(id uint) (*model.ManageNetwork, error) {
	mod := model.ManageNetwork{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountManageNetworkByNetwork(network string) (uint, error) {
	mod := model.ManageNetwork{Network: network}
	var count uint
	err := repo.db.Model(mod).Where("network = ?", network).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountManageNetworkByNetworkAndId(network string, id uint) (uint, error) {
	mod := model.ManageNetwork{}
	var count uint
	err := repo.db.Model(mod).Where("network = ? and id != ?", network, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountManageNetwork() (uint, error) {
	mod := model.ManageNetwork{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetManageNetworkListWithPage(limit uint, offset uint) ([]model.ManageNetwork, error) {
	var mods []model.ManageNetwork
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetManageNetworkById(id uint) (*model.ManageNetwork, error) {
	var mod model.ManageNetwork
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetManageNetworkIdByNetwork(network string) (uint, error) {
	mod := model.ManageNetwork{Network: network}
	err := repo.db.Where("network = ?", network).Find(&mod).Error
	return mod.ID, err
}
