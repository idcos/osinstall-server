package mysqlrepo

import (
	"fmt"
	"model"
	"server/osinstallserver/util"
	"strconv"
	"strings"
	"time"
)

//device相关
func (repo *MySQLRepo) AddDevice(batchNumber string, sn string, hostname string, ip string, manageIp string, networkId uint, manageNetworkId uint, osId uint, hardwareId uint, systemId uint, location string, locationId uint, assetNumber string, status string, isSupportVm string, userID uint) (*model.Device, error) {
	mod := model.Device{BatchNumber: batchNumber, Sn: sn, Hostname: hostname, Ip: ip, ManageIp: manageIp, NetworkID: networkId, ManageNetworkID: manageNetworkId, OsID: osId, HardwareID: hardwareId, SystemID: systemId, Location: location, LocationID: locationId, AssetNumber: assetNumber, Status: status, IsSupportVm: isSupportVm, UserID: userID}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateDeviceById(id uint, batchNumber string, sn string, hostname string, ip string, manageIp string, networkId uint, manageNetworkId uint, osId uint, hardwareId uint, systemId uint, location string, locationId uint, assetNumber string, status string, isSupportVm string, userID uint) (*model.Device, error) {
	mod := model.Device{BatchNumber: batchNumber, Sn: sn, Hostname: hostname, Ip: ip, ManageIp: manageIp, NetworkID: networkId, ManageNetworkID: manageNetworkId, OsID: osId, HardwareID: hardwareId, SystemID: systemId, Location: location, LocationID: locationId, AssetNumber: assetNumber, Status: status, IsSupportVm: isSupportVm, UserID: userID}
	//设备信息发生修改，但属主不发生变化
	err := repo.db.Unscoped().First(&mod, id).Update("batch_number", batchNumber).Update("sn", sn).Update("hostname", hostname).Update("ip", ip).Update("manage_ip", manageIp).Update("network_id", networkId).Update("manage_network_id", manageNetworkId).Update("os_id", osId).Update("hardware_id", hardwareId).Update("system_id", systemId).Update("location", location).Update("location_id", locationId).Update("asset_number", assetNumber).Update("status", status).Update("install_progress", 0.0000).Update("install_log", "").Update("is_support_vm", isSupportVm).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateInstallInfoById(id uint, status string, installProgress float64) (*model.Device, error) {
	mod := model.Device{Status: status, InstallProgress: installProgress}
	err := repo.db.Unscoped().First(&mod, id).Update("status", status).Update("install_progress", installProgress).Error
	return &mod, err
}

func (repo *MySQLRepo) ReInstallDeviceById(id uint) (*model.Device, error) {
	mod := model.Device{}
	err := repo.db.Unscoped().First(&mod, id).Update("status", "pre_install").Update("install_progress", 0.0000).Update("install_log", "").Error
	return &mod, err
}

func (repo *MySQLRepo) CancelInstallDeviceById(id uint) (*model.Device, error) {
	mod := model.Device{}
	err := repo.db.Unscoped().First(&mod, id).Update("status", "failure").Update("install_progress", 0.0000).Update("install_log", "").Error
	return &mod, err
}

//device相关
func (repo *MySQLRepo) CreateBatchNumber() (string, error) {
	date := time.Now().Format("2006-01-02")
	var batchNumber string
	//row := repo.db.DB().QueryRow("select count(*) as count from (select batch_number from devices where batch_number like ?) as t", date+"%")
	row := repo.db.DB().QueryRow("select count(*) as count from devices where created_at >= ? and created_at <= ?", date+" 00:00:00", date+" 23:59:59")
	var count = -1
	if err := row.Scan(&count); err != nil {
		return "", err
	}

	if count > 0 {
		var device model.Device
		err := repo.db.Unscoped().Where("created_at >= ? and created_at <= ?", date+" 00:00:00", date+" 23:59:59").Limit(1).Order("id DESC").Find(&device).Error
		if err != nil {
			return "", nil
		}
		fix := util.SubString(device.BatchNumber, 8, len(device.BatchNumber)-8)
		fixNum, err := strconv.Atoi(fix)
		if err != nil {
			return "", err
		}
		batchNumber = strings.Replace(date, "-", "", -1) + fmt.Sprintf("%03d", fixNum+1)
	} else {
		batchNumber = strings.Replace(date, "-", "", -1) + fmt.Sprintf("%03d", 1)
	}

	return batchNumber, nil
}

func (repo *MySQLRepo) DeleteDeviceById(id uint) (*model.Device, error) {
	mod := model.Device{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountDeviceBySn(sn string) (uint, error) {
	mod := model.Device{Sn: sn}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("sn = ?", sn).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceByHostname(hostname string) (uint, error) {
	mod := model.Device{Hostname: hostname}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("hostname = ?", hostname).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceByHostnameAndId(hostname string, id uint) (uint, error) {
	mod := model.Device{Hostname: hostname}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("hostname = ? and id != ?", hostname, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceByIp(ip string) (uint, error) {
	mod := model.Device{Ip: ip}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("ip = ?", ip).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceByManageIp(manageIp string) (uint, error) {
	mod := model.Device{ManageIp: manageIp}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("manage_ip = ?", manageIp).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceByIpAndId(ip string, id uint) (uint, error) {
	mod := model.Device{Ip: ip}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("ip = ? and id != ?", ip, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDeviceByManageIpAndId(manageIp string, id uint) (uint, error) {
	mod := model.Device{ManageIp: manageIp}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("manage_ip = ? and id != ?", manageIp, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountDevice(where string) (int, error) {
	row := repo.db.DB().QueryRow("SELECT count(t1.id) as count FROM devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join hardwares t4 on t1.hardware_id = t4.id left join system_configs t5 on t1.system_id = t5.id " + where)
	var count = -1
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) CountDeviceByWhere(where string) (int, error) {
	row := repo.db.DB().QueryRow("SELECT count(*) as count FROM devices where " + where)
	var count = -1
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) GetDeviceByWhere(where string) ([]model.Device, error) {
	var result []model.Device
	sql := "SELECT * FROM devices where " + where
	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetDeviceListWithPage(limit uint, offset uint, where string) ([]model.DeviceFull, error) {
	/*
		var mods []model.Device
		err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
		return mods, err
	*/

	var result []model.DeviceFull
	sql := "SELECT t1.*,t2.network as network_name,t6.network as manage_network_name,t3.name as os_name,concat(t4.company,'-',t4.model_name) as hardware_name,t5.name as system_name,t7.username as owner_name,t8.ip as bootos_ip,t8.oob as oob_ip FROM devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join hardwares t4 on t1.hardware_id = t4.id left join system_configs t5 on t1.system_id = t5.id left join manage_networks t6 on t1.manage_network_id = t6.id left join `users` t7 on t1.user_id = t7.id left join manufacturers t8 on t1.`sn` = t8.`sn` " + where + " order by t1.id DESC"

	if offset > 0 {
		sql += " limit " + fmt.Sprintf("%d", offset) + "," + fmt.Sprintf("%d", limit)
	} else {
		sql += " limit " + fmt.Sprintf("%d", limit)
	}

	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetFullDeviceById(id uint) (*model.DeviceFull, error) {
	var result model.DeviceFull
	err := repo.db.Raw("SELECT t1.*,t2.network as network_name,t6.network as manage_network_name,t3.name as os_name,concat(t4.company,'-',t4.model_name) as hardware_name,t5.name as system_name FROM devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join hardwares t4 on t1.hardware_id = t4.id left join system_configs t5 on t1.system_id = t5.id left join manage_networks t6 on t1.manage_network_id = t6.id where t1.id = ?", id).Scan(&result).Error
	return &result, err
}

func (repo *MySQLRepo) GetDeviceById(id uint) (*model.Device, error) {
	var mod model.Device
	err := repo.db.Unscoped().Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetDeviceBySn(sn string) (*model.Device, error) {
	var mod model.Device
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetDeviceIdBySn(sn string) (uint, error) {
	mod := model.Device{Sn: sn}
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	return mod.ID, err
}

func (repo *MySQLRepo) GetSystemBySn(sn string) (*model.SystemConfig, error) {
	var mod model.SystemConfig
	err := repo.db.Joins("inner join devices on devices.system_id = system_configs.id").Where("devices.sn = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetNetworkBySn(sn string) (*model.Network, error) {
	var mod model.Network
	err := repo.db.Joins("inner join devices on devices.network_id = networks.id").Where("devices.sn = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetInstallTimeoutDeviceList(timeout int) ([]model.Device, error) {
	var result []model.Device
	sql := "select t3.* from (select device_id,max(id) as id from device_logs where type = 'install' and device_id in (select id from devices where status = 'installing') group by device_id ) t1 inner join device_logs t2 on t1.id = t2.id and (UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(t2.created_at)) >= " + fmt.Sprintf("%d", timeout) + " inner join devices t3 on t1.device_id = t3.id"
	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) IsInstallTimeoutDevice(timeout int, deviceId uint) (bool, error) {
	sql := "select count(t3.id) as count from (select device_id,max(id) as id from device_logs where type = 'install' and device_id = " + fmt.Sprintf("%d", deviceId) + " and device_id in (select id from devices where status = 'installing') group by device_id ) t1 inner join device_logs t2 on t1.id = t2.id and (UNIX_TIMESTAMP(now()) - UNIX_TIMESTAMP(t2.created_at)) >= " + fmt.Sprintf("%d", timeout) + " inner join devices t3 on t1.device_id = t3.id"
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

func (repo *MySQLRepo) ExecDBVersionUpdateSql(sql string) error {
	err := repo.db.Exec(sql).Error
	return err
}
