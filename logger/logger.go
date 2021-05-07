package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// 临时封装 v2 的时候直接替换

func New() *Log {
	return &Log{
		logger: logrus.New(),
	}
}

type Log struct {
	logger *logrus.Logger
}

func (l *Log) Debug(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *Log) Info(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Log) Warn(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *Log) Error(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *Log) Fatal(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *Log) SetFormat(format string) {
	//l.logger.SetFormat(format)
}

func (l *Log) Close() (err error) {
	return nil
	//return l.logger.Close()
}
