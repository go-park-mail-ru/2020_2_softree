package router

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/config"
	"server/src/infrastructure/financial"
	"server/src/infrastructure/log"
	"server/src/infrastructure/persistence"
	"server/src/infrastructure/userInteraction"
	"server/src/interfaces/authorization"
	"server/src/interfaces/profile"
	"server/src/interfaces/rates"
)

func createAuthenticate() (*authorization.Authentication, error) {
	dbRepo, err := userInteraction.NewUserDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(config.SessionDatabaseConfig.AddressSessions)
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	servicesDB := application.NewUserApp(dbRepo)
	servicesAuth := application.NewUserAuth(dbAuth)
	servicesLog := log.NewLogrusLogger()

	return authorization.NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), nil
}

func createProfile() (*profile.Profile, error) {
	dbRepo, err := userInteraction.NewUserDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(config.SessionDatabaseConfig.AddressSessions)
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	servicesDB := application.NewUserApp(dbRepo)
	servicesAuth := application.NewUserAuth(dbAuth)
	servicesLog := log.NewLogrusLogger()

	return profile.NewProfile(*servicesDB, *servicesAuth, servicesLog), nil
}

func createRates() (*rates.Rates, error) {
	dbRepo, err := persistence.NewRateDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(config.SessionDatabaseConfig.AddressDayCurrency)
	if err != nil {
		return nil, err
	}
	dbAuth := financial.NewCurrencyManager(connect)

	servicesDB := application.NewRateApp(dbRepo, dbAuth)
	servicesLog := log.NewLogrusLogger()

	return rates.NewRates(*servicesDB, servicesLog), nil
}

func CreateAppStructs() (*authorization.Authentication, *profile.Profile, *rates.Rates, error) {
	userAuthenticate, err := createAuthenticate()
	if err != nil {
		fmt.Println("userAuthenticate", err)
		return nil, nil, nil, err
	}

	userProfile, err := createProfile()
	if err != nil {
		fmt.Println("userProfile", err)
		return nil, nil, nil, err
	}

	rateRates, err := createRates()
	if err != nil {
		fmt.Println("userProfile", err)
		return nil, nil, nil, err
	}

	return userAuthenticate, userProfile, rateRates, nil
}
