package middleware

import (
	"go_api/internal/util"
	"log"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := "Caught panic: %v, Stack trace: %s"
				log.Printf(msg, err, string(debug.Stack()))
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.(error).Error())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
