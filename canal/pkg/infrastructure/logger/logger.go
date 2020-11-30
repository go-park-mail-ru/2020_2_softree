package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setLevel() error {
	level, error := logrus.ParseLevel(viper.GetString("server.logLevel"))
	if error != nil {
		return error
	}
	logrus.SetLevel(level)
	return nil
}

func setOutput() error {
	var writer io.Writer
	writer = os.Stdout
	if viper.GetString("server.logFile") != "" {
		fullpath, err := filepath.Abs(viper.GetString("server.logFile"))
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
	logrus.SetOutput(writer)
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
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	return nil
}
