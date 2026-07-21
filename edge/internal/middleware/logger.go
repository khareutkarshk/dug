package middleware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id, _ := r.Context().Value(RequestIDKey).(string)

		log.Printf("[%s] %s %s", id, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
