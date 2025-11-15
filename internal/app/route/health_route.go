package route

import (
	"net/http"

	"go_api/internal/app/handler"
)

func SetupHealthRoute(mux *http.ServeMux, handler *handler.Handler) {
	mux.HandleFunc("/health", handler.HealthHandler())
}
