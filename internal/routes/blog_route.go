package routes

import (
	"net/http"

	"go_api/internal/auth"
	"go_api/internal/handlers"
)

func SetupBlogRoute(mux *http.ServeMux, handler *handlers.Handler) {
	mux.Handle("POST /blogs/", auth.AuthMiddleware(handler.CreateBlogHandler()))
	mux.Handle("GET /blogs/{id}", handler.GetBlogHandler())
	mux.Handle("DELETE /blogs/{id}", auth.AuthMiddleware(handler.DeleteBlogHandler()))
	mux.Handle("GET /blogs/", handler.ListBlogsHandler())
}
