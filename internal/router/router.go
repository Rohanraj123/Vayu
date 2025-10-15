package router

import (
	"log"
	"net/http"

	"github.com/Rohanraj123/vayu/internal/config"
	"github.com/Rohanraj123/vayu/internal/proxy"
)

func NewRouter(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

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
