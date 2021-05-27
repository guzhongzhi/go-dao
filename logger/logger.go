package logger

const (
	LevelFatal   = "fatal"
	LevelError   = "error"
	LevelWarning = "warning"
	LevelInfo    = "info"
	LevelDebug   = "debug"
	LevelTrace   = "trace"

	FormatterJSON = "json"
)

type Loggerf interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger interface {
	Loggerf
	Debug(args ...interface{})

	Info(args ...interface{})
	Warn(args ...interface{})

	Error(args ...interface{})
	Fatal(args ...interface{})
}

type CloserLogger interface {
	Close() error
}

var defaultLogger Logger

func Default() Logger {
	if defaultLogger == nil {
		defaultLogger = NewLogrus(&Config{})
	}
	return defaultLogger
}

type Config struct {
	Level        string
	FileName     string
	ReportCaller bool
	Formatter    string
}

func (s *Config) getLevel() string {
	if s.Level == "" {
		s.Level = LevelInfo
	}
	return s.Level
}

func (s *Config) getFormatter() string {
	if s.Formatter == "" {
		s.Formatter = FormatterJSON
	}
	return s.Formatter
}
