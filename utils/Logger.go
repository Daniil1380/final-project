package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	// Настройка формата
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Уровень логирования из окружения или по умолчанию
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	Logger.SetLevel(level)
}
