package mysqlrepo

import (
	"config"
	"logger"

	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"server/osinstallserver/util"
)

const (
	// ViewAction 查看操作
	ViewAction Action = "view"
	// ExecAction 执行命令操作
	ExecAction Action = "exec"
	// UpdateAction 更新操作
	UpdateAction Action = "update"
	// InsertAction 新增操作
	InsertAction Action = "insert"
	// DeleteAction 删除操作
	DeleteAction Action = "delete"

	// DataCenter 权限资源类型-数据中心
	DataCenter ResourceType = "dataCenter"
	// AppRes 权限资源类型-应用
	AppRes ResourceType = "app"
	// UnitRes 权限资源类型-部署单元
	UnitRes ResourceType = "unit"
	// NodeRes 权限资源类型-主机节点
	NodeRes ResourceType = "node"
)

// Action 操作
type Action string

// ResourceType 权限资源类型
type ResourceType string

// MySQLRepo mysql数据库实现
type MySQLRepo struct {
	conf *config.Config
	log  logger.Logger
	db   *gorm.DB
}

// NewRepo 创建mysql数据实现实例
func NewRepo(conf *config.Config, log logger.Logger) (*MySQLRepo, error) {
	var connection string
	if conf.Repo.ConnectionIsCrypted != "" {
		str, err := util.RSADecrypt(conf.Rsa.PrivateKey, conf.Repo.ConnectionIsCrypted)
		if err != nil {
			return nil, err
		}
		connection = str
	} else if conf.Repo.Connection != "" {
		connection = conf.Repo.Connection
	}

	if connection == "" {
		return nil, errors.New("请先配置好数据库连接信息!")
	}

	db, err := gorm.Open("mysql", connection)
	if err != nil {
		log.Errorf("database connection failed:%s", err.Error())
		return nil, err
	}
	db.LogMode(true)
	repo := &MySQLRepo{
		conf: conf,
		log:  log,
		db:   &db,
	}

	return repo, nil
}

// Close 关闭mysql连接
func (repo *MySQLRepo) Close() error {
	return repo.db.Close()
}

// DropDB 删除表
func (repo *MySQLRepo) DropDB() error {
	return nil
}
