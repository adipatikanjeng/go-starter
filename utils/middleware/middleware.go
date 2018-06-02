package middleware

import (
	"fmt"
	"log"
	"net/http"
	"rest-api/utils/caching"
	"strings"
)

// Define our struct
type AuthenticationMiddleware struct {
	Cache caching.Cache
}

func AuthMiddleware(c caching.Cache) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		Cache: c,
	}
}

// Middleware function, which will be called for each request
func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		if path[1] == "auth" {
			// Pass down the request to the next middleware
			next.ServeHTTP(w, r)
		}
		token := r.Header.Get("token")
		userIDStr, err := amw.Cache.Get(fmt.Sprintf("token_%s", token))
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", userIDStr)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		}

	})
}
