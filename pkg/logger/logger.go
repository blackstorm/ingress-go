package looger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
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
