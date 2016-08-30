package route

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/jakecoffman/cron"
	"io/ioutil"
	"logger"
	"model"
	"os"
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
	str = `set names utf8;
	use ~~cloudboot~~;
	alter table cloudboot.system_configs modify column ~~content~~ LONGTEXT;
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
insert into platform_configs(~~created_at~~,~~updated_at~~,~~name~~,~~content~~) values('2016-05-30 17:32:47','2016-05-30 17:32:47','Version','v1.3.1');

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
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS ~~system_configs_new~~;
CREATE TABLE IF NOT EXISTS ~~system_configs_new~~ (
  ~~id~~ int(10) unsigned NOT NULL,
  ~~created_at~~ timestamp NULL DEFAULT NULL,
  ~~updated_at~~ timestamp NULL DEFAULT NULL,
  ~~deleted_at~~ timestamp NULL DEFAULT NULL,
  ~~name~~ varchar(255) NOT NULL,
  ~~content~~ longtext
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;

INSERT INTO ~~system_configs_new~~ (~~id~~, ~~created_at~~, ~~updated_at~~, ~~deleted_at~~, ~~name~~, ~~content~~) VALUES
(1, '2015-11-25 10:24:10', '2016-06-20 07:34:25', NULL, 'centos6.7', 'install\nurl --url=http://osinstall.idcos.com/centos/6.7/os/x86_64/\nlang en_US.UTF-8\nkeyboard us\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\nrootpw  --iscrypted $6$eAdCfx9hZjVMqyS6$BYIbEu4zeKp0KLnz8rLMdU7sQ5o4hQRv55o151iLX7s2kSq.5RVsteGWJlpPMqIRJ8.WUcbZC3duqX0Rt3unK/\nfirewall --disabled\nauthconfig --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai\ntext\nreboot\nzerombr\nbootloader --location=mbr --append="console=tty0 biosdevname=0 audit=0 selinux=0"\nclearpart --all --initlabel\npart /boot --fstype=ext4 --size=256\npart swap --size=2048\npart / --fstype=ext4 --size=100 --grow\n\n%packages --ignoremissing\n@base\n@core\n@development\n\n%pre\nif dmidecode | grep -qEi ''VMware|VirtualBox|KVM|Xen|Parallels''; then\n    _sn=$(sed q /sys/class/net/*/address)\nelse\n    _sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\nfi\n\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"启动OS安装程序\\",\\"InstallProgress\\":0.6,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"分区并安装软件包\\",\\"InstallProgress\\":0.7,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n\n%post\nprogress() {\n    curl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"$1\\",\\"InstallProgress\\":$2,\\"InstallLog\\":\\"$3\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n}\n\nif dmidecode | grep -qEi ''VMware|VirtualBox|KVM|Xen|Parallels''; then\n    _sn=$(sed q /sys/class/net/*/address)\nelse\n    _sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\nfi\n\nprogress "配置主机名和网络" 0.8 "Y29uZmlnIG5ldHdvcmsK"\n\n# config network\ncurl -o /tmp/networkinfo "http://osinstall.idcos.com/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw"\nsource /tmp/networkinfo\n\ncat > /etc/sysconfig/network <<EOF\nNETWORKING=yes\nHOSTNAME=$HOSTNAME\nGATEWAY=$GATEWAY\nNOZEROCONF=yes\nNETWORKING_IPV6=no\nIPV6INIT=no\nPEERNTP=no\nEOF\n\ncat > /etc/sysconfig/network-scripts/ifcfg-eth0 <<EOF\nDEVICE=eth0\nBOOTPROTO=static\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nONBOOT=yes\nTYPE=Ethernet\nNM_CONTROLLED=no\nEOF\n\nprogress "添加用户" 0.85 "YWRkIHVzZXIgeXVuamkK"\n#useradd admin\n\nprogress "配置系统服务" 0.9 "Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg=="\n\n# config service\nservice=(crond network ntpd rsyslog sshd sysstat)\nchkconfig --list | awk ''{ print $1 }'' | xargs -n1 -I@ chkconfig @ off\necho ${service[@]} | xargs -n1 | xargs -I@ chkconfig @ on\n\nprogress "调整系统参数" 0.95 "Y29uZmlnIGJhc2ggcHJvbXB0Cg=="\n\n# custom bash prompt\ncat >> /etc/profile <<''EOF''\n\nexport LANG=en_US.UTF8\nexport PS1=''\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m~~pwd~~\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ ''\nexport HISTTIMEFORMAT=''[%F %T] ''\nEOF\n\nprogress "安装完成" 1 "aW5zdGFsbCBmaW5pc2hlZAo="'),
(2, '2015-12-08 17:45:45', '2016-03-29 09:25:11', NULL, 'sles11sp4', '<?xml version="1.0" encoding="utf-8"?>\n<!DOCTYPE profile>\n\n<profile xmlns="http://www.suse.com/1.0/yast2ns" xmlns:config="http://www.suse.com/1.0/configns">  \n  <add-on/>  \n  <bootloader> \n    <device_map config:type="list"> \n      <device_map_entry> \n        <firmware>hd0</firmware>  \n        <linux>/dev/sda</linux> \n      </device_map_entry> \n    </device_map>  \n    <global> \n      <activate>true</activate>  \n      <default>SUSE Linux Enterprise Server 11 SP4 - 3.0.101-63</default>  \n      <generic_mbr>true</generic_mbr>  \n      <gfxmenu>(hd0,0)/message</gfxmenu>  \n      <lines_cache_id>2</lines_cache_id>  \n      <timeout config:type="integer">8</timeout> \n    </global>  \n    <initrd_modules config:type="list"> \n      <initrd_module> \n        <module>megaraid_sas</module> \n      </initrd_module> \n    </initrd_modules>  \n    <loader_type>grub</loader_type>  \n    <sections config:type="list"> \n      <section> \n        <append>console=tty0 selinux=0 biosdevname=0 resume=/dev/sda2 splash=silent showopts</append>  \n        <image>(hd0,0)/vmlinuz-3.0.101-63-default</image>  \n        <initial>1</initial>  \n        <initrd>(hd0,0)/initrd-3.0.101-63-default</initrd>  \n        <lines_cache_id>0</lines_cache_id>  \n        <name>SUSE Linux Enterprise Server 11 SP4 - 3.0.101-63</name>  \n        <original_name>linux</original_name>  \n        <root>/dev/sda3</root>  \n        <type>image</type> \n      </section>  \n      <section> \n        <append>showopts ide=nodma apm=off noresume edd=off powersaved=off nohz=off highres=off processor.max_cstate=1 nomodeset x11failsafe</append>  \n        <image>(hd0,0)/vmlinuz-3.0.101-63-default</image>  \n        <initrd>(hd0,0)/initrd-3.0.101-63-default</initrd>  \n        <lines_cache_id>1</lines_cache_id>  \n        <name>Failsafe -- SUSE Linux Enterprise Server 11 SP4 - 3.0.101-63</name>  \n        <original_name>failsafe</original_name>  \n        <root>/dev/sda3</root>  \n        <type>image</type> \n      </section> \n    </sections> \n  </bootloader>  \n  <deploy_image> \n    <image_installation config:type="boolean">false</image_installation> \n  </deploy_image>  \n  <firewall> \n    <enable_firewall config:type="boolean">false</enable_firewall> \n  </firewall>  \n  <general> \n    <ask-list config:type="list"/>  \n    <mode> \n      <confirm config:type="boolean">false</confirm> \n    </mode>  \n    <mouse> \n      <id>none</id> \n    </mouse>  \n    <proposals config:type="list"/>  \n    <storage/> \n  </general>  \n  <groups config:type="list"/>  \n  <group> \n    <encrypted config:type="boolean">true</encrypted>  \n    <gid>0</gid>  \n    <group_password>x</group_password>  \n    <groupname>root</groupname>  \n    <userlist/> \n  </group>  \n  <host> \n    <hosts config:type="list">\n      <hosts_entry>\n        <host_address>127.0.0.1</host_address>  \n        <names config:type="list">\n          <name>localhost</name> \n        </names> \n      </hosts_entry> \n    </hosts> \n  </host>  \n  <kdump> \n    <add_crash_kernel config:type="boolean">false</add_crash_kernel> \n  </kdump>  \n  <keyboard> \n    <keymap>english-us</keymap> \n  </keyboard>  \n  <language> \n    <language>en_US</language>  \n    <languages>en_US</languages> \n  </language>  \n  <login_settings/>  \n  <networking>\n    <dhcp_options>\n      <dhclient_client_id></dhclient_client_id>\n      <dhclient_hostname_option>AUTO</dhclient_hostname_option>\n    </dhcp_options>\n    <dns>\n      <dhcp_hostname config:type="boolean">true</dhcp_hostname>\n      <domain>localdomain</domain>\n      <hostname>localhost</hostname>\n      <resolv_conf_policy>auto</resolv_conf_policy>\n      <write_hostname config:type="boolean">false</write_hostname>\n    </dns>\n    <interfaces config:type="list"> \n      <interface> \n        <bootproto>dhcp4</bootproto>  \n        <device>eth0</device>  \n        <startmode>auto</startmode> \n      </interface>  \n      <interface> \n        <aliases> \n          <alias> \n            <IPADDR>127.0.0.2</IPADDR>  \n            <NETMASK>255.0.0.0</NETMASK>  \n            <PREFIXLEN>8</PREFIXLEN> \n          </alias> \n        </aliases>  \n        <broadcast>127.255.255.255</broadcast>  \n        <device>lo</device>  \n        <firewall>no</firewall>  \n        <ipaddr>127.0.0.1</ipaddr>  \n        <netmask>255.0.0.0</netmask>  \n        <network>127.0.0.0</network>  \n        <prefixlen>8</prefixlen>  \n        <startmode>auto</startmode>  \n        <usercontrol>no</usercontrol> \n      </interface> \n    </interfaces>  \n    <ipv6 config:type="boolean">false</ipv6>  \n    <managed config:type="boolean">false</managed>  \n    <routing> \n      <ip_forward config:type="boolean">false</ip_forward> \n    </routing> \n  </networking>  \n  <nis> \n    <netconfig_policy>auto</netconfig_policy>  \n  </nis>  \n  <partitioning config:type="list"> \n    <drive> \n      <device>/dev/sda</device>  \n      <initialize config:type="boolean">true</initialize>  \n      <partitions config:type="list"> \n        <partition> \n          <create config:type="boolean">true</create>  \n          <crypt_fs config:type="boolean">false</crypt_fs>  \n          <filesystem config:type="symbol">ext3</filesystem>  \n          <format config:type="boolean">true</format>  \n          <fstopt>acl,user_xattr</fstopt>  \n          <loop_fs config:type="boolean">false</loop_fs>  \n          <mount>/boot</mount>  \n          <mountby config:type="symbol">id</mountby>  \n          <partition_id config:type="integer">131</partition_id>  \n          <partition_nr config:type="integer">1</partition_nr>  \n          <resize config:type="boolean">false</resize>  \n          <size>256M</size> \n        </partition>  \n        <partition> \n          <create config:type="boolean">true</create>  \n          <crypt_fs config:type="boolean">false</crypt_fs>  \n          <filesystem config:type="symbol">swap</filesystem>  \n          <format config:type="boolean">true</format>  \n          <fstopt>defaults</fstopt>  \n          <loop_fs config:type="boolean">false</loop_fs>  \n          <mount>swap</mount>  \n          <mountby config:type="symbol">id</mountby>  \n          <partition_id config:type="integer">130</partition_id>  \n          <partition_nr config:type="integer">2</partition_nr>  \n          <resize config:type="boolean">false</resize>  \n          <size>2G</size> \n        </partition>  \n        <partition> \n          <create config:type="boolean">true</create>  \n          <crypt_fs config:type="boolean">false</crypt_fs>  \n          <filesystem config:type="symbol">ext3</filesystem>  \n          <format config:type="boolean">true</format>  \n          <fstopt>acl,user_xattr</fstopt>  \n          <loop_fs config:type="boolean">false</loop_fs>  \n          <mount>/</mount>  \n          <mountby config:type="symbol">id</mountby>  \n          <partition_id config:type="integer">131</partition_id>  \n          <partition_nr config:type="integer">3</partition_nr>  \n          <resize config:type="boolean">false</resize>  \n          <size>100%</size> \n        </partition> \n      </partitions>  \n      <pesize/>  \n      <type config:type="symbol">CT_DISK</type>  \n      <use>all</use> \n    </drive> \n  </partitioning>  \n  <proxy> \n    <enabled config:type="boolean">false</enabled> \n  </proxy>  \n  <runlevel> \n    <default>3</default> \n  </runlevel>  \n  <software> \n    <patterns config:type="list"> \n      <pattern>Basis-Devel</pattern>  \n      <pattern>base</pattern>  \n      <pattern>Minimal</pattern> \n    </patterns> \n  </software>  \n  <timezone> \n    <hwclock>localtime</hwclock>  \n    <timezone>Asia/Shanghai</timezone> \n  </timezone>  \n  <user_defaults> \n    <group>100</group>  \n    <groups>video,dialout</groups>  \n    <home>/home</home>  \n    <inactive>-1</inactive>  \n    <shell>/bin/bash</shell>  \n    <skel>/etc/skel</skel>  \n    <umask>022</umask> \n  </user_defaults>  \n  <users config:type="list"/>  \n  <user> \n    <encrypted config:type="boolean">true</encrypted>  \n    <fullname>root</fullname>  \n    <gid>0</gid>  \n    <home>/root</home>  \n    <password_settings> \n      <expire/>  \n      <flag/>  \n      <inact/>  \n      <max/>  \n      <min/>  \n      <warn/> \n    </password_settings>  \n    <shell>/bin/bash</shell>  \n    <uid>0</uid>  \n    <user_password>$2y$05$P58A74K8q3STIFopY0zj/eaq9Uk.K1khj8yeuJJDq4LinaEOf1Uy.</user_password>  \n    <username>root</username> \n  </user>  \n  <scripts> \n    <pre-scripts config:type="list"> \n      <script> \n        <interpreter>shell</interpreter>  \n        <source> <![CDATA[\n#!/bin/bash\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"启动OS安装程序\\",\\"InstallProgress\\":0.6,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"分区并安装软件包\\",\\"InstallProgress\\":0.7,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n]]> </source> \n      </script> \n    </pre-scripts>  \n    <post-scripts config:type="list"> \n      <script> \n        <interpreter>shell</interpreter>  \n        <source> <![CDATA[\n#!/bin/bash\nprogress() {\n    curl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"$1\\",\\"InstallProgress\\":$2,\\"InstallLog\\":\\"$3\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n}\n\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\n\nprogress "配置主机名和网络" 0.8 "Y29uZmlnIG5ldHdvcmsK"\n\n# config network\ncurl -o /tmp/networkinfo "http://osinstall.idcos.com/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw"\nsource /tmp/networkinfo\n\nhostname $HOSTNAME\ncat > /etc/HOSTNAME <<EOF\n$HOSTNAME\nEOF\n\ncat > /etc/sysconfig/network/ifcfg-eth0 <<EOF\nBOOTPROTO=''static''\nSTARTMODE=''auto''\nNAME=''eth0''\nBROADCAST=''''\nETHTOOL_OPTIONS=''''\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nMTU=''''\nNETWORK=''''\nREMOTE_IPADDR=''''\nUSERCONTROL=''no''\nEOF\n\ncat > /etc/sysconfig/network/routes <<EOF\ndefault $GATEWAY - -\nEOF\n\nprogress "添加用户" 0.85 "YWRkIHVzZXIgeXVuamkK"\necho ''root:yunjikeji'' | chpasswd\n\nprogress "配置系统服务" 0.9 "Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg=="\n\n# config service\n\nprogress "调整系统参数" 0.95 "Y29uZmlnIGJhc2ggcHJvbXB0Cg=="\n\n# custom bash prompt\ncat >> /etc/profile <<''EOF''\n\nexport LANG=en_US.UTF8\nexport PS1=''\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m~~pwd~~\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ ''\nexport HISTTIMEFORMAT=''[%F %T] ''\nEOF\n\nprogress "安装完成" 1 "aW5zdGFsbCBmaW5pc2hlZAo="\n]]> </source> \n      </script> \n    </post-scripts> \n  </scripts> \n</profile>'),
(3, '2015-12-09 15:39:28', '2016-06-20 08:29:33', NULL, 'rhel7.2', 'install\nurl --url=http://osinstall.idcos.com/rhel/7.2/os/x86_64/\nlang en_US.UTF-8\nkeyboard --vckeymap=us --xlayouts=''us''\nnetwork --bootproto=dhcp --device=eth0 --noipv6 --activate\nrootpw --iscrypted $6$hrKVIh4.DTVDR2Fp$Q.ho5bHXzIoKmaXGJCSbBnC5PaXNe5wbrcbe70mMlZON20aX.BGySazXrfs0ePnTDrCF8JRzDmH8815CbaAVn.\nfirewall --disabled\nauth --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai --isUtc\ntext\nreboot\nzerombr\nbootloader --location=mbr --append="console=tty0 net.ifnames=0 biosdevname=0 audit=0 selinux=0"\nclearpart --all --initlabel\npart /boot --fstype=ext4 --size=256\npart swap --size=2048\npart / --fstype=ext4 --size=100 --grow\n\n%packages --ignoremissing\n@base\n@core\n@development\n%end\n\n%pre\nif dmidecode | grep -qEi ''VMware|VirtualBox|KVM|Xen|Parallels''; then\n    _sn=$(sed q /sys/class/net/*/address)\nelse\n    _sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\nfi\n\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"启动OS安装程序\\",\\"InstallProgress\\":0.6,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"分区并安装软件包\\",\\"InstallProgress\\":0.7,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n%end\n\n%post\nprogress() {\n    curl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"$1\\",\\"InstallProgress\\":$2,\\"InstallLog\\":\\"$3\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n}\n\nif dmidecode | grep -qEi ''VMware|VirtualBox|KVM|Xen|Parallels''; then\n    _sn=$(sed q /sys/class/net/*/address)\nelse\n    _sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\nfi\n\nprogress "配置主机名和网络" 0.8 "Y29uZmlnIG5ldHdvcmsK"\n\n# config network\ncurl -o /tmp/networkinfo "http://osinstall.idcos.com/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw"\nsource /tmp/networkinfo\n\necho "$HOSTNAME" > /etc/hostname\n\ncat > /etc/sysconfig/network <<EOF\nNETWORKING=yes\nGATEWAY=$GATEWAY\nNOZEROCONF=yes\nNETWORKING_IPV6=no\nIPV6INIT=no\nPEERNTP=no\nEOF\n\ncat > /etc/sysconfig/network-scripts/ifcfg-eth0 <<EOF\nDEVICE=eth0\nBOOTPROTO=static\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nONBOOT=yes\nTYPE=Ethernet\nNM_CONTROLLED=no\nEOF\n\nprogress "添加用户" 0.85 "YWRkIHVzZXIgeXVuamkK"\n#useradd admin\n\nprogress "配置系统服务" 0.9 "Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg=="\n\n# config service\nservice=(crond network ntpd rsyslog sshd sysstat)\nchkconfig --list | awk ''{ print $1 }'' | xargs -n1 -I@ chkconfig @ off\necho ${service[@]} | xargs -n1 | xargs -I@ chkconfig @ on\n\nprogress "调整系统参数" 0.95 "Y29uZmlnIGJhc2ggcHJvbXB0Cg=="\n\n# custom bash prompt\ncat >> /etc/profile <<''EOF''\n\nexport LANG=en_US.UTF8\nexport PS1=''\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m~~pwd~~\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ ''\nexport HISTTIMEFORMAT=''[%F %T] ''\nEOF\n\nprogress "安装完成" 1 "aW5zdGFsbCBmaW5pc2hlZAo="\n%end'),
(4, '2016-02-01 11:45:16', '2016-02-29 11:59:12', NULL, 'esxi6.0', 'vmaccepteula\r\nrootpw yunjikeji\r\ninstall --firstdisk --overwritevmfs\r\nreboot\r\n\r\n%include /tmp/network.ks\r\n\r\n%pre --interpreter=busybox\r\n_sn=$(localcli hardware platform get | awk ''/Serial Number/ { print $NF }'')\r\n_dns=$(awk ''/^nameserver/ { print $NF; exit }'' /etc/resolv.conf)\r\nwget -qO /tmp/networkinfo "http://osinstall.idcos.com/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw"\r\nsource /tmp/networkinfo\r\necho "network --bootproto=static --device=vmnic0 --ip=$IPADDR --netmask=$NETMASK --gateway=$GATEWAY --nameserver=$_dns --hostname=$HOSTNAME" > /tmp/network.ks\r\n\r\ncat > /tmp/progress.py <<''EOF''\r\n#!/usr/bin/env python\r\nimport sys\r\nimport traceback\r\nimport json\r\nimport urllib2\r\nimport urllib\r\n\r\nURL = "http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo"\r\n\r\ndef process_api(sn, title, process, content):\r\n\r\n    params = {\r\n        "Sn": sn,\r\n        "Title": title,\r\n        "InstallProgress": round(float(process), 2),\r\n        "InstallLog": content\r\n    }\r\n    data = json.dumps(params)\r\n    req = urllib2.Request(URL,data=data)\r\n    req.add_header("Content-Type", "application/json")\r\n    return urllib2.urlopen(req).read()\r\n\r\n\r\nif __name__ == "__main__":\r\n    try:\r\n        print process_api(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])\r\n    except:\r\n        print "error when request URL: ", URL\r\n        print traceback.print_exc()\r\nEOF\r\nchmod 755 /tmp/progress.py\r\n\r\n/tmp/progress.py $_sn "启动OS安装程序" 0.6 "SW5zdGFsbCBPUwo="\r\n\r\n%firstboot --interpreter=busybox\r\ncat > /tmp/progress.py <<''EOF''\r\n#!/usr/bin/env python\r\nimport sys\r\nimport traceback\r\nimport json\r\nimport urllib2\r\nimport urllib\r\n\r\nURL = "http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo"\r\n\r\ndef process_api(sn, title, process, content):\r\n\r\n    params = {\r\n        "Sn": sn,\r\n        "Title": title,\r\n        "InstallProgress": round(float(process), 2),\r\n        "InstallLog": content\r\n    }\r\n    data = json.dumps(params)\r\n    req = urllib2.Request(URL,data=data)\r\n    req.add_header("Content-Type", "application/json")\r\n    return urllib2.urlopen(req).read()\r\n\r\n\r\nif __name__ == "__main__":\r\n    try:\r\n        print process_api(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])\r\n    except:\r\n        print "error when request URL: ", URL\r\n        print traceback.print_exc()\r\nEOF\r\nchmod 755 /tmp/progress.py\r\n\r\n_sn=$(localcli hardware platform get | awk ''/Serial Number/ { print $NF }'')\r\n\r\n/tmp/progress.py $_sn "调整系统参数" 0.8 "ZW5hYmxlIHNoZWxsIGFuZCBzc2gK"\r\n\r\nvim-cmd hostsvc/enable_ssh\r\nvim-cmd hostsvc/start_ssh\r\n\r\nvim-cmd hostsvc/enable_esx_shell\r\nvim-cmd hostsvc/start_esx_shell\r\n\r\n/tmp/progress.py $_sn "安装完成" 1 "aW5zdGFsbCBmaW5pc2hlZAo="'),
(5, '2016-02-17 12:01:00', '2016-03-01 15:38:06', NULL, 'win2008r2_cn', '<?xml version="1.0" encoding="utf-8"?>\r<unattend xmlns="urn:schemas-microsoft-com:unattend">\r    <settings pass="generalize">\r        <component name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>\r        </component>\r        <component name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>\r        </component>\r    </settings>\r    <settings pass="specialize">\r        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <ProductKey>489J6-VHDMP-X63PK-3K798-CPX3Y</ProductKey>\r            <ShowWindowsLive>false</ShowWindowsLive>\r            <DisableAutoDaylightTimeSet>false</DisableAutoDaylightTimeSet>\r            <TimeZone>China Standard Time</TimeZone>\r            <RegisteredOwner />\r            <RegisteredOrganization />\r        </component>\r        <component name="Microsoft-Windows-TerminalServices-LocalSessionManager" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <fDenyTSConnections>false</fDenyTSConnections>\r        </component>\r        <component name="Networking-MPSSVC-Svc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <FirewallGroups>\r                <FirewallGroup wcm:action="add" wcm:keyValue="1">\r                    <Active>true</Active>\r                    <Group>远程桌面</Group>\r                    <Profile>all</Profile>\r                </FirewallGroup>\r            </FirewallGroups>\r        </component>\r    </settings>\r    <settings pass="oobeSystem">\r        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <AutoLogon>\r                <Password>\r                    <Value>eQB1AG4AagBpAGsAZQBqAGkAUABhAHMAcwB3AG8AcgBkAA==</Value>\r                    <PlainText>false</PlainText>\r                </Password>\r                <Enabled>true</Enabled>\r                <LogonCount>1</LogonCount>\r                <Username>Administrator</Username>\r            </AutoLogon>\r            <OOBE>\r                <HideEULAPage>true</HideEULAPage>\r            </OOBE>\r            <UserAccounts>\r                <AdministratorPassword>\r                    <PlainText>false</PlainText>\r                    <Value>eQB1AG4AagBpAGsAZQBqAGkAQQBkAG0AaQBuAGkAcwB0AHIAYQB0AG8AcgBQAGEAcwBzAHcAbwByAGQA</Value>\r                </AdministratorPassword>\r            </UserAccounts>\r            <FirstLogonCommands>\r                <SynchronousCommand wcm:action="add">\r                    <Order>1</Order>\r                    <CommandLine>C:\\firstboot\\winconfig.exe</CommandLine>\r                    <Description></Description>\r                    <RequiresUserInput>false</RequiresUserInput>\r                </SynchronousCommand>\r            </FirstLogonCommands>\r            <RegisteredOrganization />\r            <RegisteredOwner />\r        </component>\r        <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <SystemLocale>zh-CN</SystemLocale>\r            <UILanguage>zh-CN</UILanguage>\r            <UILanguageFallback>zh-CN</UILanguageFallback>\r            <UserLocale>zh-CN</UserLocale>\r            <InputLocale>zh-CN</InputLocale>\r        </component>\r    </settings>\r    <settings pass="windowsPE">\r        <component name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <ImageInstall>\r                <OSImage>\r                    <InstallFrom>\r                        <Path>Z:\\windows\\2008r2_cn\\sources\\install.wim</Path>\r                        <MetaData wcm:action="add">\r                            <Key>/IMAGE/NAME</Key>\r                            <Value>Windows Server 2008 R2 SERVERENTERPRISE</Value>\r                        </MetaData>\r                    </InstallFrom>\r                    <InstallTo>\r                        <DiskID>0</DiskID>\r                        <PartitionID>1</PartitionID>\r                    </InstallTo>\r                    <WillShowUI>OnError</WillShowUI>\r                </OSImage>\r            </ImageInstall>\r            <UserData>\r                <ProductKey>\r                    <WillShowUI>OnError</WillShowUI>\r                    <Key>489J6-VHDMP-X63PK-3K798-CPX3Y</Key>\r                </ProductKey>\r                <AcceptEula>true</AcceptEula>\r            </UserData>\r            <EnableFirewall>true</EnableFirewall>\r            <EnableNetwork>true</EnableNetwork>\r            <DiskConfiguration>\r                <Disk wcm:action="add">\r                    <CreatePartitions>\r                        <CreatePartition wcm:action="add">\r                            <Type>Primary</Type>\r                            <Order>1</Order>\r                            <Size>51200</Size>\r                        </CreatePartition>\r                    </CreatePartitions>\r                    <ModifyPartitions>\r                        <ModifyPartition wcm:action="add">\r                            <Active>true</Active>\r                            <Extend>false</Extend>\r                            <Format>NTFS</Format>\r                            <Order>1</Order>\r                            <Label>System</Label>\r                            <PartitionID>1</PartitionID>\r                        </ModifyPartition>\r                    </ModifyPartitions>\r                    <DiskID>0</DiskID>\r                    <WillWipeDisk>true</WillWipeDisk>\r                </Disk>\r                <WillShowUI>OnError</WillShowUI>\r            </DiskConfiguration>\r        </component>\r        <component name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <SetupUILanguage>\r                <UILanguage>zh-CN</UILanguage>\r            </SetupUILanguage>\r            <InputLocale>zh-CN</InputLocale>\r            <SystemLocale>zh-CN</SystemLocale>\r            <UILanguage>zh-CN</UILanguage>\r            <UserLocale>zh-CN</UserLocale>\r            <UILanguageFallback>zh-CN</UILanguageFallback>\r        </component>\r        <component name="Microsoft-Windows-PnpCustomizationsWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DriverPaths>\r                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\r                    <Path>Z:\\windows\\drivers\\2008r2</Path>\r                </PathAndCredentials>\r            </DriverPaths>\r        </component>\r    </settings>\r    <settings pass="offlineServicing">\r        <component name="Microsoft-Windows-PnpCustomizationsNonWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DriverPaths>\r                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\r                    <Path>Z:\\windows\\drivers\\2008r2</Path>\r                </PathAndCredentials>\r            </DriverPaths>\r        </component>\r        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <EnableLUA>false</EnableLUA>\r        </component>\r    </settings>\r    <cpi:offlineImage cpi:source="catalog://osinstall/image/windows/clg/install_windows server 2008 r2 serverenterprise.clg" xmlns:cpi="urn:schemas-microsoft-com:cpi" />\r</unattend>'),
(6, '2016-02-26 15:25:15', '2016-02-29 22:00:16', NULL, 'win2008r2_en', '<?xml version="1.0" encoding="utf-8"?>\n<unattend xmlns="urn:schemas-microsoft-com:unattend">\n    <settings pass="generalize">\n        <component name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>\n        </component>\n        <component name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>\n        </component>\n    </settings>\n    <settings pass="specialize">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ProductKey>489J6-VHDMP-X63PK-3K798-CPX3Y</ProductKey>\n            <ShowWindowsLive>false</ShowWindowsLive>\n            <DisableAutoDaylightTimeSet>false</DisableAutoDaylightTimeSet>\n            <TimeZone>China Standard Time</TimeZone>\n            <RegisteredOwner />\n            <RegisteredOrganization />\n        </component>\n        <component name="Microsoft-Windows-TerminalServices-LocalSessionManager" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <fDenyTSConnections>false</fDenyTSConnections>\n        </component>\n        <component name="Networking-MPSSVC-Svc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <FirewallGroups>\n                <FirewallGroup wcm:action="add" wcm:keyValue="1">\n                    <Active>true</Active>\n                    <Group>Remote Desktop</Group>\n                    <Profile>all</Profile>\n                </FirewallGroup>\n            </FirewallGroups>\n        </component>\n    </settings>\n    <settings pass="oobeSystem">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <AutoLogon>\n                <Password>\n                    <Value>eQB1AG4AagBpAGsAZQBqAGkAUABhAHMAcwB3AG8AcgBkAA==</Value>\n                    <PlainText>false</PlainText>\n                </Password>\n                <Enabled>true</Enabled>\n                <LogonCount>1</LogonCount>\n                <Username>Administrator</Username>\n            </AutoLogon>\n            <OOBE>\n                <HideEULAPage>true</HideEULAPage>\n            </OOBE>\n            <UserAccounts>\n                <AdministratorPassword>\n                    <PlainText>false</PlainText>\n                    <Value>eQB1AG4AagBpAGsAZQBqAGkAQQBkAG0AaQBuAGkAcwB0AHIAYQB0AG8AcgBQAGEAcwBzAHcAbwByAGQA</Value>\n                </AdministratorPassword>\n            </UserAccounts>\n            <FirstLogonCommands>\n                <SynchronousCommand wcm:action="add">\n                    <Order>1</Order>\n                    <CommandLine>C:\\firstboot\\winconfig.exe</CommandLine>\n                    <Description></Description>\n                    <RequiresUserInput>false</RequiresUserInput>\n                </SynchronousCommand>\n            </FirstLogonCommands>\n            <RegisteredOrganization />\n            <RegisteredOwner />\n        </component>\n        <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SystemLocale>en-US</SystemLocale>\n            <UILanguage>en-US</UILanguage>\n            <UILanguageFallback>en-US</UILanguageFallback>\n            <UserLocale>en-US</UserLocale>\n            <InputLocale>en-US</InputLocale>\n        </component>\n    </settings>\n    <settings pass="windowsPE">\n        <component name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ImageInstall>\n                <OSImage>\n                    <InstallFrom>\n                        <Path>Z:\\windows\\2008r2_en\\sources\\install.wim</Path>\n                        <MetaData wcm:action="add">\n                            <Key>/IMAGE/NAME</Key>\n                            <Value>Windows Server 2008 R2 SERVERENTERPRISE</Value>\n                        </MetaData>\n                    </InstallFrom>\n                    <InstallTo>\n                        <DiskID>0</DiskID>\n                        <PartitionID>1</PartitionID>\n                    </InstallTo>\n                    <WillShowUI>OnError</WillShowUI>\n                </OSImage>\n            </ImageInstall>\n            <UserData>\n                <ProductKey>\n                    <WillShowUI>OnError</WillShowUI>\n                    <Key>489J6-VHDMP-X63PK-3K798-CPX3Y</Key>\n                </ProductKey>\n                <AcceptEula>true</AcceptEula>\n            </UserData>\n            <EnableFirewall>true</EnableFirewall>\n            <EnableNetwork>true</EnableNetwork>\n            <DiskConfiguration>\n                <Disk wcm:action="add">\n                    <CreatePartitions>\n                        <CreatePartition wcm:action="add">\n                            <Type>Primary</Type>\n                            <Order>1</Order>\n                            <Size>51200</Size>\n                        </CreatePartition>\n                    </CreatePartitions>\n                    <ModifyPartitions>\n                        <ModifyPartition wcm:action="add">\n                            <Active>true</Active>\n                            <Extend>false</Extend>\n                            <Format>NTFS</Format>\n                            <Order>1</Order>\n                            <Label>System</Label>\n                            <PartitionID>1</PartitionID>\n                        </ModifyPartition>\n                    </ModifyPartitions>\n                    <DiskID>0</DiskID>\n                    <WillWipeDisk>true</WillWipeDisk>\n                </Disk>\n                <WillShowUI>OnError</WillShowUI>\n            </DiskConfiguration>\n        </component>\n        <component name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SetupUILanguage>\n                <UILanguage>en-US</UILanguage>\n            </SetupUILanguage>\n            <InputLocale>en-US</InputLocale>\n            <SystemLocale>en-US</SystemLocale>\n            <UILanguage>en-US</UILanguage>\n            <UserLocale>en-US</UserLocale>\n            <UILanguageFallback>en-US</UILanguageFallback>\n        </component>\n        <component name="Microsoft-Windows-PnpCustomizationsWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DriverPaths>\n                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\n                    <Path>Z:\\windows\\drivers\\2008r2</Path>\n                </PathAndCredentials>\n            </DriverPaths>\n        </component>\n    </settings>\n    <settings pass="offlineServicing">\n        <component name="Microsoft-Windows-PnpCustomizationsNonWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DriverPaths>\n                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\n                    <Path>Z:\\windows\\drivers\\2008r2</Path>\n                </PathAndCredentials>\n            </DriverPaths>\n        </component>\n        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <EnableLUA>false</EnableLUA>\n        </component>\n    </settings>\n    <cpi:offlineImage cpi:source="catalog://osinstall/image/windows/clg/install_windows server 2008 r2 serverenterprise.clg" xmlns:cpi="urn:schemas-microsoft-com:cpi" />\n</unattend>'),
(7, '2016-02-29 11:05:41', '2016-03-01 15:38:22', NULL, 'win2012r2_cn', '<?xml version="1.0" encoding="utf-8"?>\r<unattend xmlns="urn:schemas-microsoft-com:unattend">\r    <settings pass="generalize">\r        <component name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>\r        </component>\r        <component name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>\r        </component>\r    </settings>\r    <settings pass="specialize">\r        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <ProductKey>W3GGN-FT8W3-Y4M27-J84CP-Q3VJ9</ProductKey>\r            <ShowWindowsLive>false</ShowWindowsLive>\r            <DisableAutoDaylightTimeSet>false</DisableAutoDaylightTimeSet>\r            <TimeZone>China Standard Time</TimeZone>\r            <RegisteredOwner />\r            <RegisteredOrganization />\r        </component>\r        <component name="Microsoft-Windows-TerminalServices-LocalSessionManager" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <fDenyTSConnections>false</fDenyTSConnections>\r        </component>\r        <component name="Networking-MPSSVC-Svc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <FirewallGroups>\r                <FirewallGroup wcm:action="add" wcm:keyValue="1">\r                    <Active>true</Active>\r                    <Group>远程桌面</Group>\r                    <Profile>all</Profile>\r                </FirewallGroup>\r            </FirewallGroups>\r        </component>\r    </settings>\r    <settings pass="oobeSystem">\r        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <AutoLogon>\r                <Password>\r                    <Value>eQB1AG4AagBpAGsAZQBqAGkAUABhAHMAcwB3AG8AcgBkAA==</Value>\r                    <PlainText>false</PlainText>\r                </Password>\r                <Enabled>true</Enabled>\r                <LogonCount>1</LogonCount>\r                <Username>Administrator</Username>\r            </AutoLogon>\r            <OOBE>\r                <HideEULAPage>true</HideEULAPage>\r            </OOBE>\r            <UserAccounts>\r                <AdministratorPassword>\r                    <PlainText>false</PlainText>\r                    <Value>eQB1AG4AagBpAGsAZQBqAGkAQQBkAG0AaQBuAGkAcwB0AHIAYQB0AG8AcgBQAGEAcwBzAHcAbwByAGQA</Value>\r                </AdministratorPassword>\r            </UserAccounts>\r            <FirstLogonCommands>\r                <SynchronousCommand wcm:action="add">\r                    <Order>1</Order>\r                    <CommandLine>C:\\firstboot\\winconfig.exe</CommandLine>\r                    <Description></Description>\r                    <RequiresUserInput>false</RequiresUserInput>\r                </SynchronousCommand>\r            </FirstLogonCommands>\r            <RegisteredOrganization />\r            <RegisteredOwner />\r        </component>\r        <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <SystemLocale>zh-CN</SystemLocale>\r            <UILanguage>zh-CN</UILanguage>\r            <UILanguageFallback>zh-CN</UILanguageFallback>\r            <UserLocale>zh-CN</UserLocale>\r            <InputLocale>zh-CN</InputLocale>\r        </component>\r    </settings>\r    <settings pass="windowsPE">\r        <component name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <ImageInstall>\r                <OSImage>\r                    <InstallFrom>\r                        <Path>Z:\\windows\\2012r2_cn\\sources\\install.wim</Path>\r                        <MetaData wcm:action="add">\r                            <Key>/IMAGE/NAME</Key>\r                            <Value>Windows Server 2012 R2 SERVERDATACENTER</Value>\r                        </MetaData>\r                    </InstallFrom>\r                    <InstallTo>\r                        <DiskID>0</DiskID>\r                        <PartitionID>1</PartitionID>\r                    </InstallTo>\r                    <WillShowUI>OnError</WillShowUI>\r                </OSImage>\r            </ImageInstall>\r            <UserData>\r                <ProductKey>\r                    <WillShowUI>OnError</WillShowUI>\r                    <Key>W3GGN-FT8W3-Y4M27-J84CP-Q3VJ9</Key>\r                </ProductKey>\r                <AcceptEula>true</AcceptEula>\r            </UserData>\r            <EnableFirewall>true</EnableFirewall>\r            <EnableNetwork>true</EnableNetwork>\r            <DiskConfiguration>\r                <Disk wcm:action="add">\r                    <CreatePartitions>\r                        <CreatePartition wcm:action="add">\r                            <Type>Primary</Type>\r                            <Order>1</Order>\r                            <Size>51200</Size>\r                        </CreatePartition>\r                    </CreatePartitions>\r                    <ModifyPartitions>\r                        <ModifyPartition wcm:action="add">\r                            <Active>true</Active>\r                            <Extend>false</Extend>\r                            <Format>NTFS</Format>\r                            <Order>1</Order>\r                            <Label>System</Label>\r                            <PartitionID>1</PartitionID>\r                        </ModifyPartition>\r                    </ModifyPartitions>\r                    <DiskID>0</DiskID>\r                    <WillWipeDisk>true</WillWipeDisk>\r                </Disk>\r                <WillShowUI>OnError</WillShowUI>\r            </DiskConfiguration>\r        </component>\r        <component name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <SetupUILanguage>\r                <UILanguage>zh-CN</UILanguage>\r            </SetupUILanguage>\r            <InputLocale>zh-CN</InputLocale>\r            <SystemLocale>zh-CN</SystemLocale>\r            <UILanguage>zh-CN</UILanguage>\r            <UserLocale>zh-CN</UserLocale>\r            <UILanguageFallback>zh-CN</UILanguageFallback>\r        </component>\r        <component name="Microsoft-Windows-PnpCustomizationsWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DriverPaths>\r                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\r                    <Path>Z:\\windows\\drivers\\2012r2</Path>\r                </PathAndCredentials>\r            </DriverPaths>\r        </component>\r    </settings>\r    <settings pass="offlineServicing">\r        <component name="Microsoft-Windows-PnpCustomizationsNonWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <DriverPaths>\r                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\r                    <Path>Z:\\windows\\drivers\\2012r2</Path>\r                </PathAndCredentials>\r            </DriverPaths>\r        </component>\r        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\r            <EnableLUA>false</EnableLUA>\r        </component>\r    </settings>\r    <cpi:offlineImage cpi:source="catalog://osinstall/image/windows/clg/install_windows server 2012 r2 serverdatacenter.clg" xmlns:cpi="urn:schemas-microsoft-com:cpi" />\r</unattend>');
INSERT INTO ~~system_configs_new~~ (~~id~~, ~~created_at~~, ~~updated_at~~, ~~deleted_at~~, ~~name~~, ~~content~~) VALUES
(8, '2016-02-29 12:06:11', '2016-02-29 22:07:45', NULL, 'win2012r2_en', '<?xml version="1.0" encoding="utf-8"?>\n<unattend xmlns="urn:schemas-microsoft-com:unattend">\n    <settings pass="generalize">\n        <component name="Microsoft-Windows-OutOfBoxExperience" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenInitialConfigurationTasksAtLogon>true</DoNotOpenInitialConfigurationTasksAtLogon>\n        </component>\n        <component name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>\n        </component>\n    </settings>\n    <settings pass="specialize">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ProductKey>W3GGN-FT8W3-Y4M27-J84CP-Q3VJ9</ProductKey>\n            <ShowWindowsLive>false</ShowWindowsLive>\n            <DisableAutoDaylightTimeSet>false</DisableAutoDaylightTimeSet>\n            <TimeZone>China Standard Time</TimeZone>\n            <RegisteredOwner />\n            <RegisteredOrganization />\n        </component>\n        <component name="Microsoft-Windows-TerminalServices-LocalSessionManager" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <fDenyTSConnections>false</fDenyTSConnections>\n        </component>\n        <component name="Networking-MPSSVC-Svc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <FirewallGroups>\n                <FirewallGroup wcm:action="add" wcm:keyValue="1">\n                    <Active>true</Active>\n                    <Group>Remote Desktop</Group>\n                    <Profile>all</Profile>\n                </FirewallGroup>\n            </FirewallGroups>\n        </component>\n    </settings>\n    <settings pass="oobeSystem">\n        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <AutoLogon>\n                <Password>\n                    <Value>eQB1AG4AagBpAGsAZQBqAGkAUABhAHMAcwB3AG8AcgBkAA==</Value>\n                    <PlainText>false</PlainText>\n                </Password>\n                <Enabled>true</Enabled>\n                <LogonCount>1</LogonCount>\n                <Username>Administrator</Username>\n            </AutoLogon>\n            <OOBE>\n                <HideEULAPage>true</HideEULAPage>\n            </OOBE>\n            <UserAccounts>\n                <AdministratorPassword>\n                    <PlainText>false</PlainText>\n                    <Value>eQB1AG4AagBpAGsAZQBqAGkAQQBkAG0AaQBuAGkAcwB0AHIAYQB0AG8AcgBQAGEAcwBzAHcAbwByAGQA</Value>\n                </AdministratorPassword>\n            </UserAccounts>\n            <FirstLogonCommands>\n                <SynchronousCommand wcm:action="add">\n                    <Order>1</Order>\n                    <CommandLine>C:\\firstboot\\winconfig.exe</CommandLine>\n                    <Description></Description>\n                    <RequiresUserInput>false</RequiresUserInput>\n                </SynchronousCommand>\n            </FirstLogonCommands>\n            <RegisteredOrganization />\n            <RegisteredOwner />\n        </component>\n        <component name="Microsoft-Windows-International-Core" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SystemLocale>en-US</SystemLocale>\n            <UILanguage>en-US</UILanguage>\n            <UILanguageFallback>en-US</UILanguageFallback>\n            <UserLocale>en-US</UserLocale>\n            <InputLocale>en-US</InputLocale>\n        </component>\n    </settings>\n    <settings pass="windowsPE">\n        <component name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <ImageInstall>\n                <OSImage>\n                    <InstallFrom>\n                        <Path>Z:\\windows\\2012r2_en\\sources\\install.wim</Path>\n                        <MetaData wcm:action="add">\n                            <Key>/IMAGE/NAME</Key>\n                            <Value>Windows Server 2012 R2 SERVERDATACENTER</Value>\n                        </MetaData>\n                    </InstallFrom>\n                    <InstallTo>\n                        <DiskID>0</DiskID>\n                        <PartitionID>1</PartitionID>\n                    </InstallTo>\n                    <WillShowUI>OnError</WillShowUI>\n                </OSImage>\n            </ImageInstall>\n            <UserData>\n                <ProductKey>\n                    <WillShowUI>OnError</WillShowUI>\n                    <Key>W3GGN-FT8W3-Y4M27-J84CP-Q3VJ9</Key>\n                </ProductKey>\n                <AcceptEula>true</AcceptEula>\n            </UserData>\n            <EnableFirewall>true</EnableFirewall>\n            <EnableNetwork>true</EnableNetwork>\n            <DiskConfiguration>\n                <Disk wcm:action="add">\n                    <CreatePartitions>\n                        <CreatePartition wcm:action="add">\n                            <Type>Primary</Type>\n                            <Order>1</Order>\n                            <Size>51200</Size>\n                        </CreatePartition>\n                    </CreatePartitions>\n                    <ModifyPartitions>\n                        <ModifyPartition wcm:action="add">\n                            <Active>true</Active>\n                            <Extend>false</Extend>\n                            <Format>NTFS</Format>\n                            <Order>1</Order>\n                            <Label>System</Label>\n                            <PartitionID>1</PartitionID>\n                        </ModifyPartition>\n                    </ModifyPartitions>\n                    <DiskID>0</DiskID>\n                    <WillWipeDisk>true</WillWipeDisk>\n                </Disk>\n                <WillShowUI>OnError</WillShowUI>\n            </DiskConfiguration>\n        </component>\n        <component name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <SetupUILanguage>\n                <UILanguage>en-US</UILanguage>\n            </SetupUILanguage>\n            <InputLocale>en-US</InputLocale>\n            <SystemLocale>en-US</SystemLocale>\n            <UILanguage>en-US</UILanguage>\n            <UserLocale>en-US</UserLocale>\n            <UILanguageFallback>en-US</UILanguageFallback>\n        </component>\n        <component name="Microsoft-Windows-PnpCustomizationsWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DriverPaths>\n                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\n                    <Path>Z:\\windows\\drivers\\2012r2</Path>\n                </PathAndCredentials>\n            </DriverPaths>\n        </component>\n    </settings>\n    <settings pass="offlineServicing">\n        <component name="Microsoft-Windows-PnpCustomizationsNonWinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <DriverPaths>\n                <PathAndCredentials wcm:action="add" wcm:keyValue="1">\n                    <Path>Z:\\windows\\drivers\\2012r2</Path>\n                </PathAndCredentials>\n            </DriverPaths>\n        </component>\n        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS" xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">\n            <EnableLUA>false</EnableLUA>\n        </component>\n    </settings>\n    <cpi:offlineImage cpi:source="catalog://osinstall/image/windows/clg/install_windows server 2012 r2 serverdatacenter.clg" xmlns:cpi="urn:schemas-microsoft-com:cpi" />\n</unattend>'),
(9, '2016-03-29 08:41:04', '2016-06-20 09:17:33', NULL, 'ubuntu14.04', 'd-i debian-installer/locale string en_US.UTF-8\nd-i console-setup/ask_detect boolean false\nd-i keyboard-configuration/layoutcode string us\nd-i netcfg/choose_interface select auto\nd-i netcfg/target_network_config select ifupdown\nd-i mirror/country string manual\nd-i mirror/http/hostname string osinstall\nd-i mirror/http/directory string /ubuntu/14.04/os/x86_64\nd-i mirror/http/proxy string\nd-i live-installer/net-image string http://osinstall.idcos.com/ubuntu/14.04/os/x86_64/install/filesystem.squashfs\nd-i clock-setup/utc boolean false\nd-i time/zone string Asia/Shanghai\nd-i clock-setup/ntp boolean true\n#d-i partman-auto/disk string /dev/sda\nd-i partman-auto/method string regular\nd-i partman-lvm/device_remove_lvm boolean true\nd-i partman-md/device_remove_md boolean true\nd-i partman-lvm/confirm boolean true\nd-i partman-auto/choose_recipe select atomic\nd-i partman/default_filesystem string ext4\nd-i partman/mount_style select uuid\nd-i partman/unmount_active boolean true\nd-i partman-partitioning/confirm_write_new_label boolean true\nd-i partman/choose_partition select finish\nd-i partman/confirm boolean true\nd-i partman/confirm_nooverwrite boolean true\nd-i passwd/root-login boolean true\nd-i passwd/make-user boolean false\nd-i passwd/root-password-crypted password $1$AxPsi$GDSXEkYCIL2xfRuimCMiX1\nd-i user-setup/encrypt-home boolean false\ntasksel tasksel/first multiselect standard\nd-i pkgsel/include string openssh-server build-essential ntp vim dmidecode curl\nd-i pkgsel/update-policy select none\nd-i grub-installer/only_debian boolean true\nd-i grub-installer/with_other_os boolean true\nd-i finish-install/reboot_in_progress note\nd-i debian-installer/exit/reboot boolean true\nd-i preseed/early_command string umount /media || true\nd-i preseed/late_command string \\\nexport LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/target/usr/lib/x86_64-linux-gnu:/target/lib/x86_64-linux-gnu ; \\\n/target/usr/sbin/dmidecode | grep -qEi ''VMware|VirtualBox|KVM|Xen|Parallels'' && _sn=$(sed q /sys/class/net/*/address) || _sn=$(/target/usr/sbin/dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }''); \\\n/target/usr/bin/curl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"配置主机名和网络\\",\\"InstallProgress\\":0.8,\\"InstallLog\\":\\"Y29uZmlnIG5ldHdvcmsK\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo ; \\\n/target/usr/bin/curl -o /tmp/networkinfo "http://osinstall.idcos.com/api/osinstall/v1/device/getNetworkBySn?sn=$_sn&type=raw" ; \\\n. /tmp/networkinfo ; \\\necho -e "auto lo\\niface lo inet loopback\\nauto eth0\\niface eth0 inet static\\naddress $IPADDR\\nnetmask $NETMASK\\ngateway $GATEWAY" > /etc/network/interfaces ; \\\ncp /etc/network/interfaces /target/etc/network/interfaces ; \\\necho "$HOSTNAME" > /target/etc/hostname ; \\\n/target/usr/bin/curl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"安装完成\\",\\"InstallProgress\\":1,\\"InstallLog\\":\\"aW5zdGFsbCBmaW5pc2hlZAo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo'),
(10, '2016-06-16 06:46:19', '2016-06-20 05:55:13', NULL, 'xenserver6.5', '<?xml version="1.0" encoding="UTF-8"?>\n<installation srtype="ext">\n   <primary-disk>sda</primary-disk>\n   <keymap>us</keymap>\n   <root-password>yunjikeji</root-password>\n   <source type="url">http://osinstall.idcos.com/xenserver/6.5/</source>\n   <script stage="filesystem-populated" type="url">http://osinstall.idcos.com/scripts/xenserver.sh</script>\n   <admin-interface name="eth0" proto="dhcp" />\n   <timezone>Asia/Shanghai</timezone>\n</installation>'),
(11, '2016-06-20 05:53:16', '2016-06-23 04:39:11', NULL, 'centos6.7-kvmserver', 'install\nurl --url=http://osinstall.idcos.com/centos/6.7/os/x86_64/\nlang en_US.UTF-8\nkeyboard us\nnetwork --onboot yes --device bootif --bootproto dhcp --noipv6\nrootpw  --iscrypted $6$eAdCfx9hZjVMqyS6$BYIbEu4zeKp0KLnz8rLMdU7sQ5o4hQRv55o151iLX7s2kSq.5RVsteGWJlpPMqIRJ8.WUcbZC3duqX0Rt3unK/\nfirewall --disabled\nauthconfig --enableshadow --passalgo=sha512\nselinux --disabled\ntimezone Asia/Shanghai\ntext\nreboot\nzerombr\nbootloader --location=mbr --append="console=tty0 biosdevname=0 audit=0 selinux=0"\nclearpart --all --initlabel\npart /boot --fstype=ext4 --size=256 --ondisk=sda\npart swap --size=2048 --ondisk=sda\npart / --fstype=ext4 --size=51200 --ondisk=sda\npart pv.01 --size 100 --grow\nvolgroup VolGroup0 pv.01\n\n%packages --ignoremissing\n@base\n@core\n@development\n@virtualization\n@virtualization-platform\n\n%pre\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"启动OS安装程序\\",\\"InstallProgress\\":0.6,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\ncurl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"分区并安装软件包\\",\\"InstallProgress\\":0.7,\\"InstallLog\\":\\"SW5zdGFsbCBPUwo=\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n\n%post\nprogress() {\n    curl -H "Content-Type: application/json" -X POST -d "{\\"Sn\\":\\"$_sn\\",\\"Title\\":\\"$1\\",\\"InstallProgress\\":$2,\\"InstallLog\\":\\"$3\\"}" http://osinstall.idcos.com/api/osinstall/v1/report/deviceInstallInfo\n}\n\n_sn=$(dmidecode -s system-serial-number 2>/dev/null | awk ''/^[^#]/ { print $1 }'')\n\nprogress "配置主机名和网络" 0.8 "Y29uZmlnIG5ldHdvcmsK"\n\n# config network\ncurl -o /tmp/networkinfo "http://osinstall.idcos.com/api/osinstall/v1/device/getNetworkBySn?sn=${_sn}&type=raw"\nsource /tmp/networkinfo\n\ncat > /etc/sysconfig/network <<EOF\nNETWORKING=yes\nHOSTNAME=$HOSTNAME\nGATEWAY=$GATEWAY\nNOZEROCONF=yes\nNETWORKING_IPV6=no\nIPV6INIT=no\nPEERNTP=no\nEOF\n\ncat > /etc/sysconfig/network-scripts/ifcfg-br0 <<EOF\nDEVICE=br0\nBOOTPROTO=none\nIPADDR=$IPADDR\nNETMASK=$NETMASK\nONBOOT=yes\nTYPE=Ethernet\nNM_CONTROLLED=no\nTYPE=Bridge\nDELAY=0\nEOF\n\ncat > /etc/sysconfig/network-scripts/ifcfg-eth0 <<EOF\nDEVICE=eth0\nBOOTPROTO=none\nONBOOT=yes\nTYPE=Ethernet\nNM_CONTROLLED=no\nTYPE=Ethernet\nBRIDGE=br0\nEOF\n\nprogress "添加用户" 0.85 "YWRkIHVzZXIgeXVuamkK"\n#useradd admin\n\nprogress "配置系统服务" 0.9 "Y29uZmlnIHN5c3RlbSBzZXJ2aWNlCg=="\n\n# config service\nservice=(crond irqbalance libvirtd network ntpd rsyslog sshd sysstat)\nchkconfig --list | awk ''{ print $1 }'' | xargs -n1 -I@ chkconfig @ off\necho ${service[@]} | xargs -n1 | xargs -I@ chkconfig @ on\n\nprogress "调整系统参数" 0.95 "Y29uZmlnIGJhc2ggcHJvbXB0Cg=="\n\n# custom bash prompt\ncat >> /etc/profile <<''EOF''\n\nexport LANG=en_US.UTF8\nexport PS1=''\\n\\e[1;37m[\\e[m\\e[1;32m\\u\\e[m\\e[1;33m@\\e[m\\e[1;35m\\H\\e[m:\\e[4m~~pwd~~\\e[m\\e[1;37m]\\e[m\\e[1;36m\\e[m\\n\\$ ''\nexport HISTTIMEFORMAT=''[%F %T] ''\nEOF\n\ncat >> /etc/sysctl.conf <<''EOF''\n\nnet.ipv6.conf.all.disable_ipv6 = 1\nnet.ipv6.conf.default.disable_ipv6 = 1\nnet.ipv6.conf.lo.disable_ipv6 = 1\nnet.bridge.bridge-nf-call-arptables = 0\nnet.bridge.bridge-nf-call-ip6tables = 0\nnet.bridge.bridge-nf-call-iptables = 0\nEOF\n\nservice libvirtd start\nvirsh pool-define-as guest_images_lvm logical - - - VolGroup0 /dev/VolGroup0\nvirsh pool-autostart guest_images_lvm\n\nprogress "安装完成" 1 "aW5zdGFsbCBmaW5pc2hlZAo="');

ALTER TABLE ~~system_configs_new~~
  ADD PRIMARY KEY (~~id~~), ADD UNIQUE KEY ~~name~~ (~~name~~);

ALTER TABLE ~~system_configs_new~~
  MODIFY ~~id~~ int(10) unsigned NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=12;

update system_configs t1,system_configs_new t2
set t1.~~content~~ = t2.~~content~~,t1.updated_at = now()
where t1.~~name~~ = t2.~~name~~ and t1.~~content~~ != t2.~~content~~;

insert into system_configs(~~created_at~~, ~~updated_at~~, ~~deleted_at~~, ~~name~~, ~~content~~)
select t1.~~created_at~~, t1.~~updated_at~~, t1.~~deleted_at~~, t1.~~name~~, t1.~~content~~ from system_configs_new t1
left join system_configs t2 on t1.~~name~~ = t2.~~name~~
where t2.id is null;

DROP TABLE IF EXISTS ~~os_configs_new~~;
CREATE TABLE IF NOT EXISTS ~~os_configs_new~~ (
  ~~id~~ int(10) unsigned NOT NULL,
  ~~created_at~~ timestamp NULL DEFAULT NULL,
  ~~updated_at~~ timestamp NULL DEFAULT NULL,
  ~~deleted_at~~ timestamp NULL DEFAULT NULL,
  ~~name~~ varchar(255) NOT NULL,
  ~~pxe~~ longtext
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;

INSERT INTO ~~os_configs_new~~ (~~id~~, ~~created_at~~, ~~updated_at~~, ~~deleted_at~~, ~~name~~, ~~pxe~~) VALUES
(1, '2015-11-24 16:00:44', '2015-12-10 11:09:07', NULL, 'centos6u7-x86_64', 'DEFAULT centos6.7\nLABEL centos6.7\nKERNEL http://osinstall.idcos.com/centos/6.7/os/x86_64/images/pxeboot/vmlinuz\nAPPEND initrd=http://osinstall.idcos.com/centos/6.7/os/x86_64/images/pxeboot/initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/osinstall/v1/device/getSystemBySn?sn={sn} console=tty0 selinux=0 biosdevname=0\nIPAPPEND 2'),
(2, '2015-12-08 17:41:50', '2015-12-10 11:10:24', NULL, 'sles11sp4-x86_64', 'DEFAULT sles11sp4\nLABEL sles11sp4\nKERNEL http://osinstall.idcos.com/sles/11sp4/os/x86_64/boot/x86_64/loader/linux\nAPPEND initrd=http://osinstall.idcos.com/sles/11sp4/os/x86_64/boot/x86_64/loader/initrd netdevice=bootif install=http://osinstall.idcos.com/sles/11sp4/os/x86_64/ autoyast=http://osinstall.idcos.com/api/osinstall/v1/device/getSystemBySn?sn={sn} console=tty0 selinux=0 biosdevname=0 textmode=1\nIPAPPEND 2'),
(3, '2015-12-09 15:37:03', '2015-12-10 11:11:13', NULL, 'rhel7u2-x86_64', 'DEFAULT rhel7.2\nLABEL rhel7.2\nKERNEL http://osinstall.idcos.com/rhel/7.2/os/x86_64/images/pxeboot/vmlinuz\nAPPEND initrd=http://osinstall.idcos.com/rhel/7.2/os/x86_64/images/pxeboot/initrd.img ksdevice=bootif ks=http://osinstall.idcos.com/api/osinstall/v1/device/getSystemBySn?sn={sn} console=tty0 selinux=0 net.ifnames=0 biosdevname=0\nIPAPPEND 2'),
(4, '2016-02-01 11:35:40', '2016-02-01 11:45:49', NULL, 'esxi6.0u1-x86_64', 'DEFAULT esxi\nLABEL esxi\nMENU LABEL ^ESXi 6.0U1\nKERNEL http://osinstall.idcos.com/esxi/6.0u1/mboot.c32\nAPPEND -c http://osinstall.idcos.com/esxi/6.0u1/boot.cfg ks=http://osinstall.idcos.com/api/osinstall/v1/device/getSystemBySn?sn={sn}'),
(5, '2016-02-17 12:08:14', '2016-02-29 16:16:01', NULL, 'win2008r2-x86_64', 'DEFAULT winpe\nLABEL winpe\nMENU LABEL ^WinPE\nKERNEL memdisk\nAPPEND initrd=http://osinstall.idcos.com/winpe/winpe.iso iso raw'),
(6, '2016-02-29 11:04:57', '2016-02-29 16:16:07', NULL, 'win2012r2-x86_64', 'DEFAULT winpe\nLABEL winpe\nMENU LABEL ^WinPE\nKERNEL memdisk\nAPPEND initrd=http://osinstall.idcos.com/winpe/winpe.iso iso raw'),
(7, '2016-03-29 08:40:22', '2016-03-29 08:40:22', NULL, 'ubuntu1404-x86_64', 'DEFAULT ubuntu14.04\nLABEL ubuntu14.04\nKERNEL http://osinstall.idcos.com/ubuntu/14.04/os/x86_64/install/netboot/ubuntu-installer/amd64/linux\nAPPEND initrd=http://osinstall.idcos.com/ubuntu/14.04/os/x86_64/install/netboot/ubuntu-installer/amd64/initrd.gz auto=true priority=critical net.ifnames=1 biosdevname=0 netcfg/choose_interface=auto preseed/url=http://osinstall.idcos.com/api/osinstall/v1/device/getSystemBySn?sn={sn} --'),
(8, '2016-06-20 05:56:09', '2016-06-20 05:56:09', NULL, 'xenserver6.5-x86_64', 'DEFAULT xenserver\nLABEL xenserver\nMENU LABEL ^XenServer 6.5\nKERNEL http://osinstall.idcos.com/xenserver/6.5/boot/pxelinux/mboot.c32\nAPPEND http://osinstall.idcos.com/xenserver/6.5/boot/xen.gz dom0_max_vcpus=1-2 dom0_mem=752M,max:752M com1=115200,8n1 console=com1,vga --- http://osinstall.idcos.com/xenserver/6.5/boot/vmlinuz xencons=hvc console=hvc0 console=tty0 answerfile=http://osinstall.idcos.com/api/osinstall/v1/device/getSystemBySn?sn={sn} install --- http://osinstall.idcos.com/xenserver/6.5/install.img');

ALTER TABLE ~~os_configs_new~~
  ADD PRIMARY KEY (~~id~~), ADD UNIQUE KEY ~~name~~ (~~name~~);
ALTER TABLE ~~os_configs_new~~
  MODIFY ~~id~~ int(10) unsigned NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=9;

update os_configs t1,os_configs_new t2
set t1.~~pxe~~ = t2.~~pxe~~,t1.updated_at = now()
where t1.~~name~~ = t2.~~name~~ and t1.~~pxe~~ != t2.~~pxe~~;

insert into os_configs(~~created_at~~, ~~updated_at~~, ~~deleted_at~~, ~~name~~, ~~pxe~~)
select t1.~~created_at~~, t1.~~updated_at~~, t1.~~deleted_at~~, t1.~~name~~, t1.~~pxe~~ from os_configs_new t1
left join os_configs t2 on t1.~~name~~ = t2.~~name~~
where t2.id is null;

drop table os_configs_new;
drop table system_configs_new;`
	str = strings.Replace(str, "~~", "`", -1)
	file := "/tmp/cloudboot-v1.3-update.sql"
	bytes := []byte(str)
	errWrite := ioutil.WriteFile(file, bytes, 0644)
	if errWrite != nil {
		logger.Error(errWrite.Error())
	}

	cmd := `mysql -uroot < ` + file
	logger.Debugf("mysql exec:%s", cmd)
	result, err := util.ExecScript(cmd)
	logger.Debugf("result:%s", string(result))
	if err != nil {
		logger.Error(err.Error())
		logger.Info("db version failed upgrade to v1.3")
	} else {
		defer os.Remove(file)
		logger.Info("db version has been successfully upgraded to v1.3")
	}
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
