package route

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/jakecoffman/cron"
	"logger"
	"model"
	"regexp"
	"server/osinstallserver/util"
	"strconv"
	"strings"
)

func CloudBootCron(conf *config.Config, logger logger.Logger, repo model.Repo) {
	//db version update(v1.2.1 to v1.3)
	DBVersionUpdate(logger, repo)

	c := cron.New()
	//install timeout process
	c.AddFunc("0 */5 * * * *", func() {
		InstallTimeoutProcess(conf, logger, repo)
	}, "InstallTimeoutProcessTask")
	//init bootos ip for old data
	c.AddFunc("0 */30 * * * *", func() {
		InitBootOSIPForScanDeviceListProcess(logger, repo)
	}, "InitBootOSIPForScanDeviceListProcessTask")
	//update vm host resource
	c.AddFunc("0 1 1 * * *", func() {
		UpdateVmHostResource(logger, repo, conf, 0)
	}, "UpdateVmHostResourceTask")
	//start
	c.Start()
}

func DBVersionUpdate(logger logger.Logger, repo model.Repo) {
	_, err := repo.CountVmHost("")
	if err == nil {
		//logger.Info("db version has been upgraded to v1.3")
		return
	}

	var str string
	str = `alter table cloudboot.system_configs modify column ~~content~~ LONGTEXT;
alter table cloudboot.os_configs modify column ~~pxe~~ LONGTEXT;
alter table cloudboot.hardwares modify column ~~tpl~~ LONGTEXT;
alter table cloudboot.hardwares modify column ~~data~~ LONGTEXT;
alter table cloudboot.device_logs modify column ~~content~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~nic~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~cpu~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~memory~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~disk~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~motherboard~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~raid~~ LONGTEXT;
alter table cloudboot.manufacturers modify column ~~oob~~ LONGTEXT;

ALTER TABLE ~~manufacturers~~ ADD ~~is_vm~~ ENUM('Yes','No') NOT NULL DEFAULT 'No' ;
ALTER TABLE ~~manufacturers~~ ADD ~~is_show_in_scan_list~~ ENUM('Yes','No') NOT NULL DEFAULT 'Yes' ;
ALTER TABLE ~~manufacturers~~ ADD ~~nic_device~~ longtext NULL DEFAULT NULL;
ALTER TABLE ~~devices~~ CHANGE ~~is_support_vm~~ ~~is_support_vm~~ ENUM('Yes','No') CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'No' COMMENT '是否支持安装虚拟机';
UPDATE ~~devices~~ set ~~is_support_vm~~ = 'No';

CREATE TABLE ~~dhcp_subnets~~ (
  ~~id~~ int(10) unsigned NOT NULL AUTO_INCREMENT,
  ~~created_at~~ timestamp NULL DEFAULT NULL,
  ~~updated_at~~ timestamp NULL DEFAULT NULL,
  ~~deleted_at~~ timestamp NULL DEFAULT NULL,
  ~~start_ip~~ varchar(255) NOT NULL,
  ~~end_ip~~ varchar(255) NOT NULL,
  ~~gateway~~ varchar(255) NOT NULL,
  PRIMARY KEY (~~id~~)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

CREATE TABLE ~~platform_configs~~ (
  ~~id~~ int(10) unsigned NOT NULL AUTO_INCREMENT,
  ~~created_at~~ timestamp NULL DEFAULT NULL,
  ~~updated_at~~ timestamp NULL DEFAULT NULL,
  ~~deleted_at~~ timestamp NULL DEFAULT NULL,
  ~~name~~ varchar(255) NOT NULL,
  ~~content~~ longtext NULL DEFAULT NULL,
  PRIMARY KEY (~~id~~),
  UNIQUE KEY ~~name~~ (~~name~~)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
insert into platform_configs(~~created_at~~,~~updated_at~~,~~name~~,~~content~~) values('2016-05-30 17:32:47','2016-05-30 17:32:47','IsShowVmFunction','');

CREATE TABLE ~~vm_hosts~~ (
  ~~id~~ int(11) unsigned NOT NULL AUTO_INCREMENT,
  ~~created_at~~ timestamp NULL DEFAULT NULL,
  ~~updated_at~~ timestamp NULL DEFAULT NULL,
  ~~deleted_at~~ timestamp NULL DEFAULT NULL,
  ~~sn~~ varchar(255) DEFAULT NULL COMMENT '序列号',
  ~~cpu_sum~~ int(11) DEFAULT '0' COMMENT 'CPU总核数',
  ~~cpu_used~~ int(11) DEFAULT '0' COMMENT '已用CPU核数',
  ~~cpu_available~~ int(11) DEFAULT '0' COMMENT '可用CPU核数',
  ~~memory_sum~~ int(11) DEFAULT '0' COMMENT '内存总容量',
  ~~memory_used~~ int(11) DEFAULT '0' COMMENT '已用内存',
  ~~memory_available~~ int(11) DEFAULT '0' COMMENT '可用内存',
  ~~disk_sum~~ int(11) DEFAULT '0' COMMENT '磁盘总容量',
  ~~disk_used~~ int(11) DEFAULT '0' COMMENT '已用磁盘空间',
  ~~disk_available~~ int(11) DEFAULT '0' COMMENT '可用磁盘空间',
  ~~vm_num~~ int(11) DEFAULT '0' COMMENT '虚拟机数量',
  ~~is_available~~ ENUM('Yes','No') DEFAULT 'Yes' COMMENT '是否可用',
  ~~remark~~ text DEFAULT NULL COMMENT '备注(不可用的原因，等等)',
  PRIMARY KEY (~~id~~),
  UNIQUE KEY ~~sn~~ (~~sn~~)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

ALTER TABLE ~~vm_devices~~ ADD ~~user_id~~ INT(11) NOT NULL DEFAULT '0' ;
ALTER TABLE ~~vm_devices~~ ADD INDEX(~~user_id~~);
ALTER TABLE ~~vm_devices~~ ADD ~~system_id~~ INT(11) NOT NULL DEFAULT '0' COMMENT '系统安装模板ID';
ALTER TABLE ~~vm_devices~~ ADD ~~install_progress~~ decimal(11,4) NOT NULL DEFAULT '0.0000' COMMENT '安装进度';
ALTER TABLE ~~vm_devices~~ ADD ~~vnc_port~~ VARCHAR(256) NULL DEFAULT NULL COMMENT 'VNC端口';
ALTER TABLE ~~vm_devices~~ ADD ~~run_status~~ VARCHAR(256) NULL DEFAULT NULL COMMENT '运行状态';

CREATE TABLE ~~vm_device_logs~~ (
  ~~id~~ int(10) unsigned NOT NULL AUTO_INCREMENT,
  ~~created_at~~ timestamp NULL DEFAULT NULL,
  ~~updated_at~~ timestamp NULL DEFAULT NULL,
  ~~deleted_at~~ timestamp NULL DEFAULT NULL,
  ~~device_id~~ int(11) NOT NULL,
  ~~title~~ varchar(255) NOT NULL,
  ~~type~~ varchar(255) NOT NULL,
  ~~content~~ longtext,
  PRIMARY KEY (~~id~~),
  KEY ~~device_id~~ (~~device_id~~) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;`
	str = strings.Replace(str, "~~", "`", -1)
	sqls := strings.Split(str, ";")
	for _, sql := range sqls {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		//logger.Debug(sql)

		err := repo.ExecDBVersionUpdateSql(sql)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	logger.Info("db version has been successfully upgraded to v1.3")
}

func InstallTimeoutProcess(conf *config.Config, logger logger.Logger, repo model.Repo) {
	devices, err := repo.GetInstallTimeoutDeviceList(conf.Cron.InstallTimeout)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return
	}

	if len(devices) <= 0 {
		return
	}

	logger.Infof("install timeout config:%d", conf.Cron.InstallTimeout)
	if conf.Cron.InstallTimeout <= 0 {
		logger.Info("install timeout is not configured, don't do timeout processing")
		return
	}

	for _, device := range devices {
		isTimeout, err := repo.IsInstallTimeoutDevice(conf.Cron.InstallTimeout, device.ID)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		if !isTimeout {
			logger.Infof("the device is not timeout(SN:%s)", device.Sn)
			continue
		}

		_, errUpdate := repo.UpdateInstallInfoById(device.ID, "failure", -1)
		if errUpdate != nil {
			logger.Errorf("error:%s", errUpdate.Error())
			continue
		}

		logTitle := "安装失败(安装超时)"
		installLog := "安装超时"

		_, errAddLog := repo.AddDeviceLog(device.ID, logTitle, "install", installLog)
		if errAddLog != nil {
			logger.Errorf("error:%s", errAddLog.Error())
			continue
		}

		logger.Infof("the device timeout process success:(SN:%s)", device.Sn)
	}
	logger.Info("install timeout processing end")
	return
}

func InitBootOSIPForScanDeviceListProcess(logger logger.Logger, repo model.Repo) {
	devices, err := repo.GetManufacturerListWithPage(1000000, 0, " and (ip = '' or ip is null)")
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return
	}

	if len(devices) <= 0 {
		return
	}

	type NicInfo struct {
		Name string
		Mac  string
		Ip   string
	}

	for _, device := range devices {
		manufacturer, err := repo.GetManufacturerById(device.ID)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		if manufacturer.Ip != "" {
			continue
		}
		if manufacturer.Nic == "" {
			continue
		}
		var NicInfos []NicInfo
		errJson := json.Unmarshal([]byte(manufacturer.Nic), &NicInfos)
		if errJson != nil {
			logger.Errorf("error:%s", errJson.Error())
			continue
		}

		var ip string
		for _, nicInfo := range NicInfos {
			nicInfo.Ip = strings.TrimSpace(nicInfo.Ip)
			if nicInfo.Ip != "" {
				ip = nicInfo.Ip
				break
			}
		}
		if ip == "" {
			continue
		}
		_, errUpdate := repo.UpdateManufacturerIPById(manufacturer.ID, ip)
		if errUpdate != nil {
			logger.Errorf("error:%s", errUpdate.Error())
			continue
		}
		logger.Infof("the bootos ip init process success:(SN:%s,IP:%s)", manufacturer.Sn, ip)
	}
	logger.Info("bootos ip init processing end")
	return
}

func UpdateVmHostResource(logger logger.Logger, repo model.Repo, conf *config.Config, deviceId uint) {
	devices, err := repo.GetNeedCollectDeviceForVmHost(deviceId)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return
	}
	if len(devices) <= 0 {
		return
	}
	logger.Info("update vm host resource info")
	for _, device := range devices {
		var logTitle string
		var installLog string
		var cpuSum int
		var memorySum int
		var diskSum int
		var isAvailable = "Yes"

		_, err := RunTestConnectVmHost(repo, logger, device.ID)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			logTitle = "宿主机信息采集失败(无法SSH)"
			installLog = err.Error()
			_, errAddLog := repo.AddDeviceLog(device.ID, logTitle, "virtualization", installLog)
			if errAddLog != nil {
				logger.Errorf("error:%s", errAddLog.Error())
			}
			isAvailable = "No"
		} else {
			text, err := RunGetVmHostInfo(repo, logger, device.ID)
			if err != nil {
				logger.Errorf("error:%s", err.Error())
				isAvailable = "No"

				logTitle = "宿主机信息采集失败"
				installLog = err.Error()
				_, errAddLog := repo.AddDeviceLog(device.ID, logTitle, "virtualization", installLog)
				if errAddLog != nil {
					logger.Errorf("error:%s", errAddLog.Error())
				}
			} else {
				//cpu
				reg, _ := regexp.Compile("CPU\\(s\\):(\\s+)([\\d]+)\n")
				matchs := reg.FindStringSubmatch(text)
				cpuSum, err = strconv.Atoi(matchs[2])
				if err != nil {
					logger.Errorf("error:%s", err.Error())
				}
				//memory
				reg, _ = regexp.Compile("Memory size:(\\s+)([\\d|.]+)(\\s+)([KiB|MiB|GiB|TiB]+)")
				matchs = reg.FindStringSubmatch(text)
				float, err := strconv.ParseFloat(matchs[2], 64)
				if err != nil {
					logger.Errorf("error:%s", err.Error())
				}
				memorySum = util.FotmatNumberToMB(float, matchs[4])
			}
			//disk
			text, err = RunGetVmHostPoolInfo(repo, logger, conf, device.ID)
			reg, _ := regexp.Compile("Capacity:(\\s+)([\\d|.]+)(\\s+)([KiB|MiB|GiB|TiB]+)")
			matchs := reg.FindStringSubmatch(text)
			float, err := strconv.ParseFloat(matchs[2], 64)
			if err != nil {
				logger.Errorf("error:%s", err.Error())
				isAvailable = "No"

				logTitle = "宿主机信息采集失败"
				installLog = err.Error()
				_, errAddLog := repo.AddDeviceLog(device.ID, logTitle, "virtualization", installLog)
				if errAddLog != nil {
					logger.Errorf("error:%s", errAddLog.Error())
				}
			}
			diskSum = util.FotmatNumberToGB(float, matchs[4])
		}

		//update resource
		var infoHost model.VmHost
		where := fmt.Sprintf("device_id = %d", device.ID)
		count, err := repo.CountVmHostBySn(device.Sn)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		if count > 0 {
			vmHost, err := repo.GetVmHostBySn(device.Sn)
			if err != nil {
				logger.Errorf("error:%s", err.Error())
				continue
			}
			infoHost.ID = vmHost.ID
			infoHost.Sn = vmHost.Sn
			infoHost.CpuUsed = vmHost.CpuUsed
			infoHost.CpuAvailable = vmHost.CpuAvailable
			infoHost.MemoryUsed = vmHost.MemoryUsed
			infoHost.MemoryAvailable = vmHost.MemoryAvailable
			infoHost.DiskUsed = vmHost.DiskUsed
			infoHost.DiskAvailable = vmHost.DiskAvailable
			infoHost.VmNum = vmHost.VmNum
			infoHost.IsAvailable = isAvailable
			infoHost.Remark = vmHost.Remark
		} else {
			infoHost.Sn = device.Sn
			infoHost.CpuUsed = uint(0)
			infoHost.CpuAvailable = uint(0)
			infoHost.MemoryUsed = uint(0)
			infoHost.MemoryAvailable = uint(0)
			infoHost.DiskUsed = uint(0)
			infoHost.DiskAvailable = uint(0)
			infoHost.VmNum = uint(0)
			infoHost.IsAvailable = isAvailable
			infoHost.Remark = ""
		}
		infoHost.CpuSum = uint(cpuSum)
		infoHost.MemorySum = uint(memorySum)
		infoHost.DiskSum = uint(diskSum)
		//cpu update
		//cpu used sum
		infoHost.CpuUsed, err = repo.GetCpuUsedSum(where)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		cpuAvailable := int(infoHost.CpuSum - infoHost.CpuUsed)
		if cpuAvailable <= 0 {
			cpuAvailable = 0
		}
		infoHost.CpuAvailable = uint(cpuAvailable)
		//memory update
		infoHost.MemoryUsed, err = repo.GetMemoryUsedSum(where)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		memoryAvailable := int(infoHost.MemorySum - infoHost.MemoryUsed)
		if memoryAvailable <= 0 {
			memoryAvailable = 0
		}
		infoHost.MemoryAvailable = uint(memoryAvailable)
		//update disk
		infoHost.DiskUsed, err = repo.GetDiskUsedSum(where)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		diskAvailable := int(infoHost.DiskSum - infoHost.DiskUsed)
		if diskAvailable < 0 {
			diskAvailable = 0
		}
		infoHost.DiskAvailable = uint(diskAvailable)
		if infoHost.MemoryAvailable <= uint(0) || infoHost.DiskAvailable <= uint(0) {
			infoHost.IsAvailable = "No"
		}
		infoHost.VmNum, err = repo.CountVmDeviceByDeviceId(device.ID)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			continue
		}
		if count > 0 {
			//update host
			_, errUpdate := repo.UpdateVmHostCpuMemoryDiskVmNumById(infoHost.ID, infoHost.CpuSum, infoHost.CpuUsed, infoHost.CpuAvailable, infoHost.MemorySum, infoHost.MemoryUsed, infoHost.MemoryAvailable, infoHost.DiskSum, infoHost.DiskUsed, infoHost.DiskAvailable, infoHost.VmNum, infoHost.IsAvailable)
			if errUpdate != nil {
				logger.Errorf("error:%s", errUpdate.Error())
				continue
			}
		} else {
			_, err := repo.AddVmHost(infoHost.Sn, infoHost.CpuSum, infoHost.CpuUsed, infoHost.CpuAvailable, infoHost.MemorySum, infoHost.MemoryUsed, infoHost.MemoryAvailable, infoHost.DiskSum, infoHost.DiskUsed, infoHost.DiskAvailable, infoHost.IsAvailable, infoHost.Remark, infoHost.VmNum)
			if err != nil {
				logger.Errorf("error:%s", err.Error())
				continue
			}
		}
	}
	logger.Info("update vm host resource info end")
	return
}
