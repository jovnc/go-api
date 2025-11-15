package route

import (
	"net/http"

	"go_api/internal/app/handler"
	"go_api/internal/middleware"
)

func SetupBlogRoute(mux *http.ServeMux, handler *handler.Handler) {
	mux.Handle("POST /blogs/", middleware.AuthMiddleware(handler.CreateBlogHandler()))
	mux.Handle("GET /blogs/{id}", handler.GetBlogHandler())
	mux.Handle("DELETE /blogs/{id}", middleware.AuthMiddleware(handler.DeleteBlogHandler()))
	mux.Handle("GET /blogs/", handler.ListBlogsHandler())
}
