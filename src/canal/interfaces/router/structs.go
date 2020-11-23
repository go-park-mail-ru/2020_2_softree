package router

import (
	"google.golang.org/grpc"
	sessionService "server/src/authorization/session/gen"
	"server/src/canal/interfaces/authorization"
	"server/src/canal/interfaces/profile"
	"server/src/canal/interfaces/rates"
	profileService "server/src/profile/profile/gen"
)

func createAuthenticate(profileConn, sessionConn *grpc.ClientConn) *authorization.Authentication {
	sessionManager := sessionService.NewAuthorizationServiceClient(sessionConn)
	profileManager := profileService.NewProfileServiceClient(profileConn)

	return authorization.NewAuthenticate(profileManager, sessionManager)
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
