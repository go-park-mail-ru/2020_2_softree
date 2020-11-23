package router

import (
	"google.golang.org/grpc"
	"log"
	session "server/src/authorization/session/gen"
	"server/src/canal/application"
	"server/src/canal/infrastructure/auth"
	"server/src/canal/infrastructure/financial"
	"server/src/canal/infrastructure/persistence"
	"server/src/canal/interfaces/authorization"
	"server/src/canal/interfaces/rates"
	profile "server/src/profile/profile/gen"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func createAuthenticate(profileConn, sessionConn *grpc.ClientConn) *authorization.Authentication {
	sessionManager := session.NewAuthorizationServiceClient(sessionConn)
	profileManager := profile.NewProfileServiceClient(profileConn)

	return authorization.NewAuthenticate(profileManager, sessionManager)
}

func createProfile() (*profile.Profile, error) {
	dbRepo, err := persistence.NewUserDBManager()
	if err != nil {
		return nil, err
	}
	dbHistory, err := persistence.NewPaymentDBManager()
	if err != nil {
		return nil, err
	}
	dbWallet, err := persistence.NewWalletDBManager()
	if err != nil {
		return nil, err
	}
	dbRate, err := persistence.NewRateDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(viper.GetString("redis.sessionURL"))
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	connectRedis, err := redis.DialURL(viper.GetString("redis.currencyURL"))
	if err != nil {
		return nil, err
	}
	dbCurr := financial.NewCurrencyManager(connectRedis)

	servicesDB := application.NewUserApp(dbRepo, dbHistory, dbWallet)
	servicesAuth := application.NewUserAuth(dbAuth)
	servicesRate := application.NewRateApp(dbRate, dbCurr)

	return profile.NewProfile(*servicesDB, *servicesAuth, *servicesRate), nil
}

func createRates() (*rates.Rates, error) {
	dbRepo, err := persistence.NewRateDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(viper.GetString("redis.currencyURL"))
	if err != nil {
		return nil, err
	}
	dbAuth := financial.NewCurrencyManager(connect)

	servicesDB := application.NewRateApp(dbRepo, dbAuth)

	return rates.NewRates(*servicesDB), nil
}

func CreateAppStructs() (*authorization.Authentication, *profile.Profile, *rates.Rates, error) {
	userAuthenticate, err := createAuthenticate()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createAuthenticate",
		}).Error(err)
		return nil, nil, nil, err
	}

	userProfile, err := createProfile()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createProfile",
		}).Error(err)
		return nil, nil, nil, err
	}

	rateRates, err := createRates()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createRates",
		}).Error(err)
		return nil, nil, nil, err
	}

	return userAuthenticate, userProfile, rateRates, nil
}
