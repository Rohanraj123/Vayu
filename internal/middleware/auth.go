package middleware

import (
	"net/http"

	"github.com/Rohanraj123/vayu/internal/config"
)

func AuthMiddleware(cfg config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find route matching this path
		var route *config.RouteConfig
		for _, rt := range cfg.Routes {
			if rt.Path == r.URL.Path {
				route = &rt
				break
			}
		}

		// If route not found or auth not required, skip auth
		if route == nil || !route.AuthRequired {
			next.ServeHTTP(w, r)
			return
		}

		// Auth required, check API
		apiKey := r.Header.Get(cfg.Auth.ApiKeyHeader)
		if apiKey == "" || !contains(cfg.Auth.ValidKeys, apiKey) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
