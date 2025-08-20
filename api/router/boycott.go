package router

import (
	"canada-hires/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func BoycottRoutes(boycottController controllers.BoycottController, authMW func(http.Handler) http.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/boycotts", func(r chi.Router) {
			// Public routes
			r.Get("/top", boycottController.GetTopBoycotted)

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(authMW)
				r.Get("/stats", boycottController.GetBoycottStats)
				r.Post("/toggle", boycottController.ToggleBoycott)
				r.Get("/my", boycottController.GetUserBoycotts)
			})
		})
	}
}
