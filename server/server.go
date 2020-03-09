package server

import (
	"autograder/server/handlers"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

var log = logrus.New()

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:      false,
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02:150405",
		FieldMap: logrus.FieldMap{
			"FieldKeyTime":  "@timestamp",
			"FieldKeyLevel": "@level",
			"FieldKeyMsg":   "@message",
			"FieldKeyFunc":  "@caller",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})
	handlers.SetLogger(log)

}
