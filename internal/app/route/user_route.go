package route

import (
	"net/http"

	"go_api/internal/app/handler"
	"go_api/internal/middleware"
)

func SetupUserRoute(mux *http.ServeMux, userHandler *handler.UserHandler) {
	userMux := http.NewServeMux()

	userMux.HandleFunc("POST /register", userHandler.CreateUserHandler())
	userMux.HandleFunc("POST /login", userHandler.LoginUserHandler())
	userMux.Handle("GET /profile", middleware.AuthMiddleware(userHandler.UserProfileHandler()))
	userMux.Handle("POST /logout", middleware.AuthMiddleware(userHandler.LogoutUserHandler()))
	userMux.Handle("GET /", middleware.AuthMiddleware(userHandler.ListAllUsersHandler()))

	mux.Handle("/users/", http.StripPrefix("/users", userMux))
}
