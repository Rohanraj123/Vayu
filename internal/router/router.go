package router

import (
	"log"
	"net/http"

	"github.com/Rohanraj123/vayu/internal/apis"
	"github.com/Rohanraj123/vayu/internal/config"
	"github.com/Rohanraj123/vayu/internal/proxy"
	"k8s.io/client-go/kubernetes"
)

func NewRouter(cfg *config.Config, clientset *kubernetes.Clientset) *http.ServeMux {
	mux := http.NewServeMux()

	// add handlers
	mux.HandleFunc("/api-keys", func(w http.ResponseWriter, r *http.Request) {
		apis.CreateApiKeyHandler(w, r, clientset)
	})
	mux.HandleFunc("/healtz", func(w http.ResponseWriter, r *http.Request) {
		apis.HealthzHandler(w, r)
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		apis.ReadyzHandler(w, r)
	})

	// routes in config file
	for _, route := range cfg.Routes {
		handler, err := proxy.ProxyHandler(route.Upstream)
		if err != nil {
			log.Fatalf("Failed to set up route %s: %v", route.Path, err)
		}

		log.Printf("✅ Route registered: %s --> %s", route.Path, route.Upstream)
		mux.HandleFunc(route.Path, handler)
	}

	return mux
}
