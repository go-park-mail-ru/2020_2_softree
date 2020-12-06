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

func (l *Logrus) Info(desc entity.Description) {
	l.log.WithFields(logrus.Fields{"status": desc.Status, "function": desc.Function, "action": desc.Action}).Info(desc.Err)
}

func (l *Logrus) Warn(desc entity.Description) {
	l.log.WithFields(logrus.Fields{"status": desc.Status, "function": desc.Function, "action": desc.Action}).Warn(desc.Err)
}

func (l *Logrus) Error(desc entity.Description) {
	l.log.WithFields(logrus.Fields{"status": desc.Status, "function": desc.Function, "action": desc.Action}).Error(desc.Err)
}

func (l *Logrus) setLevel() error {
	level, error := logrus.ParseLevel(viper.GetString("server.logLevel"))
	if error != nil {
		return error
	}
	logrus.SetLevel(level)
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
