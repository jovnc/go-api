package routes

import (
	"net/http"

	"go_api/internal/handlers"
)

func SetupUserRoute(mux *http.ServeMux, handler *handlers.Handler) {
	userMux := http.NewServeMux()

	userMux.HandleFunc("POST /register", handler.CreateUserHandler())
	userMux.HandleFunc("POST /login", handler.LoginUserHandler())

	mux.Handle("/users/", http.StripPrefix("/users", userMux))
}