package main

import (
	"os"

	"github.com/urfave/cli"

	"idcos.io/osinstall/agent"
	"idcos.io/osinstall/build"
)

func main() {
	app := cli.NewApp()
	app.Name = "cloudboot-agent"
	app.Description = "cloudboot agent"
	app.Version = build.Version()

	app.Action = func(c *cli.Context) {
		runAgent(c)
	}

	app.Run(os.Args)
}

func runAgent(c *cli.Context) {
	var agent, err = agent.New()
	if err != nil {
		agent.Logger.Error(err)
		return
	}

	//run pre install script
	agent.RunPreInstallScript()

	if err = agent.ReportProductInfo(); err != nil {
		agent.Logger.Error(err)
	}

	if agent.Sn == "" {
		agent.Logger.Error("SN error:SN can not be empty!")
		return
	}

	agent.ReportProgress(0.1, "进入bootos", "正常进入bootos")
	for {
		// 状态查询（是否在装机队列中）
		agent.IsInInstallQueue()
		agent.Logger.Debug("into install queue")

		// 判断IP地址是否在使用
		if err = agent.IsIpInUse(); err != nil {
			agent.ReportProgress(-1, "IP查询错误", err.Error())
			continue
		}

		// 配置查询（15%）
		var isSkip = false
		isSkip, err = agent.IsHaveHardWareConf()
		if !isSkip {
			if err != nil {
				agent.ReportProgress(-1, "配置查询失败", "该硬件型号不存在，请打开开发者模式再尝试，错误信息："+err.Error())
				continue
			} else {
				agent.ReportProgress(0.15, "配置查询", "存在对应的硬件配置")
			}

			// 获取硬件配置模板（20%）
			if err = agent.GetHardWareConf(); err != nil {
				agent.ReportProgress(-1, "获取硬件配置模板失败", "没有对应的硬件配置模板，错误信息："+err.Error())
				continue
			} else {
				agent.ReportProgress(0.2, "获取硬件配置", "存在对应的硬件配置模板")
			}

			// 硬件初始化（30%~40%）
			if err = agent.ImplementHardConf(); err != nil {
				agent.ReportProgress(-1, "初始化硬件失败", "无法初始化硬件，错误信息："+err.Error())
				continue
			}
		}

		// 生成 PXE文件（45%）
		if err = agent.ReportMacInfo(); err != nil {
			agent.ReportProgress(-1, "生成PXE文件失败", "无法生成PXE文件，错误信息："+err.Error())
			continue
		} else {
			agent.ReportProgress(0.45, "生成PXE文件", "正常生成PXE文件")
		}

		//run post install script
		agent.RunPostInstallScript()

		// 重启系统（50%）
		agent.ReportProgress(0.5, "系统开始重启", "系统重启中... ...")
		if err = agent.Reboot(); err != nil {
			agent.ReportProgress(-1, "系统重启失败", "重启系统出错，错误信息："+err.Error())
			continue
		} else {
			break // 退出 agent
		}
	}
}
