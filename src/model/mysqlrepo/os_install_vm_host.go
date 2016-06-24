package mysqlrepo

import (
	"fmt"
	"model"
)

func (repo *MySQLRepo) AddVmHost(sn string, cpuSum uint, cpuUsed uint, cpuAvailable uint, memorySum uint, memoryUsed uint, memoryAvailable uint, diskSum uint, diskUsed uint, diskAvailable uint, isAvailable string, remark string, vmNum uint) (*model.VmHost, error) {
	mod := model.VmHost{Sn: sn, CpuSum: cpuSum, CpuUsed: cpuUsed, CpuAvailable: cpuAvailable, MemorySum: memorySum, MemoryUsed: memoryUsed, MemoryAvailable: memoryAvailable, DiskSum: diskSum, DiskUsed: diskUsed, DiskAvailable: diskAvailable, IsAvailable: isAvailable, Remark: remark, VmNum: vmNum}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateVmHostById(id uint, cpuSum uint, cpuUsed uint, cpuAvailable uint, memorySum uint, memoryUsed uint, memoryAvailable uint, diskSum uint, diskUsed uint, diskAvailable uint, isAvailable string, remark string, vmNum uint) (*model.VmHost, error) {
	mod := model.VmHost{CpuSum: cpuSum, CpuUsed: cpuUsed, CpuAvailable: cpuAvailable, MemorySum: memorySum, MemoryUsed: memoryUsed, MemoryAvailable: memoryAvailable, DiskSum: diskSum, DiskUsed: diskUsed, DiskAvailable: diskAvailable, IsAvailable: isAvailable, Remark: remark, VmNum: vmNum}
	err := repo.db.Unscoped().First(&mod, id).Update("cpu_sum", cpuSum).Update("cpu_used", cpuUsed).Update("cpu_available", cpuAvailable).Update("memory_sum", memorySum).Update("memory_used", memoryUsed).Update("memory_available", memoryAvailable).Update("disk_sum", diskSum).Update("disk_used", diskUsed).Update("disk_available", diskAvailable).Update("is_available", isAvailable).Update("remark", remark).Update("vm_num", vmNum).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateVmHostCpuMemoryDiskVmNumById(id uint, cpuSum uint, cpuUsed uint, cpuAvailable uint, memorySum uint, memoryUsed uint, memoryAvailable uint, diskSum uint, diskUsed uint, diskAvailable uint, vmNum uint, isAvailable string) (*model.VmHost, error) {
	mod := model.VmHost{CpuSum: cpuSum, CpuUsed: cpuUsed, CpuAvailable: cpuAvailable, MemorySum: memorySum, MemoryUsed: memoryUsed, MemoryAvailable: memoryAvailable, DiskSum: diskSum, DiskUsed: diskUsed, DiskAvailable: diskAvailable, VmNum: vmNum, IsAvailable: isAvailable}
	err := repo.db.Unscoped().First(&mod, id).Update("cpu_sum", cpuSum).Update("cpu_used", cpuUsed).Update("cpu_available", cpuAvailable).Update("memory_sum", memorySum).Update("memory_used", memoryUsed).Update("memory_available", memoryAvailable).Update("disk_sum", diskSum).Update("disk_used", diskUsed).Update("disk_available", diskAvailable).Update("vm_num", vmNum).Update("is_available", isAvailable).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteVmHostById(id uint) (*model.VmHost, error) {
	mod := model.VmHost{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteVmHostBySn(sn string) (*model.VmHost, error) {
	mod := model.VmHost{}
	err := repo.db.Unscoped().Where("sn = ?", sn).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountVmHostBySn(sn string) (uint, error) {
	mod := model.VmHost{Sn: sn}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("sn = ?", sn).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountVmHost(where string) (int, error) {
	row := repo.db.DB().QueryRow("SELECT count(t7.id) as count FROM devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join hardwares t4 on t1.hardware_id = t4.id left join system_configs t5 on t1.system_id = t5.id inner join vm_hosts t7 on t1.sn = t7.sn  " + where)
	var count = -1
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) GetVmHostListWithPage(limit uint, offset uint, where string) ([]model.VmHostFull, error) {
	var result []model.VmHostFull
	sql := "SELECT t7.*,t1.id as device_id,t1.sn,t1.hostname,t1.ip,t1.manage_ip,t1.network_id,t1.manage_network_id,t1.os_id,t1.hardware_id,t1.system_id,t1.location_id,t1.asset_number,t1.status,t2.network as network_name,t6.network as manage_network_name,t3.name as os_name,concat(t4.company,'-',t4.model_name) as hardware_name,t5.name as system_name FROM devices t1 left join networks t2 on t1.network_id = t2.id left join os_configs t3 on t1.os_id = t3.id left join hardwares t4 on t1.hardware_id = t4.id left join system_configs t5 on t1.system_id = t5.id left join manage_networks t6 on t1.manage_network_id = t6.id inner join vm_hosts t7 on t1.sn = t7.sn  " + where + " order by t1.id DESC"

	if offset > 0 {
		sql += " limit " + fmt.Sprintf("%d", offset) + "," + fmt.Sprintf("%d", limit)
	} else {
		sql += " limit " + fmt.Sprintf("%d", limit)
	}

	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetNeedCollectDeviceForVmHost(deviceId uint) ([]model.Device, error) {
	var where string
	if deviceId > 0 {
		where = " and t1.id = " + fmt.Sprintf("%d", deviceId)
	}
	var result []model.Device
	sql := "select t1.*,case when t2.id is null then 2 when t2.is_available = 'No' then 1 else 0 end as weight from devices t1 left join vm_hosts t2 on t1.sn = t2.sn where t1.`status` = 'success' " + where + " and t1.is_support_vm = 'Yes' order by weight desc"
	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetVmHostById(id uint) (*model.VmHost, error) {
	var mod model.VmHost
	err := repo.db.Unscoped().Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetVmHostBySn(sn string) (*model.VmHost, error) {
	var mod model.VmHost
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetCpuUsedSum(where string) (uint, error) {
	row := repo.db.DB().QueryRow("SELECT SUM(cpu_cores_number) as sum FROM `vm_devices` where " + where)
	var sum = uint(0)
	if err := row.Scan(&sum); err != nil {
		return uint(0), nil
	}
	return sum, nil
}

func (repo *MySQLRepo) GetMemoryUsedSum(where string) (uint, error) {
	row := repo.db.DB().QueryRow("SELECT SUM(memory_current) as sum FROM `vm_devices` where " + where)
	var sum = uint(0)
	if err := row.Scan(&sum); err != nil {
		return uint(0), nil
	}
	return sum, nil
}

func (repo *MySQLRepo) GetDiskUsedSum(where string) (uint, error) {
	row := repo.db.DB().QueryRow("SELECT SUM(disk_size) as sum FROM `vm_devices` where " + where)
	var sum = uint(0)
	if err := row.Scan(&sum); err != nil {
		return uint(0), nil
	}
	return sum, nil
}

func (repo *MySQLRepo) GetMaxVncPort(where string) (uint, error) {
	row := repo.db.DB().QueryRow("SELECT Max(vnc_port) as sum FROM `vm_devices` where " + where)
	var sum = uint(0)
	if err := row.Scan(&sum); err != nil {
		return uint(0), nil
	}
	return sum, nil
}

func (repo *MySQLRepo) DeleteVmInfoByDeviceSn(sn string) error {
	var mod model.Device
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	if err != nil {
		return err
	}

	deviceID := mod.ID
	//delete vm device log
	var modelVmDeviceLog = model.VmDeviceLog{}
	err = repo.db.Unscoped().Where("device_id in (select id from vm_devices where device_id = ?)", deviceID).Delete(&modelVmDeviceLog).Error
	if err != nil {
		return err
	}
	//delete vm device
	var modelVmDevice = model.VmDevice{}
	err = repo.db.Unscoped().Where("device_id = ?", deviceID).Delete(&modelVmDevice).Error
	if err != nil {
		return err
	}
	//delete vm host
	var modelVmHost = model.VmHost{}
	err = repo.db.Unscoped().Where("sn = ?", sn).Delete(&modelVmHost).Error
	if err != nil {
		return err
	}
	return nil

}
