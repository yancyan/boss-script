package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var name string

func initLogMethod(appName string) {
	name = appName
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.AddHook(&AppNameFieldHook{})
	log.AddHook(newLfsHook(10000))

}

type AppNameFieldHook struct {
}

func (hook *AppNameFieldHook) Fire(entry *log.Entry) error {
	entry.Data["appName"] = name
	return nil
}
func (hook *AppNameFieldHook) Levels() []log.Level {
	return log.AllLevels
}

func newLfsHook(maxRemainCnt uint) log.Hook {

	writer, err := rotatelogs.New(
		"log\\"+name+".%Y%m%d%H%m",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		//rotatelogs.WithLinkName(env.AppName),
		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationSize(50000),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{})
	return lfsHook
}
