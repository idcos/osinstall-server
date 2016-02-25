package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"config"

	"github.com/astaxie/beego/logs"
)

const (
	// HomeDirFlag 当前用户家目录标识符
	HomeDirFlag = "~"
)

// BeeLogger beego log实现
type BeeLogger struct {
	beeLogger *logs.BeeLogger
}

func selectLevel(level string) uint {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return 7
	case "warn":
		return 4
	case "error":
		return 3
	default:
		return 6 // default level: info
	}
}

// 将~转化为用户家目录
func rel2Abs(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, HomeDirFlag) {
		return raw, nil
	}
	user, err := user.Current()
	if err != nil {
		return raw, err
	}
	return strings.Replace(raw, HomeDirFlag, user.HomeDir, 1), nil
}

// NewBeeLogger 创建BeeLogger实例
func NewBeeLogger(conf *config.Config) *BeeLogger {
	filename := strings.TrimSpace(conf.Logger.LogFile)
	if strings.HasPrefix(filename, HomeDirFlag) {
		filename, _ = rel2Abs(filename)
	}

	var logConf struct {
		FileName string `json:"filename"`
		Level    uint   `json:"level"`
	}
	logConf.FileName = filename
	logConf.Level = selectLevel(conf.Logger.Level)

	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		fmt.Printf("MkdirAll err: %s\n", err)
	}

	logger := logs.NewLogger(1000)
	logger.EnableFuncCallDepth(true) // 输出文件名和行号
	logger.SetLogFuncCallDepth(3)

	logData, _ := json.Marshal(logConf)

	if err := logger.SetLogger("file", string(logData)); err != nil {
		fmt.Printf("SetLogger err: %s\n", err)
	}

	// 尝试重置日志文件权限为0666
	os.Chmod(filename, 0666) // 不处理error

	return &BeeLogger{logger}
}

func (log *BeeLogger) SetField(name string, value interface{}) {
	panic("not support")
}

// Debug logs a debug message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Debug(v ...interface{}) {
	log.beeLogger.Debug("%v", v...)
}

// Debug logs a debug message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Debugf(format string, v ...interface{}) {
	log.beeLogger.Debug(format, v...)
}

// Info logs a info message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Info(v ...interface{}) {
	log.beeLogger.Info("%v", v...)
}

// Info logs a info message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Infof(format string, v ...interface{}) {
	log.beeLogger.Info(format, v...)
}

// Warn logs a warning message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Warn(v ...interface{}) {
	log.beeLogger.Warn("%v", v...)
}

// Warn logs a warning message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Warnf(format string, v ...interface{}) {
	log.beeLogger.Warn(format, v...)
}

// Error logs an error message. If last parameter is a map[string]string, it's content
// is added as fields to the message.
func (log *BeeLogger) Error(v ...interface{}) {
	log.beeLogger.Error("%v", v...)
}

// Error logs an error message with format. If last parameter is a map[string]string,
// it's content is added as fields to the message.
func (log *BeeLogger) Errorf(format string, v ...interface{}) {
	log.beeLogger.Error(format, v...)
}
