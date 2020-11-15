package corsInteraction

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"server/src/infrastructure/config"
)

func enableCORS(handler http.Handler) http.Handler {
	var (
		allowedOrigins = handlers.AllowedOrigins(config.GlobalConfig.GetStringSlice("CORS.allowedOrigins"))
		allowedHeaders = handlers.AllowedHeaders(config.GlobalConfig.GetStringSlice("CORS.allowedHeaders"))
		exposedHeaders = handlers.ExposedHeaders(config.GlobalConfig.GetStringSlice("CORS.exposedHeaders"))
		allowedMethods = handlers.AllowedMethods(config.GlobalConfig.GetStringSlice("CORS.allowedMethods"))
		credentials    = handlers.AllowCredentials()
	)

	return handlers.CORS(allowedOrigins, allowedHeaders, exposedHeaders, allowedMethods, credentials)(handler)
}

func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return enableCORS(next)
	}
}
