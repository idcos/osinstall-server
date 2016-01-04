package mysqlrepo

import (
	"fmt"
	"model"
	"strings"
)

func (repo *MySQLRepo) AddLocation(pid uint, name string) (*model.Location, error) {
	name = strings.Replace(name, "-", "_", -1)
	mod := model.Location{Pid: pid, Name: name}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) UpdateLocationById(id uint, pid uint, name string) (*model.Location, error) {
	name = strings.Replace(name, "-", "_", -1)
	mod := model.Location{Pid: pid, Name: name}
	err := repo.db.First(&mod, id).Update("name", name).Update("pid", pid).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteLocationById(id uint) (*model.Location, error) {
	mod := model.Location{}
	err := repo.db.Unscoped().Where("id = ?", id).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountLocationByName(name string) (uint, error) {
	mod := model.Location{Name: name}
	var count uint
	err := repo.db.Model(mod).Where("name = ?", name).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountLocation() (uint, error) {
	mod := model.Location{}
	var count uint
	err := repo.db.Model(mod).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetLocationListWithPage(limit uint, offset uint) ([]model.Location, error) {
	var mods []model.Location
	err := repo.db.Limit(limit).Offset(offset).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) FormatLocationToTreeByPid(pid uint, content []map[string]interface{}, floor uint, selectPid uint) ([]map[string]interface{}, error) {
	var mods []model.Location
	//result := make(map[uint]interface{})
	err := repo.db.Unscoped().Where("pid = ?", pid).Find(&mods).Error
	if err != nil {
		return content, err
	}
	//result[pid] = mods
	for _, v := range mods {
		data := make(map[string]interface{})
		data["ID"] = v.ID
		data["Pid"] = v.Pid
		data["Name"] = v.Name
		//var strSelect string
		if selectPid == v.ID {
			data["Selected"] = true
		} else {
			data["Selected"] = false
		}
		data["ShowName"] = treeStrRepeat(floor) + "|-" + v.Name

		content = append(content, data)
		if v.ID != pid {
			childContent, _ := repo.FormatLocationToTreeByPid(v.ID, nil, floor+1, selectPid)
			for _, currentContent := range childContent {
				content = append(content, currentContent)
			}
		}
	}
	return content, nil
}

func (repo *MySQLRepo) GetParentLocationIdByName(name string) (uint, error) {
	list := strings.Split(name, "-")
	pid := uint(0)
	for _, location := range list {
		location = strings.TrimSpace(location)
		if location == "" {
			continue
		}
		count, err := repo.CountLocationByNameAndPid(location, pid)
		if err != nil {
			return uint(0), err
		}
		if count > 0 {
			//var mod model.Location
			mod, err := repo.GetLocationByNameAndPid(location, pid)
			if err != nil {
				return uint(0), err
			}
			pid = mod.ID
		} else {
			return uint(0), err
		}
	}
	return pid, nil
}

func (repo *MySQLRepo) ImportLocation(name string) (uint, error) {
	list := strings.Split(name, "-")
	pid := uint(0)
	for _, location := range list {
		location = strings.TrimSpace(location)
		if location == "" {
			continue
		}
		count, err := repo.CountLocationByNameAndPid(location, pid)
		if err != nil {
			return uint(0), err
		}

		if count <= 0 {
			_, err := repo.AddLocation(pid, location)
			if err != nil {
				return uint(0), err
			}
		}

		mod, err := repo.GetLocationByNameAndPid(location, pid)
		if err != nil {
			return uint(0), err
		}
		pid = mod.ID
	}
	return pid, nil
}

func (repo *MySQLRepo) FormatLocationNameById(id uint, content string, separator string) (string, error) {
	var mod model.Location
	if id <= uint(0) {
		return content, nil
	}
	//result := make(map[uint]interface{})
	err := repo.db.Unscoped().Where("id = ?", id).Find(&mod).Error
	if err != nil {
		return content, err
	}

	if content == "" {
		content = mod.Name
	} else {
		content = mod.Name + separator + content
	}

	if mod.Pid > 0 {
		parentContent, _ := repo.FormatLocationNameById(mod.Pid, "", separator)
		content = parentContent + separator + content
	}
	return content, nil
}

func (repo *MySQLRepo) FormatChildLocationIdById(id uint, content string, separator string) (string, error) {
	var mods []model.Location
	err := repo.db.Unscoped().Where("pid = ?", id).Find(&mods).Error
	if err != nil {
		return content, err
	}

	for _, v := range mods {
		//content += fmt.Sprintf("%d", v.ID) + ","
		childContent, _ := repo.FormatChildLocationIdById(v.ID, "", ",")
		if childContent != "" {
			if content != "" {
				content = content + "," + childContent
			} else {
				content = childContent
			}
		}
	}

	if content != "" {
		content = fmt.Sprintf("%d", id) + "," + content
	} else {
		content = fmt.Sprintf("%d", id)
	}

	return content, nil
}

func treeStrRepeat(floor uint) string {
	var str string
	for i := uint(0); i < (floor * uint(1)); i++ {
		str = str + "　"
	}
	return str
}

func (repo *MySQLRepo) GetLocationListByPidWithPage(limit uint, offset uint, pid uint) ([]model.Location, error) {
	var mods []model.Location
	err := repo.db.Limit(limit).Offset(offset).Where("pid = ?", pid).Find(&mods).Error
	return mods, err
}

func (repo *MySQLRepo) CountLocationByPid(pid uint) (uint, error) {
	mod := model.Location{Pid: pid}
	var count uint
	err := repo.db.Model(mod).Where("pid = ?", pid).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountLocationByNameAndPid(name string, pid uint) (uint, error) {
	mod := model.Location{}
	var count uint
	err := repo.db.Model(mod).Where("name = ? and pid = ?", name, pid).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) CountLocationByNameAndPidAndId(name string, pid uint, id uint) (uint, error) {
	mod := model.Location{}
	var count uint
	err := repo.db.Model(mod).Where("name = ? and pid = ? and id != ?", name, pid, id).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetLocationById(id uint) (*model.Location, error) {
	var mod model.Location
	err := repo.db.Where("id = ?", id).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetLocationByNameAndPid(name string, pid uint) (*model.Location, error) {
	var mod model.Location
	err := repo.db.Where("name = ? and pid = ?", name, pid).Find(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) GetLocationIdByName(name string) (uint, error) {
	mod := model.Location{Name: name}
	err := repo.db.Where("name = ?", name).Find(&mod).Error
	return mod.ID, err
}
