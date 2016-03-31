package mysqlrepo

import (
	"fmt"
	"model"
)

func (repo *MySQLRepo) UpdateHistoryDeviceStatusById(id uint, status string) (*model.DeviceHistory, error) {
	mod := model.DeviceHistory{Status: status}
	err := repo.db.First(&mod, id).Update("status", status).Error
	return &mod, err
}

func (repo *MySQLRepo) CopyDeviceToHistory(id uint) error {
	//var result model.DeviceHistory
	sqlDelete := "delete from device_histories where id = " + fmt.Sprintf("%d", id)
	errDelete := repo.db.Exec(sqlDelete).Error
	if errDelete != nil {
		return errDelete
	}

	sql := "insert into device_histories select * from devices where id = " + fmt.Sprintf("%d", id)
	err := repo.db.Exec(sql).Error
	return err
}
