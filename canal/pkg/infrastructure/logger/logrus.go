package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"server/canal/pkg/domain/entity"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Logrus struct {
	log *logrus.Logger
}

func NewLogrus() *Logrus {
	l := &Logrus{log: logrus.New()}
	if err := l.ConfigureLogger(); err != nil {
		log.Fatal(err)
	}

	return l
}

func (l *Logrus) Info(desc entity.Description, err error) {
	l.log.WithFields(logrus.Fields{"status": desc.Status, "function": desc.Function, "action": desc.Action}).Info(err)
}

func (l *Logrus) Warn(desc entity.Description, err error) {
	l.log.WithFields(logrus.Fields{"status": desc.Status, "function": desc.Function, "action": desc.Action}).Warn(err)
}

func (l *Logrus) Error(desc entity.Description, err error) {
	l.log.WithFields(logrus.Fields{"status": desc.Status, "function": desc.Function, "action": desc.Action}).Error(err)
}

func (l *Logrus) setLevel() error {
	if viper.GetString("server.logLevel") != "" {
		level, err := logrus.ParseLevel(viper.GetString("server.logLevel"))
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
		return nil
	}
	logrus.SetLevel(logrus.DebugLevel)
	return nil
}

func (l *Logrus) setOutput() error {
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
	l.log.SetOutput(writer)
	return nil
}

// ConfigureLogger is setting logging level and output dist
func (l *Logrus) ConfigureLogger() error {
	if err := l.setLevel(); err != nil {
		return err
	}
	if err := l.setOutput(); err != nil {
		return err
	}
	l.log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	return nil
}
