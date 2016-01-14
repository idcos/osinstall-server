package main

import "osinstall/agent"

func main() {
	var agent, err = agent.New()
	if err != nil {
		agent.Logger.Error(err.Error())
		return
	}
	agent.ReportProgress(0.1, "进入bootos", "正常进入bootos")
	for {
		// 状态查询（是否在装机队列中）
		if agent.IsInInstallQueue() == false {
			continue
		}
		// 判断IP地址是否在使用
		if agent.IsIpInUse() == false {
			agent.ReportProgress(-1, "IP查询错误", "IP地址冲突")
			continue
		}
		// 配置查询（15%）
		if agent.HaveHardWareConf() == false {
			agent.ReportProgress(-1, "配置查询失败", "该硬件型号不存在，请打开开发者模式再尝试")
			continue
		} else {
			agent.ReportProgress(0.15, "配置查询", "存在对应的硬件配置")
		}
		// 获取硬件配置模板（20%）
		if agent.GetHardWareConf() == false {
			agent.ReportProgress(-1, "获取硬件配置模板失败", "没有对应的硬件配置模板")
			continue
		} else {
			agent.ReportProgress(0.2, "获取硬件配置", "存在对应的硬件配置模板")
		}
		// 硬件初始化（30%~40%）
		if agent.ImplementHardConf() == false {
			agent.ReportProgress(-1, "初始化硬件失败", "无法初始化硬件")
			continue
		}
		// 生成 PXE文件（45%）
		if agent.ReportMacInfo() == false {
			agent.ReportProgress(-1, "生成PXE文件失败", "无法生成PXE文件")
			continue
		} else {
			agent.ReportProgress(0.45, "生成PXE文件", "正常生成PXE文件")
		}
		// 重启系统（50%）
		agent.ReportProgress(0.5, "系统开始重启", "系统重启中... ...")
		if agent.Reboot() == false {
			agent.ReportProgress(-1, "系统重启失败", "重启系统出错")
			continue
		} else {
			break // 退出 agent
		}

	}
}
