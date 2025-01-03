package main

import (
	"e-commerce-shop/internal/config"
	"e-commerce-shop/internal/http-server/router"
	"e-commerce-shop/internal/storage/postgres"
	"fmt"
	"net/http"
)

func main() {
	cfg := config.MustLoad()

	db := postgres.Database(cfg)
	defer db.Close()

	r := router.Router(db)

	err := http.ListenAndServe(cfg.ServerPort, r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err)
	}
}
