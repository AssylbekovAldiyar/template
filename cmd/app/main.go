package main

import (
	"log"

	"github.com/AssylbekovAldiyar/memegen/internal/config"
	"github.com/AssylbekovAldiyar/memegen/internal/controller"
)

func main() {
	cfg := config.LoadConfig()
	memeController := controller.NewMemeController(cfg)

	log.Printf("Сервер запущен на порту %s", cfg.ServerPort)
	memeController.StartServer()
}
