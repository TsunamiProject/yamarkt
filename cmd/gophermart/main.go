package main

import (
	"log"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/handler"
	"github.com/TsunamiProject/yamarkt/internal/service"
	"github.com/TsunamiProject/yamarkt/internal/storage"
)

func main() {
	//Creating config instance
	log.Println("Initializing config")
	cfg := config.New()
	log.Print(cfg)

	pStorage, _ := storage.NewPostgresStorage(cfg.DatabaseDSN)
	userService := service.NewUserService(pStorage)
	userHandler := handler.NewUserHandler(userService)
	log.Println(userHandler)
	balanceService := service.NewBalanceService(pStorage)
	balanceHandler := handler.NewBalanceHandler(balanceService)
	log.Println(balanceHandler)
	orderService := service.NewOrderService(pStorage, cfg.AccrualURL)
	orderHandler := handler.NewOrderHandler(orderService)
	log.Println(orderHandler)
}
