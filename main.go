package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rohanraj123/vayu/internal/config"
	"github.com/Rohanraj123/vayu/internal/router"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./api-gateway <config-file>")
	}

	configPath := os.Args[1]
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	mux := router.NewRouter(cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 API Gateway started on port %d", cfg.Server.Port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
