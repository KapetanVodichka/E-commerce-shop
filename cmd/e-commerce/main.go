package main

import (
	"e-commerce-shop/internal/config"
	"fmt"
	"net/http"
)

func main() {
	cfg := config.MustLoad()

	err := http.ListenAndServe(cfg.ServerPort, nil)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err)
	}
}
