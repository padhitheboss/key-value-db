package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/padhitheboss/key-value-db/pkg/controller"
)

func RegisterRoute(r chi.Router) {
	r.Use(middleware.Logger)
	r.Post("/command", controller.CommandHandler)
	r.Get("/command", controller.CommandHandler)
}
