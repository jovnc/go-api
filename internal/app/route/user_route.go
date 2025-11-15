package route

import (
	"net/http"

	"go_api/internal/app/handler"
	"go_api/internal/middleware"
)

func SetupUserRoute(mux *http.ServeMux, handler *handler.Handler) {
	userMux := http.NewServeMux()

	userMux.HandleFunc("POST /register", handler.CreateUserHandler())
	userMux.HandleFunc("POST /login", handler.LoginUserHandler())
	userMux.Handle("GET /profile", middleware.AuthMiddleware(handler.UserProfileHandler()))
	userMux.Handle("POST /logout", middleware.AuthMiddleware(handler.LogoutUserHandler()))
	userMux.Handle("GET /", middleware.AuthMiddleware(handler.ListAllUsersHandler()))

	mux.Handle("/users/", http.StripPrefix("/users", userMux))
}
