package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/TsunamiProject/yamarkt/internal/handler"
	"github.com/TsunamiProject/yamarkt/internal/middleware"
)

func NewRouter(uh *handler.UserHandler, bh *handler.BalanceHandler, oh *handler.OrderHandler) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.GzipRespWriter, middleware.GzipReqReader)

	router.Group(func(r chi.Router) {
		router.Post("/api/user/orders", oh.CreateOrder)
		router.Post("/apu/user/balance/withdrawals", bh.NewWithdrawal)
		router.Get("/api/user/orders", oh.OrderList)
		router.Get("/api/user/balance", bh.GetCurrentBalance)
		router.Get("/api/user/withdrawals", bh.GetWithdrawalList)
	})

	router.Group(func(r chi.Router) {
		router.Post("/api/user/register", uh.Register)
		router.Post("/api/user/login", uh.Auth)
	})

	return router
}
