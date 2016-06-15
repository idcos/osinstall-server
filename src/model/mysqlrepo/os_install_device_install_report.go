package mysqlrepo

import (
	"fmt"
	"model"
)

func (repo *MySQLRepo) CopyDeviceToInstallReport(id uint) error {
	sql := "insert into device_install_reports(`created_at`,`updated_at`,`sn`,`os_name`,`hardware_name`,`system_name`,`product_name`,`status`,`user_id`) select t1.`created_at`,t1.`updated_at`,t1.`sn`,t3.`name` as os_name,concat(t2.`company`,'-',t2.`model_name`) as hardware_name,t4.`name` as system_name,concat(t5.company,'-',t5.model_name) as product_name,t1.`status`,t1.`user_id` from `devices` t1 left join `hardwares` t2 on t1.`hardware_id` = t2.`id` left join `os_configs` t3 on t1.`os_id` = t3.`id` left join `system_configs` t4 on t1.`system_id` = t4.`id` left join `manufacturers` t5 on t1.`sn` = t5.`sn` where t1.`id` = " + fmt.Sprintf("%d", id)
	err := repo.db.Exec(sql).Error
	return err
}

func (repo *MySQLRepo) CopyVmDeviceToInstallReport(id uint) error {
	sql := "insert into device_install_reports(`created_at`,`updated_at`,`sn`,`os_name`,`hardware_name`,`system_name`,`product_name`,`status`,`user_id`) select t1.`created_at`,t1.`updated_at`,t1.`mac`,t3.`name` as os_name,'' as hardware_name,t4.`name` as system_name,'' as product_name,t1.`status`,t1.`user_id` from `vm_devices` t1 left join `os_configs` t3 on t1.`os_id` = t3.`id` left join `system_configs` t4 on t1.`system_id` = t4.`id` where t1.`id` = " + fmt.Sprintf("%d", id)
	err := repo.db.Exec(sql).Error
	return err
}

func (repo *MySQLRepo) CountDeviceInstallReportByWhere(where string) (uint, error) {
	mod := model.DeviceInstallReport{}
	var count uint
	err := repo.db.Model(mod).Where(where).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetDeviceHardwareNameInstallReport(where string) ([]model.DeviceHardwareNameInstallReport, error) {
	var result []model.DeviceHardwareNameInstallReport
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select hardware_name,count(*) as count from device_install_reports " + condition + " group by hardware_name order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetDeviceProductNameInstallReport(where string) ([]model.DeviceProductNameInstallReport, error) {
	var result []model.DeviceProductNameInstallReport
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select product_name,count(*) as count from device_install_reports " + condition + " group by product_name order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetDeviceCompanyNameInstallReport(where string) ([]model.DeviceProductNameInstallReport, error) {
	var result []model.DeviceProductNameInstallReport
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select substring_index(product_name,'-',1) as product_name,count(*) as count from device_install_reports " + condition + " group by substring_index(product_name,'-',1) order by count(*) desc").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetDeviceOsNameInstallReport(where string) ([]model.DeviceOsNameInstallReport, error) {
	var result []model.DeviceOsNameInstallReport
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select os_name,count(*) as count from device_install_reports " + condition + " group by os_name order by count(*) DESC").Scan(&result).Error
	return result, err
}

func (repo *MySQLRepo) GetDeviceSystemNameInstallReport(where string) ([]model.DeviceSystemNameInstallReport, error) {
	var result []model.DeviceSystemNameInstallReport
	var condition string
	if where != "" {
		condition = " where " + where
	}
	err := repo.db.Raw("select system_name,count(*) as count from device_install_reports " + condition + " group by system_name order by count(*) DESC").Scan(&result).Error
	return result, err
}
