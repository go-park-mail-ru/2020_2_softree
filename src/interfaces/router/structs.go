package router

import (
	"fmt"
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/config"
	"server/src/infrastructure/financial"
	"server/src/infrastructure/persistence"
	"server/src/interfaces/authorization"
	"server/src/interfaces/profile"
	"server/src/interfaces/rates"

	"github.com/gomodule/redigo/redis"
)

func createAuthenticate() (*authorization.Authentication, error) {
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
	connect, err := redis.DialURL(config.GlobalConfig.GetString("redis.sessionURL"))
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	servicesDB := application.NewUserApp(dbRepo, dbHistory, dbWallet)
	servicesAuth := application.NewUserAuth(dbAuth)

	return authorization.NewAuthenticate(*servicesDB, *servicesAuth), nil
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

	connect, err := redis.DialURL(config.GlobalConfig.GetString("redis.sessionURL"))
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	connectRedis, err := redis.DialURL(config.GlobalConfig.GetString("redis.currencyURL"))
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

	connect, err := redis.DialURL(config.GlobalConfig.GetString("redis.currencyURL"))
	if err != nil {
		println("4")
		return nil, err
	}
	dbAuth := financial.NewCurrencyManager(connect)

	servicesDB := application.NewRateApp(dbRepo, dbAuth)

	return rates.NewRates(*servicesDB), nil
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
