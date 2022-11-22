package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/TsunamiProject/yamarkt/internal/handler"
)

func NewRouter(h *handler.RequestHandler) chi.Router {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		router.Use() //middleware
		router.Post("/api/user/orders", nil)
		router.Post("/apu/user/balance/withdraw", nil)
		router.Get("/api/user/orders", nil)
		router.Get("/api/user/balance", nil)
		router.Get("/api/user/withdrawals", nil)
	})

	router.Group(func(r chi.Router) {
		router.Post("/api/user/register", nil)
		router.Post("/api/user/login", nil)
	})

	return router

}
