package routes

import (
	"go_api/internal/handlers"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	SetupHealthRoute(mux, handler)
}
