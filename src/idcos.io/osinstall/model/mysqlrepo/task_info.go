package mysqlrepo

import (
	"fmt"
	"idcos.io/osinstall/model"
	"strings"
	"time"
)

func (repo *MySQLRepo) GetTaskInfoByNo(taskNo string) (res []model.TaskInfo, err error) {
	err = repo.db.Model(model.TaskInfo{}).Where("task_no = ?", taskNo).Find(&res).Error
	return
}

func (repo *MySQLRepo) AddTaskInfo(info *model.TaskInfo) (res *model.TaskInfo, err error) {
	err = repo.db.Create(info).Error
	return info, err
}
func (repo *MySQLRepo) DeleteTaskInfo(id uint) (err error) {
	return repo.db.Model(model.TaskInfo{}).Where("id =?", id).Update("is_active", "0").Error
}
func (repo *MySQLRepo) GetTaskInfoPage(limit uint, offset uint, where string) (result []model.TaskInfo, err error) {
	sql := "SELECT * FROM task_info task_info " + where + " order by task_info.id DESC"
	if offset > 0 {
		sql += " limit " + fmt.Sprintf("%d", limit) + "," + fmt.Sprintf("%d", offset)
	} else {
		sql += " limit " + fmt.Sprintf("%d", limit)
	}
	err = repo.db.Raw(sql).Scan(&result).Error
	return result, err
}
func (repo *MySQLRepo) CountTaskInfo(where string) (count int, err error) {
	row := repo.db.DB().QueryRow("SELECT count(1) FROM task_info task_info " + where)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) AddTaskInfoAndResult(info *model.TaskInfo, sns []string, password string) (param *model.ConfJobIPExecParam, err error) {

	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(info).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	devices, err := repo.GetDeviceBySns(sns)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var hosts []model.ExecHost
	for _, device := range devices {
		result := model.TaskResult{
			TaskID:    info.ID,
			TaskNo:    info.TaskNo,
			SN:        device.Sn,
			HostName:  device.Hostname,
			IP:        device.Ip,
			StartTime: time.Now(),
			EndTIme:   time.Now(),
			TotalTime: 0,
			Status:    model.Init,
			Content:   "",
		}
		err = tx.Create(&result).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		os, _ := repo.GetOsConfigById(device.OsID)

		hosts = append(hosts, model.ExecHost{
			HostIP:   device.Ip,
			HostPort: 22,
			EntityID: "",
			HostID:   "",
			IdcName:  model.DefaultIDC,
			OsType:   osConvert(os.Name),
			Encoding: model.DefaultEncoding,
			ProxyID:  "",
		})
	}

	param = &model.ConfJobIPExecParam{
		ExecuteID:   string(info.TaskNo),
		ExecHosts:   hosts,
		Provider:    info.TaskChannel,
		Callback:    "",
		JobRecordID: "",
		ExecParam: model.ExecParam{
			Pattern:        info.TaskType,
			RunAs:          info.RunAs,
			Password:       password,
			Timeout:        int(info.Timeout),
			Env:            nil,
			ExtendData:     nil,
			RealTimeOutput: false,
		},
	}

	return param, tx.Commit().Error
}

func osConvert(osName string) string {
	if strings.Contains(osName, "win") {
		return model.Win
	}

	if strings.Contains(osName, "centos") || strings.Contains(osName, "rhel") || strings.Contains(osName, "ubuntu") {
		return model.Linux
	}
	return model.Aix
}
