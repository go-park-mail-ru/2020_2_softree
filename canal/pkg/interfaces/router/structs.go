package router

import (
	"google.golang.org/grpc"
	serviceSession "server/authorization/pkg/session/gen"
	"server/canal/pkg/application"
	"server/canal/pkg/infrastructure/security"
	"server/canal/pkg/interfaces/authorization"
	"server/canal/pkg/interfaces/profile"
	"server/canal/pkg/interfaces/rates"
	serviceCurrency "server/currency/pkg/currency/gen"
	serviceProfile "server/profile/pkg/profile/gen"
)

func createAuthenticate(profileApp *application.ProfileApp, authApp *application.AuthApp) *authorization.Authentication {
	return authorization.NewAuthentication(profileApp, authApp)
}

func createProfile(profileApp *application.ProfileApp, paymentApp *application.PaymentApp) *profile.Profile {
	return profile.NewProfile(profileApp, paymentApp)
}

func createRates(currencyApp *application.CurrencyApp) *rates.Rates {
	return rates.NewRates(currencyApp)
}

func createApps(profileConn, sessionConn, currencyConn *grpc.ClientConn) (*application.CurrencyApp, *application.PaymentApp, *application.ProfileApp, *application.AuthApp) {
	profileManager := serviceProfile.NewProfileServiceClient(profileConn)
	currencyManager := serviceCurrency.NewCurrencyServiceClient(currencyConn)
	securityManager := security.CreateNewSecurityUtils()
	authManager := serviceSession.NewAuthorizationServiceClient(sessionConn)

	profileApp := application.NewProfileApp(profileManager, securityManager)
	paymentApp := application.NewPaymentApp(profileManager, currencyManager, securityManager)
	currencyApp := application.NewCurrencyApp(currencyManager)
	authApp := application.NewAuthApp(profileManager, authManager, securityManager)

	return currencyApp, paymentApp, profileApp, authApp
}

func CreateAppStructs(profileConn, sessionConn, currencyConn *grpc.ClientConn) (*authorization.Authentication, *profile.Profile, *rates.Rates) {
	currencyApp, paymentApp, profileApp, authApp := createApps(profileConn, sessionConn, currencyConn)

	userAuthenticate := createAuthenticate(profileApp, authApp)
	userProfile := createProfile(profileApp, paymentApp)
	rateRates := createRates(currencyApp)

	return userAuthenticate, userProfile, rateRates
}
