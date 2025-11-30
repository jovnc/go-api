package route

import (
	"net/http"

	"go_api/internal/app/handler"
	"go_api/internal/middleware"
)

func SetupBlogRoute(mux *http.ServeMux, blogHandler *handler.BlogHandler) {
	mux.Handle("POST /blogs/", middleware.AuthMiddleware(blogHandler.CreateBlogHandler()))
	mux.Handle("GET /blogs/{id}", blogHandler.GetBlogHandler())
	mux.Handle("DELETE /blogs/{id}", middleware.AuthMiddleware(blogHandler.DeleteBlogHandler()))
	mux.Handle("GET /blogs/", blogHandler.ListBlogsHandler())
}
