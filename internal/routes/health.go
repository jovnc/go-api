package routes

import (
	"go_api/internal/handlers"
	"net/http"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("/health", handler.HealthHandler())
}
