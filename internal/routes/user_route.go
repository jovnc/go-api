package routes

import (
	"net/http"

	"go_api/internal/handlers"
)

func SetupUserRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("/user/register", handler.CreateUserHandler())
}