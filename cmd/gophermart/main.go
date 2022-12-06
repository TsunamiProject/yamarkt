package main

import (
	"log"
	"net/http"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/handler"
	appRouter "github.com/TsunamiProject/yamarkt/internal/router"
	"github.com/TsunamiProject/yamarkt/internal/service"
	"github.com/TsunamiProject/yamarkt/internal/storage"
)

func main() {
	//Creating config instance
	log.Println("Initializing config")
	cfg := config.New()
	log.Print(cfg)

	pStorage, err := storage.NewPostgresStorage(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	userService := service.NewUserService(pStorage)
	userHandler := handler.NewUserHandler(userService)
	log.Println(userHandler)

	balanceService := service.NewBalanceService(pStorage)
	balanceHandler := handler.NewBalanceHandler(balanceService)
	log.Println(balanceHandler)

	orderService := service.NewOrderService(pStorage, cfg.AccrualURL)
	orderHandler := handler.NewOrderHandler(orderService)
	log.Println(orderHandler)

	router := appRouter.NewRouter(userHandler, balanceHandler, orderHandler)

	httpServer := &http.Server{Addr: cfg.ServerAddress, Handler: router}
	log.Fatal(httpServer.ListenAndServe())

}
