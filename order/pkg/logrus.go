package pkg

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type ILogrus interface {
	RecordGRPC(endpoint string, method string, code int) *logrus.Entry
}

type Logger struct {
	*logrus.Logger
}

func Logrus() *Logger {
	name := "../data/log/order.log"
	err := os.MkdirAll(filepath.Dir(name), 0770)
	if err != nil {
		log.Println("Error while create folder")
	}
	
	file, err := os.Create(name)
	if err != nil {
		log.Println("Error while create file")
	}

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

func (l *Logger) RecordGRPC(endpoint string, method string, code int) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"type": "gprc",
		"endpoint": endpoint,
		"method": method,
		"code": code,
	})
}