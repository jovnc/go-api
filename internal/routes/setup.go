package routes

import (
	"net/http"

	"go_api/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
	SetupUserRoute(mux, handler)
}
