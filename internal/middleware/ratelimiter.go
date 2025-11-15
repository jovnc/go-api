package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"

	"go_api/internal/config"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return httprate.LimitByIP(config.GlobalConfig.RateLimit, time.Minute)(next)
}
