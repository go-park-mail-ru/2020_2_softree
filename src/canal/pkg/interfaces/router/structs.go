package router

import (
	"google.golang.org/grpc"
	sessionService "server/src/authorization/pkg/session/gen"
	"server/src/canal/pkg/infrastructure/security"
	"server/src/canal/pkg/interfaces/authorization"
	"server/src/canal/pkg/interfaces/profile"
	"server/src/canal/pkg/interfaces/rates"
	profileService "server/src/profile/pkg/profile/gen"
)

func createAuthenticate(profileConn, sessionConn *grpc.ClientConn) *authorization.Authentication {
	sessionManager := sessionService.NewAuthorizationServiceClient(sessionConn)
	profileManager := profileService.NewProfileServiceClient(profileConn)

	return authorization.NewAuthenticate(profileManager, sessionManager, security.CreateNewSecurityUtils())
}

func createProfile(profileConn, sessionConn, currencyConn *grpc.ClientConn) *profile.Profile {
	sessionManager := sessionService.NewAuthorizationServiceClient(sessionConn)
	profileManager := profileService.NewProfileServiceClient(profileConn)
	currencyManager := //

	return profile.NewProfile(profileManager, sessionManager, currencyManager)
}

func createRates(currencyConn *grpc.ClientConn) *rates.Rates {
	currencyManager := //

	return rates.NewRates(currencyManager)
}

func CreateAppStructs(
	profileConn, sessionConn, currencyConn *grpc.ClientConn) (
	*authorization.Authentication, *profile.Profile, *rates.Rates) {
	userAuthenticate := createAuthenticate(profileConn, sessionConn)
	userProfile := createProfile(profileConn, sessionConn, currencyConn)
	rateRates := createRates(currencyConn)

	return userAuthenticate, userProfile, rateRates
}
