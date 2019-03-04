package mysqlrepo

import (
	"fmt"
	"idcos.io/osinstall/model"
)

func (repo *MySQLRepo) GetTaskResultPage(limit uint, offset uint, where string) (results []model.TaskResult, err error) {
	sql := "SELECT * FROM task_result task_result " + where + " order by task_result.id DESC"
	if offset > 0 {
		sql += " limit " + fmt.Sprintf("%d", limit) + "," + fmt.Sprintf("%d", offset)
	} else {
		sql += " limit " + fmt.Sprintf("%d", limit)
	}
	err = repo.db.Raw(sql).Scan(&results).Error
	return results, err
}
func (repo *MySQLRepo) CountTaskResult(where string) (count int, err error) {
	row := repo.db.DB().QueryRow("SELECT count(1) FROM task_result task_result " + where)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MySQLRepo) GetTaskResultByTaskNo(taskNo string) (results []*model.TaskResult, err error) {
	err = repo.db.Model(model.TaskResult{}).Where("task_no = ?", taskNo).Find(&results).Error
	return
}
func (repo *MySQLRepo) GetTaskResultByTaskID(taskID uint) (results []*model.TaskResult, err error) {
	err = repo.db.Model(model.TaskResult{}).Where("task_id = ?", taskID).Find(&results).Error
	return
}

func (repo *MySQLRepo) AddTasks(info *model.TaskInfo, results []*model.TaskResult) (err error) {
	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Save(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, result := range results {
		if err := tx.Save(result).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
