package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/TsunamiProject/yamarkt/internal/handler"
)

func NewRouter(uh *handler.UserHandler) chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		router.Use() //middleware
		router.Post("/api/user/orders", nil)
		router.Post("/apu/user/balance/withdrawals", nil)
		router.Get("/api/user/orders", nil)
		router.Get("/api/user/balance", nil)
		router.Get("/api/user/withdrawals", nil)
	})

	router.Group(func(r chi.Router) {
		router.Post("/api/user/register", uh.Register)
		router.Post("/api/user/login", uh.Auth)
	})

	return router
}
