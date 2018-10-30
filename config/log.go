package config

import (
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

type Logger struct {
}

func Debug(action string, data ...lager.Data) {
	log.Debug(action, data...)
}

func Info(action string, data ...lager.Data) {
	log.Info(action, data...)
}

func Warn(action string, data ...lager.Data) {
	log.Warn(action, data...)
}

func Error(action string, err error, data ...lager.Data) {
	log.Error(action, err, data...)
}

func Fatal(action string, err error, data ...lager.Data) {
	log.Fatal(action, err, data...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(err error, format string, args ...interface{}) {
	log.Errorf(err, format, args...)
}

func Fatalf(err error, format string, args ...interface{}) {
	log.Fatalf(err, format, args...)
}
