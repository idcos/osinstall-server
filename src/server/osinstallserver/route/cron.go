package route

import (
	"config"
	"github.com/jakecoffman/cron"
	"logger"
	"model"
)

func CloudBootCron(conf *config.Config, logger logger.Logger, repo model.Repo) {
	c := cron.New()
	//install timeout process
	c.AddFunc("0 */5 * * * *", func() {
		InstallTimeoutProcess(conf, logger, repo)
	}, "InstallTimeoutProcessTask")
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
