package mysqlrepo

import (
	"encoding/json"
	"fmt"
	"model"
)

func (repo *MySQLRepo) AssignManufacturerOnwer(id uint, userId uint) (*model.Manufacturer, error) {
	mod := model.Manufacturer{UserID: userId}
	err := repo.db.First(&mod, id).Update("user_id", userId).Error
	return &mod, err
}

func (repo *MySQLRepo) AssignManufacturerNewOnwer(newOwnerID uint, oldOwnerID uint) error {
	sql := "UPDATE manufacturers SET `user_id` = " + fmt.Sprintf("%d", newOwnerID) + " where `user_id` = " + fmt.Sprintf("%d", oldOwnerID)
	err := repo.db.Exec(sql).Error
	return err
}

func (repo *MySQLRepo) GetManufacturerCompanyByGroup(where string) ([]model.Manufacturer, error) {
	var result []model.Manufacturer
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select company from manufacturers " + condition + " group by company order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetManufacturerProductByGroup(where string) ([]model.Manufacturer, error) {
	var result []model.Manufacturer
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select product from manufacturers " + condition + " group by product order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetManufacturerModelNameByGroup(where string) ([]model.Manufacturer, error) {
	var result []model.Manufacturer
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select model_name from manufacturers " + condition + " group by model_name order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetManufacturerListWithPage(limit uint, offset uint, where string) ([]model.ManufacturerFull, error) {
	var result []model.ManufacturerFull
	sql := "SELECT t1.*,t2.username as owner_name from manufacturers t1 left join users t2 on t1.user_id = t2.id where t1.id > 0 " + where + " order by t1.id DESC"

	if offset > 0 {
		sql += " limit " + fmt.Sprintf("%d", offset) + "," + fmt.Sprintf("%d", limit)
	} else {
		sql += " limit " + fmt.Sprintf("%d", limit)
	}

	err := repo.db.Raw(sql).Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) CountManufacturerByWhere(where string) (int, error) {
	row := repo.db.DB().QueryRow("SELECT count(t1.id) as c from manufacturers t1 left join users t2 on t1.user_id = t2.id where t1.id > 0 " + where)
	var count = -1
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) AddManufacturer(deviceId uint, company string, product string, modelName string, sn string, ip string, mac string, nic string, cpu string, cpuSum uint, memory string, memorySum uint, disk string, diskSum uint, motherboard string, raid string, oob string, isVm string, nicDevice string, isShowInScanList string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{DeviceID: deviceId, Company: company, Product: product, ModelName: modelName, Sn: sn, Ip: ip, Mac: mac, Nic: nic, Cpu: cpu, CpuSum: cpuSum, Memory: memory, MemorySum: memorySum, Disk: disk, DiskSum: diskSum, Motherboard: motherboard, Raid: raid, Oob: oob, IsVm: isVm, NicDevice: nicDevice, IsShowInScanList: isShowInScanList}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteManufacturerById(id uint) (*model.Manufacturer, error) {
	mod := model.Manufacturer{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteManufacturerBySn(sn string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{}
	err := repo.db.Unscoped().Where("sn = ?", sn).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountManufacturerByDeviceID(deviceId uint) (uint, error) {
	mod := model.Manufacturer{DeviceID: deviceId}
	var count uint
	err := repo.db.Model(mod).Where("device_id = ?", deviceId).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetManufacturerById(id uint) (*model.Manufacturer, error) {
	var mod model.Manufacturer
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetManufacturerBySn(sn string) (*model.Manufacturer, error) {
	var mod model.Manufacturer
	err := repo.db.Where("sn = ?", sn).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetManufacturerByDeviceId(deviceId uint) (*model.Manufacturer, error) {
	var mod model.Manufacturer
	err := repo.db.Where("device_id = ?", deviceId).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetManufacturerByDeviceID(deviceId uint) (*model.Manufacturer, error) {
	var mod model.Manufacturer
	err := repo.db.Where("device_id = ?", deviceId).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateManufacturerById(id uint, company string, product string, modelName string, sn string, ip string, mac string, nic string, cpu string, cpuSum uint, memory string, memorySum uint, disk string, diskSum uint, motherboard string, raid string, oob string, isVm string, nicDevice string, isShowInScanList string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{Company: company, Product: product, ModelName: modelName}
	err := repo.db.First(&mod, id).Update("company", company).Update("product", product).Update("model_name", modelName).Update("sn", sn).Update("ip", ip).Update("mac", mac).Update("nic", nic).Update("cpu", cpu).Update("cpu_sum", cpuSum).Update("memory", memory).Update("memory_sum", memorySum).Update("disk", disk).Update("disk_sum", diskSum).Update("motherboard", motherboard).Update("raid", raid).Update("oob", oob).Update("is_vm", isVm).Update("nic_device", nicDevice).Update("is_show_in_scan_list", isShowInScanList).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateManufacturerIsShowInScanListById(id uint, isShowInScanList string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{IsShowInScanList: isShowInScanList}
	err := repo.db.First(&mod, id).Update("is_show_in_scan_list", isShowInScanList).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateManufacturerIPById(id uint, ip string) (*model.Manufacturer, error) {
	mod := model.Manufacturer{Ip: ip}
	err := repo.db.First(&mod, id).Update("ip", ip).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateManufacturerDeviceIdById(id uint, deviceId uint) (*model.Manufacturer, error) {
	mod := model.Manufacturer{DeviceID: deviceId}
	err := repo.db.First(&mod, id).Update("device_id", deviceId).Update("is_show_in_scan_list", "No").Error
	return &mod, err
}

func (repo *MySQLRepo) CountManufacturerBySn(sn string) (uint, error) {
	mod := model.Manufacturer{Sn: sn}
	var count uint
	err := repo.db.Unscoped().Model(mod).Where("sn = ?", sn).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetManufacturerIdBySn(sn string) (uint, error) {
	mod := model.Manufacturer{Sn: sn}
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	return mod.ID, err
}

func (repo *MySQLRepo) GetManufacturerSnByNicMacForVm(mac string) (string, error) {
	var result model.Manufacturer
	sql := "SELECT * from manufacturers where is_vm = 'Yes' and (sn = '" + mac + "' or nic like '%%\"Mac\":\"" + mac + "\"%%')"
	err := repo.db.Raw(sql).Scan(&result).Error
	return result.Sn, err

}

func (repo *MySQLRepo) GetManufacturerMacBySn(sn string) (string, error) {
	mod := model.Manufacturer{Sn: sn}
	err := repo.db.Unscoped().Where("sn = ?", sn).Find(&mod).Error
	if err != nil {
		return "", err
	}
	var mac string
	if mod.Nic != "" {
		type Nic struct {
			Name string `json:"Name"`
			Mac  string `json:"Mac"`
			Ip   string `json:"Ip"`
		}
		var nics []Nic

		err := json.Unmarshal([]byte(mod.Nic), &nics)
		if err != nil {
			return "", err
		}
		for _, nic := range nics {
			if nic.Ip != "" {
				mac = nic.Mac
				break
			}
		}
	}

	return mac, nil
}
