package logger

import (
	"io"
	"os"
	"path/filepath"
	"server/src/infrastructure/config"

	"github.com/sirupsen/logrus"
)

// GlobalLogger instance
var GlobalLogger = logrus.New()

func setLevel() error {
	level, error := logrus.ParseLevel(config.GlobalConfig.GetString("server.logLevel"))
	if error != nil {
		return error
	}
	GlobalLogger.SetLevel(level)
	return nil
}

func setOutput() error {
	var writer io.Writer
	writer = os.Stdout
	if config.GlobalConfig.GetString("server.logFile") != "" {
		fullpath, err := filepath.Abs(config.GlobalConfig.GetString("server.logFile"))
		if err != nil {
			return err
		}
		logFile, err := os.OpenFile(fullpath, os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		logrus.RegisterExitHandler(func() {
			if logFile == nil {
				return
			}
			logFile.Close()
		})
		writer = logFile
	}
	GlobalLogger.SetOutput(writer)
	return nil
}

// ConfigureLogger is setting logging level and output dist
func ConfigureLogger() error {
	if err := setLevel(); err != nil {
		return err
	}
	if err := setOutput(); err != nil {
		return err
	}
	GlobalLogger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	return nil
}
