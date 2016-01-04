package mysqlrepo

import (
	"model"
)

//SystemConfig相关
func (repo *MySQLRepo) AddNetwork(network string, netmask string, gateway string, vlan string, trunk string, bonding string) (*model.Network, error) {
	mod := model.Network{Network: network, Netmask: netmask, Gateway: gateway, Vlan: vlan, Trunk: trunk, Bonding: bonding}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateNetworkById(id uint, network string, netmask string, gateway string, vlan string, trunk string, bonding string) (*model.Network, error) {
	mod := model.Network{Network: network, Netmask: netmask, Gateway: gateway, Vlan: vlan, Trunk: trunk, Bonding: bonding}
	err := repo.db.First(&mod, id).Update("network", network).Update("netmask", netmask).Update("gateway", gateway).Update("vlan", vlan).Update("trunk", trunk).Update("bonding", bonding).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteNetworkById(id uint) (*model.Network, error) {
	mod := model.Network{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountNetworkByNetwork(network string) (uint, error) {
	mod := model.Network{Network: network}
	var count uint
	err := repo.db.Model(mod).Where("network = ?", network).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountNetworkByNetworkAndId(network string, id uint) (uint, error) {
	mod := model.Network{}
	var count uint
	err := repo.db.Model(mod).Where("network = ? and id != ?", network, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountNetwork() (uint, error) {
	mod := model.Network{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetNetworkListWithPage(limit uint, offset uint) ([]model.Network, error) {
	var mods []model.Network
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) GetNetworkById(id uint) (*model.Network, error) {
	var mod model.Network
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetNetworkIdByNetwork(network string) (uint, error) {
	mod := model.Network{Network: network}
	err := repo.db.Where("network = ?", network).Find(&mod).Error
	return mod.ID, err
}

func (repo *MySQLRepo) AssignNewIpByNetworkId(network_id uint) (string, error) {
	var result model.Ip
	err := repo.db.Raw("select t1.* from ips t1 left join devices t2 on t1.ip = t2.ip left join vm_devices t3 on t1.ip = t3.ip where t1.network_id = ? and t2.id is null and t3.id is null order by t1.id asc limit 1", network_id).Scan(&result).Error
	if err != nil {
		return "", err
	}
	return result.Ip, nil
}

func (repo *MySQLRepo) GetNotUsedIPListByNetworkId(network_id uint) ([]model.Ip, error) {
	var result []model.Ip
	err := repo.db.Raw("select t1.* from ips t1 left join devices t2 on t1.ip = t2.ip left join vm_devices t3 on t1.ip = t3.ip where t1.network_id = ? and t2.id is null and t3.id is null order by t1.id asc limit 200", network_id).Scan(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}
