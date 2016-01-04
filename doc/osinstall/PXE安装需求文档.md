# PXE安装需求文档

> 面对大规模的系统部署，PXE安装操作系统是一种十分高效的部署方式，同时也为标准化运维奠定了基础。下面针对PXE安装的细节展开一些讨论。

## 什么是PXE

PXE(preboot execute environment，预启动执行环境)是由Intel公司开发的最新技术，工作于Client/Server的网络模式，支持工作站通过网络从远端服务器下载映像，并由此支持通过网络启动操作系统，在启动过程中，通过DHCP服务器分配IP地址，再用TFTP协议下载一个启动软件到本机内存中执行。

## PXE工作原理

![image](http://linux.vbird.org/linux_enterprise/0120installation/pxe.jpg)

![image](http://diddy.boot-land.net/pxe/files/img/pxeboot1.jpg)

## 利用PXE安装操作系统

- 启动BootOS环境

```
DEFAULT bootos
TIMEOUT 100

LABEL bootos
    KERNEL bootos/vmlinuz
    APPEND initrd=bootos/initrd.img console=tty0 console=ttyS0,115200n8
```

- 收集mac地址生成pxe安装配置

```
DEFAULT centos6
TIMEOUT 100

LABEL centos6
    KERNEL centos6/vmlinuz
    APPEND initrd=centos6/initrd.img ks=nfs:10.0.1.1:/srv/centos6/centos6.ks ip=dhcp console=tty0 console=ttyS0,115200n8
```

- 生成kickstart文件

```
install
nfs --server=10.0.1.1 --dir=/srv/centos6
lang en_US.UTF-8
keyboard us
network --onboot yes --device eth0 --bootproto dhcp --noipv6
rootpw  --iscrypted $6$7V7NyfAKwWTjydPj$ay7WH0Oj6ew4xCwQRYOE/NNgHDIizFye0pswgC1KnXqT34KWdD0r1Zm4ZFDcj.U4yo1NFW8rqH96M4CwRBkWl1
firewall --disabled
authconfig --enableshadow --passalgo=sha512
selinux --disabled
timezone Asia/Shanghai
reboot
text
zerombr
bootloader --location=mbr
clearpart --all --initlabel
part /boot --fstype=ext4 --size=256 --ondisk=sda
part / --fstype=ext4 --size=51200 --ondisk=sda
part swap --size=2048 --ondisk=sda
part /home --fstype=ext4 --size=100 --grow --ondisk=sda
%packages
@base
@core
@development
%end
```

## 平台设计

![image](http://www.ibm.com/developerworks/cn/linux/l-cobbler/fig-1.png)

功能模块：

- 系统镜像：明确各个发行版以及对应版本的支持，如CentOS，RedHat，SUSE等，需要能够通过前台导入并管理每一套系统安装镜像。
- 安装模板：对应每个发行版提供一套自动安装模板，前台提供模板自动生成，关键配置可以由用户自定义。
- 软件仓库：提供并管理对应各个版本的软件仓库，方便用户安装软件包，方便快速部署应用和升级安全漏洞。
- 装机配置：系统安装所需基础配置包括主机配置，网络配置等。系统安装好以后，提供接口方便用户自定义配置，可以通过脚本或者配管工具的方式进行。
- 权限管理：验证和识别不同用户，区分权限，防止误操作重新系统。提供装机日志查看审计功能。

## 需求

1. 依赖BootOS内置的agent收集网卡mac地址，服务端通过收集的mac地址创建pxe安装文件。
2. 若要配置硬件，需在BootOS阶段通过服务端发送命令控制agent完成，重启后进入PXE安装。
3. kickstart模块化配置，支持用户自定义需求，如分区，安装包，执行post初始化脚本等。
4. 前台界面可以提供安装进度，安装界面或者安装日志，方便出现问题的时候进行排查分析。