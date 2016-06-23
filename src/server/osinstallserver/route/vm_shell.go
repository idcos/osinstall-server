package route

import (
	"config"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"logger"
	"middleware"
	"model"
	"os"
	"server/osinstallserver/util"
	"strings"
)

func RunCreateVol(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	storage := conf.Vm.Storage
	if storage == "" {
		storage = "guest_images_lvm"
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system vol-create-as %s %s %dG`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		storage,
		vmDevice.Hostname,
		vmDevice.DiskSize)
	logger.Debugf("vm create vol:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("create result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("create error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

func RunCreateVm(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	var cmdFormat = `LANG=C virt-install --connect qemu+ssh://root@%s/system \
--name=%s \
--os-type=windows \
--vcpus=%d \
--ram=%d \
--hvm \
--boot hd,network,menu=on \
--accelerate \
--graphics vnc,listen=0.0.0.0,port=%s \
--noautoconsole \
--autostart \
--network bridge=br0,model=virtio,mac=%s \
--disk path=/dev/VolGroup0/%s,device=disk,bus=virtio,cache=none,sparse=false,format=raw`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname,
		vmDevice.CpuCoresNumber,
		vmDevice.MemoryCurrent,
		vmDevice.VncPort,
		vmDevice.Mac,
		vmDevice.Hostname)
	logger.Debugf("create vm:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("create result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("create error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

func RunDestroyVol(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	conf, ok := middleware.ConfigFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	storage := conf.Vm.Storage
	if storage == "" {
		storage = "guest_images_lvm"
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system vol-delete %s --pool %s`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname,
		storage)
	logger.Debugf("vm destroy vol:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("destroy result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("destroy error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

func RunDestroyVm(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system 'destroy %s; undefine %s'`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname,
		vmDevice.Hostname)
	logger.Debugf("destroy vm:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("destroy result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("destroy error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

//start vm
func RunStartVm(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system start %s`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname)
	logger.Debugf("start vm:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("start result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("start error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

//stop vm
func RunStopVm(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system destroy %s`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname)
	logger.Debugf("stop vm:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("stop result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("stop error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

//restart vm
func RunReStartVm(ctx context.Context, vmDeviceId uint) error {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system 'destroy %s; start %s'`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname,
		vmDevice.Hostname)
	logger.Debugf("restart vm:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("restart result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("restart error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return errors.New(runResult)
	}
	return nil
}

//get vm info
func RunGetVmInfo(repo model.Repo, logger logger.Logger, vmDeviceId uint) (string, error) {
	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return "", err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return "", err
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system dominfo %s`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		vmDevice.Hostname)
	logger.Debugf("get vm info:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return "", errors.New(runResult)
	}
	return string(bytes), nil
}

//get vm host info
func RunGetVmHostInfo(repo model.Repo, logger logger.Logger, deviceId uint) (string, error) {
	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		return "", err
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system nodeinfo`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip)
	logger.Debugf("get vm host info:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return "", errors.New(runResult)
	}
	return string(bytes), nil
}

//get vm pool info
func RunGetVmHostPoolInfo(repo model.Repo, logger logger.Logger, conf *config.Config, deviceId uint) (string, error) {
	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		return "", err
	}

	storage := conf.Vm.Storage
	if storage == "" {
		storage = "guest_images_lvm"
	}

	var cmdFormat = `LANG=C virsh --connect qemu+ssh://root@%s/system pool-info %s`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip,
		storage)
	logger.Debugf("get vm host pool info:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return "", errors.New(runResult)
	}
	return string(bytes), nil
}

//test connect vm host
func RunTestConnectVmHost(repo model.Repo, logger logger.Logger, deviceId uint) (string, error) {
	device, err := repo.GetDeviceById(deviceId)
	if err != nil {
		return "", err
	}

	var cmdFormat = `LANG=C ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o BatchMode=yes -o ConnectTimeout=3 root@%s 'w'`
	var cmd = fmt.Sprintf(cmdFormat,
		device.Ip)
	logger.Debugf("test connect vm host:%s", cmd)
	var runResult = "执行脚本:\n" + cmd
	bytes, err := util.ExecScript(cmd)
	logger.Debugf("result:%s", string(bytes))
	runResult += "\n\n" + "执行结果:\n" + string(bytes)
	if err != nil {
		logger.Errorf("error:%s", err.Error())
		runResult += "\n\n" + "错误信息:\n" + err.Error()
		return "", errors.New(runResult)
	}
	return string(bytes), nil
}

//create noVNC token file
func RunCreateVmNoVncTokenFile(repo model.Repo, logger logger.Logger, vmDeviceId uint) error {
	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	device, err := repo.GetDeviceById(vmDevice.DeviceID)
	if err != nil {
		return err
	}

	var newHostname = strings.Replace(vmDevice.Hostname, ":", "_", -1)
	var strFormat = `%s: %s:%s`
	var str = fmt.Sprintf(strFormat,
		newHostname,
		device.Ip,
		vmDevice.VncPort)

	var dir = "/etc/cloudboot-server/novnc-tokens"
	var file = vmDevice.Hostname

	if !util.FileExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	//文件已存在,先删除
	if util.FileExist(dir + "/" + file) {
		err := os.Remove(dir + "/" + file)
		if err != nil {
			return err
		}
	}
	logger.Debugf("create vm novnc token file %s:%s", dir+"/"+file, str)
	bytes := []byte(str)
	errCreate := ioutil.WriteFile(dir+"/"+file, bytes, 0644)
	if errCreate != nil {
		logger.Errorf("error:%s", errCreate.Error())
		return errCreate
	}
	return nil
}

//delete noVNC token file
func RunDeleteVmNoVncTokenFile(repo model.Repo, logger logger.Logger, vmDeviceId uint) error {
	vmDevice, err := repo.GetVmDeviceById(vmDeviceId)
	if err != nil {
		return err
	}

	var dir = "/etc/cloudboot-server/novnc-tokens"
	var file = vmDevice.Hostname
	logger.Debugf("delete vm novnc token file:%s", dir+"/"+file)
	if util.FileExist(dir + "/" + file) {
		err := os.Remove(dir + "/" + file)
		if err != nil {
			logger.Errorf("error:%s", err.Error())
			return err
		}
	}
	return nil
}
