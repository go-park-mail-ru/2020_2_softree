package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"server/canal/pkg/infrastructure/config"
	"server/canal/pkg/infrastructure/logger"
	"server/canal/pkg/interfaces/router"
	"time"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("config", "c", "", "path to config file")
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
				"port":     8000,
				"domain":   "localhost",
				"secure":   false,
				"logLevel": "Info",
				"logFile":  "",
				"timeout":  10,
			},

			"session": map[string]interface{}{
				"ip":   "127.0.0.1",
				"port": 8001,
			},

			"profile": map[string]interface{}{
				"ip":   "127.0.0.1",
				"port": 8002,
			},

			"currency": map[string]interface{}{
				"ip":   "127.0.0.1",
				"port": 8003,
			},

			"CORS": map[string]interface{}{
				"allowedOrigins": []string{
					"http://localhost",
					"http://localhost:3000",
					"https://softree.group",
				},
				"allowedHeaders": []string{
					"If-Modified-Since",
					"Cache-Control",
					"Content-Type",
					"Range",
				},
				"allowedMethods": []string{
					"GET",
					"POST",
					"OPTIONS",
					"PUT",
					"PATCH",
					"DELETE",
				},
				"exposedHeaders": []string{
					"Content-Length",
					"Content-Range",
				},
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
	sessionConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", viper.GetString("session.ip"), viper.GetInt("session.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
		}).Fatalln("Can't connect to session grpc", err)
	}
	defer sessionConn.Close()

	profileConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", viper.GetString("profile.ip"), viper.GetInt("profile.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
		}).Fatalln("Can't connect to profile grpc", err)
	}
	defer profileConn.Close()

	currencyConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", viper.GetString("currency.ip"), viper.GetInt("currency.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
		}).Fatalln("Can't connect to currency grpc", err)
	}
	defer currencyConn.Close()

	userAuthenticate, userProfile, rateRates := router.CreateAppStructs(profileConn, sessionConn, currencyConn)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("server.ip"), viper.GetInt("server.port")),
		Handler:      router.NewRouter(userAuthenticate, userProfile, rateRates),
		WriteTimeout: time.Duration(viper.GetInt("server.timeout")) * time.Second,
		ReadTimeout:  time.Duration(viper.GetInt("server.timeout")) * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "main",
		}).Fatal("Server cannot start", err)
	}
}
