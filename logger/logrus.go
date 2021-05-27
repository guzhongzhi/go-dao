package logger

import "github.com/sirupsen/logrus"

func NewLogrus(cfg *Config) Logger {
	log := logrus.New()
	level, err := logrus.ParseLevel(cfg.getLevel())
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
	log.SetReportCaller(cfg.ReportCaller)
	switch cfg.Formatter {
	case FormatterJSON:
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log
}
