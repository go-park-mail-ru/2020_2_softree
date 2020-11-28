package router

import (
	"google.golang.org/grpc"
	sessionService "server/authorization/pkg/session/gen"
	"server/canal/pkg/infrastructure/security"
	"server/canal/pkg/interfaces/authorization"
	"server/canal/pkg/interfaces/profile"
	"server/canal/pkg/interfaces/rates"
	currencyService "server/currency/pkg/currency/gen"
	profileService "server/profile/pkg/profile/gen"
)

func createAuthenticate(profileConn, sessionConn *grpc.ClientConn) *authorization.Authentication {
	sessionManager := sessionService.NewAuthorizationServiceClient(sessionConn)
	profileManager := profileService.NewProfileServiceClient(profileConn)

	return authorization.NewAuthenticate(profileManager, sessionManager, security.CreateNewSecurityUtils())
}

func createProfile(profileConn, currencyConn *grpc.ClientConn) *profile.Profile {
	profileManager := profileService.NewProfileServiceClient(profileConn)
	currencyManager := currencyService.NewCurrencyServiceClient(currencyConn)
	securityManager := security.CreateNewSecurityUtils()

	return profile.NewProfile(profileManager, securityManager, currencyManager)
}

func createRates(currencyConn *grpc.ClientConn) *rates.Rates {
	currencyManager := currencyService.NewCurrencyServiceClient(currencyConn)

	return rates.NewRates(currencyManager)
}

func CreateAppStructs(
	profileConn, sessionConn, currencyConn *grpc.ClientConn) (
	*authorization.Authentication, *profile.Profile, *rates.Rates) {
	userAuthenticate := createAuthenticate(profileConn, sessionConn)
	userProfile := createProfile(profileConn, currencyConn)
	rateRates := createRates(currencyConn)

	return userAuthenticate, userProfile, rateRates
}
