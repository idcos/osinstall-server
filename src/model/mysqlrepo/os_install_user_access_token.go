package mysqlrepo

import (
	//"fmt"
	"model"
)

func (repo *MySQLRepo) AddUserAccessToken(userId uint, accessToken string) (*model.UserAccessToken, error) {
	mod := model.UserAccessToken{UserID: userId, AccessToken: accessToken}
	err := repo.db.Create(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) DeleteUserAccessTokenByToken(accessToken string) (*model.UserAccessToken, error) {
	mod := model.UserAccessToken{}
	err := repo.db.Unscoped().Where("access_token = ?", accessToken).Delete(&mod).Error
	return &mod, err
}

func (repo *MySQLRepo) CountUserAccessTokenByToken(accessToken string) (uint, error) {
	mod := model.UserAccessToken{AccessToken: accessToken}
	var count uint
	err := repo.db.Model(mod).Where("access_token = ?", accessToken).Count(&count).Error
	return count, err
}

func (repo *MySQLRepo) GetUserByAccessToken(accessToken string) (*model.UserWithToken, error) {
	var result model.UserWithToken
	err := repo.db.Raw("SELECT t1.access_token,t2.username,t2.role,t2.name,t2.id from user_access_tokens t1 inner join users t2 on t1.user_id = t2.id where t1.access_token = ?", accessToken).Scan(&result).Error
	return &result, err
}
