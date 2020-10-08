package corsInteraction

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"server/src/infrastructure/config"
)

func enableCORS(cfg *config.CORSConfig, handler http.Handler) http.Handler {
	var (
		allowedOrigins = handlers.AllowedOrigins(cfg.AllowedOrigins)
		allowedHeaders = handlers.AllowedHeaders(cfg.AllowedHeaders)
		exposedHeaders = handlers.ExposedHeaders(cfg.ExposedHeaders)
		allowedMethods = handlers.AllowedMethods(cfg.AllowedMethods)
		credentials    = handlers.AllowCredentials()
	)

	return handlers.CORS(allowedOrigins, allowedHeaders, exposedHeaders, allowedMethods, credentials)(handler)
}

func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return enableCORS(&config.GlobalCORSConfig, next)
	}
}
