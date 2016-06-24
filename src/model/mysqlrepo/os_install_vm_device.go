package mysqlrepo

import (
	"fmt"
	"model"
)

func (repo *MySQLRepo) GetVmDeviceIdByMac(mac string) (uint, error) {
	mod := model.VmDevice{Mac: mac}
	err := repo.db.Unscoped().Where("mac = ?", mac).Find(&mod).Error
	return mod.ID, err
}

func (repo *MySQLRepo) GetVmDeviceById(id uint) (*model.VmDevice, error) {
	var mod model.VmDevice
	err := repo.db.Unscoped().Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetVmDeviceByMac(mac string) (*model.VmDevice, error) {
	var mod model.VmDevice
	err := repo.db.Unscoped().Where("mac = ?", mac).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetSystemByVmMac(sn string) (*model.SystemConfig, error) {
	var mod model.SystemConfig
	err := repo.db.Joins("inner join vm_devices on vm_devices.system_id = system_configs.id").Where("vm_devices.mac = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetNetworkByVmMac(sn string) (*model.Network, error) {
	var mod model.Network
	err := repo.db.Joins("inner join vm_devices on vm_devices.network_id = networks.id").Where("vm_devices.mac = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetFullVmDeviceById(id uint) (*model.VmDeviceFull, error) {
	var result model.VmDeviceFull
	err := repo.db.Raw("SELECT t1.*,t2.network as network_name,t3.name as os_name,t5.name as system_name,t4.sn as device_sn FROM vm_devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join devices t4 on t1.device_id = t4.id left join system_configs t5 on t1.system_id = t5.id where t1.id = ?", id).Scan(&result).Error
	return &result, err
}

func (repo *MySQLRepo) UpdateVmInstallInfoById(id uint, status string, installProgress float64) (*model.VmDevice, error) {
	mod := model.VmDevice{Status: status, InstallProgress: installProgress}
	err := repo.db.Unscoped().First(&mod, id).Update("status", status).Update("install_progress", installProgress).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateVmRunStatusById(id uint, runStatus string) (*model.VmDevice, error) {
	mod := model.VmDevice{RunStatus: runStatus}
	err := repo.db.Unscoped().First(&mod, id).Update("run_status", runStatus).Error
	return &mod, err
}

func (repo *MySQLRepo) CountVmDeviceByHostname(hostname string) (uint, error) {
	mod := model.VmDevice{Hostname: hostname}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("hostname = ?", hostname).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceByHostnameAndId(hostname string, id uint) (uint, error) {
	mod := model.VmDevice{Hostname: hostname}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("hostname = ? and id != ?", hostname, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) ReInstallVmDeviceById(id uint) (*model.VmDevice, error) {
	mod := model.VmDevice{}
	err := repo.db.First(&mod, id).Update("status", "pre_install").Error
	return &mod, err
}

func (repo *MySQLRepo) CountVmDeviceByMac(mac string) (uint, error) {
	mod := model.VmDevice{Mac: mac}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("mac = ?", mac).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceByDeviceId(deviceId uint) (uint, error) {
	mod := model.VmDevice{DeviceID: deviceId}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("device_id = ?", deviceId).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceByMacAndId(mac string, id uint) (uint, error) {
	mod := model.VmDevice{Mac: mac}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("mac = ? and id != ?", mac, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceByIp(ip string) (uint, error) {
	mod := model.VmDevice{Ip: ip}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("ip = ?", ip).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDeviceByIpAndId(ip string, id uint) (uint, error) {
	mod := model.VmDevice{Ip: ip}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("ip = ? and id != ?", ip, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmDevice(where string) (int, error) {
	row := repo.db.DB().QueryRow("SELECT count(t1.id) as count FROM vm_devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join devices t4 on t1.device_id = t4.id " + where)
	var count = -1
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) GetVmDeviceListWithPage(limit uint, offset uint, where string) ([]model.VmDeviceFull, error) {
	var result []model.VmDeviceFull
	sql := "SELECT t1.*,t2.network as network_name,t3.name as os_name,t5.name as system_name,t4.sn as device_sn FROM vm_devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join devices t4 on t1.device_id = t4.id left join system_configs t5 on t1.system_id = t5.id " + where + " order by t1.id DESC"
	if offset > 0 {
		sql += " limit " + fmt.Sprintf("%d", offset) + "," + fmt.Sprintf("%d", limit)
	} else {
		sql += " limit " + fmt.Sprintf("%d", limit)
	}

	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) DeleteVmDeviceById(id uint) (*model.VmDevice, error) {
	mod := model.VmDevice{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) AddVmDevice(
	DeviceID uint,
	Hostname string,
	Mac string,
	Ip string,
	NetworkID uint,
	OsID uint,
	SystemID uint,
	CpuCoresNumber uint,
	CpuHotPlug string,
	CpuPassthrough string,
	CpuTopSockets uint,
	CpuTopCores uint,
	CpuTopThreads uint,
	CpuPinning string,
	MemoryCurrent uint,
	MemoryMax uint,
	MemoryKsm string,
	DiskType string,
	DiskSize uint,
	DiskBusType string,
	DiskCacheMode string,
	DiskIoMode string,
	NetworkType string,
	NetworkDeviceType string,
	DisplayType string,
	DisplayPassword string,
	DisplayUpdatePassword string,
	Status string,
	UserID uint,
	VncPort string,
	RunStatus string) (model.VmDevice, error) {
	mod := model.VmDevice{
		DeviceID:              DeviceID,
		Hostname:              Hostname,
		Mac:                   Mac,
		Ip:                    Ip,
		NetworkID:             NetworkID,
		OsID:                  OsID,
		SystemID:              SystemID,
		CpuCoresNumber:        CpuCoresNumber,
		CpuHotPlug:            CpuHotPlug,
		CpuPassthrough:        CpuPassthrough,
		CpuTopSockets:         CpuTopSockets,
		CpuTopCores:           CpuTopCores,
		CpuTopThreads:         CpuTopThreads,
		CpuPinning:            CpuPinning,
		MemoryCurrent:         MemoryCurrent,
		MemoryMax:             MemoryMax,
		MemoryKsm:             MemoryKsm,
		DiskType:              DiskType,
		DiskSize:              DiskSize,
		DiskBusType:           DiskBusType,
		DiskCacheMode:         DiskCacheMode,
		DiskIoMode:            DiskIoMode,
		NetworkType:           NetworkType,
		NetworkDeviceType:     NetworkDeviceType,
		DisplayType:           DisplayType,
		DisplayPassword:       DisplayPassword,
		DisplayUpdatePassword: DisplayUpdatePassword,
		UserID:                UserID,
		VncPort:               VncPort,
		RunStatus:             RunStatus,
		Status:                Status}
	err := repo.db.Create(&mod).Error
	return mod, err
}

func (repo *MySQLRepo) UpdateVmDeviceById(
	Id uint,
	DeviceID uint,
	Hostname string,
	Mac string,
	Ip string,
	NetworkID uint,
	OsID uint,
	SystemID uint,
	CpuCoresNumber uint,
	CpuHotPlug string,
	CpuPassthrough string,
	CpuTopSockets uint,
	CpuTopCores uint,
	CpuTopThreads uint,
	CpuPinning string,
	MemoryCurrent uint,
	MemoryMax uint,
	MemoryKsm string,
	DiskType string,
	DiskSize uint,
	DiskBusType string,
	DiskCacheMode string,
	DiskIoMode string,
	NetworkType string,
	NetworkDeviceType string,
	DisplayType string,
	DisplayPassword string,
	DisplayUpdatePassword string,
	Status string,
	UserID uint,
	VncPort string,
	RunStatus string) (model.VmDevice, error) {
	mod := model.VmDevice{
		DeviceID:              DeviceID,
		Hostname:              Hostname,
		Mac:                   Mac,
		Ip:                    Ip,
		NetworkID:             NetworkID,
		OsID:                  OsID,
		SystemID:              SystemID,
		CpuCoresNumber:        CpuCoresNumber,
		CpuHotPlug:            CpuHotPlug,
		CpuPassthrough:        CpuPassthrough,
		CpuTopSockets:         CpuTopSockets,
		CpuTopCores:           CpuTopCores,
		CpuTopThreads:         CpuTopThreads,
		CpuPinning:            CpuPinning,
		MemoryCurrent:         MemoryCurrent,
		MemoryMax:             MemoryMax,
		MemoryKsm:             MemoryKsm,
		DiskType:              DiskType,
		DiskSize:              DiskSize,
		DiskBusType:           DiskBusType,
		DiskCacheMode:         DiskCacheMode,
		DiskIoMode:            DiskIoMode,
		NetworkType:           NetworkType,
		NetworkDeviceType:     NetworkDeviceType,
		DisplayType:           DisplayType,
		DisplayPassword:       DisplayPassword,
		DisplayUpdatePassword: DisplayUpdatePassword,
		UserID:                UserID,
		VncPort:               VncPort,
		RunStatus:             RunStatus,
		Status:                Status}
	err := repo.db.Unscoped().First(&mod, Id).Update("device_id", DeviceID).
		Update("hostname", Hostname).
		Update("mac", Mac).
		Update("ip", Ip).
		Update("network_id", NetworkID).
		Update("os_id", OsID).
		Update("cpu_cores_number", CpuCoresNumber).
		Update("cpu_hot_plug", CpuHotPlug).
		Update("cpu_passthrough", CpuPassthrough).
		Update("cpu_top_sockets", CpuTopSockets).
		Update("cpu_top_cores", CpuTopCores).
		Update("cpu_top_threads", CpuTopThreads).
		Update("cpu_pinning", CpuPinning).
		Update("memory_current", MemoryCurrent).
		Update("memory_max", MemoryMax).
		Update("memory_ksm", MemoryKsm).
		Update("disk_type", DiskType).
		Update("disk_size", DiskSize).
		Update("disk_bus_type", DiskBusType).
		Update("disk_cache_mode", DiskCacheMode).
		Update("disk_io_mode", DiskIoMode).
		Update("network_type", NetworkType).
		Update("network_device_type", NetworkDeviceType).
		Update("display_type", DisplayType).
		Update("display_password", DisplayPassword).
		Update("display_update_password", DisplayUpdatePassword).
		Update("UserID", UserID).
		Update("VncPort", VncPort).
		Update("RunStatus", RunStatus).
		Update("status", Status).Error
	return mod, err
}
