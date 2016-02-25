package initdb

import (
	"config"
	"model"

	"github.com/jinzhu/gorm"
)

var osInstallTables = []interface{}{
	&model.Device{},
	&model.Network{},
	&model.OsConfig{},
	&model.SystemConfig{},
	&model.Hardware{},
	&model.Location{},
	&model.Ip{},
	&model.DeviceLog{},
	&model.DeviceHistory{},
	&model.Mac{},
	&model.Manufacturer{},
	&model.VmDevice{},
	&model.User{},
}

func InitOsInstallTables(db *gorm.DB, conf *config.Config) error {
	//db.DropTableIfExists(osInstallTables...)
	db.CreateTable(osInstallTables...)
	return nil
}

func DropOsInstallTables(db *gorm.DB, conf *config.Config) error {
	db.DropTableIfExists(osInstallTables...)
	return nil
}
