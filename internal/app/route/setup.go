package route

import (
	"net/http"

	"go_api/internal/app/handler"
	"go_api/internal/middleware"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(mux *http.ServeMux, db *gorm.DB, redis *redis.Client) http.Handler {

	// Create handlers
	blogHandler := handler.NewBlogHandler(db, redis)
	handler := handler.NewHandler(db, redis)

	// Register routes
	SetupHealthRoute(mux, handler)
	SetupUserRoute(mux, handler)
	SetupBlogRoute(mux, blogHandler)

	// Create middleware chain
	middlewares := []func(http.Handler) http.Handler{
		middleware.LoggerMiddleware,
		middleware.RecoveryMiddleware,
		middleware.RateLimiterMiddleware,
	}

	// Wrap the mux with all middlewares in order
	var wrappedHandler http.Handler = mux
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrappedHandler = middlewares[i](wrappedHandler)
	}

	return wrappedHandler
}
