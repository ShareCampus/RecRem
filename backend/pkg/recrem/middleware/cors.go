package middleware

import (
	"net/http"

	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
	"github.com/rs/cors"
)

// AllowCors is a middleware to allow CORS
func AllowCors(next http.Handler) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins:   GetCORSSettings(),
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		//Only GET, POST, and HEAD requests are allowed by default, with the addition of PUT and DELETE requests
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodPut, http.MethodDelete},
	}
	logger.Infof("Allowed CORS origins: %v", corsOptions.AllowedOrigins)

	cors := cors.New(corsOptions)
	return cors.Handler(next)
}
