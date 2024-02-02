package console

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func Log(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
