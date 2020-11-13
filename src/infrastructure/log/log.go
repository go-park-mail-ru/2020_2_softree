package log

import (
	"io"
	"os"
	"server/src/infrastructure/config"

	"github.com/sirupsen/logrus"
)

type LoggerLogrus struct {
	logger logrus.Logger
}

func NewLogrusLogger() *LoggerLogrus {
	return &LoggerLogrus{logger: *logrus.New()}
}

func (l *LoggerLogrus) setLevel() error {
	level, error := logrus.ParseLevel(config.GlobalServerConfig.LogLevel)
	if error != nil {
		return error
	}
	l.logger.SetLevel(level)
	return nil
}

func (l *LoggerLogrus) setOutput() error {
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
	l.logger.SetOutput(io.MultiWriter(os.Stderr, writer))
	return nil
}

// ConfigureLogger is setting logging level and output dist
func (l *LoggerLogrus) ConfigureLogger() error {
	if err := l.setLevel(); err != nil {
		return err
	}

	if err := l.setOutput(); err != nil {
		return err
	}
	return nil
}

func (l *LoggerLogrus) Print(err ...interface{}) {
	l.logger.Println(err)
}

func (l *LoggerLogrus) Info(msg ...interface{}) {
	l.logger.Info(msg)
}

func (l *LoggerLogrus) Debug(msg ...interface{}) {
	l.logger.Debug(msg)
}
