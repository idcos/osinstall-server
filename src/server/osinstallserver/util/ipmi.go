package util

import (
	"errors"
	"fmt"
	"logger"
	"model"
	"regexp"
	"strings"
)

func GetDeviceSnFromIpmitool(repo model.Repo, logger logger.Logger, ip string, username string, password string) (string, error) {
	if ip == "" || username == "" || password == "" {
		return "", errors.New("OOB IP、用户名、密码不能为空")
	}

	cmd := fmt.Sprintf(`ipmitool -I lanplus -H %s -U %s -P %s fru list 0`, ip, username, password)
	//cmd = fmt.Sprintf(`/usr/local/Cellar/ipmitool/1.8.13/bin/ipmitool -I lanplus -H %s -U %s -P %s fru list 0`, ip, username, password)
	logger.Debug(cmd)
	bytes, err := ExecScript(cmd)
	logger.Debug(string(bytes))
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return "", fmt.Errorf("校验设备实际SN失败: %s", err.Error())
	}

	reg, _ := regexp.Compile("Product Serial(\\s+):(.*)\n")
	matchs := reg.FindStringSubmatch(string(bytes))
	if len(matchs) != 3 {
		return "", errors.New("校验设备实际SN失败")
	}
	sn := strings.TrimSpace(matchs[2])
	return sn, nil
}

func GetDevicePowerStatusFromIpmitool(repo model.Repo, logger logger.Logger, ip string, username string, password string) (string, error) {
	if ip == "" || username == "" || password == "" {
		return "", errors.New("OOB IP、用户名、密码不能为空")
	}

	cmd := fmt.Sprintf(`ipmitool -I lanplus -H %s -U %s -P %s power status`, ip, username, password)
	//cmd = fmt.Sprintf(`/usr/local/Cellar/ipmitool/1.8.13/bin/ipmitool -I lanplus -H %s -U %s -P %s power status`, ip, username, password)
	logger.Debug(cmd)
	bytes, err := ExecScript(cmd)
	logger.Debug(string(bytes))
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return "", errors.New("通过IPMI获取设备电源状态失败" + ":" + err.Error())
	}

	reg, _ := regexp.Compile("Chassis Power is(\\s)(\\w+)")
	matchs := reg.FindStringSubmatch(string(bytes))
	if len(matchs) != 3 {
		return "", errors.New("通过IPMI获取设备电源状态失败")
	}
	result := strings.TrimSpace(matchs[2])
	if result != "on" && result != "off" {
		result = ""
	}
	return result, nil
}

func RestartDeviceFromIpmitool(repo model.Repo, logger logger.Logger, ip string, username string, password string) error {
	if ip == "" || username == "" || password == "" {
		return errors.New("OOB IP、用户名、密码不能为空")
	}

	cmd := fmt.Sprintf(`ipmitool -I lanplus -H %s -U %s -P %s power reset`, ip, username, password)
	//cmd = fmt.Sprintf(`/usr/local/Cellar/ipmitool/1.8.13/bin/ipmitool -I lanplus -H %s -U %s -P %s power reset`, ip, username, password)
	logger.Debug(cmd)
	bytes, err := ExecScript(cmd)
	logger.Debug(string(bytes))
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		return fmt.Errorf("通过IPMI重启失败: %s", err.Error())
	}
	return nil
}
