package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DeviceFull struct {
	ID                uint
	BatchNumber       string
	Sn                string
	Hostname          string
	Ip                string
	ManageIp          string
	NetworkID         uint
	ManageNetworkID   uint
	OsID              uint
	HardwareID        uint
	SystemID          uint
	Location          string
	LocationID        uint
	AssetNumber       string
	Status            string
	InstallProgress   float64
	InstallLog        string
	NetworkName       string
	ManageNetworkName string
	OsName            string
	HardwareName      string
	SystemName        string
	IsSupportVm       string
	UserID            uint
	OwnerName         string
	Callback          string
	BootosIp          string
	OobIp             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Device struct {
	gorm.Model
	BatchNumber     string  `sql:"not null;"`        //录入批次号
	Sn              string  `sql:"not null;unique;"` //序列号
	Hostname        string  `sql:"not null;"`        //主机名
	Ip              string  `sql:"not null;unique;"` //IP
	ManageIp        string  `sql:"unique;"`          //IP
	NetworkID       uint    `sql:"not null;"`        //网段模板ID
	ManageNetworkID uint    ``                       //管理网段模板ID
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
	UserID          uint    `sql:"not null;default:0;"`
}

// IDevice 设备操作接口
type IDevice interface {
	GetDeviceBySnAndStatus(sn string, status string) (*Device, error)
	CountDeviceBySn(sn string) (uint, error)
	CountDeviceByHostname(hostname string) (uint, error)
	CountDeviceByHostnameAndId(hostname string, id uint) (uint, error)
	CountDeviceByIp(ip string) (uint, error)
	CountDeviceByManageIp(ManageIp string) (uint, error)
	CountDeviceByIpAndId(ip string, id uint) (uint, error)
	CountDeviceByManageIpAndId(ManageIp string, id uint) (uint, error)
	GetDeviceIdBySn(sn string) (uint, error)
	GetDeviceBySn(sn string) (*Device, error)
	CountDevice(where string) (int, error)
	GetDeviceListWithPage(Limit uint, Offset uint, where string) ([]DeviceFull, error)
	GetDeviceById(Id uint) (*Device, error)
	DeleteDeviceById(Id uint) (*Device, error)
	ReInstallDeviceById(Id uint) (*Device, error)
	CancelInstallDeviceById(Id uint) (*Device, error)
	CreateBatchNumber() (string, error)
	AddDevice(BatchNumber string, Sn string, Hostname string, Ip string, ManageIp string, NetworkID uint, ManageNetworkID uint, OsID uint, HardwareID uint, SystemID uint, Location string, LocationID uint, AssetNumber string, Status string, IsSupportVm string, UserID uint) (*Device, error)
	UpdateDeviceById(ID uint, BatchNumber string, Sn string, Hostname string, Ip string, ManageIp string, NetworkID uint, ManageNetworkID uint, OsID uint, HardwareID uint, SystemID uint, Location string, LocationID uint, AssetNumber string, Status string, IsSupportVm string, UserID uint) (*Device, error)
	UpdateInstallInfoById(ID uint, status string, installProgress float64) (*Device, error)
	GetSystemBySn(sn string) (*SystemConfig, error)
	GetNetworkBySn(sn string) (*Network, error)
	GetFullDeviceById(id uint) (*DeviceFull, error)
	CountDeviceByWhere(where string) (int, error)
	GetDeviceByWhere(where string) ([]Device, error)
	GetInstallTimeoutDeviceList(timeout int) ([]Device, error)
	IsInstallTimeoutDevice(timeout int, deviceId uint) (bool, error)
	ExecDBVersionUpdateSql(sql string) error
}

type DeviceHistory struct {
	gorm.Model
	BatchNumber     string  `sql:"not null;"`        //录入批次号
	Sn              string  `sql:"not null;unique;"` //序列号
	Hostname        string  `sql:"not null;"`        //主机名
	Ip              string  `sql:"not null;unique;"` //IP
	ManageIp        string  `sql:"unique;"`          //ManageIP
	NetworkID       uint    `sql:"not null;"`        //网段模板ID
	ManageNetworkID uint    ``                       //管理网段模板ID
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

type DeviceInstallReport struct {
	gorm.Model
	Sn           string `sql:"not null;unique;"` //序列号
	OsName       string
	HardwareName string
	SystemName   string
	Status       string
	UserID       uint
}

type DeviceHardwareNameInstallReport struct {
	HardwareName string
	Count        uint
}

type DeviceProductNameInstallReport struct {
	ProductName string
	Count       uint
}

type DeviceOsNameInstallReport struct {
	OsName string
	Count  uint
}

type DeviceSystemNameInstallReport struct {
	SystemName string
	Count      uint
}

// IDevice 设备操作接口
type IDeviceInstallReport interface {
	CopyDeviceToInstallReport(ID uint) error
	CopyVmDeviceToInstallReport(ID uint) error
	CountDeviceInstallReportByWhere(Where string) (uint, error)
	GetDeviceHardwareNameInstallReport(Where string) ([]DeviceHardwareNameInstallReport, error)
	GetDeviceProductNameInstallReport(Where string) ([]DeviceProductNameInstallReport, error)
	GetDeviceCompanyNameInstallReport(Where string) ([]DeviceProductNameInstallReport, error)
	GetDeviceOsNameInstallReport(Where string) ([]DeviceOsNameInstallReport, error)
	GetDeviceSystemNameInstallReport(Where string) ([]DeviceSystemNameInstallReport, error)
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

// ManageNetwork 网络
type ManageNetwork struct {
	gorm.Model
	Network string `sql:"not null;unique;"` //网段
	Netmask string `sql:"not null;`         //掩码
	Gateway string `sql:"not null;"`        //网关
	Vlan    string //vlan
	Trunk   string //trunk
	Bonding string //bonding
}

// INetwork 网络操作接口
type IManageNetwork interface {
	CountManageNetworkByNetwork(Network string) (uint, error)
	GetManageNetworkIdByNetwork(Network string) (uint, error)
	CountManageNetworkByNetworkAndId(Network string, ID uint) (uint, error)
	CountManageNetwork() (uint, error)
	GetManageNetworkListWithPage(Limit uint, Offset uint) ([]ManageNetwork, error)
	GetManageNetworkById(Id uint) (*ManageNetwork, error)
	UpdateManageNetworkById(Id uint, Network string, Netmask string, Gateway string, Vlan string, Trunk string, Bonding string) (*ManageNetwork, error)
	DeleteManageNetworkById(Id uint) (*ManageNetwork, error)
	AddManageNetwork(Network string, Netmask string, Gateway string, Vlan string, Trunk string, Bonding string) (*ManageNetwork, error)
	GetManufacturerMacBySn(Sn string) (string, error)
}

// Network 网络
type ManageIp struct {
	gorm.Model
	NetworkID uint   `sql:"not null;"`
	Ip        string `sql:"not null;"`
}

// INetwork 网络操作接口
type IManageIp interface {
	DeleteManageIpByNetworkId(NetworkID uint) (*ManageIp, error)
	AddManageIp(NetworkID uint, Ip string) (*ManageIp, error)
	CountManageIpByIp(Ip string) (uint, error)
	GetManageIpByIp(Ip string) (*ManageIp, error)
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
	Source      string //来源
	Version     string //版本
	Status      string `sql:"enum('Pending','Success','Failure');NOT NULL;DEFAULT 'Success'"` //状态
}

// IHardware 硬件配置操作接口
type IHardware interface {
	GetHardwareBySn(sn string) (*Hardware, error)
	CountHardwareByCompanyAndProductAndName(Company string, Product string, ModelName string) (uint, error)
	CountHardwareByCompanyAndProductAndNameAndId(Company string, Product string, ModelName string, ID uint) (uint, error)
	CountHardwareWithSeparator(Name string) (uint, error)
	GetHardwareIdByCompanyAndProductAndName(Company string, Product string, ModelName string) (uint, error)
	CountHardware(where string) (uint, error)
	GetHardwareListWithPage(Limit uint, Offset uint, where string) ([]Hardware, error)
	GetHardwareById(Id uint) (*Hardware, error)
	UpdateHardwareById(Id uint, Company string, Product string, ModelName string, Raid string, Oob string, Bios string, Tpl string, Data string, Source string, Version string, Status string) (*Hardware, error)
	DeleteHardwareById(Id uint) (*Hardware, error)
	AddHardware(Company string, Product string, ModelName string, Raid string, Oob string, Bios string, IsSystemAdd string, Tpl string, Data string, Source string, Version string, Status string) (*Hardware, error)
	GetCompanyByGroup() ([]Hardware, error)
	GetProductByWhereAndGroup(where string) ([]Hardware, error)
	GetModelNameByWhereAndGroup(where string) ([]Hardware, error)
	GetHardwareBySeaprator(Name string) (*Hardware, error)
	ValidateHardwareProductModel(Company string, Product string, ModelName string) (bool, error)
	CountHardwareByWhere(Where string) (uint, error)
	GetHardwareByWhere(Where string) (*Hardware, error)
	GetLastestVersionHardware() (Hardware, error)
	CreateHardwareBackupTable(Fix string) error
	RollbackHardwareFromBackupTable(Fix string) error
	DropHardwareBackupTable(Fix string) error
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
	DeviceID         uint   `sql:"not null;"`
	Company          string `sql:"not null;"`
	Product          string
	ModelName        string
	Sn               string
	Ip               string
	Mac              string
	Nic              string
	Cpu              string
	CpuSum           uint `sql:"type:int(11);default:0;"`
	Memory           string
	MemorySum        uint `sql:"type:int(11);default:0;"`
	Disk             string
	DiskSum          uint `sql:"type:int(11);default:0;"`
	Motherboard      string
	Raid             string
	Oob              string
	UserID           uint   `sql:"not null;default:0;"`
	IsVm             string `sql:"enum('Yes','No');NOT NULL;DEFAULT 'Yes'"`
	IsShowInScanList string `sql:"enum('Yes','No');NOT NULL;DEFAULT 'Yes'"`
	NicDevice        string
}

type ManufacturerFull struct {
	ID               uint
	DeviceID         uint
	Company          string
	Product          string
	ModelName        string
	Sn               string
	Ip               string
	Mac              string
	Nic              string
	Cpu              string
	CpuSum           uint
	Memory           string
	MemorySum        uint
	Disk             string
	DiskSum          uint
	Motherboard      string
	Raid             string
	Oob              string
	UserID           uint
	OwnerName        string
	IsVm             string
	NicDevice        string
	IsShowInScanList string
}

type IManufacturer interface {
	CountManufacturerByDeviceID(DeviceID uint) (uint, error)
	GetManufacturerById(Id uint) (*Manufacturer, error)
	GetManufacturerBySn(Sn string) (*Manufacturer, error)
	GetManufacturerByDeviceId(DeviceID uint) (*Manufacturer, error)
	GetManufacturerByDeviceID(DeviceID uint) (*Manufacturer, error)
	DeleteManufacturerById(Id uint) (*Manufacturer, error)
	DeleteManufacturerBySn(Sn string) (*Manufacturer, error)
	AddManufacturer(DeviceID uint, Company string, Product string, ModelName string, Sn string, Ip string, Mac string, Nic string, Cpu string, CpuSum uint, Memory string, MemorySum uint, Disk string, DiskSum uint, Motherboard string, Raid string, Oob string, IsVm string, NicDevice string, IsShowInScanList string) (*Manufacturer, error)
	UpdateManufacturerById(Id uint, Company string, Product string, ModelName string, Sn string, Ip string, Mac string, Nic string, Cpu string, CpuSum uint, Memory string, MemorySum uint, Disk string, DiskSum uint, Motherboard string, Raid string, Oob string, IsVm string, NicDevice string, IsShowInScanList string) (*Manufacturer, error)
	UpdateManufacturerIsShowInScanListById(id uint, IsShowInScanList string) (*Manufacturer, error)
	UpdateManufacturerDeviceIdById(id uint, deviceId uint) (*Manufacturer, error)
	UpdateManufacturerIPById(id uint, ip string) (*Manufacturer, error)
	GetManufacturerListWithPage(Limit uint, Offset uint, Where string) ([]ManufacturerFull, error)
	CountManufacturerByWhere(Where string) (int, error)
	GetManufacturerCompanyByGroup(Where string) ([]Manufacturer, error)
	GetManufacturerProductByGroup(Where string) ([]Manufacturer, error)
	GetManufacturerModelNameByGroup(Where string) ([]Manufacturer, error)
	CountManufacturerBySn(Sn string) (uint, error)
	GetManufacturerIdBySn(Sn string) (uint, error)
	AssignManufacturerOnwer(Id uint, UserID uint) (*Manufacturer, error)
	AssignManufacturerNewOnwer(NewUserID uint, OldUserID uint) error
	GetManufacturerSnByNicMacForVm(Mac string) (string, error)
}

type VmDevice struct {
	gorm.Model
	DeviceID              uint
	Hostname              string
	Mac                   string
	Ip                    string
	NetworkID             uint
	OsID                  uint
	SystemID              uint
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
	UserID                uint
	VncPort               string
	InstallProgress       float64
	RunStatus             string
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
	SystemID              uint
	OsName                string
	SystemName            string
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
	UserID                uint
	VncPort               string
	InstallProgress       float64
	RunStatus             string
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
	UpdateVmInstallInfoById(ID uint, status string, installProgress float64) (*VmDevice, error)
	UpdateVmRunStatusById(ID uint, runStatus string) (*VmDevice, error)
	GetSystemByVmMac(mac string) (*SystemConfig, error)
	GetNetworkByVmMac(mac string) (*Network, error)
	AddVmDevice(DeviceID uint,
		Hostname string,
		Mac string,
		Ip string,
		NetworkID uint,
		OsID uint,
		SystemID uint,
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
		Status string,
		UserID uint,
		VncPort string,
		RunStatus string) (VmDevice, error)
	UpdateVmDeviceById(ID uint,
		DeviceID uint,
		Hostname string,
		Mac string,
		Ip string,
		NetworkID uint,
		OsID uint,
		SystemID uint,
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
		Status string,
		UserID uint,
		VncPort string,
		RunStatus string) (VmDevice, error)
}

// Mac mac地址
type User struct {
	gorm.Model
	Username    string `sql:"not null;unique;"`
	Password    string `sql:"not null;"`
	Name        string
	PhoneNumber string
	Permission  string
	Status      string `sql:"enum('Enable','Disable');NOT NULL;DEFAULT 'Enable'"`
	Role        string `sql:"enum('Administrator','User');NOT NULL;DEFAULT 'User'"`
}

type IUser interface {
	CountUserByUsername(Username string) (uint, error)
	GetUserByUsername(Username string) (*User, error)
	GetUserById(Id uint) (*User, error)
	CountUserById(Id uint) (uint, error)
	CountUserByWhere(Where string) (uint, error)
	GetUserByWhere(Where string) (*User, error)
	CountUser(Where string) (uint, error)
	DeleteUserById(Id uint) (*User, error)
	AddUser(Username string, Password string, Name string, PhoneNumber string, Permission string, Status string, Role string) (*User, error)
	UpdateUserById(Id uint, Password string, Name string, PhoneNumber string, Permission string, Status string, Role string) (*User, error)
	GetUserListWithPage(Limit uint, Offset uint, Where string) ([]User, error)
}

// Mac mac地址
type UserWithToken struct {
	ID          uint
	Username    string
	Name        string
	PhoneNumber string
	Status      string
	Role        string
	AccessToken string
}

// Mac mac地址
type UserAccessToken struct {
	gorm.Model
	UserID      uint   `sql:"not null;"`
	AccessToken string `sql:"not null;"`
}

type IUserAccessToken interface {
	CountUserAccessTokenByToken(AccessToken string) (uint, error)
	GetUserByAccessToken(AccessToken string) (*UserWithToken, error)
	DeleteUserAccessTokenByToken(AccessToken string) (*UserAccessToken, error)
	AddUserAccessToken(UserID uint, AccessToken string) (*UserAccessToken, error)
}

type DeviceInstallCallback struct {
	gorm.Model
	DeviceID     uint   `sql:"not null;"`
	CallbackType string `sql:"not null;"`
	Content      string `sql:"not null;"`
	RunTime      string
	RunResult    string
	RunStatus    string
}

type IDeviceInstallCallback interface {
	CountDeviceInstallCallbackByDeviceIDAndType(DeviceID uint, CallbackType string) (uint, error)
	GetDeviceInstallCallbackByWhere(Where string, Order string) ([]DeviceInstallCallback, error)
	GetDeviceInstallCallbackByDeviceIDAndType(DeviceID uint, CallbackType string) (*DeviceInstallCallback, error)
	DeleteDeviceInstallCallbackByID(Id uint) (*DeviceInstallCallback, error)
	DeleteDeviceInstallCallbackByDeviceID(DeviceID uint) (*DeviceInstallCallback, error)
	AddDeviceInstallCallback(DeviceID uint, CallbackType string, Content string, RunTime string, RunResult string, RunStatus string) (*DeviceInstallCallback, error)
	UpdateDeviceInstallCallbackByID(Id uint, DeviceID uint, CallbackType string, Content string, RunTime string, RunResult string, RunStatus string) (*DeviceInstallCallback, error)
	UpdateDeviceInstallCallbackRunInfoByID(Id uint, RunTime string, RunResult string, RunStatus string) (*DeviceInstallCallback, error)
}

type DhcpSubnet struct {
	gorm.Model
	StartIp string `sql:"not null;"`
	EndIp   string `sql:"not null;"`
	Gateway string `sql:"not null;"`
}

type IDhcpSubnet interface {
	CountDhcpSubnet() (uint, error)
	GetDhcpSubnetListWithPage(Limit uint, Offset uint) ([]DhcpSubnet, error)
	GetDhcpSubnetById(Id uint) (*DhcpSubnet, error)
	UpdateDhcpSubnetById(Id uint, StartIp string, EndIp string, Gateway string) (*DhcpSubnet, error)
	DeleteDhcpSubnetById(Id uint) (*DhcpSubnet, error)
	AddDhcpSubnet(StartIp string, EndIp string, Gateway string) (*DhcpSubnet, error)
}

type PlatformConfig struct {
	gorm.Model
	Name    string `sql:"not null;unique;"`
	Content string `sql:"type:longtext;"`
}

type IPlatformConfig interface {
	CountPlatformConfigByName(Name string) (uint, error)
	CountPlatformConfigByNameAndId(Name string, ID uint) (uint, error)
	CountPlatformConfig() (uint, error)
	GetPlatformConfigListWithPage(Limit uint, Offset uint) ([]PlatformConfig, error)
	GetPlatformConfigIdByName(Name string) (uint, error)
	GetPlatformConfigById(Id uint) (*PlatformConfig, error)
	UpdatePlatformConfigById(Id uint, Name string, Pxe string) (*PlatformConfig, error)
	DeletePlatformConfigById(Id uint) (*PlatformConfig, error)
	AddPlatformConfig(Name string, Content string) (*PlatformConfig, error)
	GetPlatformConfigByName(Name string) (*PlatformConfig, error)
}

type VmHost struct {
	gorm.Model
	Sn              string `sql:"not null;"`
	CpuSum          uint   `sql:"type:int(11);default:0;"`
	CpuUsed         uint   `sql:"type:int(11);default:0;"`
	CpuAvailable    uint   `sql:"type:int(11);default:0;"`
	MemorySum       uint   `sql:"type:int(11);default:0;"`
	MemoryUsed      uint   `sql:"type:int(11);default:0;"`
	MemoryAvailable uint   `sql:"type:int(11);default:0;"`
	DiskSum         uint   `sql:"type:int(11);default:0;"`
	DiskUsed        uint   `sql:"type:int(11);default:0;"`
	DiskAvailable   uint   `sql:"type:int(11);default:0;"`
	VmNum           uint   `sql:"type:int(11);default:0;"`
	IsAvailable     string `sql:"enum('Yes','No');NOT NULL;DEFAULT 'Yes'"`
	Remark          string `sql:"type:text"`
}

type VmHostFull struct {
	ID                uint
	DeviceID          uint
	Sn                string
	Hostname          string
	Ip                string
	ManageIp          string
	NetworkID         uint
	ManageNetworkID   uint
	OsID              uint
	HardwareID        uint
	SystemID          uint
	LocationID        uint
	AssetNumber       string
	Status            string
	NetworkName       string
	ManageNetworkName string
	OsName            string
	SystemName        string
	HardwareName      string
	IsSupportVm       string
	CpuSum            uint
	CpuUsed           uint
	CpuAvailable      uint
	MemorySum         uint
	MemoryUsed        uint
	MemoryAvailable   uint
	DiskSum           uint
	DiskUsed          uint
	DiskAvailable     uint
	VmNum             uint
	IsAvailable       string
	Remark            string
}

type IVmHost interface {
	CountVmHostBySn(Sn string) (uint, error)
	CountVmHost(Where string) (int, error)
	GetVmHostListWithPage(Limit uint, Offset uint, Where string) ([]VmHostFull, error)
	GetVmHostById(Id uint) (*VmHost, error)
	UpdateVmHostById(Id uint, CpuSum uint, CpuUsed uint, CpuAvailable uint, MemorySum uint, MemoryUsed uint, MemoryAvailable uint, DiskSum uint, DiskUsed uint, DiskAvailable uint, IsAvailable string, Remark string, VmNum uint) (*VmHost, error)
	UpdateVmHostCpuMemoryDiskVmNumById(Id uint, CpuSum uint, CpuUsed uint, CpuAvailable uint, MemorySum uint, MemoryUsed uint, MemoryAvailable uint, DiskSum uint, DiskUsed uint, DiskAvailable uint, VmNum uint, IsAvailable string) (*VmHost, error)
	DeleteVmHostById(Id uint) (*VmHost, error)
	DeleteVmHostBySn(Sn string) (*VmHost, error)
	AddVmHost(Sn string, CpuSum uint, CpuUsed uint, CpuAvailable uint, MemorySum uint, MemoryUsed uint, MemoryAvailable uint, DiskSum uint, DiskUsed uint, DiskAvailable uint, IsAvailable string, Remark string, VmNum uint) (*VmHost, error)
	GetVmHostBySn(Sn string) (*VmHost, error)
	GetCpuUsedSum(Where string) (uint, error)
	GetMemoryUsedSum(Where string) (uint, error)
	GetDiskUsedSum(Where string) (uint, error)
	CountVmDeviceByDeviceId(DeviceID uint) (uint, error)
	GetMaxVncPort(Where string) (uint, error)
	GetNeedCollectDeviceForVmHost(DeviceID uint) ([]Device, error)
	DeleteVmInfoByDeviceSn(Sn string) error
}

type VmDeviceLog struct {
	gorm.Model
	DeviceID  uint   `sql:"not null;"`
	Title     string `sql:"not null;"`
	Type      string `sql:"not null;default:'install';"`
	Content   string `sql:"type:text;"` //pxe信息
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IVmDeviceLog interface {
	CountVmDeviceLogByDeviceID(DeviceID uint) (uint, error)
	CountVmDeviceLogByDeviceIDAndType(DeviceID uint, Type string) (uint, error)
	CountVmDeviceLog() (uint, error)
	GetVmDeviceLogListByDeviceID(DeviceID uint, Order string) ([]VmDeviceLog, error)
	GetLastVmDeviceLogByDeviceID(DeviceID uint) (VmDeviceLog, error)
	GetVmDeviceLogListByDeviceIDAndType(DeviceID uint, Type string, Order string, MaxID uint) ([]VmDeviceLog, error)
	GetVmDeviceLogById(Id uint) (*VmDeviceLog, error)
	DeleteVmDeviceLogById(Id uint) (*VmDeviceLog, error)
	DeleteVmDeviceLogByDeviceIDAndType(DeviceID uint, Type string) (*VmDeviceLog, error)
	DeleteVmDeviceLogByDeviceID(DeviceID uint) (*VmDeviceLog, error)
	AddVmDeviceLog(DeviceID uint, Title string, Type string, Content string) (*VmDeviceLog, error)
	UpdateVmDeviceLogTypeByDeviceIdAndType(deviceID uint, Type string, NewType string) ([]VmDeviceLog, error)
}
