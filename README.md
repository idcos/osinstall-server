## 依赖
- MySQL（5.6+）
- Git
- Go
- gb
- gbb

## 安装
### 拉取源代码

``` shell
$ cd $YOUR_WORK_SPACE && git clone https://github.com/idcos/osinstall-server.git
```

### *nix下安装编译环境
1. 登录[golang官网](https://golang.org/dl/)或者[golang中国](http://golangtc.com/download)下载最新的稳定版本的go安装包并安装。

	```
	$ wget https://storage.googleapis.com/golang/go1.7.4.darwin-amd64.tar.gz
	# 解压缩后go被安装在/usr/local/go
	$ sudo tar -xzv -f ./go1.7.4.darwin-amd64.tar.gz -C /usr/local/
	```

1. 配置go环境变量

	``` shell
	$ vim ~/.bashrc
	export GOROOT=/usr/local/go
	export GOPATH=$YOUR_GO_LIB_DIR
	export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
	```
配置后需要让配置生效👇

	``` shell
	$ source ~/.bashrc
	```

1. 安装[gb](https://getgb.io/)

	``` shell
	$ go get -u -v github.com/constabulary/gb/...
	```

1. 安装[gbb](https://github.com/voidint/gbb)

	``` shell
	$ go get -u -v github.com/voidint/gbb
	```

如果以上工具都安装完毕，并且`$GOROOT/bin`和`$GOPATH/bin`都已经加入到`$PATH`环境变量下，那么执行操作后是否有类似输出👇

```
$ gbb version
gbb version 0.2.0
date: 2016-12-30T13:57:50+08:00
commit: bf4f04baac5c2c20c875884672e3ed01d2962e54
```

### 编译
进入源代码根目录后执行`gbb --debug`

``` shell
$ cd $YOUR_WORK_SPACE && gbb --debug
```
编译完毕后，项目根目录`$YOUR_WORK_SPACE`下多了`bin`和`pkg`两个目录，其中`bin`目录下包含了项目的可执行文件。

``` shell
$ ls -l bin
total 105032
-rwxr-xr-x  1 voidint  staff   7.8M 12 30 17:48 cloudboot-agent
-rwxr-xr-x  1 voidint  staff   8.0M 12 30 17:48 cloudboot-encrypt-generator
-rwxr-xr-x  1 voidint  staff   5.9M 12 30 17:48 cloudboot-initdb
-rwxr-xr-x  1 voidint  staff    14M 12 30 17:48 cloudboot-server
-rwxr-xr-x  1 voidint  staff   7.4M 12 30 17:48 pe-agent
-rwxr-xr-x  1 voidint  staff   8.5M 12 30 17:48 win-agent
```

查看编译得到的可执行文件的版本信息，可以看到编译的时间戳信息-`date`和源代码的版本信息-`commit`都已经烙印在了这个二进制可执行文件的版本信息中。这类信息对于`追溯`有重要作用。

``` shell
$ ./bin/cloudboot-server -v
cloudboot-server version v1.3.1
date: 2016-12-30T17:48:34+08:00
commit: 8d614f11902927b0790a0491b1a4e192fe053657
```

详情，请移步[gbb](https://github.com/voidint/gbb)。

### 初始化数据
1. 导入SQL文件初始化数据库
将`$osinstall_server/doc/db/cloudboot.sql`导入MySQL。

1. 配置文件`/etc/cloudboot-server/cloudboot-server.conf`

``` JSON
{
    "repo": {
        "connection": "root:mypassword@tcp(localhost:3306)/cloudboot?charset=utf8&parseTime=True&loc=Local"
    },
    "osInstall": {
        "httpPort": 8081,
        "pxeConfigDir": "/etc/osinstall-server/pxelinux.cfg"
    },
    "logger": {
        "logFile": "~/logs/osinstall.log",
        "level": "debug"
    },
    "vm": {
        "storage": "guest_images_lvm"
    },
    "rsa": {
        "publicKey": "/etc/cloudboot-server/rsa/public.pem",
        "privateKey": "/etc/cloudboot-server/rsa/private.pem"
    },
    "cron": {
        "installTimeout": 3600
    },
    "activeMQ": {
        "server": "activemq.dev.idcos.net:61614"
    },
    "device": {
        "maxBatchOperateNum": 5,
        "maxOperateNumIn5Minutes": 5
    }
}
```

## 运行

``` shell
$ cd $osinstall_server && ./bin/cloudboot-server -c /etc/cloudboot-server/cloudboot-server.conf
```