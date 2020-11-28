package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"server/authorization/pkg/infrastructure/persistence"
	session "server/authorization/pkg/session/gen"
	"server/canal/pkg/infrastructure/config"
)

func init() {
	pflag.StringP("viper", "c", "", "path to viper file")
	pflag.BoolP("help", "h", false, "usage info")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if err := config.ParseConfig(
		viper.GetString("viper"),
		map[string]interface{}{
			"redis": map[string]interface{}{
				"host": "127.0.0.1",
				"port": 6379,
				"user": "user",
			},
		}); err != nil {
		log.Fatalln("Error during parse defaults", err)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	fmt.Println(viper.GetString("redis.URL"))
	connect, err := redis.DialURL(viper.GetString("redis.URL"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"infrastructure": "session",
			"action":         "connect to redis",
		}).Error(err)
	}
	session.RegisterAuthorizationServiceServer(server, persistence.NewSessionManager(connect))

	fmt.Println("starting server at :8081")
	if err := server.Serve(lis); err != nil {
		logrus.WithFields(logrus.Fields{
			"infrastructure": "session",
			"action":         "Serve",
		}).Error(err)
	}
}
