package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"server/authorization/pkg/infrastructure/persistence"
	session "server/authorization/pkg/session/gen"
	"server/canal/pkg/infrastructure/config"
	"server/canal/pkg/infrastructure/logger"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	pflag.StringP("config", "c", "", "path to viper file")
	pflag.BoolP("help", "h", false, "usage info")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalln(err)
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if err := config.ParseConfig(
		viper.GetString("config"),
		map[string]interface{}{
			"server": map[string]interface{}{
				"ip":       "127.0.0.1",
				"port":     8001,
				"logLevel": "Info",
			},

			"redis": map[string]interface{}{
				"host": "127.0.0.1",
				"port": 6379,
				"user": "user",
			},
		}); err != nil {
		log.Fatalln("Error during parse defaults", err)
	}

	if err := logger.ConfigureLogger(); err != nil {
		log.Fatalln(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	connect, err := redis.DialURL(fmt.Sprintf("redis://%s:%s:%d",
		viper.GetString("redis.user"),
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
	))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
			"action":   "connect to redis",
		}).Fatalln(err)
	}

	server := grpc.NewServer()

	session.RegisterAuthorizationServiceServer(server, persistence.NewSessionManager(connect))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		viper.GetString("server.ip"),
		viper.GetInt("server.port"),
	))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
			"action":   "starting listening tcp port",
		}).Fatalln(err)
	}

	if err := server.Serve(lis); err != nil {
		logrus.WithFields(logrus.Fields{
			"action": "Starting server",
		}).Fatalln(err)
	}
}
