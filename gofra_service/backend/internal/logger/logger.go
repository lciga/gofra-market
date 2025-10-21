package logger

import (
	"Gofra_Market/internal/config"
	"sync"

	"github.com/sirupsen/logrus"
)

var once sync.Once

func InitFromConfig() {
	once.Do(func() {
		cfg := config.Load()

		level := logrus.InfoLevel

		if l, e := logrus.ParseLevel(cfg.LogLevel); e == nil {
			level = l
		}

		logrus.SetLevel(level)
		logrus.SetFormatter(&ginLikeFormatter{})
		logrus.SetReportCaller(false)
	})
}

func Debug(msg string, f ...logrus.Fields) { InitFromConfig(); with(f).Debug(msg) }
func Info(msg string, f ...logrus.Fields)  { InitFromConfig(); with(f).Info(msg) }
func Warn(msg string, f ...logrus.Fields)  { InitFromConfig(); with(f).Warn(msg) }
func Error(err error, f ...logrus.Fields)  { InitFromConfig(); with(f).Error(err) }

func Debugf(fmt string, a ...any) { InitFromConfig(); logrus.Debugf(fmt, a...) }
func Infof(fmt string, a ...any)  { InitFromConfig(); logrus.Infof(fmt, a...) }
func Warnf(fmt string, a ...any)  { InitFromConfig(); logrus.Warnf(fmt, a...) }
func Errorf(fmt string, a ...any) { InitFromConfig(); logrus.Errorf(fmt, a...) }

func with(f []logrus.Fields) *logrus.Entry {
	if len(f) > 0 {
		return logrus.WithFields(f[0])
	}
	return logrus.NewEntry(logrus.StandardLogger())
}
