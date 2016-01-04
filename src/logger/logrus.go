package logger

import (
	"config"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

type LogrusLogger struct {
	*logrus.Entry
}

func (log *LogrusLogger) SetField(key string, value interface{}) {
	log.WithField(key, value)
}

func NewLogger() *LogrusLogger {
	return &LogrusLogger{logrus.NewEntry(logrus.New())}
}

func NewLogrusLogger(conf *config.Config) (logrusLogger *LogrusLogger) {
	var logger = logrus.New()

	formater := &logrus.TextFormatter{
		DisableColors: true,
	}
	logger.Formatter = formater

	if conf == nil {
		return
	}

	if conf.Logger.Color {
		formater.ForceColors = true
		formater.DisableColors = false
	} else {
		formater.ForceColors = false
		formater.DisableColors = true
	}

	level, err := logrus.ParseLevel(conf.Logger.Level)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logger.Level = level

	if conf.Logger.LogFile != "" {
		f, err := os.OpenFile(conf.Logger.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		logger.Out = f
	}

	logrusLogger = &LogrusLogger{logrus.NewEntry(logger)}

	// if pc, file, line, ok := runtime.Caller(1); ok {
	// 	fName := runtime.FuncForPC(pc).Name()
	// 	logrusLogger.Entry = logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	// }

	return
}
