package mysqlrepo

import (
	"fmt"
)

func (repo *MySQLRepo) CopyDeviceToInstallReport(id uint) error {
	sql := "insert into device_install_reports(`created_at`,`updated_at`,`sn`,`os_name`,`hardware_name`,`system_name`,`status`,`user_id`) select t1.`created_at`,t1.`updated_at`,t1.`sn`,t3.`name` as os_name,concat(t2.`company`,'-',t2.`model_name`) as hardware_name,t4.`name` as system_name,t1.`status`,t1.`user_id` from `devices` t1 left join `hardwares` t2 on t1.`hardware_id` = t2.`id` left join `os_configs` t3 on t1.`os_id` = t3.`id` left join `system_configs` t4 on t1.`system_id` = t4.`id` where t1.`id` = " + fmt.Sprintf("%d", id)
	err := repo.db.Exec(sql).Error
	return err
}
