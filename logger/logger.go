package logger

import (
	"os"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var l log.Logger
var once sync.Once

func newLogger() (logger log.Logger) {
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, allowedLevels())
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "time", log.DefaultTimestampUTC)
	return
}

func Info() log.Logger {
	once.Do(func() {
		l = newLogger()
	})
	return level.Info(l)
}

func Warn() log.Logger {
	once.Do(func() {
		l = newLogger()
	})
	return level.Warn(l)
}

func Error() log.Logger {
	once.Do(func() {
		l = newLogger()
	})
	return level.Error(l)
}

func allowedLevels() level.Option {
	switch config := os.Getenv("LOG_LEVEL"); {
	case config == "DEBUG":
		return level.AllowDebug()
	case config == "WARN":
		return level.AllowWarn()
	case config == "NONE":
		return level.AllowNone()
	default:
		return level.AllowInfo()
	}
}
