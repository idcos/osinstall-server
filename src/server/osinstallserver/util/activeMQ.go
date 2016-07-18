package util

import (
	"config"
	"encoding/json"
	"github.com/go-stomp/stomp"
	"logger"
	"net"
	"time"
)

type ActiveMQMessageAttributes struct {
	Company     string
	Product     string
	ModelName   string
	Sn          string
	Ip          string
	Mac         string
	Nic         string
	Cpu         string
	CpuSum      uint
	Memory      string
	MemorySum   uint
	Disk        string
	DiskSum     uint
	Motherboard string
	Raid        string
	Oob         string
	IsVm        string
	NicDevice   string
}

func SendActiveMQMessage(conf *config.Config, logger logger.Logger, attributes map[string]interface{}) bool {
	if conf.ActiveMQ.Server == "" {
		logger.Error("send activeMQ message error:the server url is null")
		return false
	}

	conn, err := net.DialTimeout("tcp", conf.ActiveMQ.Server, 10*time.Second)
	if err != nil {
		logger.Errorf("connect error:%s", err.Error())
		return false
	}
	stompConn, err := stomp.Connect(conn,
		stomp.ConnOpt.Login("admin", "idcos.net"),
		stomp.ConnOpt.Host("/"),
	)
	if err != nil {
		logger.Errorf("stomp connect error:%s", err.Error())
		return false
	}
	defer stompConn.Disconnect()

	message := make(map[string]interface{})
	message["domain"] = "egfbank"
	message["source"] = "cloudboot"
	message["ci_class_name"] = "x86_phycalserver"
	message["attributes"] = attributes

	str, err := json.Marshal(message)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	logger.Debugf("send activemq:\n%s", str)
	err = stompConn.Send("XMDB_CI_Queue", "", []byte(str))
	if err != nil {
		logger.Errorf("stomp connect error:%s", err.Error())
		return false
	}
	return true
}
