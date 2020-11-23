package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	session "server/src/authorizationService/session/gen"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	connect, err := redis.DialURL(viper.GetString("redis.sessionURL"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"service": "session",
			"action":  "connect to redis",
		}).Error(err)
	}
	session.RegisterAuthorizationServiceServer(server, NewSessionManager(connect))

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}
