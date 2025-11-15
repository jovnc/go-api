package route

import (
	"net/http"

	"go_api/internal/app/handler"
)

func SetupRoutes(mux *http.ServeMux, handler *handler.Handler) {
	SetupHealthRoute(mux, handler)
	SetupUserRoute(mux, handler)
	SetupBlogRoute(mux, handler)
}
