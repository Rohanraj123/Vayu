package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rohanraj123/vayu/internal/config"
	"github.com/Rohanraj123/vayu/internal/middleware"
	"github.com/Rohanraj123/vayu/internal/router"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./api-gateway <config-file>")
	}

	// builds config from service account credentials mounted in pod
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("failed to get incluster-config: %v", err)
	}

	// creates a clientset - an object gives access to different k8s resources
	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalf("failed to create k8s clientset: %v", err)
	}

	configPath := os.Args[1]
	cfg, err := config.LoadConfig(configPath)

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// create a new mux
	mux := router.NewRouter(cfg, clientset)

	// adding different layers of middlewares
	handler := middleware.LoggingMiddleware(
		middleware.AuthMiddleware(
			*cfg, mux, clientset))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 API Gateway started on port %d", cfg.Server.Port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
