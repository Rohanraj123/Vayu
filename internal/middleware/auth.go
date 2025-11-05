package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/Rohanraj123/vayu/internal/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// authentication-middleware
func AuthMiddleware(cfg config.Config, next http.Handler, clientset *kubernetes.Clientset) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Auth.Enabled {
			next.ServeHTTP(w, r)
			return
		}

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

		// extract the key and compare it with request header
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			http.Error(w, "missing api-key", http.StatusUnauthorized)
			return
		}

		secret, _ := clientset.CoreV1().Secrets("vayu-system").Get(context.TODO(), "vayu-api-keys", metav1.GetOptions{})
		storedHashKey := string(secret.Data[route.Service])

		if !compareHash(storedHashKey, apiKey) {
			http.Error(w, "invalid api-key", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func compareHash(hashKey, apiKey string) bool {
	storedHashBytes, err := hex.DecodeString(hashKey)
	if err != nil {
		return false
	}

	newHashKey := sha256.Sum256([]byte(apiKey))

	return subtle.ConstantTimeCompare(newHashKey[:], storedHashBytes) == 1
}
