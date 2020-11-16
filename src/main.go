package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/src/infrastructure/config"
	"server/src/infrastructure/logger"
	"server/src/interfaces/router"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("config", "c", "", "path to config file")
	pflag.BoolP("help", "h", false, "usage info")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if viper.GetString("config") == "" {
		fmt.Fprint(os.Stderr, "There is must explicitly specify the config file`\n")
		pflag.Usage()
		os.Exit(1)
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

			"postgres": map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     5432,
				"db":       "db",
				"user":     "user",
				"password": "",
			},

			"redis": map[string]interface{}{
				"host":         "127.0.0.1",
				"port":         6379,
				"sessionPath":  "/1",
				"currencyPath": "/2",
				"user":         "user",
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
			"finnhub-api": map[string]interface{}{
				"token": "",
			},
		}); err != nil {
		fmt.Fprint(os.Stderr, "Error during parse defaults \n", err, "\n")
		os.Exit(1)
	}

	logger.ConfigureLogger()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	userAuthenticate, userProfile, rateRates, err := router.CreateAppStructs()
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"function": "main",
		}).Error(err)
		return
	}

	go rateRates.GetRatesFromApi()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.GlobalConfig.GetString("server.ip"), config.GlobalConfig.GetInt("server.port")),
		Handler:      router.NewRouter(userAuthenticate, userProfile, rateRates),
		WriteTimeout: time.Duration(config.GlobalConfig.GetInt("server.timeout")) * time.Second,
		ReadTimeout:  time.Duration(config.GlobalConfig.GetInt("server.timeout")) * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"function": "main",
		}).Error("Server cannot start", err)
		os.Exit(1)
	}
	logger.GlobalLogger.Info("Server listening")
}
