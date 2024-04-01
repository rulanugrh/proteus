package pkg

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type ILogrus interface {
	Record(endpoint string, code int, method string) *logrus.Entry
}

type Logger struct {
	*logrus.Logger
}

func Logrus() *Logger {
	file, _ := os.Create("../../data/log/product.log")
	defer file.Close()

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}

	log.SetOutput(logger.Writer())
	logger.SetOutput(io.MultiWriter(os.Stdout))
	logger.SetOutput(file)

	return &Logger{logger}
}

func (l *Logger) Record(endpoint string, code int, method string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"endpoint": endpoint,
		"code":     code,
		"method":   method,
	})
}
