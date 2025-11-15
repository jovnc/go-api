package route

import (
	"net/http"

	"go_api/internal/app/handler"
	"go_api/internal/middleware"
)

func SetupRoutes(mux *http.ServeMux, handler *handler.Handler) http.Handler {
	// Register routes
	SetupHealthRoute(mux, handler)
	SetupUserRoute(mux, handler)
	SetupBlogRoute(mux, handler)

	// Create middleware chain
	middlewares := []func(http.Handler) http.Handler{
		middleware.LoggerMiddleware,
		middleware.RecoveryMiddleware,
	}

	// Wrap the mux with all middlewares in order
	var wrappedHandler http.Handler = mux
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrappedHandler = middlewares[i](wrappedHandler)
	}

	return wrappedHandler
}
