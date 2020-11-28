package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"os"
	"server/canal/pkg/infrastructure/config"
	"server/canal/pkg/infrastructure/logger"
	"server/canal/pkg/interfaces/router"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("viper", "c", "", "path to viper file")
	pflag.BoolP("help", "h", false, "usage info")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if err := config.ParseConfig(
		viper.GetString("viper"),
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
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
}

func main() {
	sessionConn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cant connect to session grpc")
	}
	profileConn, err := grpc.Dial("127.0.0.1:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cant connect to profile grpc")
	}
	currencyConn, err := grpc.Dial("127.0.0.1:8083", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cant connect to currency grpc")
	}

	defer sessionConn.Close()
	defer profileConn.Close()
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
			"function": "canal",
		}).Fatal("Server cannot start", err)
	}
}
