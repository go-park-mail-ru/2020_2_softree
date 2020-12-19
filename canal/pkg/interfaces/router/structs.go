package router

import (
	"google.golang.org/grpc"
	sessionService "server/authorization/pkg/session/gen"
	"server/canal/pkg/application"
	"server/canal/pkg/infrastructure/security"
	"server/canal/pkg/interfaces/authorization"
	"server/canal/pkg/interfaces/profile"
	"server/canal/pkg/interfaces/rates"
	currencyService "server/currency/pkg/currency/gen"
	profileService "server/profile/pkg/profile/gen"
)

func createAuthenticate(profileConn, sessionConn *grpc.ClientConn) *authorization.Authentication {
	authService := sessionService.NewAuthorizationServiceClient(sessionConn)
	profileService := profileService.NewProfileServiceClient(profileConn)
	securityManager := security.CreateNewSecurityUtils()

	authApp := application.NewAuthApp(profileService, authService, securityManager)
	profileApp := application.NewProfileApp(profileService, securityManager)

	return authorization.NewAuthentication(profileApp, authApp)
}

func createProfile(profileConn, currencyConn *grpc.ClientConn) *profile.Profile {
	profileManager := profileService.NewProfileServiceClient(profileConn)
	currencyManager := currencyService.NewCurrencyServiceClient(currencyConn)
	securityManager := security.CreateNewSecurityUtils()

	profileApp := application.NewProfileApp(profileManager, securityManager)
	paymentApp := application.NewPaymentApp(profileManager, currencyManager, securityManager)

	return profile.NewProfile(profileApp, paymentApp)
}

func createRates(currencyConn *grpc.ClientConn) *rates.Rates {
	currencyManager := currencyService.NewCurrencyServiceClient(currencyConn)

	currencyApp := application.NewCurrencyApp(currencyManager)

	return rates.NewRates(currencyApp)
}

func CreateAppStructs(
	profileConn, sessionConn, currencyConn *grpc.ClientConn) (
	*authorization.Authentication, *profile.Profile, *rates.Rates) {
	userAuthenticate := createAuthenticate(profileConn, sessionConn)
	userProfile := createProfile(profileConn, currencyConn)
	rateRates := createRates(currencyConn)

	return userAuthenticate, userProfile, rateRates
}
