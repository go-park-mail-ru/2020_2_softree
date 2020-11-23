package main

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	profile "server/src/profileService/profile/gen"
)

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	db, err := sql.Open("postgres", viper.GetString("postgres.URL"))
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	db.SetMaxOpenConns(10)

	profile.RegisterProfileServiceServer(server, NewUserDBManager(db))

	fmt.Println("starting server at :8082")
	server.Serve(lis)
}
