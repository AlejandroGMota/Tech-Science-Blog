package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlejandroGMota/Tech-Science-Blog/backend/api"
	"github.com/AlejandroGMota/Tech-Science-Blog/backend/internal/config"
)

func main() {
	cfg := config.Load()

	router := api.NewRouter(cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Tech Blog API running on %s", addr)
	log.Printf("Storage: %s", cfg.DBType)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
