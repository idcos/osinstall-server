package mysqlrepo

import (
	"idcos.io/osinstall/model"
)

func (repo *MySQLRepo) AddTaskResult(info *model.TaskResult) (err error) {
	return
}
func (repo *MySQLRepo) GetTaskResultPage(limit uint, offset uint, where string) (results []model.TaskResult, err error) {
	return
}
func (repo *MySQLRepo) CountTaskResult(where string) (count int, err error) {
	return
}

func (repo *MySQLRepo) GetTaskResultByTaskNo(taskNo string) (results []*model.TaskResult, err error) {
	err = repo.db.Model(model.TaskResult{}).Where("task_no = ?", taskNo).Find(&results).Error
	return
}

func (repo *MySQLRepo) AddTasks(info *model.TaskInfo, results []*model.TaskResult) (err error) {
	tx := repo.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
