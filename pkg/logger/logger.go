package looger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func Info(message string, values ...interface{}) {
	logrus.Info(fmt.Sprintf(message, values...))
}

func InfoWithFields(message string, fields map[string]interface{}) {
	logrus.WithFields(fields).Info(message)
}

func Warn(message string, values ...interface{}) {
	logrus.Warn(fmt.Sprintf(message, values...))
}

func WarnWithFields(message string, fields map[string]interface{}) {
	logrus.WithFields(fields).Warn(message)
}
