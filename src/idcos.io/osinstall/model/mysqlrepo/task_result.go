package mysqlrepo

import "idcos.io/osinstall/model"

func (repo *MySQLRepo) AddTaskResult(info *model.TaskResult) (err error) {
	return
}
func (repo *MySQLRepo) GetTaskResultPage(limit uint, offset uint, where string) (results []model.TaskResult, err error) {
	return
}
func (repo *MySQLRepo) CountTaskResult(where string) (count int, err error) {
	return
}
