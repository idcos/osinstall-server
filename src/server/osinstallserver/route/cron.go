package route

import (
	"config"
	"encoding/json"
	"github.com/jakecoffman/cron"
	"logger"
	"model"
	"strings"
)

func CloudBootCron(conf *config.Config, logger logger.Logger, repo model.Repo) {
	c := cron.New()
	//install timeout process
	c.AddFunc("0 */5 * * * *", func() {
		InstallTimeoutProcess(conf, logger, repo)
	}, "InstallTimeoutProcessTask")
	//init bootos ip for old data
	c.AddFunc("0 */30 * * * *", func() {
		InitBootOSIPForScanDeviceListProcess(logger, repo)
	}, "InitBootOSIPForScanDeviceListProcessTask")
	//start
	c.Start()
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
