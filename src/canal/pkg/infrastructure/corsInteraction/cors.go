package corsInteraction

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func enableCORS(handler http.Handler) http.Handler {
	var (
		allowedOrigins = handlers.AllowedOrigins(viper.GetStringSlice("CORS.allowedOrigins"))
		allowedHeaders = handlers.AllowedHeaders(viper.GetStringSlice("CORS.allowedHeaders"))
		exposedHeaders = handlers.ExposedHeaders(viper.GetStringSlice("CORS.exposedHeaders"))
		allowedMethods = handlers.AllowedMethods(viper.GetStringSlice("CORS.allowedMethods"))
		credentials    = handlers.AllowCredentials()
	)

	return handlers.CORS(allowedOrigins, allowedHeaders, exposedHeaders, allowedMethods, credentials)(handler)
}

func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return enableCORS(next)
	}
}
