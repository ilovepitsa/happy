package app

import (
	"os"

	"github.com/ilovepitsa/happy/backend/pkg/config"
	"github.com/sirupsen/logrus"
)

func SetLogrusParams(cfg *config.Config) {
	logrusLevel, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrusLevel)
	}

	if cfg.Log.JSONEnable {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	if cfg.Log.Report {
		logrus.SetReportCaller(true)
	}

	logrus.SetOutput(os.Stdout)
}
