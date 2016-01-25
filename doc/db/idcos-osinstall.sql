-- MySQL dump 10.13  Distrib 5.5.42, for Linux (x86_64)
--
-- Host: localhost    Database: idcos-osinstall
-- ------------------------------------------------------
-- Server version	5.5.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `idcos-osinstall`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `idcos-osinstall` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `idcos-osinstall`;

--
-- Table structure for table `device_histories`
--

DROP TABLE IF EXISTS `device_histories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `device_histories` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `batch_number` varchar(255) NOT NULL,
  `sn` varchar(255) NOT NULL,
  `hostname` varchar(255) NOT NULL,
  `ip` varchar(255) NOT NULL,
  `network_id` int(10) unsigned NOT NULL,
  `os_id` int(10) unsigned NOT NULL,
  `hardware_id` int(10) unsigned DEFAULT NULL,
  `system_id` int(10) unsigned NOT NULL,
  `location` varchar(255) NOT NULL,
  `location_id` int(11) NOT NULL,
  `asset_number` varchar(255) DEFAULT NULL,
  `status` varchar(255) NOT NULL,
  `install_progress` decimal(11,4) DEFAULT '0.0000',
  `install_log` text,
  `is_support_vm` enum('Yes','No') DEFAULT 'Yes' COMMENT '是否支持安装虚拟机',
  PRIMARY KEY (`id`),
  KEY `batch_number` (`batch_number`),
  KEY `status` (`status`) USING BTREE,
  KEY `sn` (`sn`) USING BTREE,
  KEY `hostname` (`hostname`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `device_histories`
--

LOCK TABLES `device_histories` WRITE;
/*!40000 ALTER TABLE `device_histories` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_histories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `device_logs`
--

DROP TABLE IF EXISTS `device_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `device_logs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `device_id` int(11) NOT NULL,
  `title` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `content` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `device_id` (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `device_logs`
--

LOCK TABLES `device_logs` WRITE;
/*!40000 ALTER TABLE `device_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `devices`
--

DROP TABLE IF EXISTS `devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `devices` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `batch_number` varchar(255) NOT NULL,
  `sn` varchar(255) NOT NULL,
  `hostname` varchar(255) NOT NULL,
  `ip` varchar(255) NOT NULL,
  `network_id` int(10) unsigned NOT NULL,
  `os_id` int(10) unsigned NOT NULL,
  `hardware_id` int(10) unsigned DEFAULT NULL,
  `system_id` int(10) unsigned NOT NULL,
  `location` varchar(255) NOT NULL,
  `location_id` int(11) NOT NULL,
  `asset_number` varchar(255) DEFAULT NULL,
  `status` varchar(255) NOT NULL,
  `install_progress` decimal(11,4) DEFAULT '0.0000',
  `install_log` text,
  `is_support_vm` enum('Yes','No') DEFAULT 'Yes' COMMENT '是否支持安装虚拟机',
  PRIMARY KEY (`id`),
  UNIQUE KEY `sn` (`sn`),
  UNIQUE KEY `ip` (`ip`),
  KEY `batch_number` (`batch_number`),
  KEY `status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `devices`
--

LOCK TABLES `devices` WRITE;
/*!40000 ALTER TABLE `devices` DISABLE KEYS */;
/*!40000 ALTER TABLE `devices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hardwares`
--

DROP TABLE IF EXISTS `hardwares`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hardwares` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `company` varchar(255) NOT NULL,
  `product` varchar(255) NOT NULL,
  `model_name` varchar(255) NOT NULL,
  `raid` text,
  `oob` text,
  `bios` text,
  `is_system_add` enum('Yes','No') NOT NULL DEFAULT 'Yes' COMMENT '是否是系统添加的配置',
  `tpl` text COMMENT '厂商提交的JSON信息',
  `data` text COMMENT '最终要执行的信息',
  `source` varchar(255) DEFAULT NULL,
  `version` varchar(255) DEFAULT NULL,
  `status` enum('Pending','Success','Failure') DEFAULT 'Success',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hardwares`
--

LOCK TABLES `hardwares` WRITE;
/*!40000 ALTER TABLE `hardwares` DISABLE KEYS */;
INSERT INTO `hardwares` VALUES (1,'2015-11-20 11:41:50','2015-12-11 01:41:40',NULL,'Dell','PowerEdge','R420','/opt/yunji/osinstall/dell/raid/raid.sh','/opt/yunji/osinstall/dell/oob/oob.sh','/opt/yunji/osinstall/dell/bios/bios.sh','Yes','[{\"name\":\"RAID\",\"data\":[{\"name\":\"RAID\",\"type\":\"select\",\"data\":[{\"name\":\"RAID0\",\"value\":\"/opt/yunji/osinstall/dell/raid.sh -c -l 0\",\"checked\":false},{\"name\":\"RAID1\",\"value\":\"/opt/yunji/osinstall/dell/raid.sh -c -l 1\",\"checked\":false},{\"name\":\"RAID5\",\"value\":\"/opt/yunji/osinstall/dell/raid.sh -c -l 5\",\"checked\":false},{\"name\":\"RAID10\",\"value\":\"/opt/yunji/osinstall/dell/raid.sh -c -l 10\",\"checked\":true}],\"default\":\"/opt/yunji/osinstall/dell/raid.sh -c -l 10\"}]},{\"name\":\"OOB\",\"data\":[{\"name\":\"网络类型\",\"type\":\"select\",\"data\":[{\"name\":\"DHCP\",\"value\":\"/opt/yunji/osinstall/dell/oob.sh -n dhcp\",\"checked\":false},{\"name\":\"静态IP\",\"value\":\"/opt/yunji/osinstall/dell/oob.sh -n static\",\"checked\":true}],\"default\":\"/opt/yunji/osinstall/dell/oob.sh -n dhcp\"},{\"name\":\"用户名\",\"type\":\"input\",\"tpl\":\"/opt/yunji/osinstall/dell/oob.sh -u <{##}>\",\"default\":\"/opt/yunji/osinstall/dell/oob.sh -u root\",\"input\":\"root\"},{\"name\":\"密码\",\"type\":\"input\",\"tpl\":\"/opt/yunji/osinstall/dell/oob.sh -p <{##}>\",\"default\":\"/opt/yunji/osinstall/dell/oob.sh -p calvin\",\"input\":\"calvin\"}]},{\"name\":\"BIOS\",\"data\":[{\"name\":\"VT\",\"type\":\"select\",\"data\":[{\"name\":\"ON\",\"value\":\"/opt/yunji/osinstall/dell/bios.sh -t enable\",\"checked\":false},{\"name\":\"OFF\",\"value\":\"/opt/yunji/osinstall/dell/bios.sh -t disable\",\"checked\":true}],\"default\":\"/opt/yunji/osinstall/dell/bios.sh -t enable\"},{\"name\":\"C-States\",\"type\":\"select\",\"data\":[{\"name\":\"ON\",\"value\":\"/opt/yunji/osinstall/dell/bios.sh -c enable\",\"checked\":false},{\"name\":\"OFF\",\"value\":\"/opt/yunji/osinstall/dell/bios.sh -c disable\",\"checked\":true}],\"default\":\"/opt/yunji/osinstall/dell/bios.sh -c disable\"}]}]','[{\"Name\":\"RAID\",\"Data\":[{\"Name\":\"RAID\",\"Value\":\"/opt/yunji/osinstall/dell/raid.sh -c -l 10\"}]},{\"Name\":\"OOB\",\"Data\":[{\"Name\":\"网络类型\",\"Value\":\"/opt/yunji/osinstall/dell/oob.sh -n dhcp\"},{\"Name\":\"用户名\",\"Value\":\"/opt/yunji/osinstall/dell/oob.sh -u root\"},{\"Name\":\"密码\",\"Value\":\"/opt/yunji/osinstall/dell/oob.sh -p calvin\"}]},{\"Name\":\"BIOS\",\"Data\":[{\"Name\":\"VT\",\"Value\":\"/opt/yunji/osinstall/dell/bios.sh -t enable\"},{\"Name\":\"C-States\",\"Value\":\"/opt/yunji/osinstall/dell/bios.sh -c disable\"}]}]',NULL,NULL,'Success');
/*!40000 ALTER TABLE `hardwares` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ips`
--

DROP TABLE IF EXISTS `ips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ips` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `network_id` int(10) unsigned NOT NULL,
  `ip` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `network_id` (`network_id`) USING BTREE,
  KEY `ip` (`ip`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ips`
--

LOCK TABLES `ips` WRITE;
/*!40000 ALTER TABLE `ips` DISABLE KEYS */;
/*!40000 ALTER TABLE `ips` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `locations`
--

DROP TABLE IF EXISTS `locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `locations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `pid` int(10) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `pid` (`pid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locations`
--

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
/*!40000 ALTER TABLE `locations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `macs`
--

DROP TABLE IF EXISTS `macs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `macs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `device_id` int(11) unsigned NOT NULL,
  `mac` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mac` (`mac`),
  KEY `device_id` (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `macs`
--

LOCK TABLES `macs` WRITE;
/*!40000 ALTER TABLE `macs` DISABLE KEYS */;
/*!40000 ALTER TABLE `macs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `manufacturers`
--

DROP TABLE IF EXISTS `manufacturers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `manufacturers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `device_id` int(11) unsigned NOT NULL,
  `company` varchar(255) NOT NULL,
  `product` varchar(255) DEFAULT NULL,
  `model_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `device_id` (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `manufacturers`
--

LOCK TABLES `manufacturers` WRITE;
/*!40000 ALTER TABLE `manufacturers` DISABLE KEYS */;
/*!40000 ALTER TABLE `manufacturers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `networks`
--

DROP TABLE IF EXISTS `networks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `networks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `network` varchar(255) NOT NULL,
  `netmask` varchar(255) NOT NULL,
  `gateway` varchar(255) NOT NULL,
  `vlan` varchar(255) DEFAULT NULL,
  `trunk` varchar(255) DEFAULT NULL,
  `bonding` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `network` (`network`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `networks`
--

LOCK TABLES `networks` WRITE;
/*!40000 ALTER TABLE `networks` DISABLE KEYS */;
/*!40000 ALTER TABLE `networks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `os_configs`
--

DROP TABLE IF EXISTS `os_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `os_configs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `pxe` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `os_configs`
--

LOCK TABLES `os_configs` WRITE;
/*!40000 ALTER TABLE `os_configs` DISABLE KEYS */;
INSERT INTO `os_configs` VALUES (2,'2015-11-24 08:00:44','2015-12-10 03:09:07',NULL,'centos6u7-x86_64','DEFAULT centos6.7\nLABEL centos6.7\nKERNEL http://osinstall.idcos.net/centos/6.7/os/x86_64/images/pxeboot/vmlinuz\nAPPEND initrd=http://osinstall.idcos.net/centos/6.7/os/x86_64/images/pxeboot/initrd.img ksdevice=bootif ks=http://osinstall.idcos.net/api/osinstall/v1/device/getSystemBySn?sn={sn} console=tty0 selinux=0 biosdevname=0\nIPAPPEND 2'),(30,'2015-12-08 09:41:50','2015-12-10 03:10:24',NULL,'sles11sp4-x86_64','DEFAULT sles11sp4\nLABEL sles11sp4\nKERNEL http://osinstall.idcos.net/sles/11sp4/os/x86_64/boot/x86_64/loader/linux\nAPPEND initrd=http://osinstall.idcos.net/sles/11sp4/os/x86_64/boot/x86_64/loader/initrd netdevice=bootif install=http://osinstall.idcos.net/sles/11sp4/os/x86_64/ autoyast=http://osinstall.idcos.net/api/osinstall/v1/device/getSystemBySn?sn={sn} console=tty0 selinux=0 biosdevname=0 textmode=1\nIPAPPEND 2'),(31,'2015-12-09 07:37:03','2015-12-10 03:11:13',NULL,'rhel7u2-x86_64','DEFAULT rhel7.2\nLABEL rhel7.2\nKERNEL http://osinstall.idcos.net/rhel/7.2/os/x86_64/images/pxeboot/vmlinuz\nAPPEND initrd=http://osinstall.idcos.net/rhel/7.2/os/x86_64/images/pxeboot/initrd.img ksdevice=bootif ks=http://osinstall.idcos.net/api/osinstall/v1/device/getSystemBySn?sn={sn} console=tty0 selinux=0 net.ifnames=0 biosdevname=0\nIPAPPEND 2');
/*!40000 ALTER TABLE `os_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `system_configs`
--

DROP TABLE IF EXISTS `system_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `system_configs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `content` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `system_configs`
--

LOCK TABLES `system_configs` WRITE;
/*!40000 ALTER TABLE `system_configs` DISABLE KEYS */;
INSERT INTO `system_configs` VALUES (1,'2015-11-25 02:24:10','2015-12-08 06:35:55',NULL,'centos6','install\nurl --url=http://mirror.idcos.net/centos/6.7/os/x86_64/\nlang en_US.UTF-8\nkeyboard us\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\nrootpw  --iscrypted $6$eAdCfx9hZjVMqyS6$BYIbEu4zeKp0KLnz8rLMdU7sQ5o4hQRv55o151iLX7s2kSq.5RVsteGWJlpPMqIRJ8.WUcbZC3duqX0Rt3unK/\nfirewall --disabled\nauthconfig --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai\ntext\nreboot\nzerombr\nbootloader --location=mbr --append=\"console=tty0 biosdevname=0 audit=0 selinux=0\"\nclearpart --all --initlabel\npart /boot --fstype=ext4 --size=256 --ondisk=sda\npart swap --size=2048 --ondisk=sda\npart / --fstype=ext4 --size=100 --grow --ondisk=sda\n\n%packages --ignoremissing\n@base\n@core\n@development\n\n%pre\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\ncurl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"启动OS安装程序\\\",\\\"InstallProgress\\\":0.6,\\\"InstallLog\\\":\\\"SW5zdGFsbCBPUwo=\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\ncurl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"分区并安装软件包\\\",\\\"InstallProgress\\\":0.7,\\\"InstallLog\\\":\\\"SW5zdGFsbCBPUwo=\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\n\n%post\nprogress() {\n    curl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"$1\\\",\\\"InstallProgress\\\":$2,\\\"InstallLog\\\":\\\"$3\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\n}\n\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\n\n#exec &>/tmp/install.log\n\nprogress \"配置主机名和网络\" 0.8 \"Y29uZmlnIG5ldHdvcmsK\"\n\n# config network\ncat > /etc/modprobe.d/disable_ipv6.conf <<EOF\ninstall ipv6 /bin/true\nEOF\n\ncurl -o /tmp/networkinfo \"http://osinstall.idcos.net/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw\"\nsource /tmp/networkinfo\n\ncat > /etc/sysconfig/network <<EOF\nNETWORKING=yes\nHOSTNAME=$HOSTNAME\nGATEWAY=$GATEWAY\nNOZEROCONF=yes\nNETWORKING_IPV6=no\nIPV6INIT=no\nPEERNTP=no\nEOF\n\ncat > /etc/sysconfig/network-scripts/ifcfg-eth0 <<EOF\nDEVICE=eth0\nBOOTPROTO=static\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nONBOOT=yes\nTYPE=Ethernet\nNM_CONTROLLED=no\nEOF\n\nprogress \"添加用户\" 0.85 \"YWRkIHVzZXIgeXVuamkK\"\nuseradd yunji\n\nprogress \"配置系统服务\" 0.9 \"Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg==\"\n\n# config service\nservice=(crond network ntpd rsyslog sshd sysstat)\nchkconfig --list | awk \'{ print $1 }\' | xargs -n1 -I@ chkconfig @ off\necho ${service[@]} | xargs -n1 | xargs -I@ chkconfig @ on\n\nprogress \"调整系统参数\" 0.95 \"Y29uZmlnIGJhc2ggcHJvbXB0Cg==\"\n\n# custom bash prompt\ncat >> /etc/profile <<\'EOF\'\n\nexport LANG=en_US.UTF8\nexport PS1=\'\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m`pwd`\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ \'\nexport HISTTIMEFORMAT=\'[%F %T] \'\nEOF\n\nsed -i -r -e \'/^(serial|terminal)/d\' /boot/grub/grub.conf\n\n#_log=$(base64 -w 0 /tmp/install.log)\nprogress \"安装完成\" 1 \"aW5zdGFsbCBmaW5pc2hlZAo=\"'),(11,'2015-12-08 09:45:45','2015-12-10 07:54:26',NULL,'sles11','<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<!DOCTYPE profile>\n\n<profile xmlns=\"http://www.suse.com/1.0/yast2ns\" xmlns:config=\"http://www.suse.com/1.0/configns\">  \n  <add-on/>  \n  <bootloader> \n    <device_map config:type=\"list\"> \n      <device_map_entry> \n        <firmware>hd0</firmware>  \n        <linux>/dev/sda</linux> \n      </device_map_entry> \n    </device_map>  \n    <global> \n      <activate>true</activate>  \n      <default>SUSE Linux Enterprise Server 11 SP4 - 3.0.101-63</default>  \n      <generic_mbr>true</generic_mbr>  \n      <gfxmenu>(hd0,0)/message</gfxmenu>  \n      <lines_cache_id>2</lines_cache_id>  \n      <timeout config:type=\"integer\">8</timeout> \n    </global>  \n    <initrd_modules config:type=\"list\"> \n      <initrd_module> \n        <module>megaraid_sas</module> \n      </initrd_module> \n    </initrd_modules>  \n    <loader_type>grub</loader_type>  \n    <sections config:type=\"list\"> \n      <section> \n        <append>console=tty0 selinux=0 biosdevname=0 resume=/dev/sda2 splash=silent showopts</append>  \n        <image>(hd0,0)/vmlinuz-3.0.101-63-default</image>  \n        <initial>1</initial>  \n        <initrd>(hd0,0)/initrd-3.0.101-63-default</initrd>  \n        <lines_cache_id>0</lines_cache_id>  \n        <name>SUSE Linux Enterprise Server 11 SP4 - 3.0.101-63</name>  \n        <original_name>linux</original_name>  \n        <root>/dev/sda3</root>  \n        <type>image</type> \n      </section>  \n      <section> \n        <append>showopts ide=nodma apm=off noresume edd=off powersaved=off nohz=off highres=off processor.max_cstate=1 nomodeset x11failsafe</append>  \n        <image>(hd0,0)/vmlinuz-3.0.101-63-default</image>  \n        <initrd>(hd0,0)/initrd-3.0.101-63-default</initrd>  \n        <lines_cache_id>1</lines_cache_id>  \n        <name>Failsafe -- SUSE Linux Enterprise Server 11 SP4 - 3.0.101-63</name>  \n        <original_name>failsafe</original_name>  \n        <root>/dev/sda3</root>  \n        <type>image</type> \n      </section> \n    </sections> \n  </bootloader>  \n  <deploy_image> \n    <image_installation config:type=\"boolean\">false</image_installation> \n  </deploy_image>  \n  <firewall> \n    <enable_firewall config:type=\"boolean\">false</enable_firewall> \n  </firewall>  \n  <general> \n    <ask-list config:type=\"list\"/>  \n    <mode> \n      <confirm config:type=\"boolean\">false</confirm> \n    </mode>  \n    <mouse> \n      <id>none</id> \n    </mouse>  \n    <proposals config:type=\"list\"/>  \n    <storage/> \n  </general>  \n  <groups config:type=\"list\"/>  \n  <group> \n    <encrypted config:type=\"boolean\">true</encrypted>  \n    <gid>0</gid>  \n    <group_password>x</group_password>  \n    <groupname>root</groupname>  \n    <userlist/> \n  </group>  \n  <host> \n    <hosts config:type=\"list\">\n      <hosts_entry>\n        <host_address>127.0.0.1</host_address>  \n        <names config:type=\"list\">\n          <name>localhost</name> \n        </names> \n      </hosts_entry> \n    </hosts> \n  </host>  \n  <kdump> \n    <add_crash_kernel config:type=\"boolean\">false</add_crash_kernel> \n  </kdump>  \n  <keyboard> \n    <keymap>english-us</keymap> \n  </keyboard>  \n  <language> \n    <language>en_US</language>  \n    <languages>en_US</languages> \n  </language>  \n  <login_settings/>  \n  <networking>\n    <dhcp_options>\n      <dhclient_client_id></dhclient_client_id>\n      <dhclient_hostname_option>AUTO</dhclient_hostname_option>\n    </dhcp_options>\n    <dns>\n      <dhcp_hostname config:type=\"boolean\">true</dhcp_hostname>\n      <domain>localdomain</domain>\n      <hostname>localhost</hostname>\n      <resolv_conf_policy>auto</resolv_conf_policy>\n      <write_hostname config:type=\"boolean\">false</write_hostname>\n    </dns>\n    <interfaces config:type=\"list\"> \n      <interface> \n        <bootproto>dhcp4</bootproto>  \n        <device>eth0</device>  \n        <startmode>auto</startmode> \n      </interface>  \n      <interface> \n        <aliases> \n          <alias> \n            <IPADDR>127.0.0.2</IPADDR>  \n            <NETMASK>255.0.0.0</NETMASK>  \n            <PREFIXLEN>8</PREFIXLEN> \n          </alias> \n        </aliases>  \n        <broadcast>127.255.255.255</broadcast>  \n        <device>lo</device>  \n        <firewall>no</firewall>  \n        <ipaddr>127.0.0.1</ipaddr>  \n        <netmask>255.0.0.0</netmask>  \n        <network>127.0.0.0</network>  \n        <prefixlen>8</prefixlen>  \n        <startmode>auto</startmode>  \n        <usercontrol>no</usercontrol> \n      </interface> \n    </interfaces>  \n    <ipv6 config:type=\"boolean\">false</ipv6>  \n    <managed config:type=\"boolean\">false</managed>  \n    <routing> \n      <ip_forward config:type=\"boolean\">false</ip_forward> \n    </routing> \n  </networking>  \n  <nis> \n    <netconfig_policy>auto</netconfig_policy>  \n  </nis>  \n  <partitioning config:type=\"list\"> \n    <drive> \n      <device>/dev/sda</device>  \n      <initialize config:type=\"boolean\">true</initialize>  \n      <partitions config:type=\"list\"> \n        <partition> \n          <create config:type=\"boolean\">true</create>  \n          <crypt_fs config:type=\"boolean\">false</crypt_fs>  \n          <filesystem config:type=\"symbol\">ext3</filesystem>  \n          <format config:type=\"boolean\">true</format>  \n          <fstopt>acl,user_xattr</fstopt>  \n          <loop_fs config:type=\"boolean\">false</loop_fs>  \n          <mount>/boot</mount>  \n          <mountby config:type=\"symbol\">id</mountby>  \n          <partition_id config:type=\"integer\">131</partition_id>  \n          <partition_nr config:type=\"integer\">1</partition_nr>  \n          <resize config:type=\"boolean\">false</resize>  \n          <size>256M</size> \n        </partition>  \n        <partition> \n          <create config:type=\"boolean\">true</create>  \n          <crypt_fs config:type=\"boolean\">false</crypt_fs>  \n          <filesystem config:type=\"symbol\">swap</filesystem>  \n          <format config:type=\"boolean\">true</format>  \n          <fstopt>defaults</fstopt>  \n          <loop_fs config:type=\"boolean\">false</loop_fs>  \n          <mount>swap</mount>  \n          <mountby config:type=\"symbol\">id</mountby>  \n          <partition_id config:type=\"integer\">130</partition_id>  \n          <partition_nr config:type=\"integer\">2</partition_nr>  \n          <resize config:type=\"boolean\">false</resize>  \n          <size>2G</size> \n        </partition>  \n        <partition> \n          <create config:type=\"boolean\">true</create>  \n          <crypt_fs config:type=\"boolean\">false</crypt_fs>  \n          <filesystem config:type=\"symbol\">ext3</filesystem>  \n          <format config:type=\"boolean\">true</format>  \n          <fstopt>acl,user_xattr</fstopt>  \n          <loop_fs config:type=\"boolean\">false</loop_fs>  \n          <mount>/</mount>  \n          <mountby config:type=\"symbol\">id</mountby>  \n          <partition_id config:type=\"integer\">131</partition_id>  \n          <partition_nr config:type=\"integer\">3</partition_nr>  \n          <resize config:type=\"boolean\">false</resize>  \n          <size>100%</size> \n        </partition> \n      </partitions>  \n      <pesize/>  \n      <type config:type=\"symbol\">CT_DISK</type>  \n      <use>all</use> \n    </drive> \n  </partitioning>  \n  <proxy> \n    <enabled config:type=\"boolean\">false</enabled> \n  </proxy>  \n  <runlevel> \n    <default>3</default> \n  </runlevel>  \n  <software> \n    <patterns config:type=\"list\"> \n      <pattern>Basis-Devel</pattern>  \n      <pattern>base</pattern>  \n      <pattern>Minimal</pattern> \n    </patterns> \n  </software>  \n  <timezone> \n    <hwclock>localtime</hwclock>  \n    <timezone>Asia/Shanghai</timezone> \n  </timezone>  \n  <user_defaults> \n    <group>100</group>  \n    <groups>video,dialout</groups>  \n    <home>/home</home>  \n    <inactive>-1</inactive>  \n    <shell>/bin/bash</shell>  \n    <skel>/etc/skel</skel>  \n    <umask>022</umask> \n  </user_defaults>  \n  <users config:type=\"list\"/>  \n  <user> \n    <encrypted config:type=\"boolean\">true</encrypted>  \n    <fullname>root</fullname>  \n    <gid>0</gid>  \n    <home>/root</home>  \n    <password_settings> \n      <expire/>  \n      <flag/>  \n      <inact/>  \n      <max/>  \n      <min/>  \n      <warn/> \n    </password_settings>  \n    <shell>/bin/bash</shell>  \n    <uid>0</uid>  \n    <user_password>$2y$05$P58A74K8q3STIFopY0zj/eaq9Uk.K1khj8yeuJJDq4LinaEOf1Uy.</user_password>  \n    <username>root</username> \n  </user>  \n  <scripts> \n    <pre-scripts config:type=\"list\"> \n      <script> \n        <interpreter>shell</interpreter>  \n        <source> <![CDATA[\n#!/bin/bash\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\ncurl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"启动OS安装程序\\\",\\\"InstallProgress\\\":0.6,\\\"InstallLog\\\":\\\"SW5zdGFsbCBPUwo=\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\ncurl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"分区并安装软件包\\\",\\\"InstallProgress\\\":0.7,\\\"InstallLog\\\":\\\"SW5zdGFsbCBPUwo=\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\n]]> </source> \n      </script> \n    </pre-scripts>  \n    <post-scripts config:type=\"list\"> \n      <script> \n        <interpreter>shell</interpreter>  \n        <source> <![CDATA[\n#!/bin/bash\nprogress() {\n    curl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"$1\\\",\\\"InstallProgress\\\":$2,\\\"InstallLog\\\":\\\"$3\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\n}\n\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\n\nexec &>/tmp/install.log\n\nprogress \"配置主机名和网络\" 0.8 \"Y29uZmlnIG5ldHdvcmsK\"\n\n# config network\ncurl -o /tmp/networkinfo \"http://osinstall.idcos.net/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw\"\nsource /tmp/networkinfo\n\nhostname $HOSTNAME\ncat > /etc/HOSTNAME <<EOF\n$HOSTNAME\nEOF\n\ncat > /etc/sysconfig/network/ifcfg-eth0 <<EOF\nBOOTPROTO=\'static\'\nSTARTMODE=\'auto\'\nNAME=\'eth0\'\nBROADCAST=\'\'\nETHTOOL_OPTIONS=\'\'\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nMTU=\'\'\nNETWORK=\'\'\nREMOTE_IPADDR=\'\'\nUSERCONTROL=\'no\'\nEOF\n\ncat > /etc/sysconfig/network/routes <<EOF\ndefault $GATEWAY - -\nEOF\n\nprogress \"添加用户\" 0.85 \"YWRkIHVzZXIgeXVuamkK\"\necho \'root:yunjikeji\' | chpasswd\n\nprogress \"配置系统服务\" 0.9 \"Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg==\"\n\n# config service\n\nprogress \"调整系统参数\" 0.95 \"Y29uZmlnIGJhc2ggcHJvbXB0Cg==\"\n\n# custom bash prompt\ncat >> /etc/profile <<\'EOF\'\n\nexport LANG=en_US.UTF8\nexport PS1=\'\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m`pwd`\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ \'\nexport HISTTIMEFORMAT=\'[%F %T] \'\nEOF\n\n#_log=$(base64 -w 0 /tmp/install.log)\nprogress \"安装完成\" 1 \"aW5zdGFsbCBmaW5pc2hlZAo=\"\n]]> </source> \n      </script> \n    </post-scripts> \n  </scripts> \n</profile>'),(12,'2015-12-09 07:39:28','2015-12-09 13:43:17',NULL,'rhel7','install\nurl --url=http://mirror.idcos.net/rhel/7.2/os/x86_64/\nlang en_US.UTF-8\nkeyboard --vckeymap=us --xlayouts=\'us\'\nnetwork --bootproto=dhcp --device=eth0 --noipv6 --activate\nrootpw --iscrypted $6$hrKVIh4.DTVDR2Fp$Q.ho5bHXzIoKmaXGJCSbBnC5PaXNe5wbrcbe70mMlZON20aX.BGySazXrfs0ePnTDrCF8JRzDmH8815CbaAVn.\nfirewall --disabled\nauth --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai --isUtc\ntext\nreboot\nzerombr\nbootloader --location=mbr --append=\"console=tty0 net.ifnames=0 biosdevname=0 audit=0 selinux=0\"\nclearpart --all --initlabel\npart /boot --fstype=ext4 --size=256 --ondisk=sda\npart swap --size=2048 --ondisk=sda\npart / --fstype=ext4 --size=100 --grow --ondisk=sda\n\n%packages --ignoremissing\n@base\n@core\n@development\n%end\n\n%pre\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\ncurl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"启动OS安装程序\\\",\\\"InstallProgress\\\":0.6,\\\"InstallLog\\\":\\\"SW5zdGFsbCBPUwo=\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\ncurl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"分区并安装软件包\\\",\\\"InstallProgress\\\":0.7,\\\"InstallLog\\\":\\\"SW5zdGFsbCBPUwo=\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\n%end\n\n%post\nprogress() {\n    curl -H \"Content-Type: application/json\" -X POST -d \"{\\\"Sn\\\":\\\"$_sn\\\",\\\"Title\\\":\\\"$1\\\",\\\"InstallProgress\\\":$2,\\\"InstallLog\\\":\\\"$3\\\"}\" http://osinstall.idcos.net/api/osinstall/v1/report/deviceInstallInfo\n}\n\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk \'/^[^#]/ { print $1 }\')\n\n#exec &>/tmp/install.log\n\nprogress \"配置主机名和网络\" 0.8 \"Y29uZmlnIG5ldHdvcmsK\"\n\n# config network\ncurl -o /tmp/networkinfo \"http://osinstall.idcos.net/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw\"\nsource /tmp/networkinfo\n\necho \"$HOSTNAME\" > /etc/hostname\n\ncat > /etc/sysconfig/network <<EOF\nNETWORKING=yes\nGATEWAY=$GATEWAY\nNOZEROCONF=yes\nNETWORKING_IPV6=no\nIPV6INIT=no\nPEERNTP=no\nEOF\n\ncat > /etc/sysconfig/network-scripts/ifcfg-eth0 <<EOF\nDEVICE=eth0\nBOOTPROTO=static\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nONBOOT=yes\nTYPE=Ethernet\nNM_CONTROLLED=no\nEOF\n\nprogress \"添加用户\" 0.85 \"YWRkIHVzZXIgeXVuamkK\"\nuseradd yunji\n\nprogress \"配置系统服务\" 0.9 \"Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg==\"\n\n# config service\nservice=(crond network ntpd rsyslog sshd sysstat)\nchkconfig --list | awk \'{ print $1 }\' | xargs -n1 -I@ chkconfig @ off\necho ${service[@]} | xargs -n1 | xargs -I@ chkconfig @ on\n\nprogress \"调整系统参数\" 0.95 \"Y29uZmlnIGJhc2ggcHJvbXB0Cg==\"\n\n# custom bash prompt\ncat >> /etc/profile <<\'EOF\'\n\nexport LANG=en_US.UTF8\nexport PS1=\'\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m`pwd`\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ \'\nexport HISTTIMEFORMAT=\'[%F %T] \'\nEOF\n\n#_log=$(base64 -w 0 /tmp/install.log)\nprogress \"安装完成\" 1 \"aW5zdGFsbCBmaW5pc2hlZAo=\"\n%end');
/*!40000 ALTER TABLE `system_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vm_devices`
--

DROP TABLE IF EXISTS `vm_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vm_devices` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `device_id` int(10) NOT NULL,
  `hostname` varchar(255) NOT NULL,
  `mac` varchar(255) NOT NULL,
  `ip` varchar(255) NOT NULL,
  `network_id` int(10) NOT NULL,
  `os_id` int(10) NOT NULL,
  `cpu_cores_number` int(10) NOT NULL,
  `cpu_hot_plug` enum('Yes','No') NOT NULL DEFAULT 'No',
  `cpu_passthrough` enum('Yes','No') NOT NULL DEFAULT 'No',
  `cpu_top_sockets` int(10) DEFAULT '0',
  `cpu_top_cores` int(10) DEFAULT '0',
  `cpu_top_threads` int(10) DEFAULT '0',
  `cpu_pinning` text,
  `memory_current` int(10) DEFAULT '0',
  `memory_max` int(10) DEFAULT '0',
  `memory_ksm` enum('Yes','No') NOT NULL DEFAULT 'No',
  `disk_type` varchar(255) NOT NULL,
  `disk_size` int(10) DEFAULT '0',
  `disk_bus_type` varchar(255) DEFAULT NULL,
  `disk_cache_mode` varchar(255) DEFAULT NULL,
  `disk_io_mode` varchar(255) DEFAULT NULL,
  `network_type` varchar(255) NOT NULL,
  `network_device_type` varchar(255) NOT NULL,
  `display_type` varchar(255) NOT NULL,
  `display_password` varchar(255) DEFAULT NULL,
  `display_update_password` enum('Yes','No') NOT NULL DEFAULT 'No',
  `status` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `hostname` (`hostname`),
  UNIQUE KEY `ip` (`ip`),
  UNIQUE KEY `mac` (`mac`),
  KEY `status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vm_devices`
--

LOCK TABLES `vm_devices` WRITE;
/*!40000 ALTER TABLE `vm_devices` DISABLE KEYS */;
/*!40000 ALTER TABLE `vm_devices` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-01-21 16:02:25
