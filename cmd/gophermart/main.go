package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

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
		log.Fatalf("error while initializing posgtres storage: %s", err)
	}

	userService := service.NewUserService(pStorage)
	userHandler := handler.NewUserHandler(userService)

	balanceService := service.NewBalanceService(pStorage)
	balanceHandler := handler.NewBalanceHandler(balanceService)

	orderService := service.NewOrderService(pStorage, cfg.AccrualURL)
	orderHandler := handler.NewOrderHandler(orderService)

	router := appRouter.NewRouter(userHandler, balanceHandler, orderHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	var wg sync.WaitGroup
	httpServer := &http.Server{Addr: cfg.ServerAddress, Handler: router}
	wg.Add(1)
	go gracefulShutdown(ctx, &wg, httpServer)
	err = httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("listen and serve error: %s", err)
	}
	wg.Wait()
	stop()
}

func gracefulShutdown(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
	<-ctx.Done()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("error while shutting down http server: %s", err)
	}
	wg.Done()
}
