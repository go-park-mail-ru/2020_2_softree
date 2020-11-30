package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"server/canal/pkg/infrastructure/config"
	"server/canal/pkg/infrastructure/logger"
	"server/profile/pkg/infrastructure/persistence"
	profile "server/profile/pkg/profile/gen"

	_ "github.com/lib/pq"
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
				"port":     8002,
				"logLevel": "Info",
			},

			"postgres": map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     5432,
				"db":       "db",
				"user":     "user",
				"password": "",
			},
		}); err != nil {
		log.Fatalln("Error during parse defaults", err)
	}

	if err := logger.ConfigureLogger(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.db"),
	))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
		}).Fatalln("Can't connecting to postgres", err)
	}
	err = db.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
		}).Fatalln("Postgres is unavailable", err)
	}
	db.SetMaxOpenConns(10)

	server := grpc.NewServer()

	profile.RegisterProfileServiceServer(server, persistence.NewUserDBManager(db))

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
			"infrastructure": "session",
			"action":         "Serve",
		}).Fatalln(err)
	}
}
