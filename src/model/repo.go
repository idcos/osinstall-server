package model

// Repo 数据仓库
type Repo interface {
	Close() error
	DropDB() error // 测试时使用

	//装机相关
	IDevice
	INetwork
	IOsConfig
	ISystemConfig
	IHardware
	ILocation
	IIp
	IMac
	IManufacturer
	IDeviceLog
	IDeviceHistory
	IVmDevice
}
