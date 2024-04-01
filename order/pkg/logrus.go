package pkg

import "github.com/sirupsen/logrus"

type ILogrust interface {
	Record(endpoint string, method string, code int)
}

type Logger struct {
	*logrus.Logger
}
