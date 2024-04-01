package pkg

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type ILogrus interface {
	Record(endpoint string, code int, method string) *logrus.Entry
}

type Logger struct {
	*logrus.Logger
}

func Logrus() *Logger {
	name := "../data/log/product.log"
	os.MkdirAll(filepath.Dir(name), 0770)
	file, _ := os.Create(name)
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
