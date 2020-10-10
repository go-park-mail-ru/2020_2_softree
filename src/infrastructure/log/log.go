package log

import (
	"io"
	"os"
	"server/src/infrastructure/config"

	"github.com/sirupsen/logrus"
)

// GlobalLogger instance.
var GlobalLogger = logrus.New()

func setLevel() error {
	level, error := logrus.ParseLevel(config.GlobalServerConfig.LogLevel)
	if error != nil {
		return error
	}
	GlobalLogger.SetLevel(level)
	return nil
}

func setOutput() error {
	var writer io.Writer
	writer = os.Stdout
	if config.GlobalServerConfig.LogFile != "" {
		logFile, err := os.Open(config.GlobalServerConfig.LogFile)
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
	GlobalLogger.SetOutput(io.MultiWriter(os.Stderr, writer))
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
	return nil
}
