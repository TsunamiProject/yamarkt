package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/TsunamiProject/yamarkt/internal/config"
	"github.com/TsunamiProject/yamarkt/internal/handler"
	"github.com/TsunamiProject/yamarkt/internal/middleware"
)

func NewRouter(uh *handler.UserHandler, bh *handler.BalanceHandler, oh *handler.OrderHandler) chi.Router {
	router := chi.NewRouter()
	router.Use(chiMiddleware.Recoverer)
	//router.Use(chiMiddleware.Logger)
	router.Use(middleware.GzipRespWriter, middleware.GzipReqReader)

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		router.Post("/api/user/orders", oh.CreateOrder)
		router.Post("/api/user/balance/withdraw", bh.CreateWithdrawal)
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
