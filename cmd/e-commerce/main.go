package main

import (
	"e-commerce-shop/internal/config"
	"e-commerce-shop/internal/storage/postgres"
	"fmt"
	"net/http"
)

func main() {
	cfg := config.MustLoad()

	db := postgres.Database(cfg)
	defer db.Close()

	err := http.ListenAndServe(cfg.ServerPort, nil)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err)
	}
}
