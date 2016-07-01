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
	&model.ManageNetwork{},
	&model.ManageIp{},
	&model.UserAccessToken{},
	&model.DeviceInstallReport{},
	&model.DeviceInstallCallback{},
	&model.DhcpSubnet{},
	&model.PlatformConfig{},
	&model.VmHost{},
	&model.VmDeviceLog{},
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
