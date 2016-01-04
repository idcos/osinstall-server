package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DeviceFull struct {
	ID              uint
	BatchNumber     string
	Sn              string
	Hostname        string
	Ip              string
	NetworkID       uint
	OsID            uint
	HardwareID      uint
	SystemID        uint
	Location        string
	LocationID      uint
	AssetNumber     string
	Status          string
	InstallProgress float64
	InstallLog      string
	NetworkName     string
	OsName          string
	HardwareName    string
	SystemName      string
	IsSupportVm     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Device struct {
	gorm.Model
	BatchNumber     string  `sql:"not null;"`        //录入批次号
	Sn              string  `sql:"not null;unique;"` //序列号
	Hostname        string  `sql:"not null;"`        //主机名
	Ip              string  `sql:"not null;unique;"` //IP
	NetworkID       uint    `sql:"not null;"`        //网段模板ID
	OsID            uint    `sql:"not null;"`        //操作系统ID
	HardwareID      uint    ``                       //硬件配置模板ID
	SystemID        uint    `sql:"not null;"`        //系统配置模板ID
	Location        string  `sql:"not null;"`        //位置
	LocationID      uint    `sql:"not null;"`
	AssetNumber     string  //财编
	Status          string  `sql:"not null;"`                     //状态 'pre_run' 待安装,'running' 安装中,'success' 安装成功,'failure' 安装失败
	InstallProgress float64 `sql:"type:decimal(11,4);default:0;"` //安装进度
	InstallLog      string  `sql:"type:text;"`                    //安装日志
	IsSupportVm     string  `sql:"enum('Yes','No');NOT NULL;DEFAULT 'Yes'"`
}

// IDevice 设备操作接口
type IDevice interface {
	GetDeviceBySnAndStatus(sn string, status string) (*Device, error)
	CountDeviceBySn(sn string) (uint, error)
	CountDeviceByHostname(hostname string) (uint, error)
	CountDeviceByHostnameAndId(hostname string, id uint) (uint, error)
	CountDeviceByIp(ip string) (uint, error)
	CountDeviceByIpAndId(ip string, id uint) (uint, error)
	GetDeviceIdBySn(sn string) (uint, error)
	CountDevice(where string) (int, error)
	GetDeviceListWithPage(Limit uint, Offset uint, where string) ([]DeviceFull, error)
	GetDeviceById(Id uint) (*Device, error)
	DeleteDeviceById(Id uint) (*Device, error)
	ReInstallDeviceById(Id uint) (*Device, error)
	CreateBatchNumber() (string, error)
	AddDevice(BatchNumber string, Sn string, Hostname string, Ip string, NetworkID uint, OsID uint, HardwareID uint, SystemID uint, Location string, LocationID uint, AssetNumber string, Status string, IsSupportVm string) (*Device, error)
	UpdateDeviceById(ID uint, BatchNumber string, Sn string, Hostname string, Ip string, NetworkID uint, OsID uint, HardwareID uint, SystemID uint, Location string, LocationID uint, AssetNumber string, Status string, IsSupportVm string) (*Device, error)
	UpdateInstallInfoById(ID uint, status string, installProgress float64) (*Device, error)
	GetSystemBySn(sn string) (*SystemConfig, error)
	GetNetworkBySn(sn string) (*Network, error)
	GetFullDeviceById(id uint) (*DeviceFull, error)
	CountDeviceByWhere(where string) (int, error)
	GetDeviceByWhere(where string) ([]Device, error)
}

type DeviceHistory struct {
	gorm.Model
	BatchNumber     string  `sql:"not null;"`        //录入批次号
	Sn              string  `sql:"not null;unique;"` //序列号
	Hostname        string  `sql:"not null;"`        //主机名
	Ip              string  `sql:"not null;unique;"` //IP
	NetworkID       uint    `sql:"not null;"`        //网段模板ID
	OsID            uint    `sql:"not null;"`        //操作系统ID
	HardwareID      uint    ``                       //硬件配置模板ID
	SystemID        uint    `sql:"not null;"`        //系统配置模板ID
	Location        string  `sql:"not null;"`        //位置
	LocationID      uint    `sql:"not null;"`
	AssetNumber     string  //财编
	Status          string  `sql:"not null;"`                     //状态 'pre_run' 待安装,'running' 安装中,'success' 安装成功,'failure' 安装失败
	InstallProgress float64 `sql:"type:decimal(11,4);default:0;"` //安装进度
	InstallLog      string  `sql:"type:text;"`                    //安装日志
	IsSupportVm     string
}

// IDevice 设备操作接口
type IDeviceHistory interface {
	UpdateHistoryDeviceStatusById(ID uint, status string) (*DeviceHistory, error)
	CopyDeviceToHistory(ID uint) error
}

// Network 网络
type Network struct {
	gorm.Model
	Network string `sql:"not null;unique;"` //网段
	Netmask string `sql:"not null;`         //掩码
	Gateway string `sql:"not null;"`        //网关
	Vlan    string //vlan
	Trunk   string //trunk
	Bonding string //bonding
}

// INetwork 网络操作接口
type INetwork interface {
	CountNetworkByNetwork(Network string) (uint, error)
	GetNetworkIdByNetwork(Network string) (uint, error)
	CountNetworkByNetworkAndId(Network string, ID uint) (uint, error)
	CountNetwork() (uint, error)
	GetNetworkListWithPage(Limit uint, Offset uint) ([]Network, error)
	GetNetworkById(Id uint) (*Network, error)
	UpdateNetworkById(Id uint, Network string, Netmask string, Gateway string, Vlan string, Trunk string, Bonding string) (*Network, error)
	DeleteNetworkById(Id uint) (*Network, error)
	AddNetwork(Network string, Netmask string, Gateway string, Vlan string, Trunk string, Bonding string) (*Network, error)
}

// Network 网络
type Ip struct {
	gorm.Model
	NetworkID uint   `sql:"not null;"`
	Ip        string `sql:"not null;"`
}

// INetwork 网络操作接口
type IIp interface {
	DeleteIpByNetworkId(NetworkID uint) (*Ip, error)
	AddIp(NetworkID uint, Ip string) (*Ip, error)
	CountIpByIp(Ip string) (uint, error)
	GetIpByIp(Ip string) (*Ip, error)
	AssignNewIpByNetworkId(NetworkID uint) (string, error)
	GetNotUsedIPListByNetworkId(NetworkID uint) ([]Ip, error)
}

// OS 操作系统
type OsConfig struct {
	gorm.Model
	Name string `sql:"not null;unique;"`    //操作系统名称
	Pxe  string `sql:"type:text;not null;"` //pxe信息
}

// IOS 操作系统操作接口
type IOsConfig interface {
	//GetOSByID(ID uint) (*OsConfig, error)
	CountOsConfigByName(Name string) (uint, error)
	CountOsConfigByNameAndId(Name string, ID uint) (uint, error)
	CountOsConfig() (uint, error)
	GetOsConfigListWithPage(Limit uint, Offset uint) ([]OsConfig, error)
	GetOsConfigIdByName(Name string) (uint, error)
	GetOsConfigById(Id uint) (*OsConfig, error)
	UpdateOsConfigById(Id uint, Name string, Pxe string) (*OsConfig, error)
	DeleteOsConfigById(Id uint) (*OsConfig, error)
	AddOsConfig(Name string, Pxe string) (*OsConfig, error)
	GetOsConfigByName(Name string) (*OsConfig, error)
}

// OS 操作系统
type DeviceLog struct {
	gorm.Model
	DeviceID  uint   `sql:"not null;"`
	Title     string `sql:"not null;"`
	Type      string `sql:"not null;default:'install';"`
	Content   string `sql:"type:text;"` //pxe信息
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IDeviceLog interface {
	CountDeviceLogByDeviceID(DeviceID uint) (uint, error)
	CountDeviceLogByDeviceIDAndType(DeviceID uint, Type string) (uint, error)
	CountDeviceLog() (uint, error)
	GetDeviceLogListByDeviceID(DeviceID uint, Order string) ([]DeviceLog, error)
	GetLastDeviceLogByDeviceID(DeviceID uint) (DeviceLog, error)
	GetDeviceLogListByDeviceIDAndType(DeviceID uint, Type string, Order string, MaxID uint) ([]DeviceLog, error)
	GetDeviceLogById(Id uint) (*DeviceLog, error)
	DeleteDeviceLogById(Id uint) (*DeviceLog, error)
	DeleteDeviceLogByDeviceIDAndType(DeviceID uint, Type string) (*DeviceLog, error)
	DeleteDeviceLogByDeviceID(DeviceID uint) (*DeviceLog, error)
	AddDeviceLog(DeviceID uint, Title string, Type string, Content string) (*DeviceLog, error)
	UpdateDeviceLogTypeByDeviceIdAndType(deviceID uint, Type string, NewType string) ([]DeviceLog, error)
}

// System 系统配置
type SystemConfig struct {
	gorm.Model
	Name    string `sql:"not null;unique;"`    //操作系统名称
	Content string `sql:"type:text;not null;"` //信息
}

// ISystemConfg 操作系统操作接口
type ISystemConfig interface {
	//GetSystemByID(ID uint) (*SystemConfig, error)
	CountSystemConfigByName(Name string) (uint, error)
	CountSystemConfigByNameAndId(Name string, ID uint) (uint, error)
	GetSystemConfigIdByName(Name string) (uint, error)
	CountSystemConfig() (uint, error)
	GetSystemConfigListWithPage(Limit uint, Offset uint) ([]SystemConfig, error)
	GetSystemConfigById(Id uint) (*SystemConfig, error)
	UpdateSystemConfigById(Id uint, Name string, Content string) (*SystemConfig, error)
	DeleteSystemConfigById(Id uint) (*SystemConfig, error)
	AddSystemConfig(Name string, Content string) (*SystemConfig, error)
}

// Hardware 硬件配置
type Hardware struct {
	gorm.Model
	Company     string `sql:"not null;"`  //企业名称
	Product     string `sql:"not null;"`  //产品
	ModelName   string `sql:"not null;"`  //型号
	Raid        string `sql:"type:text;"` //raid配置
	Oob         string `sql:"type:text;"` //oob配置
	Bios        string `sql:"type:text;"` //bios配置
	IsSystemAdd string `sql:"enum('Yes','No');NOT NULL;DEFAULT 'Yes'"`
	Tpl         string //厂商提交的JSON信息
	Data        string //最终要执行的脚本信息
}

// IHardware 硬件配置操作接口
type IHardware interface {
	GetHardwareBySn(sn string) (*Hardware, error)
	CountHardwareByCompanyAndProductAndName(Company string, Product string, ModelName string) (uint, error)
	CountHardwareByCompanyAndProductAndNameAndId(Company string, Product string, ModelName string, ID uint) (uint, error)
	CountHardwarrWithSeparator(Name string) (uint, error)
	GetHardwareIdByCompanyAndProductAndName(Company string, Product string, ModelName string) (uint, error)
	CountHardware(where string) (uint, error)
	GetHardwareListWithPage(Limit uint, Offset uint, where string) ([]Hardware, error)
	GetHardwareById(Id uint) (*Hardware, error)
	UpdateHardwareById(Id uint, Company string, Product string, ModelName string, Raid string, Oob string, Bios string, Tpl string, Data string) (*Hardware, error)
	DeleteHardwareById(Id uint) (*Hardware, error)
	AddHardware(Company string, Product string, ModelName string, Raid string, Oob string, Bios string, IsSystemAdd string, Tpl string, Data string) (*Hardware, error)
	GetCompanyByGroup() ([]Hardware, error)
	GetProductByWhereAndGroup(where string) ([]Hardware, error)
	GetModelNameByWhereAndGroup(where string) ([]Hardware, error)
	GetHardwareBySeaprator(Name string) (*Hardware, error)
}

// Location 位置
type Location struct {
	gorm.Model
	Pid  uint   `sql:"not null;"` //父级ID
	Name string `sql:"not null;"` //位置名
}

// ILocation 位置操作接口
type ILocation interface {
	CountLocationByName(Name string) (uint, error)
	GetLocationIdByName(Name string) (uint, error)
	CountLocation() (uint, error)
	GetLocationListWithPage(Limit uint, Offset uint) ([]Location, error)
	//FormatLocationToTreeByPid(Pid uint, Content string, Floor uint, SelectPid uint) (string, error)
	FormatLocationToTreeByPid(Pid uint, Content []map[string]interface{}, Floor uint, SelectPid uint) ([]map[string]interface{}, error)
	FormatLocationNameById(id uint, content string, separator string) (string, error)
	GetLocationListByPidWithPage(Limit uint, Offset uint, pid uint) ([]Location, error)
	CountLocationByPid(Pid uint) (uint, error)
	CountLocationByNameAndPid(Name string, Pid uint) (uint, error)
	CountLocationByNameAndPidAndId(Name string, Pid uint, ID uint) (uint, error)
	GetLocationById(Id uint) (*Location, error)
	UpdateLocationById(Id uint, Pid uint, Name string) (*Location, error)
	DeleteLocationById(Id uint) (*Location, error)
	AddLocation(Pid uint, Name string) (*Location, error)
	GetLocationByNameAndPid(Name string, Pid uint) (*Location, error)
	ImportLocation(Name string) (uint, error)
	FormatChildLocationIdById(id uint, content string, separator string) (string, error)
}

// Mac mac地址
type Mac struct {
	gorm.Model
	DeviceID uint   `sql:"not null;"`
	Mac      string `sql:"not null;unique;"` //位置名
}

type IMac interface {
	CountMacByMac(Mac string) (uint, error)
	CountMacByMacAndDeviceID(Mac string, DeviceID uint) (uint, error)
	GetMacById(Id uint) (*Mac, error)
	DeleteMacById(Id uint) (*Mac, error)
	AddMac(DeviceID uint, Mac string) (*Mac, error)
	GetMacListByDeviceID(DeviceID uint) ([]Mac, error)
	DeleteMacByDeviceId(deviceId uint) (*Mac, error)
}

type Manufacturer struct {
	gorm.Model
	DeviceID  uint   `sql:"not null;"`
	Company   string `sql:"not null;"`
	Product   string
	ModelName string
}

type IManufacturer interface {
	CountManufacturerByDeviceID(DeviceID uint) (uint, error)
	GetManufacturerById(Id uint) (*Manufacturer, error)
	GetManufacturerByDeviceID(DeviceID uint) (*Manufacturer, error)
	DeleteManufacturerById(Id uint) (*Manufacturer, error)
	AddManufacturer(DeviceID uint, Company string, Product string, ModelName string) (*Manufacturer, error)
	UpdateManufacturerById(Id uint, Company string, Product string, ModelName string) (*Manufacturer, error)
}

type VmDevice struct {
	gorm.Model
	DeviceID              uint
	Hostname              string
	Mac                   string
	Ip                    string
	NetworkID             uint
	OsID                  uint
	CpuCoresNumber        uint
	CpuHotPlug            string
	CpuPassthrough        string
	CpuTopSockets         uint
	CpuTopCores           uint
	CpuTopThreads         uint
	CpuPinning            string
	MemoryCurrent         uint
	MemoryMax             uint
	MemoryKsm             string
	DiskType              string
	DiskSize              uint
	DiskBusType           string
	DiskCacheMode         string
	DiskIoMode            string
	NetworkType           string
	NetworkDeviceType     string
	DisplayType           string
	DisplayPassword       string
	DisplayUpdatePassword string
	Status                string
}

type VmDeviceFull struct {
	gorm.Model
	DeviceID              uint
	DeviceSn              string
	Hostname              string
	Mac                   string
	Ip                    string
	NetworkID             uint
	NetworkName           string
	OsID                  uint
	OsName                string
	CpuCoresNumber        uint
	CpuHotPlug            string
	CpuPassthrough        string
	CpuTopSockets         uint
	CpuTopCores           uint
	CpuTopThreads         uint
	CpuPinning            string
	MemoryCurrent         uint
	MemoryMax             uint
	MemoryKsm             string
	DiskType              string
	DiskSize              uint
	DiskBusType           string
	DiskCacheMode         string
	DiskIoMode            string
	NetworkType           string
	NetworkDeviceType     string
	DisplayType           string
	DisplayPassword       string
	DisplayUpdatePassword string
	Status                string
}

type IVmDevice interface {
	CountVmDeviceByHostname(Hostname string) (uint, error)
	CountVmDeviceByMac(Mac string) (uint, error)
	CountVmDeviceByIp(Ip string) (uint, error)
	CountVmDeviceByHostnameAndId(Hostname string, ID uint) (uint, error)
	CountVmDeviceByMacAndId(Mac string, ID uint) (uint, error)
	CountVmDeviceByIpAndId(Ip string, ID uint) (uint, error)
	CountVmDevice(Where string) (int, error)
	GetVmDeviceListWithPage(Limit uint, Offset uint, Where string) ([]VmDeviceFull, error)
	GetVmDeviceById(Id uint) (*VmDevice, error)
	GetFullVmDeviceById(Id uint) (*VmDeviceFull, error)
	GetVmDeviceByMac(Mac string) (*VmDevice, error)
	GetVmDeviceIdByMac(Mac string) (uint, error)
	DeleteVmDeviceById(Id uint) (*VmDevice, error)
	ReInstallVmDeviceById(Id uint) (*VmDevice, error)
	AddVmDevice(DeviceID uint,
		Hostname string,
		Mac string,
		Ip string,
		NetworkID uint,
		OsID uint,
		CpuCoresNumber uint,
		CpuHotPlug string,
		CpuPassthrough string,
		CpuTopSockets uint,
		CpuTopCores uint,
		CpuTopThreads uint,
		CpuPinning string,
		MemoryCurrent uint,
		MemoryMax uint,
		MemoryKsm string,
		DiskType string,
		DiskSize uint,
		DiskBusType string,
		DiskCacheMode string,
		DiskIoMode string,
		NetworkType string,
		NetworkDeviceType string,
		DisplayType string,
		DisplayPassword string,
		DisplayUpdatePassword string,
		Status string) (VmDevice, error)
	UpdateVmDeviceById(ID uint,
		DeviceID uint,
		Hostname string,
		Mac string,
		Ip string,
		NetworkID uint,
		OsID uint,
		CpuCoresNumber uint,
		CpuHotPlug string,
		CpuPassthrough string,
		CpuTopSockets uint,
		CpuTopCores uint,
		CpuTopThreads uint,
		CpuPinning string,
		MemoryCurrent uint,
		MemoryMax uint,
		MemoryKsm string,
		DiskType string,
		DiskSize uint,
		DiskBusType string,
		DiskCacheMode string,
		DiskIoMode string,
		NetworkType string,
		NetworkDeviceType string,
		DisplayType string,
		DisplayPassword string,
		DisplayUpdatePassword string,
		Status string) (VmDevice, error)
}
